package db

import (
	"context"
	"core/internal/models"
	"fmt"

	"github.com/lib/pq"
)

func (db *DB) queryGeoEntities(ctx context.Context, query string, args ...interface{}) ([]models.GeoEntity, error) {
	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot query geographic entities: %w", err)
	}
	defer rows.Close()

	entities := []models.GeoEntity{}
	for rows.Next() {
		var entity models.GeoEntity
		if err := rows.Scan(
			&entity.Id,
			&entity.Code,
			&entity.Name,
			&entity.Type,
		); err != nil {
			return nil, fmt.Errorf("cannot scan row: %w", err)
		}
		entities = append(entities, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot read rows: %w", err)
	}

	return entities, nil
}

func (db *DB) GetGeoEntitiesPaginated(ctx context.Context, limit, offset int) ([]models.GeoEntity, bool, error) {
	query := `
        SELECT id, code, name, type
        FROM geo_entities
        ORDER BY code
        LIMIT $1 OFFSET $2
    `
	entities, err := db.queryGeoEntities(ctx, query, limit+1, offset)
	if err != nil {
		return nil, false, err
	}

	hasMore := false
	if len(entities) > limit {
		hasMore = true
		entities = entities[:limit]
	}

	return entities, hasMore, nil
}

func (db *DB) GetGeoEntitiesByCodes(ctx context.Context, codes []string) ([]models.GeoEntity, error) {
	query := `
        SELECT id, code, name, type
        FROM geo_entities
        WHERE code = ANY($1)
        ORDER BY code
    `
	return db.queryGeoEntities(ctx, query, pq.Array(codes))
}
