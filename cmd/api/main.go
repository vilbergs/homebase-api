package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/vilbergs/homebase-api/db"
	"github.com/vilbergs/homebase-api/handlers"
)

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

func main() {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
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

	postgresErr := db.PostgresInit()
	if postgresErr != nil {
		log.Fatalf("Could not set up Postgres: %v", postgresErr)
	}

	db.InfluxInit()

	envType := os.Getenv("API_ENV")

	if envType == "DEV" {
		router.Handle("/telemetry/{zoneId:[0-9]+}", handlers.AddTelemetry).Methods("POST")
		router.Handle("/zones", handlers.AddZone).Methods("POST")
		router.Handle("/zones", handlers.GetALLZones).Methods("GET")
		router.Handle("/zones/{zoneId:[0-9]+}", handlers.GetZone).Methods("GET")
	} else {
		router.Handle("/telemetry/{zoneId:[0-9]+}", jwtMiddleware.Handler(handlers.AddTelemetry)).Methods("POST")
		router.Handle("/zones", jwtMiddleware.Handler(handlers.AddZone)).Methods("POST")
		router.Handle("/zones", jwtMiddleware.Handler(handlers.GetALLZones)).Methods("GET")
		router.Handle("/zones/{zoneId:[0-9]+}", jwtMiddleware.Handler(handlers.GetZone)).Methods("GET")
	}

	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	log.Fatal(http.ListenAndServe(":8080", corsWrapper.Handler(router)))

	// always close client at the end
	defer db.Postgres.Close()
	defer db.Influx.Close()
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

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}
