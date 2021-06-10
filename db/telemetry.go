package db

import (
	"net/http"
	"strconv"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api/write"

	"github.com/vilbergs/homebase-api/models"
)

func AddTelemetry(t *models.Telemetry, zoneId int) error {
	writeAPI := Influx.WriteAPI(INFLUX_ORG, INFLUX_BUCKET)

	writeAPI.WritePoint(write.NewPoint(
		"telemetry",
		map[string]string{
			"zoneId": strconv.Itoa(zoneId),
		},
		map[string]interface{}{
			"temperature": t.Temperature,
			"humidity":    t.Humidity,
		},
		time.Now(),
	))

	return nil
}

func GetTelemetry(w http.ResponseWriter, r *http.Request) error {
	// TODO: Implement

	return nil
}
