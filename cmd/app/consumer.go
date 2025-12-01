package app

import (
	"fmt"
	"log/slog"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/processor"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/database"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/messagebroker"
)

const (
	queueName = "payments"
	workers   = 3
)

// StartConsumer initializes and starts the message consumer
// Returns after setup is complete. Message consumption runs in background goroutines.
func StartConsumer(db *database.DB, conn *messagebroker.Connection, walletClient interface{}, gatewayClient interface{}) error {

	// Create channel for consumer
	channel, err := conn.NewChannel()
	if err != nil {
		return fmt.Errorf("consumer: failed to create channel: %w", err)
	}

	// Consumer configuration
	config := messagebroker.ConsumerConfig{
		QueueName: queueName,
		Workers:   workers,
	}

	// Create infrastructure consumer
	consumer, err := messagebroker.NewConsumer(channel, config)
	if err != nil {
		return fmt.Errorf("consumer: failed to create consumer: %w", err)
	}

	// Create payment processor handler using the vertical pattern
	handler, err := processor.Build(db, walletClient, gatewayClient)
	if err != nil {
		return fmt.Errorf("consumer: failed to create processor: %w", err)
	}

	slog.Info("Starting consumer", "queue", queueName, "workers", workers)

	// Start consuming
	return consumer.Start(handler)
}
