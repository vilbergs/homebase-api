package db

import "github.com/vilbergs/homebase-api/models"

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
