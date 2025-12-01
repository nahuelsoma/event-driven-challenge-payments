package messagebroker

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/streadway/amqp"
)

// MessageHandler handles incoming messages
type MessageHandler interface {
	HandleMessage(body []byte) error
}

// ConsumerConfig configures a consumer
type ConsumerConfig struct {
	Exchange   string // Exchange name for topic-based routing
	QueueName  string // Queue name to consume from
	RoutingKey string // Routing key for binding (usually same as queue name)
	Workers    int
}

// Consumer consumes messages from RabbitMQ
type Consumer struct {
	channel *Channel
	config  ConsumerConfig
}

// NewConsumer creates a new consumer
func NewConsumer(channel *Channel, config ConsumerConfig) (*Consumer, error) {
	if channel == nil || channel.ch == nil {
		return nil, errors.New("consumer: channel cannot be nil")
	}
	if config.Exchange == "" {
		return nil, errors.New("consumer: exchange name cannot be empty")
	}
	if config.QueueName == "" {
		return nil, errors.New("consumer: queue name cannot be empty")
	}
	if config.RoutingKey == "" {
		config.RoutingKey = config.QueueName // Default routing key to queue name
	}
	if config.Workers < 1 {
		config.Workers = 1
	}

	return &Consumer{
		channel: channel,
		config:  config,
	}, nil
}

// Start starts consuming messages
func (c *Consumer) Start(handler MessageHandler) error {
	// Declare topic exchange for flexible routing
	err := c.channel.ch.ExchangeDeclare(
		c.config.Exchange, // exchange name
		"topic",           // topic exchange allows pattern-based routing
		true,              // durable
		false,             // auto-delete
		false,             // internal
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare queue
	_, err = c.channel.ch.QueueDeclare(
		c.config.QueueName, // queue name
		true,               // durable
		false,              // auto-delete
		false,              // exclusive
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange with routing key
	err = c.channel.ch.QueueBind(
		c.config.QueueName,  // queue name
		c.config.RoutingKey, // routing key
		c.config.Exchange,   // exchange
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue to exchange: %w", err)
	}

	// Set prefetch count
	if err := c.channel.ch.Qos(c.config.Workers*2, 0, false); err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// Start consuming
	msgs, err := c.channel.ch.Consume(
		c.config.QueueName,
		"",    // consumer tag (auto-generated)
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	// Start workers
	for i := 0; i < c.config.Workers; i++ {
		go c.worker(i, msgs, handler)
	}

	return nil
}

func (c *Consumer) worker(id int, msgs <-chan amqp.Delivery, handler MessageHandler) {
	for msg := range msgs {
		if err := handler.HandleMessage(msg.Body); err != nil {
			slog.Error("Worker failed to handle message", "worker_id", id, "error", err)
			msg.Nack(false, true) // requeue
		} else {
			msg.Ack(false)
		}
	}
}
