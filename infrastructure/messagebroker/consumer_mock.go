package messagebroker

import "github.com/stretchr/testify/mock"

// MockConsumer is a mock implementation of MessageConsumer for testing
type MockConsumer struct {
	mock.Mock
}

// Start mocks the Start method
func (m *MockConsumer) Start(handler MessageHandler) error {
	args := m.Called(handler)
	return args.Error(0)
}
