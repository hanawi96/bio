package config

import (
	"os"
	"strings"
)

type Config struct {
	DatabaseURL    string
	JWTSecret      string
	AllowedOrigins string
	Environment    string
}

func New() *Config {
	return &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://linkbio:linkbio_dev_password@localhost:5432/linkbio?sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:5173"),
		Environment:    getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.TrimSpace(value)
}
