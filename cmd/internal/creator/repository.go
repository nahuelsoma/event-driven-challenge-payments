package creator

import (
	"context"
	"errors"
	"fmt"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
)

// MessageBroker interface for the infrastructure publisher
type MessageBroker interface {
	Publish(body []byte) error
}

type PaymentPublisherRepository struct {
	messageBroker MessageBroker
}

func NewPaymentPublisherRepository(mb MessageBroker) (*PaymentPublisherRepository, error) {
	if mb == nil {
		return nil, errors.New("payment publisher: message broker cannot be nil")
	}

	return &PaymentPublisherRepository{
		messageBroker: mb,
	}, nil
}

func (r *PaymentPublisherRepository) Publish(ctx context.Context, payment *domain.Payment) error {
	data, err := payment.Marshal()
	if err != nil {
		return fmt.Errorf("publisher: failed to marshal payment: %w", err)
	}

	if err := r.messageBroker.Publish(data); err != nil {
		return fmt.Errorf("publisher: failed to publish payment: %w", err)
	}

	return nil
}
