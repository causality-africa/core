package handler

import (
	repository "core/internal/repositories"
	"core/internal/utils"
	"log"
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch locations"})
	}

	return c.JSON(http.StatusOK, locations)
}

func (h *LocationHandler) GetLocationByISOHandler(c echo.Context) error {
	isoCode := c.Param("iso")

	location, err := h.repo.GetLocationByISO(isoCode)

	if err != nil {
		if err.Error() == "location not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, "Location not found.")
		}

		log.Println("Error fetching location:", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error.")
	}

	return c.JSON(http.StatusOK, location)
}
