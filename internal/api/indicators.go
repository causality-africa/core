package api

import (
	"core/internal/db"
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (api *API) GetIndicators(c echo.Context) error {
	p := getPaginationParams(c)

	ctx := c.Request().Context()
	offset := (p.Page - 1) * p.Size
	indicators, more, err := api.db.GetIndicators(ctx, p.Size, offset)
	if err != nil {
		slog.Error("Error getting indicators", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "error querying database"},
		)
	}

	resp := map[string]any{"results": indicators, "more": more}
	return c.JSON(http.StatusOK, resp)
}

func (api *API) GetIndicatorByCode(c echo.Context) error {
	code := c.Param("code")
	if len(code) == 0 {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "code must be provided"},
		)
	}

	ctx := c.Request().Context()
	indicator, err := api.db.GetIndicatorByCode(ctx, code)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return c.JSON(
				http.StatusNotFound,
				map[string]string{"error": "indicator not found"},
			)
		}

		slog.Error("Error getting indicator", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "error querying database"},
		)
	}

	return c.JSON(http.StatusOK, indicator)
}
