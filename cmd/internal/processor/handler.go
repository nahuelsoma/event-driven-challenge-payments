package processor

import (
	"context"
	"errors"
	"fmt"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
)

// PaymentProcessor defines the interface for payment processing business logic
type PaymentProcessor interface {
	Process(ctx context.Context, payment *domain.Payment) error
}

// Handler handles incoming payment messages from the queue
type Handler struct {
	paymentProcessor PaymentProcessor
}

// NewHandler creates a new processor handler
func NewHandler(pp PaymentProcessor) (*Handler, error) {
	if pp == nil {
		return nil, errors.New("processor handler: payment processor cannot be nil")
	}

	return &Handler{
		paymentProcessor: pp,
	}, nil
}

// HandleMessage handles incoming messages from the queue
// Implements the MessageHandler interface from infrastructure/messagebroker
func (h *Handler) HandleMessage(body []byte) error {
	ctx := context.Background()

	// Parse message
	var payment domain.Payment
	if err := payment.Parse(body); err != nil {
		return fmt.Errorf("processor handler: failed to parse payment message: %w", err)
	}

	// Validate message
	if err := payment.Validate(); err != nil {
		return fmt.Errorf("processor handler: invalid payment message: %w", err)
	}

	// Process payment
	if err := h.paymentProcessor.Process(ctx, &payment); err != nil {
		return fmt.Errorf("processor handler: failed to process payment: %w", err)
	}

	return nil
}
