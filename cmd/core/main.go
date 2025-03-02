package main

import (
	"core/internal/db"
	handler "core/internal/handlers"
	repository "core/internal/repositories"
	"core/internal/routes"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {

	// initialize the db connection
	db.InitDB()

	// get the database instance
	database := db.GetDB()
	defer database.Close()

	// initialize the repositories
	locationRepo := repository.NewLocationRepository(database)

	//initialize handlers
	locationHandler := handler.NewLocationHandler(locationRepo)

	// create an echo instance
	e := echo.New()

	// register routes
	routes.RegisterRoutes(e, locationHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
