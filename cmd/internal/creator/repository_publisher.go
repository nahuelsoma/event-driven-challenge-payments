package creator

import (
	"context"
	"errors"

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
	// TODO: Implement the logic to publish the payment
	return nil
}
