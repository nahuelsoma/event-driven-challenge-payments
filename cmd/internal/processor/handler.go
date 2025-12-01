package processor

import (
	"context"
	"errors"
	"log/slog"

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

	slog.InfoContext(ctx, "Processing payment message", "body", string(body))

	// Parse message
	var payment domain.Payment
	if err := payment.Parse(body); err != nil {
		slog.ErrorContext(ctx, "Failed to parse payment message", "error", err)
		return err
	}

	// Validate message
	if err := payment.Validate(); err != nil {
		slog.WarnContext(ctx, "Invalid payment message", "error", err, "payment_id", payment.ID)
		return err
	}

	slog.InfoContext(ctx, "Processing payment", "payment_id", payment.ID, "user_id", payment.UserID, "amount", payment.Amount)

	// Process payment
	if err := h.paymentProcessor.Process(ctx, &payment); err != nil {
		slog.ErrorContext(ctx, "Failed to process payment", "error", err, "payment_id", payment.ID)
		return err
	}

	slog.InfoContext(ctx, "Payment processed successfully", "payment_id", payment.ID)
	return nil
}
