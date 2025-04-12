package api

import (
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
	GeoCodes  []string
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

	geoCodes := c.QueryParam("geo_codes")
	if len(geoCodes) == 0 {
		return nil, errors.New("geo_codes must be provided")
	}

	q.GeoCodes = strings.Split(geoCodes, ",")
	if len(q.GeoCodes) > pageSize {
		return nil, fmt.Errorf("maximum of %d geo codes allowed", pageSize)
	}

	return &q, nil
}

func (api *API) query(c echo.Context) error {
	q, err := validateQuery(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ctx := c.Request().Context()
	dataPoints, err := api.db.GetDataPointsByGeoCodes(
		ctx, q.Indicator, q.GeoCodes, q.StartDate, q.EndDate,
	)
	if err != nil {
		slog.Error("Error getting data points for geo codes", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "error querying database"},
		)
	}

	return c.JSON(http.StatusOK, dataPoints)
}
