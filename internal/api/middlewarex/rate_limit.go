package middlewarex

import (
	"context"
	"core/internal/cache"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4/middleware"
)

var _ middleware.RateLimiterStore = (*RateLimiterCacheStore)(nil)

type RateLimiterCacheStore struct {
	rate  float64
	every time.Duration

	cache *cache.Cache
}

func NewRateLimiterCacheStore(
	rate float64,
	every time.Duration,
	cache *cache.Cache,
) *RateLimiterCacheStore {
	return &RateLimiterCacheStore{
		rate:  rate,
		every: every,
		cache: cache,
	}
}

type LimiterState struct {
	Remaining  float64   `json:"remaining"`
	LastRefill time.Time `json:"last_refill"`
}

func LimiterCacheKey(identifier string) string {
	hasher := sha256.New()
	hasher.Write([]byte(identifier))
	return fmt.Sprintf("core:rate-limits:%s", hex.EncodeToString(hasher.Sum(nil)))
}

func (store *RateLimiterCacheStore) Allow(identifier string) (bool, error) {
	key := LimiterCacheKey(identifier)

	ctx := context.Background()
	state, err := cache.Get[LimiterState](store.cache, ctx, key)
	if err != nil {
		if !errors.Is(err, cache.ErrNotFound) {
			slog.Error("cannot check rate limit", "error", err)
			return false, fmt.Errorf("cannot check rate limit: %w", err)
		}
	}

	if time.Since(state.LastRefill) > store.every {
		state.Remaining = store.rate
		state.LastRefill = time.Now()
	}

	if state.Remaining < 1 {
		return false, nil // Rate limit exceeded
	}

	state.Remaining -= 1
	elapsed := time.Since(state.LastRefill)

	err = cache.Set(store.cache, ctx, key, &state, store.every-elapsed)
	if err != nil {
		slog.Error("cannot update rate limit", "error", err)
		return false, fmt.Errorf("cannot update rate limit: %w", err)
	}

	return true, nil
}
