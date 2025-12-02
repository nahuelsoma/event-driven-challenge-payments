package processor

import (
	"context"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
	"github.com/stretchr/testify/mock"
)

// MockPaymentProcessorService is a mock implementation of PaymentProcessor for testing
type MockPaymentProcessorService struct {
	mock.Mock
}

// Process mocks the Process method
func (m *MockPaymentProcessorService) Process(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}
