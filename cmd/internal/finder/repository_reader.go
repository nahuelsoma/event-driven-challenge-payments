package finder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
)

// paymentReaderDB defines the database operations required by PaymentReaderRepository
type paymentReaderDB interface {
	Conn() *sql.DB
}

type PaymentReaderRepository struct {
	db paymentReaderDB
}

func NewPaymentReaderRepository(db paymentReaderDB) (*PaymentReaderRepository, error) {
	if db == nil {
		return nil, errors.New("payment reader: database cannot be nil")
	}

	return &PaymentReaderRepository{db: db}, nil
}

func (r *PaymentReaderRepository) GetByID(ctx context.Context, paymentID string) (*domain.Payment, error) {
	slog.DebugContext(ctx, "[DEBUG] PaymentReaderRepository.GetByID called", "payment_id", paymentID)

	query := `
		SELECT id, idempotency_key, user_id, amount, currency, status, created_at, updated_at
		FROM payments
		WHERE id = $1
	`

	var payment domain.Payment
	err := r.db.Conn().QueryRowContext(ctx, query, paymentID).Scan(
		&payment.ID,
		&payment.IdempotencyKey,
		&payment.UserID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPaymentNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get payment by id: %w", err)
	}

	return &payment, nil
}
