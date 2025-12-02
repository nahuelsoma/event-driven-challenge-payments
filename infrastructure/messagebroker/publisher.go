package messagebroker

import (
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
func NewPublisher(conn *Connection, config PublisherConfig) (*Publisher, error) {
	if err := conn.Validate(); err != nil {
		return nil, err
	}

	channel, err := conn.NewChannel()
	if err != nil {
		return nil, err
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
