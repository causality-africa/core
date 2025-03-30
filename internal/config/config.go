package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type API struct {
	Port int
}

type DB struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

type Cache struct {
	Host string
	Port int
}

type Config struct {
	API
	DB
	Cache
}

func Load() (*Config, error) {
	k := koanf.New(".")
	k.Load(env.Provider("CORE_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "CORE_")), "_", ".", -1)
	}), nil)

	cfg := &Config{}
	if err := k.Unmarshal("", cfg); err != nil {
		return nil, fmt.Errorf("cannot decode config: %w", err)
	}

	return cfg, nil
}
