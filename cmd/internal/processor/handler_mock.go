package processor

import "github.com/stretchr/testify/mock"

// MockHandler is a mock implementation of Handler for testing
type MockHandler struct {
	mock.Mock
}

// HandleMessage mocks the HandleMessage method
func (m *MockHandler) HandleMessage(body []byte) error {
	args := m.Called(body)
	return args.Error(0)
}
