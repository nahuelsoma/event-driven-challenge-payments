package app

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/processor"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/database"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/messagebroker"
)

const (
	workers = 3
)

// StartConsumer initializes and starts the message consumer
// Returns after setup is complete. Message consumption runs in background goroutines.
func StartConsumer(db *database.DB, walletClient *http.Client, conn *messagebroker.Connection, exchange, queueName string) error {
	// Create channel for consumer
	channel, err := conn.NewChannel()
	if err != nil {
		return fmt.Errorf("consumer: failed to create channel: %w", err)
	}

	// Consumer configuration with topic exchange for flexible routing
	config := messagebroker.ConsumerConfig{
		Exchange:   exchange,
		QueueName:  queueName,
		RoutingKey: queueName, // Use queue name as routing key
		Workers:    workers,
	}

	// Create infrastructure consumer
	consumer, err := messagebroker.NewConsumer(channel, config)
	if err != nil {
		return fmt.Errorf("consumer: failed to create consumer: %w", err)
	}

	// Create payment processor handler using the vertical pattern
	handler, err := processor.Build(db, walletClient)
	if err != nil {
		return fmt.Errorf("consumer: failed to create processor: %w", err)
	}

	// Start consuming
	if err := consumer.Start(handler); err != nil {
		return fmt.Errorf("consumer: failed to start consumer: %w", err)
	}

	slog.Info("Consumer started", "queue", queueName, "workers", workers)

	return nil
}
