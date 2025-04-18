package cache

import (
	"context"
	"core/internal/config"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go"
)

var ErrNotFound = errors.New("key not found in cache")

type Cache struct {
	client valkey.Client
}

func New(cfg *config.Cache) (*Cache, error) {
	option := valkey.ClientOption{
		InitAddress: []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
	}
	client, err := valkey.NewClient(option)
	if err != nil {
		return nil, fmt.Errorf("cannot init cache: %w", err)
	}

	return &Cache{client: client}, nil
}

func (c *Cache) Close() error {
	c.client.Close()
	return nil
}

func Set[T any](c *Cache, ctx context.Context, key string, value T, ttl time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cannot marshal value: %w", err)
	}

	if ttl == 0 {
		ttl = 5 * time.Minute
	}

	cmd := c.client.B().Set().Key(key).Value(string(bytes)).Ex(ttl).Build()
	err = c.client.Do(ctx, cmd).Error()
	if err != nil {
		return fmt.Errorf("cannot set value: %w", err)
	}

	return nil
}

func Get[T any](c *Cache, ctx context.Context, key string) (T, error) {
	var value T

	cmd := c.client.B().Get().Key(key).Build()
	resp, err := c.client.Do(ctx, cmd).AsBytes()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return value, ErrNotFound
		}

		return value, fmt.Errorf("cannot get key: %w", err)
	}

	if err := json.Unmarshal(resp, &value); err != nil {
		return value, fmt.Errorf("cannot unmarshal value: %w", err)
	}

	return value, nil
}

func Delete(c *Cache, ctx context.Context, key string) error {
	cmd := c.client.B().Del().Key(key).Build()
	err := c.client.Do(ctx, cmd).Error()
	if err != nil {
		return fmt.Errorf("cannot delete key: %w", err)
	}

	return nil
}
