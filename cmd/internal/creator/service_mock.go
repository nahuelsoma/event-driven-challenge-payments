package creator

import (
	"context"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
	"github.com/stretchr/testify/mock"
)

// MockPaymentCreatorService is a mock implementation of PaymentCreatorService
type MockPaymentCreatorService struct {
	mock.Mock
}

// Create creates a new payment
func (m *MockPaymentCreatorService) Create(ctx context.Context, idempotencyKey string, pr *PaymentRequest) (*domain.Payment, error) {
	args := m.Called(ctx, idempotencyKey, pr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}
