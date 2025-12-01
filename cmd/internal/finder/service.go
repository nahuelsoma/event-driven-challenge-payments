package finder

import (
	"context"
	"errors"
	"fmt"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
)

// ErrPaymentNotFound is returned when a payment is not found
var ErrPaymentNotFound = errors.New("payment not found")

// PaymentReader interface for reading payments
type PaymentReader interface {
	GetByID(ctx context.Context, paymentID string) (*domain.Payment, error)
}

// PaymentFinderService handles payment finding business logic
type PaymentFinderService struct {
	paymentReader PaymentReader
}

// NewPaymentFinderService creates a new PaymentFinderService
func NewPaymentFinderService(pr PaymentReader) (*PaymentFinderService, error) {
	if pr == nil {
		return nil, errors.New("payment finder: reader cannot be nil")
	}

	return &PaymentFinderService{
		paymentReader: pr,
	}, nil
}

// Find finds a payment by ID
func (pfs *PaymentFinderService) Find(ctx context.Context, filter *PaymentFilter) (*domain.Payment, error) {
	payment, err := pfs.paymentReader.GetByID(ctx, filter.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("payment finder: failed to get payment: %w", err)
	}

	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	return payment, nil
}
