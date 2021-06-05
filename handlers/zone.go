package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vilbergs/homebase-api/db"
	"github.com/vilbergs/homebase-api/models"
)

var AddZone = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var newZone models.Zone

	reqBody, readError := ioutil.ReadAll(r.Body)
	if readError != nil {
		fmt.Print("Error!")
	}

	json.Unmarshal(reqBody, &newZone)
	db.AddZone(&newZone)

	w.WriteHeader(http.StatusCreated)

	jZone, _ := json.Marshal(newZone)

	fmt.Printf("Created zone %s", string(jZone))

	json.NewEncoder(w).Encode(newZone)
})

var GetALLZones = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	zones, err := db.GetALLZones()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(zones)
})

var GetZone = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	zoneId, parseErr := strconv.Atoi(mux.Vars(r)["zoneId"])

	if parseErr != nil {
		http.Error(w, parseErr.Error(), http.StatusInternalServerError)
	}

	zones, err := db.GetZone(zoneId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(zones)
})
