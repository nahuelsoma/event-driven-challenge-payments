package main

import (
	"log"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/app"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/database"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/http"
	messagebroker "github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/message_broker"
)

func main() {
	// Create database connection
	dbConfig := map[string]string{
		"host":     "localhost",
		"port":     "5432",
		"user":     "postgres",
		"password": "postgres",
		"dbname":   "payments",
	}

	dbConn, err := database.NewSQLDatabase(dbConfig)
	if err != nil {
		log.Fatalf("api: failed to create database connection: %v", err)
	}

	// Create HTTP client
	httpConfig := map[string]string{
		"host": "localhost",
		"port": "3000",
	}

	httpClient, err := http.NewHTTPClient(httpConfig)
	if err != nil {
		log.Fatalf("api: failed to create HTTP client: %v", err)
	}

	// Create RabbitMQ connection
	messageBrokerConfig := map[string]string{
		"host":     "localhost",
		"port":     "5672",
		"user":     "guest",
		"password": "guest",
		"vhost":    "/",
	}

	messageBrokerConn, err := messagebroker.NewMessageBrokerConnection(messageBrokerConfig)
	if err != nil {
		log.Fatalf("api: failed to create RabbitMQ connection: %v", err)
	}

	if err := app.StartAPI(dbConn, httpClient, messageBrokerConn); err != nil {
		log.Fatalf("api: failed to start API: %v", err)
	}
}
