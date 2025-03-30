package db

import (
	"context"
	"core/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

func (db *DB) GetSources(
	ctx context.Context,
	limit, offset int,
) ([]models.DataSource, error) {
	query := `
		SELECT id, name, url, description, date
		FROM data_sources
		ORDER BY name LIMIT $1 OFFSET $2
	`
	rows, err := db.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("cannot query data sources: %w", err)
	}
	defer rows.Close()

	sources := []models.DataSource{}
	for rows.Next() {
		var src models.DataSource
		if err := rows.Scan(
			&src.Id,
			&src.Name,
			&src.URL,
			&src.Description,
			&src.Date,
		); err != nil {
			return nil, fmt.Errorf("cannot scan row: %w", err)
		}

		sources = append(sources, src)
	}

	return sources, nil
}

func (db *DB) GetSourceById(
	ctx context.Context, id int,
) (*models.DataSource, error) {
	query := `
		SELECT id, name, url, description, date
		FROM data_sources
		WHERE id = $1 LIMIT 1
	`
	row := db.pool.QueryRow(ctx, query, id)

	var src models.DataSource
	err := row.Scan(
		&src.Id,
		&src.Name,
		&src.URL,
		&src.Description,
		&src.Date,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("cannot scan row: %w", err)
	}

	return &src, nil
}
