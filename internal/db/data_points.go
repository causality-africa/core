package db

import (
	"context"
	"core/internal/models"
	"fmt"
	"strings"
	"time"
)

func (db *DB) GetDataPointsByGeoCodes(
	ctx context.Context,
	indicator string,
	geoCodes []string,
	startDate, endDate time.Time,
) (map[string][]models.DataPoint, error) {
	if len(geoCodes) == 0 {
		return nil, nil
	}

	geoParams := make([]string, len(geoCodes))
	args := make([]any, len(geoCodes)+3) // +3 for indicator, start/end dates

	args[0] = indicator
	args[1] = startDate
	args[2] = endDate

	for i, code := range geoCodes {
		geoParams[i] = fmt.Sprintf("$%d", i+4) // +4 because indicator, dates are $1, $2, $3
		args[i+3] = code
	}

	query := fmt.Sprintf(`
        SELECT dp.id, ge.code, dp.source_id, dp.date, dp.numeric_value, dp.text_value
        FROM data_points dp
        JOIN geo_entities ge ON dp.geo_entity_id = ge.id
        JOIN indicators i ON dp.indicator_id = i.id
        WHERE i.code = $1
        AND ge.code IN (%s)
        AND dp.date >= $2 AND dp.date <= $3
        ORDER BY ge.code, dp.date
    `, strings.Join(geoParams, ","))

	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot query data points: %w", err)
	}
	defer rows.Close()

	results := make(map[string][]models.DataPoint)
	for _, code := range geoCodes {
		results[code] = []models.DataPoint{}
	}

	for rows.Next() {
		var dp models.DataPoint
		var geoCode string
		if err := rows.Scan(
			&dp.Id,
			&geoCode,
			&dp.SourceId,
			&dp.Date,
			&dp.NumericValue,
			&dp.TextValue,
		); err != nil {
			return nil, fmt.Errorf("cannot scan row: %w", err)
		}

		results[geoCode] = append(results[geoCode], dp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot read rows: %w", err)
	}

	return results, nil
}
