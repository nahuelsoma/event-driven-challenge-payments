package creator

import (
	"context"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
	"github.com/stretchr/testify/mock"
)

// MockPaymentPublisherRepository is a mock implementation of PaymentPublisherRepository
type MockPaymentPublisherRepository struct {
	mock.Mock
}

// Publish publishes a payment
func (m *MockPaymentPublisherRepository) Publish(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}
