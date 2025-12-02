package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageBrokerConfig_URL(t *testing.T) {
	tests := []struct {
		name           string
		config         MessageBrokerConfig
		expectedURL    string
		expectedFormat string
	}{
		{
			name: "when all fields are provided it should return valid RabbitMQ connection URL",
			config: MessageBrokerConfig{
				Host:     "localhost",
				Port:     "5672",
				User:     "guest",
				Password: "guest",
			},
			expectedURL:    "amqp://guest:guest@localhost:5672",
			expectedFormat: "amqp://%s:%s@%s:%s",
		},
		{
			name: "when using custom host and port it should return URL with custom values",
			config: MessageBrokerConfig{
				Host:     "rabbitmq.example.com",
				Port:     "5673",
				User:     "admin",
				Password: "secret123",
			},
			expectedURL:    "amqp://admin:secret123@rabbitmq.example.com:5673",
			expectedFormat: "amqp://%s:%s@%s:%s",
		},
		{
			name: "when password contains special characters it should return URL with password",
			config: MessageBrokerConfig{
				Host:     "localhost",
				Port:     "5672",
				User:     "guest",
				Password: "p@ssw0rd#123",
			},
			expectedURL:    "amqp://guest:p@ssw0rd#123@localhost:5672",
			expectedFormat: "amqp://%s:%s@%s:%s",
		},
		{
			name: "when using default RabbitMQ port it should return URL with port 5672",
			config: MessageBrokerConfig{
				Host:     "localhost",
				Port:     "5672",
				User:     "guest",
				Password: "guest",
			},
			expectedURL:    "amqp://guest:guest@localhost:5672",
			expectedFormat: "amqp://%s:%s@%s:%s",
		},
		{
			name: "when user contains special characters it should return URL with user",
			config: MessageBrokerConfig{
				Host:     "localhost",
				Port:     "5672",
				User:     "user@domain",
				Password: "password",
			},
			expectedURL:    "amqp://user@domain:password@localhost:5672",
			expectedFormat: "amqp://%s:%s@%s:%s",
		},
		{
			name: "when using production credentials it should return URL with production values",
			config: MessageBrokerConfig{
				Host:     "rabbitmq.prod.example.com",
				Port:     "5672",
				User:     "prod_user",
				Password: "secure_password_123",
			},
			expectedURL:    "amqp://prod_user:secure_password_123@rabbitmq.prod.example.com:5672",
			expectedFormat: "amqp://%s:%s@%s:%s",
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
			assert.Contains(t, result, "amqp://")
			assert.Contains(t, result, tt.config.User)
			assert.Contains(t, result, tt.config.Password)
			assert.Contains(t, result, tt.config.Host)
			assert.Contains(t, result, tt.config.Port)
		})
	}
}

