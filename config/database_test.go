package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseConfig_URL(t *testing.T) {
	tests := []struct {
		name           string
		config         DatabaseConfig
		expectedURL    string
		expectedFormat string
	}{
		{
			name: "when all fields are provided it should return valid PostgreSQL connection URL",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "postgres",
				Password: "password",
				DBName:   "payments_db",
			},
			expectedURL:    "postgres://postgres:password@localhost:5432/payments_db?sslmode=disable",
			expectedFormat: "postgres://%s:%s@%s:%s/%s?sslmode=disable",
		},
		{
			name: "when using custom host and port it should return URL with custom values",
			config: DatabaseConfig{
				Host:     "db.example.com",
				Port:     "5433",
				User:     "admin",
				Password: "secret123",
				DBName:   "production_db",
			},
			expectedURL:    "postgres://admin:secret123@db.example.com:5433/production_db?sslmode=disable",
			expectedFormat: "postgres://%s:%s@%s:%s/%s?sslmode=disable",
		},
		{
			name: "when password contains special characters it should return URL with encoded password",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "postgres",
				Password: "p@ssw0rd#123",
				DBName:   "payments_db",
			},
			expectedURL:    "postgres://postgres:p@ssw0rd#123@localhost:5432/payments_db?sslmode=disable",
			expectedFormat: "postgres://%s:%s@%s:%s/%s?sslmode=disable",
		},
		{
			name: "when database name contains underscores it should return URL with underscores",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "postgres",
				Password: "password",
				DBName:   "payment_service_db",
			},
			expectedURL:    "postgres://postgres:password@localhost:5432/payment_service_db?sslmode=disable",
			expectedFormat: "postgres://%s:%s@%s:%s/%s?sslmode=disable",
		},
		{
			name: "when using default PostgreSQL port it should return URL with port 5432",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "postgres",
				Password: "password",
				DBName:   "test_db",
			},
			expectedURL:    "postgres://postgres:password@localhost:5432/test_db?sslmode=disable",
			expectedFormat: "postgres://%s:%s@%s:%s/%s?sslmode=disable",
		},
		{
			name: "when user contains special characters it should return URL with user",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "user@domain",
				Password: "password",
				DBName:   "payments_db",
			},
			expectedURL:    "postgres://user@domain:password@localhost:5432/payments_db?sslmode=disable",
			expectedFormat: "postgres://%s:%s@%s:%s/%s?sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (Config already prepared in test struct)

			// Act
			result := tt.config.URL()

			// Assert
			assert.Equal(t, tt.expectedURL, result)
			assert.Contains(t, result, "postgres://")
			assert.Contains(t, result, "?sslmode=disable")
			assert.Contains(t, result, tt.config.User)
			assert.Contains(t, result, tt.config.Password)
			assert.Contains(t, result, tt.config.Host)
			assert.Contains(t, result, tt.config.Port)
			assert.Contains(t, result, tt.config.DBName)
		})
	}
}
