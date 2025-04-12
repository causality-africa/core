package db

import (
	"context"
	"core/internal/models"
	"fmt"

	"github.com/lib/pq"
)

func (db *DB) querySources(ctx context.Context, query string, args ...interface{}) ([]models.DataSource, error) {
	rows, err := db.pool.Query(ctx, query, args...)
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
			&src.LastUpdated,
		); err != nil {
			return nil, fmt.Errorf("cannot scan row: %w", err)
		}
		sources = append(sources, src)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot read rows: %w", err)
	}

	return sources, nil
}

func (db *DB) GetSourcesPaginated(ctx context.Context, limit, offset int) ([]models.DataSource, bool, error) {
	query := `
        SELECT id, name, url, description, last_updated
        FROM data_sources
        ORDER BY name
        LIMIT $1 OFFSET $2
    `
	sources, err := db.querySources(ctx, query, limit+1, offset)
	if err != nil {
		return nil, false, err
	}

	hasMore := false
	if len(sources) > limit {
		hasMore = true
		sources = sources[:limit]
	}

	return sources, hasMore, nil
}

func (db *DB) GetSourcesByIds(ctx context.Context, ids []int) ([]models.DataSource, error) {
	query := `
        SELECT id, name, url, description, last_updated
        FROM data_sources
        WHERE id = ANY($1)
        ORDER BY name
    `
	return db.querySources(ctx, query, pq.Array(ids))
}
