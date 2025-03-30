package api

import (
	"core/internal/api/middlewarex"
	"core/internal/cache"
	"core/internal/config"
	"core/internal/db"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	rateLimit         = 100.0
	rateLimitDuration = 5 * time.Minute

	cacheTTL = 24 * time.Hour
)

type API struct {
	echo  *echo.Echo
	db    *db.DB
	cache *cache.Cache
}

func New(
	database *db.DB,
	cache *cache.Cache,
	observabilityCfg *config.Observability,
) *API {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              observabilityCfg.SentryDSN,
		TracesSampleRate: 1.0,
		SendDefaultPII:   false,
	}); err != nil {
		slog.Error("Failed to initialized Sentry", "error", err)
	}

	e := echo.New()
	api := &API{echo: e, db: database, cache: cache}

	// Routes
	e.GET("/v1/query", api.query)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	rateLimiterCfg := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store:   middlewarex.NewRateLimiterCacheStore(rateLimit, rateLimitDuration, cache),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusInternalServerError, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			errStr := fmt.Sprintf("rate limit exceeded, try again in %s", rateLimitDuration)
			return context.JSON(
				http.StatusTooManyRequests,
				map[string]string{"error": errStr},
			)
		},
	}
	e.Use(middleware.RateLimiterWithConfig(rateLimiterCfg))

	e.Use(middlewarex.CacheMiddleware(cache, cacheTTL))

	e.Use(sentryecho.New(sentryecho.Options{}))

	return api
}

func (api *API) Start(cfg *config.API) error {
	addr := fmt.Sprintf(":%d", cfg.Port)
	return api.echo.Start(addr)
}
