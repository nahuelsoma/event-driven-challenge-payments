package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Required environment variables
	requiredVars := []string{
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"BROKER_HOST",
		"BROKER_PORT",
		"BROKER_USER",
		"BROKER_PASSWORD",
	}

	tests := []struct {
		name              string
		envVars           map[string]string
		expectedError     string
		expectedExchange  string
		expectedQueueName string
	}{
		{
			name: "when all environment variables are provided it should load config successfully and no error",
			envVars: map[string]string{
				"DB_HOST":        "localhost",
				"DB_PORT":        "5432",
				"DB_USER":        "postgres",
				"DB_PASSWORD":    "password",
				"DB_NAME":        "payments_db",
				"BROKER_HOST":    "rabbitmq",
				"BROKER_PORT":    "5672",
				"BROKER_USER":    "guest",
				"BROKER_PASSWORD": "guest",
			},
			expectedError:     "",
			expectedExchange:  "payments",
			expectedQueueName: "payments.created",
		},
		{
			name: "when DB_HOST is missing it should return error with missing variable",
			envVars: map[string]string{
				"DB_PORT":        "5432",
				"DB_USER":        "postgres",
				"DB_PASSWORD":    "password",
				"DB_NAME":        "payments_db",
				"BROKER_HOST":    "rabbitmq",
				"BROKER_PORT":    "5672",
				"BROKER_USER":    "guest",
				"BROKER_PASSWORD": "guest",
			},
			expectedError: "missing required environment variables: DB_HOST",
		},
		{
			name: "when DB_PORT is missing it should return error with missing variable",
			envVars: map[string]string{
				"DB_HOST":        "localhost",
				"DB_USER":        "postgres",
				"DB_PASSWORD":    "password",
				"DB_NAME":        "payments_db",
				"BROKER_HOST":    "rabbitmq",
				"BROKER_PORT":    "5672",
				"BROKER_USER":    "guest",
				"BROKER_PASSWORD": "guest",
			},
			expectedError: "missing required environment variables: DB_PORT",
		},
		{
			name: "when DB_USER is missing it should return error with missing variable",
			envVars: map[string]string{
				"DB_HOST":        "localhost",
				"DB_PORT":        "5432",
				"DB_PASSWORD":    "password",
				"DB_NAME":        "payments_db",
				"BROKER_HOST":    "rabbitmq",
				"BROKER_PORT":    "5672",
				"BROKER_USER":    "guest",
				"BROKER_PASSWORD": "guest",
			},
			expectedError: "missing required environment variables: DB_USER",
		},
		{
			name: "when DB_PASSWORD is missing it should return error with missing variable",
			envVars: map[string]string{
				"DB_HOST":        "localhost",
				"DB_PORT":        "5432",
				"DB_USER":        "postgres",
				"DB_NAME":        "payments_db",
				"BROKER_HOST":    "rabbitmq",
				"BROKER_PORT":    "5672",
				"BROKER_USER":    "guest",
				"BROKER_PASSWORD": "guest",
			},
			expectedError: "missing required environment variables: DB_PASSWORD",
		},
		{
			name: "when DB_NAME is missing it should return error with missing variable",
			envVars: map[string]string{
				"DB_HOST":        "localhost",
				"DB_PORT":        "5432",
				"DB_USER":        "postgres",
				"DB_PASSWORD":    "password",
				"BROKER_HOST":    "rabbitmq",
				"BROKER_PORT":    "5672",
				"BROKER_USER":    "guest",
				"BROKER_PASSWORD": "guest",
			},
			expectedError: "missing required environment variables: DB_NAME",
		},
		{
			name: "when BROKER_HOST is missing it should return error with missing variable",
			envVars: map[string]string{
				"DB_HOST":        "localhost",
				"DB_PORT":        "5432",
				"DB_USER":        "postgres",
				"DB_PASSWORD":    "password",
				"DB_NAME":        "payments_db",
				"BROKER_PORT":    "5672",
				"BROKER_USER":    "guest",
				"BROKER_PASSWORD": "guest",
			},
			expectedError: "missing required environment variables: BROKER_HOST",
		},
		{
			name: "when BROKER_PORT is missing it should return error with missing variable",
			envVars: map[string]string{
				"DB_HOST":        "localhost",
				"DB_PORT":        "5432",
				"DB_USER":        "postgres",
				"DB_PASSWORD":    "password",
				"DB_NAME":        "payments_db",
				"BROKER_HOST":    "rabbitmq",
				"BROKER_USER":    "guest",
				"BROKER_PASSWORD": "guest",
			},
			expectedError: "missing required environment variables: BROKER_PORT",
		},
		{
			name: "when BROKER_USER is missing it should return error with missing variable",
			envVars: map[string]string{
				"DB_HOST":        "localhost",
				"DB_PORT":        "5432",
				"DB_USER":        "postgres",
				"DB_PASSWORD":    "password",
				"DB_NAME":        "payments_db",
				"BROKER_HOST":    "rabbitmq",
				"BROKER_PORT":    "5672",
				"BROKER_PASSWORD": "guest",
			},
			expectedError: "missing required environment variables: BROKER_USER",
		},
		{
			name: "when BROKER_PASSWORD is missing it should return error with missing variable",
			envVars: map[string]string{
				"DB_HOST":        "localhost",
				"DB_PORT":        "5432",
				"DB_USER":        "postgres",
				"DB_PASSWORD":    "password",
				"DB_NAME":        "payments_db",
				"BROKER_HOST":    "rabbitmq",
				"BROKER_PORT":    "5672",
				"BROKER_USER":    "guest",
			},
			expectedError: "missing required environment variables: BROKER_PASSWORD",
		},
		{
			name: "when multiple environment variables are missing it should return error with all missing variables",
			envVars: map[string]string{
				"DB_HOST":        "localhost",
				"DB_PORT":        "5432",
				"BROKER_HOST":    "rabbitmq",
				"BROKER_PORT":    "5672",
			},
			expectedError: "missing required environment variables:",
		},
		{
			name: "when all environment variables are empty it should return error with all missing variables",
			envVars: map[string]string{
				"DB_HOST":        "",
				"DB_PORT":        "",
				"DB_USER":        "",
				"DB_PASSWORD":    "",
				"DB_NAME":        "",
				"BROKER_HOST":    "",
				"BROKER_PORT":    "",
				"BROKER_USER":    "",
				"BROKER_PASSWORD": "",
			},
			expectedError: "missing required environment variables:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// Save current environment variables
			savedEnv := make(map[string]string)
			for _, key := range requiredVars {
				if val, exists := os.LookupEnv(key); exists {
					savedEnv[key] = val
				}
			}

			// Clear all required environment variables
			for _, key := range requiredVars {
				os.Unsetenv(key)
			}

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Cleanup: restore original environment
			defer func() {
				for _, key := range requiredVars {
					os.Unsetenv(key)
				}
				for key, value := range savedEnv {
					os.Setenv(key, value)
				}
			}()

			// Act
			result, err := Load()

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				if strings.Contains(tt.expectedError, ",") || strings.HasSuffix(tt.expectedError, ":") {
					// Multiple missing variables or all missing
					assert.Contains(t, err.Error(), tt.expectedError)
				} else {
					assert.Equal(t, tt.expectedError, err.Error())
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedExchange, result.Exchange)
				assert.Equal(t, tt.expectedQueueName, result.QueueName)
				assert.Equal(t, tt.envVars["DB_HOST"], result.Database.Host)
				assert.Equal(t, tt.envVars["DB_PORT"], result.Database.Port)
				assert.Equal(t, tt.envVars["DB_USER"], result.Database.User)
				assert.Equal(t, tt.envVars["DB_PASSWORD"], result.Database.Password)
				assert.Equal(t, tt.envVars["DB_NAME"], result.Database.DBName)
				assert.Equal(t, tt.envVars["BROKER_HOST"], result.MessageBroker.Host)
				assert.Equal(t, tt.envVars["BROKER_PORT"], result.MessageBroker.Port)
				assert.Equal(t, tt.envVars["BROKER_USER"], result.MessageBroker.User)
				assert.Equal(t, tt.envVars["BROKER_PASSWORD"], result.MessageBroker.Password)
			}
		})
	}
}

func TestGetRequiredEnv(t *testing.T) {
	tests := []struct {
		name          string
		envKey        string
		envValue      string
		expectedValue string
		shouldTrack   bool
	}{
		{
			name:          "when environment variable is set it should return value and not track as missing",
			envKey:         "TEST_VAR",
			envValue:       "test_value",
			expectedValue:  "test_value",
			shouldTrack:    false,
		},
		{
			name:          "when environment variable is empty it should return empty and track as missing",
			envKey:         "TEST_VAR",
			envValue:       "",
			expectedValue:  "",
			shouldTrack:    true,
		},
		{
			name:          "when environment variable is not set it should return empty and track as missing",
			envKey:         "TEST_VAR_NOT_SET",
			envValue:       "", // Not set
			expectedValue:  "",
			shouldTrack:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var missingVars []string

			// Set or unset environment variable
			if tt.envValue != "" || tt.envKey == "TEST_VAR" {
				if tt.envValue == "" {
					os.Unsetenv(tt.envKey)
				} else {
					os.Setenv(tt.envKey, tt.envValue)
				}
			} else {
				os.Unsetenv(tt.envKey)
			}

			// Cleanup
			defer os.Unsetenv(tt.envKey)

			// Act
			result := getRequiredEnv(tt.envKey, &missingVars)

			// Assert
			assert.Equal(t, tt.expectedValue, result)
			if tt.shouldTrack {
				assert.Contains(t, missingVars, tt.envKey)
			} else {
				assert.NotContains(t, missingVars, tt.envKey)
			}
		})
	}
}

