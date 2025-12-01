package finder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
		return nil, fmt.Errorf("payment reader: get by id: %w", err)
	}

	return &payment, nil
}

func (r *PaymentReaderRepository) GetEventsByPaymentID(ctx context.Context, paymentID string) ([]*domain.Event, error) {
	query := `
		SELECT id, payment_id, sequence, event_type, payload, created_at
		FROM payment_events
		WHERE payment_id = $1
		ORDER BY sequence ASC
	`

	rows, err := r.db.Conn().QueryContext(ctx, query, paymentID)
	if err != nil {
		return nil, fmt.Errorf("payment reader: get events by payment id: %w", err)
	}
	defer rows.Close()

	events := []*domain.Event{}
	for rows.Next() {
		var event domain.Event
		err := rows.Scan(
			&event.ID,
			&event.PaymentID,
			&event.Sequence,
			&event.EventType,
			&event.Payload,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("payment reader: scan event: %w", err)
		}
		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("payment reader: iterate events: %w", err)
	}

	return events, nil
}
