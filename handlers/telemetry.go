package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/vilbergs/homebase-api/db"
	"github.com/vilbergs/homebase-api/models"
)

var AddTelemetry = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var newTelemetry models.Telemetry

	reqBody, readError := ioutil.ReadAll(r.Body)
	if readError != nil {
		fmt.Print("Error!")
	}

	json.Unmarshal(reqBody, &newTelemetry)
	db.AddTelemetry(&newTelemetry)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encodedTelemetry, _ := json.Marshal(newTelemetry)

	log.Printf("Created telemetry %s", string(encodedTelemetry))

	json.NewEncoder(w).Encode(newTelemetry)
})
