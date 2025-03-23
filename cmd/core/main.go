package main

import (
	"core/internal/api"
	"core/internal/config"
	"core/internal/db"
	"log"

	_ "github.com/lib/pq"
)

var (
	Version = "dev"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	db, err := db.New(&cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v\n", err)
	}
	defer db.Close()

	server := api.New(db)
	err = server.Start(&cfg.API)
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
