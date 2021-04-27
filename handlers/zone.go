package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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