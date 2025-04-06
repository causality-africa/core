package api

import (
	"core/internal/db"
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (api *API) GetLocations(c echo.Context) error {
	p := getPaginationParams(c)

	ctx := c.Request().Context()
	offset := (p.Page - 1) * p.Size
	locations, more, err := api.db.GetLocations(ctx, p.Size, offset)

	if err != nil {
		slog.Error("Error getting locations", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "error querying database"},
		)
	}

	resp := map[string]any{"results": locations, "more": more}
	return c.JSON(http.StatusOK, resp)
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

func (api *API) GetRegions(c echo.Context) error {
	p := getPaginationParams(c)

	ctx := c.Request().Context()
	offset := (p.Page - 1) * p.Size
	regions, more, err := api.db.GetRegions(ctx, p.Size, offset)
	if err != nil {
		slog.Error("Error getting regions", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "error querying database"},
		)
	}

	resp := map[string]any{"results": regions, "more": more}
	return c.JSON(http.StatusOK, resp)
}

func (api *API) GetRegionByCode(c echo.Context) error {
	code := c.Param("code")
	if len(code) == 0 {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "code must be provided"},
		)
	}

	ctx := c.Request().Context()
	region, err := api.db.GetRegionByCode(ctx, code)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return c.JSON(
				http.StatusNotFound,
				map[string]string{"error": "region not found"},
			)
		}

		slog.Error("Error getting region", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "error querying database"},
		)
	}

	return c.JSON(http.StatusOK, region)
}
