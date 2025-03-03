package handler

import (
	repository "core/internal/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LocationHandler struct {
	repo *repository.LocationRepository
}

func NewLocationHandler(repo *repository.LocationRepository) *LocationHandler {
	return &LocationHandler{repo: repo}
}

func (h *LocationHandler) GetLocationsHandler(c echo.Context) error {
	isoCode := c.QueryParam("iso")

	locations, err := h.repo.GetLocations(isoCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Locations not found"})
	}

	return c.JSON(http.StatusOK, locations)
}

func (h *LocationHandler) GetLocationByISOHandler(c echo.Context) error {
	isoCode := c.Param("iso")

	location, err := h.repo.GetLocationByISO(isoCode)

	if err != nil {
		if err.Error() == "location not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Location not found"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, location)
}
