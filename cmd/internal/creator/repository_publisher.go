package creator

import (
	"context"
	"errors"
	"log/slog"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
)

type PaymentPublisherRepository struct {
	connection interface{}
}

func NewPaymentPublisherRepository(conn interface{}) (*PaymentPublisherRepository, error) {
	if conn == nil {
		return nil, errors.New("payment publisher: connection cannot be nil")
	}

	return &PaymentPublisherRepository{
		connection: conn,
	}, nil
}

func (r *PaymentPublisherRepository) Publish(ctx context.Context, payment *domain.Payment) error {
	slog.DebugContext(ctx, "[DEBUG] PaymentPublisherRepository.Publish called", "payment_id", payment.ID, "status", payment.Status)
	// TODO: Implement the logic to publish the payment
	return nil
}
