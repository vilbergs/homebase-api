package db

import (
	"github.com/vilbergs/homebase-api/models"
)

func AddZone(z *models.Zone) error {
	var id int
	var createdAt string

	query := `INSERT INTO zones (name) VALUES ($1) RETURNING id, created_at`
	err := Postgres.QueryRow(query, z.Name).Scan(&id, &createdAt)

	if err != nil {
		return err
	}

	z.ID = id
	z.CreatedAt = createdAt

	return nil
}

func GetZone(zoneId int) (*models.Zone, error) {
	var z models.Zone

	query := "SELECT id, name, created_at FROM zones WHERE id = $1"
	err := Postgres.QueryRow(query, zoneId).Scan(&z.ID, &z.Name, &z.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &z, nil
}

func GetALLZones() (*[]models.Zone, error) {
	query := `SELECT id, name, created_at FROM zones`
	rows, err := Postgres.Query(query)

	if err != nil {
		return nil, err
	}

	zones := []models.Zone{}
	for rows.Next() {
		var z models.Zone

		err := rows.Scan(&z.ID, &z.Name, &z.CreatedAt)

		if err != nil {
			return nil, err
		}

		zones = append(zones, z)
	}

	defer rows.Close()

	return &zones, nil
}
