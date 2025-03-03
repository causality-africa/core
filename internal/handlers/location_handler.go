package handler

import (
	repository "core/internal/repositories"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type LocationHandler struct {
	repo *repository.LocationRepository
}

func NewLocationHandler(repo *repository.LocationRepository) *LocationHandler {
	return &LocationHandler{repo: repo}
}

func (h *LocationHandler) GetLocationsHandler(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	page, _ := strconv.Atoi(c.QueryParam("page"))

	if limit <= 0 {
		limit = 10
	}

	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	locations, err := h.repo.GetLocations(limit, offset)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Locations not found"})
	}

	return c.JSON(http.StatusOK, locations)
}

func (h *LocationHandler) GetLocationByISOHandler(c echo.Context) error {
	code := c.Param("code")

	location, err := h.repo.GetLocationByISO(code)

	if err != nil {
		if err.Error() == "location not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Location not found"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, location)
}
