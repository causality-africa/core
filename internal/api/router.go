package api

import (
	"core/internal/config"
	"core/internal/db"
	"fmt"

	"github.com/labstack/echo/v4"
)

type API struct {
	echo *echo.Echo
	db   *db.DB
}

func New(database *db.DB) *API {
	e := echo.New()

	api := &API{echo: e, db: database}

	e.GET("/v1/query", api.query)

	return api
}

func (api *API) Start(cfg *config.API) error {
	addr := fmt.Sprintf(":%d", cfg.Port)
	return api.echo.Start(addr)
}
