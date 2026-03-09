package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds application configuration from environment variables
type Config struct {
	Port          int
	PublicAPIPort int
	DatabaseURL   string
	JWTSecret     string
}

// Load reads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Port:          getEnvInt("PORT", 3000),
		PublicAPIPort: getEnvInt("PUBLIC_API_PORT", 0),
		DatabaseURL:   getEnvStr("DATABASE_URL", "postgres://localhost:5432/go_webserver?sslmode=disable"),
		JWTSecret:     getEnvStr("JWT_SECRET", "dev-secret-change-in-production"),
	}
}

// Addr returns the server listen address
func (c *Config) Addr() string {
	return fmt.Sprintf(":%d", c.Port)
}

func getEnvStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
