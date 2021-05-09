package db

import (
	"fmt"
	"net/http"

	"github.com/vilbergs/homebase-api/models"
)

func AddTelemetry(t *models.Telemetry) error {
  writeAPI := Influx.WriteAPI(INFLUX_ORG, INLFUX_BUCKET)

	// write line protocol
	writeAPI.WriteRecord(fmt.Sprintf("telemetry,zoneId=%d temperature=%f,humidity=%f", t.ZoneId, t.Temperature, t.Humidity))
	// Flush writes
	writeAPI.Flush()

	return nil
}

func GetTelemetry(w http.ResponseWriter, r *http.Request) error {
	// TODO: Implement

	return nil
}