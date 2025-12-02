package messagebroker

import "github.com/stretchr/testify/mock"

// MockPublisher is a mock implementation of Publisher for testing
type MockPublisher struct {
	mock.Mock
}

// Publish publishes a message to the broker
func (m *MockPublisher) Publish(body []byte) error {
	args := m.Called(body)
	return args.Error(0)
}

