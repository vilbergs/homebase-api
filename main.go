package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/rs/cors"
)

type Zone struct {
	Name    	string  `json:"name"`
	Temp      float64 `json:"temperature"`
	Humidity	float64 `json:"humidity"`
}

type Stat struct {
	ZoneId    string  `json:"zoneId"`
	Temp      float64 `json:"temperature"`
	Humidity	float64 `json:"humidity"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

const token = "bmLbrO-AGJjXMrlYV1UbYqdJWctdd-yz8RjiouLNlqINXAjYH3QtB0jfrW3G4C0nOz_hubbSl2n_Q25R-YLW_Q=="
const bucket = "homebase"
const org = "Homebase"

var influxClient influxdb2.Client
var zones = []Zone{}


func createStat(w http.ResponseWriter, r *http.Request) {
  var newEvent Stat

	reqBody, readError := ioutil.ReadAll(r.Body)
	if readError != nil {
		fmt.Print("Error!")
	}

  json.Unmarshal(reqBody, &newEvent)

  writeAPI := influxClient.WriteAPI(org, bucket)

// write line protocol
writeAPI.WriteRecord(fmt.Sprintf("stat,zone=%s temperature=%f,humidity=%f", newEvent.ZoneId, newEvent.Temp, newEvent.Humidity))
// Flush writes
writeAPI.Flush()


query := fmt.Sprintf("from(bucket:\"%v\")|> range(start: -1h) |> filter(fn: (r) => r._measurement == \"stat\")", bucket)
// Get query client
queryAPI := influxClient.QueryAPI(org)
// get QueryTableResult
result, err := queryAPI.Query(context.Background(), query)
if err == nil {
  // Iterate over query response
 
  w.WriteHeader(http.StatusCreated)


  json.NewEncoder(w).Encode(newEvent)

  // check for an error
  if result.Err() != nil {
    fmt.Printf("query parsing error: %s\n", result.Err().Error())
  }
} else {
  panic(err)
}
	

}

var ZonePostHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var newZone Zone

	reqBody, readError := ioutil.ReadAll(r.Body)
	if readError != nil {
		fmt.Print("Error!")
	}

	json.Unmarshal(reqBody, &newZone)

	zones = append(zones, newZone)

	w.WriteHeader(http.StatusCreated)


  json.NewEncoder(w).Encode(newZone)
})

var ZoneGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {



  json.NewEncoder(w).Encode(zones)
})

func main() {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options {
    ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
      // Verify 'aud' claim
      aud := "https://api.homebase-app.com"
      checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
      if !checkAud {
          return token, errors.New("invalid audience")
      }
      // Verify 'iss' claim
      iss := "https://homebase-app.eu.auth0.com/"
      checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
      if !checkIss {
          return token, errors.New("invalid issuer")
      }

      cert, err := getPemCert(token)
      if err != nil {
          panic(err.Error())
      }

      result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
      return result, nil
    },
    SigningMethod: jwt.SigningMethodRS256,
  })


  router := mux.NewRouter().StrictSlash(true)
  influxClient = influxdb2.NewClient("http://localhost:8086", token)

	router.HandleFunc("/stat", createStat).Methods("POST")
	router.Handle("/zones", jwtMiddleware.Handler(ZonePostHandler)).Methods("POST")
	router.Handle("/zones", jwtMiddleware.Handler(ZoneGetHandler)).Methods("GET")

	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
})

	log.Fatal(http.ListenAndServe(":8080", corsWrapper.Handler(router)))
  
  // always close client at the end
  defer influxClient.Close()
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://homebase-app.eu.auth0.com/.well-known/jwks.json")

	if err != nil {
			return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
			return cert, err
	}

	for k, _ := range jwks.Keys {
			if token.Header["kid"] == jwks.Keys[k].Kid {
					cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
			}
	}

	if cert == "" {
			err := errors.New("Unable to find appropriate key.")
			return cert, err
	}

	return cert, nil
}