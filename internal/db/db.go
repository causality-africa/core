package db

import (
	"context"
	"core/internal/config"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("row not found")

type DB struct {
	pool *pgxpool.Pool
}

func New(cfg *config.DB) (*DB, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("cannot create db pool: %w", err)
	}

	return &DB{pool: pool}, nil
}

func (s *DB) Close() error {
	s.pool.Close()
	return nil
}
