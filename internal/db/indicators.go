package db

import (
	"context"
	"core/internal/models"
	"fmt"
	"strings"
	"time"
)

func (db *DB) GetDataPointsForLocations(
	ctx context.Context,
	indicator string,
	locationIds []int,
	startDate, endDate time.Time,
) (map[int][]models.DataPoint, error) {
	results := map[int][]models.DataPoint{}
	if len(locationIds) == 0 {
		return results, nil
	}

	locationParams := make([]string, len(locationIds))
	args := make([]interface{}, len(locationIds)+3) // +3 for indicator code, start and end dates

	args[0] = indicator
	args[1] = startDate
	args[2] = endDate

	for i, id := range locationIds {
		locationParams[i] = fmt.Sprintf("$%d", i+4) // +4 because indicator, dates are $1, $2, $3
		args[i+3] = id

		results[id] = []models.DataPoint{}
	}

	query := fmt.Sprintf(`
        SELECT dp.id, entity_type, entity_id, indicator_id, source_id, date, numeric_value, text_value
        FROM data_points dp
		JOIN indicators i ON dp.indicator_id = i.id
        WHERE i.code = $1
        AND dp.entity_type = 'location'
        AND dp.entity_id IN (%s)
        AND dp.date >= $2 AND dp.date <= $3
    `, strings.Join(locationParams, ","))

	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot query data points: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dp models.DataPoint

		if err := rows.Scan(
			&dp.Id,
			&dp.EntityType,
			&dp.EntityID,
			&dp.IndicatorId,
			&dp.SourceId,
			&dp.Date,
			&dp.NumericValue,
			&dp.TextValue,
		); err != nil {
			return nil, fmt.Errorf("cannot scan row: %w", err)
		}

		results[dp.EntityID] = append(results[dp.EntityID], dp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot read rows: %w", err)
	}

	return results, nil
}
