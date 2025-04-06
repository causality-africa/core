package db

import (
	"context"
	"core/internal/models"
	"fmt"
	"sort"
	"strings"
	"time"
)

func (db *DB) GetLocations(
	ctx context.Context,
	limit, offset int,
) ([]models.Location, bool, error) {
	query := `
        SELECT id, name, code, admin_level, parent_id, map
        FROM locations
        ORDER BY code LIMIT $1 OFFSET $2
	`
	rows, err := db.pool.Query(ctx, query, limit+1, offset)
	if err != nil {
		return nil, false, fmt.Errorf("cannot query locations: %w", err)
	}
	defer rows.Close()

	locations := []models.Location{}
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
			return nil, false, fmt.Errorf("cannot scan row: %w", err)
		}
		locations = append(locations, loc)
	}

	if err := rows.Err(); err != nil {
		return nil, false, fmt.Errorf("cannot read rows: %w", err)
	}

	hasMore := false
	if len(locations) > limit {
		hasMore = true
		locations = locations[:limit]
	}

	return locations, hasMore, nil
}

func (db *DB) GetLocationsByCodes(
	ctx context.Context,
	codes []string,
) ([]models.Location, error) {
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

	locations := []models.Location{}
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

func (db *DB) GetRegions(
	ctx context.Context,
	limit, offset int,
) ([]models.Region, bool, error) {
	query := `
		SELECT r.id, r.name, r.code, r.description,
		       lr.location_id, l.code, lr.region_id, lr.join_date, lr.exit_date
		FROM regions r
		LEFT JOIN location_in_region lr ON r.id = lr.region_id
		LEFT JOIN locations l ON lr.location_id = l.id
		ORDER BY r.name LIMIT $1 OFFSET $2
	`
	rows, err := db.pool.Query(ctx, query, limit+1, offset)
	if err != nil {
		return nil, false, fmt.Errorf("cannot query regions: %w", err)
	}
	defer rows.Close()

	regionsMap := make(map[int]models.Region)
	for rows.Next() {
		var regionId int
		var name, code string
		var description *string
		var locationId *int
		var locationCode *string
		var regionIdFromJoin *int
		var joinDate *time.Time
		var exitDate *time.Time

		if err := rows.Scan(
			&regionId,
			&name,
			&code,
			&description,
			&locationId,
			&locationCode,
			&regionIdFromJoin,
			&joinDate,
			&exitDate,
		); err != nil {
			return nil, false, fmt.Errorf("cannot scan row: %w", err)
		}

		region, exists := regionsMap[regionId]
		if !exists {
			region = models.Region{
				Id:          regionId,
				Name:        name,
				Code:        code,
				Description: description,
				Locations:   []models.LocationInRegion{},
			}
		}

		if locationId != nil {
			location := models.LocationInRegion{
				LocationId:   *locationId,
				LocationCode: *locationCode,
				RegionId:     *regionIdFromJoin,
				JoinDate:     *joinDate,
				ExitDate:     exitDate,
			}
			region.Locations = append(region.Locations, location)
		}

		regionsMap[regionId] = region
	}

	regions := make([]models.Region, 0, len(regionsMap))
	for _, region := range regionsMap {
		regions = append(regions, region)
	}

	// Sort by name to ensure consistent ordering
	sort.Slice(regions, func(i, j int) bool {
		return regions[i].Name < regions[j].Name
	})

	hasMore := false
	if len(regions) > limit {
		hasMore = true
		regions = regions[:limit]
	}

	return regions, hasMore, nil
}

func (db *DB) GetRegionByCode(
	ctx context.Context, code string,
) (*models.Region, error) {
	query := `
		SELECT r.id, r.name, r.code, r.description,
		       lr.location_id, l.code, lr.region_id, lr.join_date, lr.exit_date
		FROM regions r
		LEFT JOIN location_in_region lr ON r.id = lr.region_id
		LEFT JOIN locations l ON lr.location_id = l.id
		WHERE r.code = $1
	`
	rows, err := db.pool.Query(ctx, query, code)
	if err != nil {
		return nil, fmt.Errorf("cannot query region: %w", err)
	}
	defer rows.Close()

	var region *models.Region

	for rows.Next() {
		var regionId int
		var name, regionCode string
		var description *string
		var locationId *int
		var locationCode *string
		var regionIdFromJoin *int
		var joinDate *time.Time
		var exitDate *time.Time

		if err := rows.Scan(
			&regionId,
			&name,
			&regionCode,
			&description,
			&locationId,
			&locationCode,
			&regionIdFromJoin,
			&joinDate,
			&exitDate,
		); err != nil {
			return nil, fmt.Errorf("cannot scan row: %w", err)
		}

		if region == nil {
			region = &models.Region{
				Id:          regionId,
				Name:        name,
				Code:        regionCode,
				Description: description,
				Locations:   []models.LocationInRegion{},
			}
		}

		if locationId != nil {
			location := models.LocationInRegion{
				LocationId:   *locationId,
				LocationCode: *locationCode,
				RegionId:     *regionIdFromJoin,
				JoinDate:     *joinDate,
				ExitDate:     exitDate,
			}
			region.Locations = append(region.Locations, location)
		}
	}

	if region == nil {
		return nil, ErrNotFound
	}

	return region, nil
}
