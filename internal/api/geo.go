package api

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (api *API) GetGeoEntities(c echo.Context) error {
	p := getPaginationParams(c)

	ctx := c.Request().Context()
	offset := (p.Page - 1) * p.Size
	geoEntities, more, err := api.db.GetGeoEntitiesPaginated(ctx, p.Size, offset)

	if err != nil {
		slog.Error("Error getting geographic entities", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "error querying database"},
		)
	}

	resp := map[string]any{"results": geoEntities, "more": more}
	return c.JSON(http.StatusOK, resp)
}

func (api *API) GetGeoEntityByCode(c echo.Context) error {
	code := c.Param("code")
	if len(code) == 0 {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "code must be provided"},
		)
	}

	ctx := c.Request().Context()
	geoEntities, err := api.db.GetGeoEntitiesByCodes(ctx, []string{code})
	if err != nil {
		slog.Error("Error getting geo entities", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "error querying database"},
		)
	}

	if len(geoEntities) == 0 {
		return c.JSON(
			http.StatusNotFound,
			map[string]string{"error": "geographic entity not found"},
		)
	}

	return c.JSON(http.StatusOK, geoEntities[0])
}
