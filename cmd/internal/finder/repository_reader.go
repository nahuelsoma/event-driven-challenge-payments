package finder

import (
	"context"
	"errors"
	"log/slog"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
)

type PaymentReaderRepository struct {
	databaseConn interface{}
}

func NewPaymentReaderRepository(databaseConn interface{}) (*PaymentReaderRepository, error) {
	if databaseConn == nil {
		return nil, errors.New("payment reader: database cannot be nil")
	}

	return &PaymentReaderRepository{databaseConn: databaseConn}, nil
}

func (r *PaymentReaderRepository) GetByID(ctx context.Context, paymentID string) (*domain.Payment, error) {
	slog.DebugContext(ctx, "[DEBUG] PaymentReaderRepository.GetByID called", "payment_id", paymentID)
	// TODO: Implement the logic to get the payment by ID
	// TODO: Return nil, ErrPaymentNotFound if the payment is not found
	// TODO: Return the error if there is an error different from sql.ErrNoRows
	return nil, nil
}
