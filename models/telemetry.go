package models

type Telemetry struct {
	ZoneId      int     `json:"zoneId"`
	Temperature float64 `json:"temperature"`
	Humidity	  float64 `json:"humidity"`
}