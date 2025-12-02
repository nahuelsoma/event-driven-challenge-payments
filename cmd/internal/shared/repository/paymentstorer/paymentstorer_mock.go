package paymentstorer

import (
	"context"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
	"github.com/stretchr/testify/mock"
)

// MockPaymentRepository is a mock implementation of PaymentRepository for external testing
type MockPaymentRepository struct {
	mock.Mock
}

// GetByID retrieves a payment by ID
func (m *MockPaymentRepository) GetByID(ctx context.Context, paymentID string) (*domain.Payment, error) {
	args := m.Called(ctx, paymentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

// GetByIDempotencyKey retrieves a payment by idempotency key
func (m *MockPaymentRepository) GetByIDempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Payment, error) {
	args := m.Called(ctx, idempotencyKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

// Save saves a new payment with its initial event
func (m *MockPaymentRepository) Save(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

// UpdateStatus updates the payment status with optional gateway reference
func (m *MockPaymentRepository) UpdateStatus(ctx context.Context, paymentID string, status domain.Status, gatewayRef string) error {
	args := m.Called(ctx, paymentID, status, gatewayRef)
	return args.Error(0)
}

// GetEventsByPaymentID retrieves all events for a payment
func (m *MockPaymentRepository) GetEventsByPaymentID(ctx context.Context, paymentID string) ([]*domain.Event, error) {
	args := m.Called(ctx, paymentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Event), args.Error(1)
}
