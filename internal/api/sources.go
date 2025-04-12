package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (api *API) GetSources(c echo.Context) error {
	p := getPaginationParams(c)

	ctx := c.Request().Context()
	offset := (p.Page - 1) * p.Size
	sources, more, err := api.db.GetSourcesPaginated(ctx, p.Size, offset)
	if err != nil {
		slog.Error("Error getting sources", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "error querying database"},
		)
	}

	resp := map[string]any{"results": sources, "more": more}
	return c.JSON(http.StatusOK, resp)
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
	sources, err := api.db.GetSourcesByIds(ctx, []int{id})
	if err != nil {
		slog.Error("Error getting source", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "error querying database"},
		)
	}

	if len(sources) == 0 {
		return c.JSON(
			http.StatusNotFound,
			map[string]string{"error": "source not found"},
		)
	}

	return c.JSON(http.StatusOK, sources[0])
}
