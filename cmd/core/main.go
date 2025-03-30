package main

import (
	"core/internal/api"
	"core/internal/cache"
	"core/internal/config"
	"core/internal/db"
	"log/slog"
	"os"
)

var (
	Version = "dev"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Starting Core...", "version", Version)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		return
	}

	db, err := db.New(&cfg.DB)
	if err != nil {
		slog.Error("Failed to connect to DB", "error", err)
		return
	}
	defer db.Close()

	cache, err := cache.New(&cfg.Cache)
	if err != nil {
		slog.Error("Failed to connect to cache", "error", err)
		return
	}

	server := api.New(db, cache)
	err = server.Start(&cfg.API)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
		return
	}
}
