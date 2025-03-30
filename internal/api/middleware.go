package api

import (
	"context"
	"core/internal/cache"
	"encoding/json"
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

type limiterState struct {
	Remaining  float64   `json:"remaining"`
	LastRefill time.Time `json:"last_refill"`
}

func (store *RateLimiterCacheStore) Allow(identifier string) (bool, error) {
	key := fmt.Sprintf("core:rate-limit:%s", identifier)
	var state limiterState

	ctx := context.Background()
	rawState, err := store.cache.Get(ctx, key)
	if err != nil {
		if !errors.Is(err, cache.ErrNotFound) {
			slog.Error("cannot check rate limit", "error", err)
			return false, fmt.Errorf("cannot check rate limit: %w", err)
		}
	} else {
		stateBytes, _ := json.Marshal(rawState)
		if err := json.Unmarshal(stateBytes, &state); err != nil {
			slog.Error("invalid rate limiter state", "error", err)
			return false, fmt.Errorf("invalid rate limiter state: %w", err)
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

	err = store.cache.Set(ctx, key, state, store.every-elapsed)
	if err != nil {
		slog.Error("cannot update rate limit", "error", err)
		return false, fmt.Errorf("cannot update rate limit: %w", err)
	}

	return true, nil
}
