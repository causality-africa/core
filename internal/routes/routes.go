package routes

import (
	handler "core/internal/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, locationHandler *handler.LocationHandler, indicatorHandler *handler.IndicatorHandler) {
	e.GET("/v1/api/locations", locationHandler.GetLocationsHandler)
	e.GET("/v1/api/locations/:code", locationHandler.GetLocationByISOHandler)

	e.GET("/v1/api/indicators", indicatorHandler.GetIndicatorsHandler)
	e.GET("/v1/api/indicators/:id", indicatorHandler.GetIndicatorByIdHandler)
}
