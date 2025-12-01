package messagebroker

import (
	"errors"

	"github.com/streadway/amqp"
)

// PublisherConfig configures a publisher
type PublisherConfig struct {
	Exchange   string
	RoutingKey string
}

// Publisher publishes messages to RabbitMQ
type Publisher struct {
	channel *Channel
	config  PublisherConfig
}

// NewPublisher creates a new publisher
func NewPublisher(channel *Channel, config PublisherConfig) (*Publisher, error) {
	if channel == nil || channel.ch == nil {
		return nil, errors.New("publisher: channel cannot be nil")
	}

	return &Publisher{
		channel: channel,
		config:  config,
	}, nil
}

// Publish publishes a JSON message
func (p *Publisher) Publish(body []byte) error {
	return p.channel.ch.Publish(
		p.config.Exchange,
		p.config.RoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}
