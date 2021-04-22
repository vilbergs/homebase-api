package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Stat struct {
	ZoneId       string `json:"ZoneId"`
	Temp       float64 `json:"Temp"`
	Humidity       float64 `json:"Humidity"`
}

const token = "bmLbrO-AGJjXMrlYV1UbYqdJWctdd-yz8RjiouLNlqINXAjYH3QtB0jfrW3G4C0nOz_hubbSl2n_Q25R-YLW_Q=="
const bucket = "homebase"
const org = "Homebase"

var influxClient influxdb2.Client

func createTemperature(w http.ResponseWriter, r *http.Request) {
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

func main() {
  router := mux.NewRouter().StrictSlash(true)
  influxClient = influxdb2.NewClient("http://localhost:8086", token)

	router.HandleFunc("/", homeLink)
	router.HandleFunc("/zones", createTemperature).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
  // You can generate a Token from the "Tokens Tab" in the UI
  
	// get non-blocking write client

  // always close client at the end
  defer influxClient.Close()
}