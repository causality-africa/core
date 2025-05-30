package api

import (
	"context"
	"core/internal/api/middlewarex"
	"core/internal/cache"
	"core/internal/config"
	"core/internal/db"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	rateLimit         = 300.0
	rateLimitDuration = 15 * time.Minute

	cacheTTL = 24 * time.Hour
)

type API struct {
	echo  *echo.Echo
	db    *db.DB
	cache *cache.Cache
}

func New(database *db.DB, cacheStore *cache.Cache, version string) *API {
	e := echo.New()
	api := &API{echo: e, db: database, cache: cacheStore}

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("Causality Core %s", version))
	})

	e.GET("/v1/geo", api.GetGeoEntities)
	e.GET("/v1/geo/:code", api.GetGeoEntityByCode)

	e.GET("/v1/indicators", api.GetIndicators)
	e.GET("/v1/indicators/:code", api.GetIndicatorByCode)

	e.GET("/v1/sources", api.GetSources)
	e.GET("/v1/sources/:id", api.GetSourceById)

	e.GET("/v1/query", api.query)

	// Middleware
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.Recover())

	rateLimiterCfg := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store:   middlewarex.NewRateLimiterCacheStore(rateLimit, rateLimitDuration, cacheStore),
		IdentifierExtractor: func(c echo.Context) (string, error) {
			return c.RealIP(), nil
		},
		ErrorHandler: func(c echo.Context, err error) error {
			slog.Warn("cannot identify client", "error", err)
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{"error": "cannot identify client"},
			)
		},
		DenyHandler: func(c echo.Context, identifier string, _ error) error {
			var retryAfter = rateLimitDuration
			key := middlewarex.LimiterCacheKey(identifier)
			ctx := context.Background()
			state, err := cache.Get[middlewarex.LimiterState](cacheStore, ctx, key)
			if err == nil {
				elapsed := time.Since(state.LastRefill)
				if elapsed < rateLimitDuration {
					retryAfter = rateLimitDuration - elapsed
				}
			}

			retryAfter = retryAfter.Truncate(time.Second)
			c.Response().Header().Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))

			errStr := fmt.Sprintf("rate limit exceeded, try again in %v", retryAfter)
			return c.JSON(
				http.StatusTooManyRequests,
				map[string]string{"error": errStr},
			)
		},
	}
	e.Use(middleware.RateLimiterWithConfig(rateLimiterCfg))

	e.Use(sentryecho.New(sentryecho.Options{
		Repanic:         true,
		WaitForDelivery: false,
	}))

	e.Use(middlewarex.CacheMiddleware(cacheStore, cacheTTL))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodOptions, http.MethodHead},
		AllowHeaders: []string{"*"},
	}))

	return api
}

func (api *API) Start(cfg *config.API) error {
	addr := fmt.Sprintf(":%d", cfg.Port)
	return api.echo.Start(addr)
}
