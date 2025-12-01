package messagebroker

import "errors"

// MessageBrokerConnection represents a Message Broker connection
type MessageBrokerConnection struct {
	config interface{}
}

// NewMessageBrokerConnection creates a new Message Broker connection
// TODO: Implement the logic to create a new Message Broker connection
func NewMessageBrokerConnection(config interface{}) (*MessageBrokerConnection, error) {
	if config == nil {
		return nil, errors.New("message broker connection: config cannot be nil")
	}
	return &MessageBrokerConnection{config: config}, nil
}
