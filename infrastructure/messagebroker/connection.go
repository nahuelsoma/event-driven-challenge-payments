package messagebroker

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

// Connection represents a TCP connection to RabbitMQ
type Connection struct {
	conn *amqp.Connection
}

// Validate validates the connection
func (c *Connection) Validate() error {
	if c.conn == nil {
		return errors.New("connection: connection cannot be nil")
	}
	return nil
}

// Connect establishes a TCP connection to RabbitMQ
func Connect(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return &Connection{conn: conn}, nil
}

// Close closes the TCP connection
func (c *Connection) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// NewChannel creates a new channel from this connection
func (c *Connection) NewChannel() (*Channel, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Channel{ch: ch}, nil
}

// Channel represents a virtual connection multiplexed over a TCP connection
type Channel struct {
	ch *amqp.Channel
}

// Close closes the channel
func (c *Channel) Close() error {
	if c.ch != nil {
		return c.ch.Close()
	}
	return nil
}
