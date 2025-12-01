package config

import (
	"errors"
	"os"
	"strings"
)

// Config holds all application configuration
type Config struct {
	Database      DatabaseConfig
	MessageBroker MessageBrokerConfig
}

// Load reads environment variables and returns the application configuration
func Load() (*Config, error) {
	var missingVars []string

	dbConfig := loadDatabaseConfig(&missingVars)
	messageBrokerConfig := loadMessageBrokerConfig(&missingVars)

	if len(missingVars) > 0 {
		return nil, errors.New("missing required environment variables: " + strings.Join(missingVars, ", "))
	}

	return &Config{
		Database:      dbConfig,
		MessageBroker: messageBrokerConfig,
	}, nil
}

// getRequiredEnv retrieves an environment variable and tracks if it's missing
func getRequiredEnv(key string, missingVars *[]string) string {
	value := os.Getenv(key)
	if value == "" {
		*missingVars = append(*missingVars, key)
	}
	return value
}
