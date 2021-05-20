package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	POSTGRES_HOST = "postgres-db"
	POSTGRES_PORT = 5432
)

var (
	POSTGRES_USER = os.Getenv("POSTGRES_USER")
	POSTGRES_PW   = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DB   = os.Getenv("POSTGRES_DB_NAME")
)

func main() {
	m, err := migrate.New(
		"file://db/migrations",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			POSTGRES_USER,
			POSTGRES_PW,
			POSTGRES_HOST,
			POSTGRES_PORT,
			POSTGRES_DB,
		))
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
