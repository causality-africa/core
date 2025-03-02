package routes

import (
	handler "core/internal/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, locationHandler *handler.LocationHandler) {
	e.GET("/locations", locationHandler.GetLocationsHandler)
	e.GET("/locations/{iso}", locationHandler.GetLocationByISOHandler)
}
