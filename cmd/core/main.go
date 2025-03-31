package main

import (
	"core/internal/api"
	"core/internal/cache"
	"core/internal/config"
	"core/internal/db"
	"log/slog"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentryslog "github.com/getsentry/sentry-go/slog"
	slogmulti "github.com/samber/slog-multi"
)

var (
	Version = "dev"
)

func main() {
	slog.Info("Starting Core...", "version", Version)

	// Config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		return
	}

	// Logging
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.Observability.SentryDSN,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		SendDefaultPII:   false,
	}); err != nil {
		slog.Warn("Failed to initialized Sentry", "error", err)
	}
	defer sentry.Flush(2 * time.Second)

	sentryHandler := sentryslog.Option{Level: slog.LevelWarn}.NewSentryHandler()
	stdoutHandler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)
	logger := slog.New(slogmulti.Fanout(stdoutHandler, sentryHandler))
	logger = logger.With("release", Version)
	slog.SetDefault(logger)

	// Database
	db, err := db.New(&cfg.DB)
	if err != nil {
		slog.Error("Failed to connect to DB", "error", err)
		return
	}
	defer db.Close()

	// Cache
	cache, err := cache.New(&cfg.Cache)
	if err != nil {
		slog.Error("Failed to connect to cache", "error", err)
		return
	}

	// API
	server := api.New(db, cache, Version)
	err = server.Start(&cfg.API)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
		return
	}
}
