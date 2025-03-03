package repository

import (
	"core/internal/models"
	"database/sql"
	"errors"
)

type LocationRepository struct {
	DB *sql.DB
}

func NewLocationRepository(db *sql.DB) *LocationRepository {
	return &LocationRepository{DB: db}
}

func (r *LocationRepository) GetLocations(limit int, offset int) ([]models.Location, error) {
	query := "SELECT * FROM locations"
	args := []interface{}{}

	query += " ORDER BY id LIMIT $1 OFFSET $2"
	args = append(args, limit, offset)

	rows, err := r.DB.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	locations := []models.Location{}

	for rows.Next() {
		var loc models.Location

		if err := rows.Scan(&loc.Id, &loc.Name, &loc.AdminLevel, &loc.ParentId, &loc.Code, &loc.Map); err != nil {
			return nil, err
		}

		locations = append(locations, loc)
	}

	return locations, nil
}

func (r *LocationRepository) GetLocationByISO(code string) (*models.Location, error) {
	query := "SELECT * FROM locations WHERE code = $1 LIMIT 1"
	row := r.DB.QueryRow(query, code)

	var loc models.Location

	err := row.Scan(&loc.Id, &loc.Name, &loc.AdminLevel, &loc.ParentId, &loc.Code, &loc.Map)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("location not found")
		}

		return nil, err
	}

	return &loc, nil
}
