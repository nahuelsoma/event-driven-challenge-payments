package paymentstorer

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/database"
)

// PaymentDB defines the database operations required by PaymentRepository
type PaymentDB interface {
	QueryRowContext(ctx context.Context, query string, args ...any) database.RowScanner
	QueryContext(ctx context.Context, query string, args ...any) (database.Rows, error)
	WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error
}

// PaymentRepository handles all database operations for payments
type PaymentRepository struct {
	db PaymentDB
}

// NewStorer creates a new PaymentRepository
func NewStorer(db PaymentDB) (*PaymentRepository, error) {
	if db == nil {
		return nil, errors.New("payment repository: database cannot be nil")
	}

	return &PaymentRepository{db: db}, nil
}

// GetByID retrieves a payment by ID
func (r *PaymentRepository) GetByID(ctx context.Context, paymentID string) (*domain.Payment, error) {
	query := `
		SELECT id, idempotency_key, user_id, amount, currency, status, created_at, updated_at
		FROM payments
		WHERE id = $1
	`

	var payment domain.Payment
	err := r.db.QueryRowContext(ctx, query, paymentID).Scan(
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
		return nil, domain.ErrPaymentNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("payment repository: get by id: %w", err)
	}

	return &payment, nil
}

// GetByIDempotencyKey retrieves a payment by idempotency key
func (r *PaymentRepository) GetByIDempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Payment, error) {
	query := `
		SELECT id, idempotency_key, user_id, amount, currency, status, created_at, updated_at
		FROM payments
		WHERE idempotency_key = $1
	`

	var payment domain.Payment
	err := r.db.QueryRowContext(ctx, query, idempotencyKey).Scan(
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
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("payment repository: get by idempotency key: %w", err)
	}

	return &payment, nil
}

// Save saves a new payment with its initial event
func (r *PaymentRepository) Save(ctx context.Context, payment *domain.Payment) error {
	payload, err := json.Marshal(map[string]interface{}{
		"payment_id":      payment.ID,
		"idempotency_key": payment.IdempotencyKey,
		"user_id":         payment.UserID,
		"amount":          payment.Amount,
		"currency":        payment.Currency,
		"status":          payment.Status,
	})
	if err != nil {
		return fmt.Errorf("payment repository: marshal payload: %w", err)
	}

	err = r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		// Insert into Event Store (source of truth)
		eventQuery := `
			INSERT INTO payment_events (id, payment_id, sequence, event_type, payload, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		_, err := tx.ExecContext(ctx, eventQuery,
			uuid.New().String(),
			payment.ID,
			1,
			"created",
			payload,
			time.Now(),
		)
		if err != nil {
			return fmt.Errorf("insert event: %w", err)
		}

		// Insert into Read Model (for queries)
		paymentQuery := `
			INSERT INTO payments (id, idempotency_key, user_id, amount, currency, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`
		_, err = tx.ExecContext(ctx, paymentQuery,
			payment.ID,
			payment.IdempotencyKey,
			payment.UserID,
			payment.Amount,
			payment.Currency,
			payment.Status,
			payment.CreatedAt,
			payment.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("insert payment: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("payment repository: save: %w", err)
	}

	return nil
}

// UpdateStatus updates the payment status with optional gateway reference
func (r *PaymentRepository) UpdateStatus(ctx context.Context, paymentID string, status domain.Status, gatewayRef string) error {
	payload, err := json.Marshal(map[string]interface{}{
		"payment_id":  paymentID,
		"status":      status,
		"gateway_ref": gatewayRef,
	})
	if err != nil {
		return fmt.Errorf("payment repository: marshal payload: %w", err)
	}

	now := time.Now()

	err = r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		// Get next sequence number for this payment
		var nextSequence int
		sequenceQuery := `
			SELECT COALESCE(MAX(sequence), 0) + 1
			FROM payment_events
			WHERE payment_id = $1
		`
		err := tx.QueryRowContext(ctx, sequenceQuery, paymentID).Scan(&nextSequence)
		if err != nil {
			return fmt.Errorf("get sequence: %w", err)
		}

		// Insert into Event Store
		eventQuery := `
			INSERT INTO payment_events (id, payment_id, sequence, event_type, payload, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		_, err = tx.ExecContext(ctx, eventQuery,
			uuid.New().String(),
			paymentID,
			nextSequence,
			string(status),
			payload,
			now,
		)
		if err != nil {
			return fmt.Errorf("insert event: %w", err)
		}

		// Update Read Model
		updateQuery := `
			UPDATE payments
			SET status = $1, gateway_ref = $2, updated_at = $3
			WHERE id = $4
		`
		result, err := tx.ExecContext(ctx, updateQuery, status, gatewayRef, now, paymentID)
		if err != nil {
			return fmt.Errorf("update status: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("rows affected: %w", err)
		}

		if rowsAffected == 0 {
			return fmt.Errorf("payment not found: %s", paymentID)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("payment repository: update status: %w", err)
	}

	return nil
}

// GetEventsByPaymentID retrieves all events for a payment
func (r *PaymentRepository) GetEventsByPaymentID(ctx context.Context, paymentID string) ([]*domain.Event, error) {
	query := `
		SELECT id, payment_id, sequence, event_type, payload, created_at
		FROM payment_events
		WHERE payment_id = $1
		ORDER BY sequence ASC
	`

	rows, err := r.db.QueryContext(ctx, query, paymentID)
	if err != nil {
		return nil, fmt.Errorf("payment repository: get events by payment id: %w", err)
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
			return nil, fmt.Errorf("payment repository: scan event: %w", err)
		}
		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("payment repository: iterate events: %w", err)
	}

	return events, nil
}
