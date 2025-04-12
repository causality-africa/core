package db

import (
	"context"
	"core/internal/models"
	"fmt"

	"github.com/lib/pq"
)

func (db *DB) queryIndicators(ctx context.Context, query string, args ...interface{}) ([]models.Indicator, error) {
	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot query indicators: %w", err)
	}
	defer rows.Close()

	indicators := []models.Indicator{}
	for rows.Next() {
		var ind models.Indicator
		if err := rows.Scan(
			&ind.Id,
			&ind.Code,
			&ind.Name,
			&ind.Category,
			&ind.Unit,
			&ind.Description,
			&ind.DataType,
		); err != nil {
			return nil, fmt.Errorf("cannot scan row: %w", err)
		}
		indicators = append(indicators, ind)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot read rows: %w", err)
	}

	return indicators, nil
}

func (db *DB) GetIndicatorsPaginated(ctx context.Context, limit, offset int) ([]models.Indicator, bool, error) {
	query := `
        SELECT id, code, name, category, unit, description, data_type
        FROM indicators
        ORDER BY name
        LIMIT $1 OFFSET $2
    `
	indicators, err := db.queryIndicators(ctx, query, limit+1, offset)
	if err != nil {
		return nil, false, err
	}

	hasMore := false
	if len(indicators) > limit {
		hasMore = true
		indicators = indicators[:limit]
	}

	return indicators, hasMore, nil
}

func (db *DB) GetIndicatorsByCodes(ctx context.Context, codes []string) ([]models.Indicator, error) {
	query := `
        SELECT id, code, name, category, unit, description, data_type
        FROM indicators
        WHERE code = ANY($1)
        ORDER BY name
    `
	return db.queryIndicators(ctx, query, pq.Array(codes))
}
