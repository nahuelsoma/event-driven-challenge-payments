package finder

import (
	"context"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
	"github.com/stretchr/testify/mock"
)

// MockPaymentFinder is a mock implementation of PaymentFinder
type MockPaymentFinder struct {
	mock.Mock
}

// Find mocks the Find method
func (m *MockPaymentFinder) Find(ctx context.Context, filter *PaymentFilter) (*domain.Payment, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

// FindEvents mocks the FindEvents method
func (m *MockPaymentFinder) FindEvents(ctx context.Context, paymentID string) ([]*domain.Event, error) {
	args := m.Called(ctx, paymentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Event), args.Error(1)
}

