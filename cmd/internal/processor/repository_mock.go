package processor

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockGatewayProcessor is a mock implementation of GatewayProcessor for testing
type MockGatewayProcessor struct {
	mock.Mock
}

// Process mocks the Process method
func (m *MockGatewayProcessor) Process(ctx context.Context, paymentID string, amount float64) (string, error) {
	args := m.Called(ctx, paymentID, amount)
	return args.String(0), args.Error(1)
}
