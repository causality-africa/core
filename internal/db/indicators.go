package db

import (
	"context"
	"core/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

func (db *DB) GetIndicators(
	ctx context.Context,
	limit, offset int,
) ([]models.Indicator, error) {
	query := `
		SELECT id, name, code, category, description, unit, data_type
		FROM indicators
		ORDER BY name LIMIT $1 OFFSET $2
	`
	rows, err := db.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("cannot query indicators: %w", err)
	}
	defer rows.Close()

	indicators := []models.Indicator{}
	for rows.Next() {
		var ind models.Indicator
		if err := rows.Scan(
			&ind.Id,
			&ind.Name,
			&ind.Code,
			&ind.Category,
			&ind.Description,
			&ind.Unit,
			&ind.DataType,
		); err != nil {
			return nil, fmt.Errorf("cannot scan row: %w", err)
		}

		indicators = append(indicators, ind)
	}

	return indicators, nil
}

func (db *DB) GetIndicatorByCode(
	ctx context.Context, code string,
) (*models.Indicator, error) {
	query := `
		SELECT id, name, code, category, description, unit, data_type
		FROM indicators
		WHERE code = $1 LIMIT 1
	`
	row := db.pool.QueryRow(ctx, query, code)

	var ind models.Indicator
	err := row.Scan(
		&ind.Id,
		&ind.Name,
		&ind.Code,
		&ind.Category,
		&ind.Description,
		&ind.Unit,
		&ind.DataType,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("cannot scan row: %w", err)
	}

	return &ind, nil
}
