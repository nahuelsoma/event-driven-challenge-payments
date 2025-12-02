package config

import "fmt"

// MessageBrokerConfig holds RabbitMQ configuration
type MessageBrokerConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

// URL returns the RabbitMQ connection URL
func (c MessageBrokerConfig) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", c.User, c.Password, c.Host, c.Port)
}

// loadMessageBrokerConfig reads message broker configuration from environment variables
func loadMessageBrokerConfig(missingVars *[]string) MessageBrokerConfig {
	return MessageBrokerConfig{
		Host:     getRequiredEnv("BROKER_HOST", missingVars),
		Port:     getRequiredEnv("BROKER_PORT", missingVars),
		User:     getRequiredEnv("BROKER_USER", missingVars),
		Password: getRequiredEnv("BROKER_PASSWORD", missingVars),
	}
}
