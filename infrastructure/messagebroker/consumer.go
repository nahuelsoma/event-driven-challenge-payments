package messagebroker

import (
	"errors"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// MessageHandler handles incoming messages
type MessageHandler interface {
	HandleMessage(body []byte) error
}

// ConsumerConfig configures a consumer
type ConsumerConfig struct {
	QueueName string
	Workers   int
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
	if config.QueueName == "" {
		return nil, errors.New("consumer: queue name cannot be empty")
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
	// Declare queue
	_, err := c.channel.ch.QueueDeclare(
		c.config.QueueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
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
	log.Printf("Worker %d started", id)
	for msg := range msgs {
		if err := handler.HandleMessage(msg.Body); err != nil {
			log.Printf("Worker %d: error handling message: %v", id, err)
			msg.Nack(false, true) // requeue
		} else {
			msg.Ack(false)
		}
	}
	log.Printf("Worker %d stopped (channel closed)", id)
}
