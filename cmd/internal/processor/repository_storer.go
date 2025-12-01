package processor

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
)

// ErrPaymentNotFound is returned when a payment is not found
var ErrPaymentNotFound = errors.New("payment not found")

// paymentStorerDB defines the database operations required by PaymentStorerRepository
type paymentStorerDB interface {
	Conn() *sql.DB
	WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error
}

// PaymentStorerRepository handles all database operations for the processor
type PaymentStorerRepository struct {
	db paymentStorerDB
}

// NewPaymentStorerRepository creates a new PaymentStorerRepository
func NewPaymentStorerRepository(db paymentStorerDB) (*PaymentStorerRepository, error) {
	if db == nil {
		return nil, errors.New("payment storer: database cannot be nil")
	}

	return &PaymentStorerRepository{db: db}, nil
}

// GetByID retrieves a payment by ID
func (r *PaymentStorerRepository) GetByID(ctx context.Context, paymentID string) (*domain.Payment, error) {
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

// UpdateStatus updates the payment status with optional gateway reference
func (r *PaymentStorerRepository) UpdateStatus(ctx context.Context, paymentID string, status domain.Status, gatewayRef string) error {
	// Build event payload
	payload, err := json.Marshal(map[string]interface{}{
		"payment_id":  paymentID,
		"status":      status,
		"gateway_ref": gatewayRef,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
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
			return fmt.Errorf("failed to get next sequence: %w", err)
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
			return fmt.Errorf("failed to insert payment event: %w", err)
		}

		// Update Read Model
		updateQuery := `
			UPDATE payments
			SET status = $1, gateway_ref = $2, updated_at = $3
			WHERE id = $4
		`
		result, err := tx.ExecContext(ctx, updateQuery, status, gatewayRef, now, paymentID)
		if err != nil {
			return fmt.Errorf("failed to update payment status: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}

		if rowsAffected == 0 {
			return fmt.Errorf("payment not found: %s", paymentID)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}
