package routes

import (
	handler "core/internal/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, locationHandler *handler.LocationHandler) {
	e.GET("/v1/api/locations", locationHandler.GetLocationsHandler)
	e.GET("/v1/api/locations/:code", locationHandler.GetLocationByISOHandler)
}
