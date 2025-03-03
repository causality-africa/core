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

func (r *LocationRepository) GetLocations(isoCode string) ([]models.Location, error) {
	query := "SELECT * FROM locations"

	var rows *sql.Rows
	var err error

	if isoCode != "" {
		query += ` WHERE iso_code = $1`
		rows, err = r.DB.Query(query, isoCode)
	} else {
		rows, err = r.DB.Query(query)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	locations := []models.Location{}

	for rows.Next() {
		var loc models.Location

		if err := rows.Scan(&loc.Id, &loc.Name, &loc.AdminLevel, &loc.ParentId, &loc.IsoCode, &loc.Map); err != nil {
			return nil, err
		}

		locations = append(locations, loc)
	}

	return locations, nil
}

func (r *LocationRepository) GetLocationByISO(isoCode string) (*models.Location, error) {
	query := "SELECT * FROM locations WHERE iso_code = $1 LIMIT 1"
	row := r.DB.QueryRow(query, isoCode)

	var loc models.Location

	err := row.Scan(&loc.Id, &loc.Name, &loc.AdminLevel, &loc.ParentId, &loc.IsoCode, &loc.Map)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("location not found")
		}

		return nil, err
	}

	return &loc, nil
}
