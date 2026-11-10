package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr        string
	ShutdownTimeout time.Duration
	DatabaseURL     string
}

func Load() (Config, error) {
	if err := godotenv.Load("config.env"); err != nil {
		return Config{}, fmt.Errorf("read config.env: %w", err)
	}

	shutdownTimeout, err := getEnvDuration("SHUTDOWN_TIMEOUT", 10*time.Second)
	if err != nil {
		return Config{}, fmt.Errorf("read SHUTDOWN_TIMEOUT: %w", err)
	}

	cfg := Config{
		HTTPAddr:        strings.TrimSpace(os.Getenv("HTTP_ADDR")),
		ShutdownTimeout: shutdownTimeout,
		DatabaseURL:     strings.TrimSpace(os.Getenv("DATABASE_URL")),
	}

	if cfg.HTTPAddr == "" {
		return Config{}, fmt.Errorf("HTTP_ADDR is required")
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	if cfg.ShutdownTimeout <= 0 {
		return Config{}, fmt.Errorf("SHUTDOWN_TIMEOUT must be positive")
	}

	return cfg, nil
}

func getEnvDuration(key string, defaultValue time.Duration) (time.Duration, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue, nil
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}

	return duration, nil
}
