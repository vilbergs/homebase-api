package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vilbergs/homebase-api/db"
	"github.com/vilbergs/homebase-api/models"
)

var AddTelemetry = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var newTelemetry models.Telemetry
	zoneId, parseErr := strconv.Atoi(mux.Vars(r)["zoneId"])

	if parseErr != nil {
		http.Error(w, parseErr.Error(), http.StatusInternalServerError)
	}

	newTelemetry.ZoneId = zoneId
	reqBody, readError := ioutil.ReadAll(r.Body)
	if readError != nil {
		log.Print("Error!")
	}

	json.Unmarshal(reqBody, &newTelemetry)
	db.AddTelemetry(&newTelemetry, zoneId)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encodedTelemetry, _ := json.Marshal(newTelemetry)

	log.Printf("Created telemetry %s", string(encodedTelemetry))

	json.NewEncoder(w).Encode(newTelemetry)
})
