package config

import (
	"errors"
	"os"
	"strings"
)

// Config holds all application configuration
type Config struct {
	Database DatabaseConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Load reads environment variables and returns the application configuration
func Load() (*Config, error) {
	var missingVars []string

	dbHost := getEnv("DB_HOST", &missingVars)
	dbPort := getEnv("DB_PORT", &missingVars)
	dbUser := getEnv("DB_USER", &missingVars)
	dbPassword := getEnv("DB_PASSWORD", &missingVars)
	dbName := getEnv("DB_NAME", &missingVars)
	dbSSLMode := getEnvOrDefault("DB_SSL_MODE", "disable")

	if len(missingVars) > 0 {
		return nil, errors.New("missing required environment variables: " + strings.Join(missingVars, ", "))
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
			SSLMode:  dbSSLMode,
		},
	}, nil
}

// getEnv retrieves an environment variable and tracks if it's missing
func getEnv(key string, missingVars *[]string) string {
	value := os.Getenv(key)
	if value == "" {
		*missingVars = append(*missingVars, key)
	}
	return value
}

// getEnvOrDefault retrieves an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
