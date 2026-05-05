package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	defaultHTTPPort = "8080"
	defaultDBPath   = "./data/app.db"
)

type Config struct {
	HTTPPort string
	DBPath   string
}

func Load() (Config, error) {
	cfg := Config{
		HTTPPort: valueOrDefault(os.Getenv("HTTP_PORT"), defaultHTTPPort),
		DBPath:   valueOrDefault(os.Getenv("DB_PATH"), defaultDBPath),
	}

	port, err := strconv.Atoi(cfg.HTTPPort)
	if err != nil || port < 1 || port > 65535 {
		return Config{}, fmt.Errorf("HTTP_PORT must be a number between 1 and 65535")
	}

	return cfg, nil
}

func valueOrDefault(value string, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}
