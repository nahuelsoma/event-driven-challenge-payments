package main

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/app"
	"github.com/nahuelsoma/event-driven-challenge-payments/config"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/database"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/messagebroker"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/restclient"
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
		log.Fatalf("main: failed to create database connection: %v", err)
	}
	defer dbConn.Close()

	// Create wallet HTTP client
	httpConfig := &restclient.Config{
		BaseURL: "http://test-wallet-service.com",
		Timeout: 300 * time.Millisecond,
	}

	walletClient, err := restclient.NewRestClient(httpConfig)
	if err != nil {
		log.Fatalf("main: failed to create HTTP client: %v", err)
	}

	// Create RabbitMQ connection
	messageBrokerConn, err := messagebroker.Connect(cfg.MessageBroker.URL())
	if err != nil {
		log.Fatalf("main: failed to create RabbitMQ connection: %v", err)
	}
	defer messageBrokerConn.Close()

	// Start consumer first (runs in background goroutines)
	if err := app.StartConsumer(dbConn, walletClient, messageBrokerConn, cfg.Exchange, cfg.QueueName); err != nil {
		log.Fatalf("main: failed to start consumer: %v", err)
	}

	// Start API server (blocks to keep the application running)
	if err := app.StartAPI(dbConn, walletClient, messageBrokerConn, cfg.Exchange, cfg.QueueName); err != nil {
		log.Fatalf("main: failed to start API: %v", err)
	}
}
