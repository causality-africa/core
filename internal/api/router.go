package api

import (
	"core/internal/api/middlewarex"
	"core/internal/cache"
	"core/internal/config"
	"core/internal/db"
	"fmt"
	"net/http"
	"time"

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

func New(database *db.DB, cache *cache.Cache) *API {
	e := echo.New()
	api := &API{echo: e, db: database, cache: cache}

	e.GET("/v1/query", api.query)

	// Rate limiting middleware
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

	// Caching middleware
	e.Use(middlewarex.CacheMiddleware(cache, cacheTTL))

	return api
}

func (api *API) Start(cfg *config.API) error {
	addr := fmt.Sprintf(":%d", cfg.Port)
	return api.echo.Start(addr)
}
