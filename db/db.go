package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	_ "github.com/lib/pq"
)

/**
 * The POSTGRES_HOST and INFLUX_HOST hostnames are set up by docker, which defaults the hostname
 * to the container's ID
 * - https://docs.docker.com/config/containers/container-networking/#ip-address-and-hostname
 */
const (
	POSTGRES_HOST = "postgres_db"
	POSTGRES_PORT = 5432
	INFLUX_HOST   = "influx_db"
	INFLUX_PORT   = 8086
)

var (
	POSTGRES_USER = os.Getenv("POSTGRES_USER")
	POSTGRES_PW   = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DB   = os.Getenv("POSTGRES_DB_NAME")
	INFLUX_TOKEN  = os.Getenv("INFLUX_TOKEN")
	INLFUX_BUCKET = os.Getenv("INFLUX_BUCKET")
	INFLUX_ORG    = os.Getenv("INFLUX_ORG")
)

// ErrNoMatch is returned when we request a row that doesn't exist
var ErrNoMatch = fmt.Errorf("no matching record")
var Postgres *sql.DB
var Influx influxdb2.Client

func PostgresInit() error {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		POSTGRES_HOST,
		POSTGRES_PORT,
		POSTGRES_USER,
		POSTGRES_PW,
		POSTGRES_DB,
	)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	Postgres = conn
	err = Postgres.Ping()
	if err != nil {
		return err
	}

	log.Println("Postgres connection established")

	return nil
}

func InfluxInit() error {

	Influx = influxdb2.NewClient(fmt.Sprintf("http://%s:%d", INFLUX_HOST, INFLUX_PORT), INFLUX_TOKEN)

	log.Println("Influx connection established")

	return nil
}
