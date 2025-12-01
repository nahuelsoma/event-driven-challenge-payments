package creator

import (
	"context"
	"errors"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
)

type PaymentStorerRepository struct {
	databaseConn interface{}
}

func NewPaymentStorerRepository(databaseConn interface{}) (*PaymentStorerRepository, error) {
	if databaseConn == nil {
		return nil, errors.New("payment storer: database cannot be nil")
	}

	return &PaymentStorerRepository{databaseConn: databaseConn}, nil
}

func (r *PaymentStorerRepository) GetByIDempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Payment, error) {
	// TODO: Implement the logic to get the payment by idempotency key
	// TODO: Return nil if the payment is not found, otherwise return the payment
	// TODO: Return the error if there is an error different from sql.ErrNoRows
	return nil, nil
}

func (r *PaymentStorerRepository) Save(ctx context.Context, payment *domain.Payment) error {
	// TODO: Implement the logic to save the payment in Event Store and Read Model
	return nil
}

func (r *PaymentStorerRepository) UpdateStatus(ctx context.Context, paymentID string, status domain.Status) error {
	// TODO: Implement the logic to update the status of the payment in Event Store and Read Model
	return nil
}
