package api

import (
	"core/internal/db"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (api *API) GetSources(c echo.Context) error {
	p := getPaginationParams(c)

	ctx := c.Request().Context()
	offset := (p.Page - 1) * p.Size
	sources, err := api.db.GetSources(ctx, p.Size, offset)
	if err != nil {
		slog.Error("Error getting sources", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "error querying database"},
		)
	}

	return c.JSON(http.StatusOK, sources)
}

func (api *API) GetSourceById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "invalid id provided"},
		)
	}

	ctx := c.Request().Context()
	source, err := api.db.GetSourceById(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return c.JSON(
				http.StatusNotFound,
				map[string]string{"error": "source not found"},
			)
		}

		slog.Error("Error getting source", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "error querying database"},
		)
	}

	return c.JSON(http.StatusOK, source)
}
