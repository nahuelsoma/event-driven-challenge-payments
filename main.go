package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/app"
	"github.com/nahuelsoma/event-driven-challenge-payments/config"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/database"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/http"
	messagebroker "github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/message_broker"
)

func main() {
	// Configure slog to show DEBUG level logs
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// Create database connection
	dbConn, err := database.NewPostgresConnection(cfg.Database.URL())
	if err != nil {
		log.Fatalf("api: failed to create database connection: %v", err)
	}
	defer dbConn.Close()

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
	messageBrokerConn, err := messagebroker.Connect(cfg.MessageBroker.URL())
	if err != nil {
		log.Fatalf("api: failed to create RabbitMQ connection: %v", err)
	}

	if err := app.StartAPI(dbConn, httpClient, messageBrokerConn); err != nil {
		log.Fatalf("api: failed to start API: %v", err)
	}
}
