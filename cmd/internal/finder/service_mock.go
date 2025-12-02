package finder

import (
	"context"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
	"github.com/stretchr/testify/mock"
)

// MockPaymentFinderService is a mock implementation of PaymentFinderService
type MockPaymentFinderService struct {
	mock.Mock
}

// Find finds a payment by ID
func (m *MockPaymentFinderService) Find(ctx context.Context, filter *PaymentFilter) (*domain.Payment, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

// FindEvents finds all events for a payment by ID
func (m *MockPaymentFinderService) FindEvents(ctx context.Context, paymentID string) ([]*domain.Event, error) {
	args := m.Called(ctx, paymentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Event), args.Error(1)
}
