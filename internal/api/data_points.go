package api

import (
	"core/internal/models"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type Query struct {
	Indicator string
	StartDate time.Time
	EndDate   time.Time
	Locations []string
}

func validateQuery(c echo.Context) (*Query, error) {
	var q Query
	var err error

	q.Indicator = c.QueryParam("indicator")
	if len(q.Indicator) == 0 {
		return nil, errors.New("indicator must be provided")
	}

	startDate := c.QueryParam("start")
	q.StartDate, err = time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}

	endDate := c.QueryParam("end")
	q.EndDate, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, errors.New("invalid end date format")
	}

	locations := c.QueryParam("locations")
	if len(locations) == 0 {
		return nil, errors.New("locations must be provided")
	}

	q.Locations = strings.Split(locations, ",")
	if len(q.Locations) > pageSize {
		return nil, fmt.Errorf("maximum of %d locations allowed", pageSize)
	}

	return &q, nil
}

func (api *API) query(c echo.Context) error {
	q, err := validateQuery(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ctx := c.Request().Context()
	locations, err := api.db.GetLocationsByCodes(ctx, q.Locations)
	if err != nil {
		slog.Error("Error getting locations by codes", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "error querying database"},
		)
	}

	locationIds := make([]int, len(locations))
	for i, location := range locations {
		locationIds[i] = location.Id
	}

	dataPoints, err := api.db.GetDataPointsForLocations(
		ctx, q.Indicator, locationIds, q.StartDate, q.EndDate,
	)
	if err != nil {
		slog.Error("Error getting data points for locations", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "error querying database"},
		)
	}

	result := make(map[string][]models.DataPoint)
	for _, loc := range locations {
		if len(dataPoints[loc.Id]) > 0 {
			result[loc.Code] = dataPoints[loc.Id]
		}
	}

	return c.JSON(http.StatusOK, result)
}
