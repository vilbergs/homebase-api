package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	_ "github.com/lib/pq"
)

const (
    POSTGRES_PORT = 5432
    INLFUX_BUCKET = "homebase"
    INFLUX_ORG = "Homebase"
)

var (
	POSTGRES_USER = os.Getenv("POSTGRES_USER")
	POSTGRES_PW = os.Getenv("POSTGRES_PASSWORD")
    POSTGRES_DB = os.Getenv("POSTGRES_DB_NAME")
    INFLUX_TOKEN = os.Getenv("INFLUX_TOKEN")
)

// ErrNoMatch is returned when we request a row that doesn't exist
var ErrNoMatch = fmt.Errorf("no matching record")
var Postgres *sql.DB
var Influx influxdb2.Client


func PostgresInit() error {
    dsn := "host=database port=5432 user=homebase password=homebase_pass dbname=homebase sslmode=disable"
    conn, err := sql.Open("postgres", dsn)
    if err != nil {
        return err
    }
    Postgres = conn
    err = Postgres.Ping()
    if err != nil {
        return  err
    }

    log.Println("Database connection established")

    return nil
}

func InfluxInit() error {
    hostname := "influx_db"

    Influx = influxdb2.NewClient(fmt.Sprintf("http://%s:8086", hostname), INFLUX_TOKEN)

    return nil
}