package db

import (
	"context"
	"core/internal/models"
	"fmt"
	"strings"
)

func (db *DB) GetLocationsByCodes(
	ctx context.Context,
	codes []string,
) ([]models.Location, error) {
	locations := []models.Location{}

	rows, err := db.pool.Query(
		ctx,
		`
        SELECT id, name, code, admin_level, parent_id, map
        FROM locations
        WHERE code = ANY($1)
        ORDER BY code
		`,
		"{"+strings.Join(codes, ",")+"}")
	if err != nil {
		return nil, fmt.Errorf("cannot query locations by codes: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var loc models.Location
		if err := rows.Scan(
			&loc.Id,
			&loc.Name,
			&loc.Code,
			&loc.AdminLevel,
			&loc.ParentID,
			&loc.Map,
		); err != nil {
			return nil, fmt.Errorf("cannot scan row: %w", err)
		}
		locations = append(locations, loc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot read rows: %w", err)
	}

	return locations, nil
}
