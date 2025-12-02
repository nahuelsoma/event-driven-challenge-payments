package config

import "fmt"

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// URL returns the PostgreSQL connection URL
func (c DatabaseConfig) URL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

// loadDatabaseConfig reads database configuration from environment variables
func loadDatabaseConfig(missingVars *[]string) DatabaseConfig {
	return DatabaseConfig{
		Host:     getRequiredEnv("DB_HOST", missingVars),
		Port:     getRequiredEnv("DB_PORT", missingVars),
		User:     getRequiredEnv("DB_USER", missingVars),
		Password: getRequiredEnv("DB_PASSWORD", missingVars),
		DBName:   getRequiredEnv("DB_NAME", missingVars),
	}
}
