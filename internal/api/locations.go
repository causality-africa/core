package api

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (api *API) GetLocations(c echo.Context) error {
	p, err := validatePagination(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ctx := c.Request().Context()
	offset := (p.Page - 1) * p.Size
	locations, err := api.db.GetLocations(ctx, p.Size, offset)

	if err != nil {
		slog.Error("Error getting locations", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "error querying database"},
		)
	}

	return c.JSON(http.StatusOK, locations)
}

func (api *API) GetLocationByCode(c echo.Context) error {
	code := c.Param("code")
	if len(code) == 0 {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "code must be provided"},
		)
	}

	ctx := c.Request().Context()
	locations, err := api.db.GetLocationsByCodes(ctx, []string{code})
	if err != nil {
		slog.Error("Error getting locations", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "error querying database"},
		)
	}

	if len(locations) == 0 {
		return c.JSON(
			http.StatusNotFound,
			map[string]string{"error": "location not found"},
		)
	}

	return c.JSON(http.StatusOK, locations[0])
}
