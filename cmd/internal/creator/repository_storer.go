package creator

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

// paymentStorerDB defines the database operations required by PaymentStorerRepository
type paymentStorerDB interface {
	Conn() *sql.DB
	WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error
}

type PaymentStorerRepository struct {
	db paymentStorerDB
}

func NewPaymentStorerRepository(db paymentStorerDB) (*PaymentStorerRepository, error) {
	if db == nil {
		return nil, errors.New("payment storer: database cannot be nil")
	}

	return &PaymentStorerRepository{db: db}, nil
}

func (r *PaymentStorerRepository) GetByIDempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Payment, error) {
	query := `
		SELECT id, idempotency_key, user_id, amount, currency, status, created_at, updated_at
		FROM payments
		WHERE idempotency_key = $1
	`

	var payment domain.Payment
	err := r.db.Conn().QueryRowContext(ctx, query, idempotencyKey).Scan(
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
		return nil, fmt.Errorf("failed to get payment by idempotency key: %w", err)
	}

	return &payment, nil
}

func (r *PaymentStorerRepository) Save(ctx context.Context, payment *domain.Payment) error {
	// Build event payload
	payload, err := json.Marshal(map[string]interface{}{
		"payment_id":      payment.ID,
		"idempotency_key": payment.IdempotencyKey,
		"user_id":         payment.UserID,
		"amount":          payment.Amount,
		"currency":        payment.Currency,
		"status":          payment.Status,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
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
			return fmt.Errorf("failed to insert payment event: %w", err)
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
			return fmt.Errorf("failed to insert payment: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to save payment: %w", err)
	}

	return nil
}

func (r *PaymentStorerRepository) UpdateStatus(ctx context.Context, paymentID string, status domain.Status) error {
	// Build event payload
	payload, err := json.Marshal(map[string]interface{}{
		"payment_id": paymentID,
		"status":     status,
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
			SET status = $1, updated_at = $2
			WHERE id = $3
		`
		result, err := tx.ExecContext(ctx, updateQuery, status, now, paymentID)
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
