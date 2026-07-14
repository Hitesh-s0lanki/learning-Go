// Package config loads and validates runtime configuration from the
// environment. All configuration lives in ONE place so the rest of the app
// never reads os.Getenv directly.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds everything the server needs to run.
type Config struct {
	Port        string        // HTTP port to listen on
	DatabaseURL string        // Postgres DSN, e.g. postgres://user:pass@host:5432/db?sslmode=disable
	JWTSecret   string        // secret used to sign JWTs
	JWTTTL      time.Duration // how long an access token stays valid
	Env         string        // "development" | "production"
}

// Load reads configuration from environment variables, applies sensible
// defaults, and fails fast if a required value is missing.
func Load() (*Config, error) {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		Env:         getEnv("ENV", "development"),
		JWTTTL:      getEnvDuration("JWT_TTL", 24*time.Hour),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		// Allow either a Go duration ("24h") or a plain number of seconds.
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
		if secs, err := strconv.Atoi(v); err == nil {
			return time.Duration(secs) * time.Second
		}
	}
	return fallback
}
