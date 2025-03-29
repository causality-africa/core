package api

import (
	"core/internal/models"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type QueryResponse struct {
	Data map[string][]models.DataPoint `json:"data"`
	Next string                        `json:"next,omitempty"` // URL for next page
}

func (api *API) query(c echo.Context) error {
	indicator := c.QueryParam("indicator")
	if len(indicator) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "indicator must be provided"})
	}

	startDate := c.QueryParam("start")
	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start date format"})
	}

	endDate := c.QueryParam("end")
	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end date format"})
	}

	afterLocation := c.QueryParam("after")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	ctx := c.Request().Context()
	locations, err := api.db.GetLocations(ctx, afterLocation, limit)
	if err != nil {
		slog.Error("Error getting locations", "error", err)
		return err
	}

	locationIds := make([]int, len(locations))
	for i, location := range locations {
		locationIds[i] = location.Id
	}

	dataPoints, err := api.db.GetDataPointsForLocations(ctx, indicator, locationIds, startDateParsed, endDateParsed)
	if err != nil {
		slog.Error("Error getting data points for locations", "error", err)
		return err
	}

	results := make(map[string][]models.DataPoint)
	for _, loc := range locations {
		if len(dataPoints[loc.Id]) > 0 {
			results[loc.Code] = dataPoints[loc.Id]
		}
	}

	var nextPageURL string
	if len(locations) > 0 {
		lastLocation := locations[len(locations)-1].Code
		nextPageURL = fmt.Sprintf(
			"/v1/query?indicator=%s&start=%s&end=%s&after=%s&limit=%d",
			indicator, startDate, endDate, lastLocation, limit,
		)
	}

	return c.JSON(http.StatusOK, QueryResponse{Data: results, Next: nextPageURL})
}
