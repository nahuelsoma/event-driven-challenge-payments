package database

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// RowScanner is an interface for scanning a single row
type RowScanner interface {
	Scan(dest ...any) error
}

// Rows is an interface for iterating over multiple rows
type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

// DB represents a PostgreSQL database with retry capabilities
type DB struct {
	conn       *sql.DB
	maxRetries int
	baseDelay  time.Duration
}

// NewPostgresConnection creates a new PostgreSQL database connection with connection pooling
func NewPostgresConnection(url string) (*DB, error) {
	if url == "" {
		return nil, errors.New("database: url cannot be empty")
	}

	conn, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	// Connection pool settings
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &DB{
		conn:       conn,
		maxRetries: 3,
		baseDelay:  100 * time.Millisecond,
	}, nil
}

// Conn returns the underlying sql.DB connection
func (db *DB) Conn() *sql.DB {
	return db.conn
}

// QueryRowContext executes a query that returns a single row
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...any) RowScanner {
	return db.conn.QueryRowContext(ctx, query, args...)
}

// QueryContext executes a query that returns multiple rows
func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
	return db.conn.QueryContext(ctx, query, args...)
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// WithTransaction executes a function within a transaction with retry logic for transient errors
func (db *DB) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	var lastErr error

	for attempt := 0; attempt < db.maxRetries; attempt++ {
		tx, err := db.conn.BeginTx(ctx, nil)
		if err != nil {
			if isTransientError(err) {
				lastErr = err
				time.Sleep(db.exponentialBackoff(attempt))
				continue
			}
			return err
		}

		if err := fn(tx); err != nil {
			tx.Rollback()
			if isTransientError(err) {
				lastErr = err
				time.Sleep(db.exponentialBackoff(attempt))
				continue
			}
			return err
		}

		if err := tx.Commit(); err != nil {
			if isTransientError(err) {
				lastErr = err
				time.Sleep(db.exponentialBackoff(attempt))
				continue
			}
			return err
		}

		return nil
	}

	return lastErr
}

// isTransientError checks if an error is transient and can be retried
func isTransientError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, sql.ErrConnDone) ||
		strings.Contains(err.Error(), "deadlock") ||
		strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "connection reset")
}

// exponentialBackoff calculates the delay with jitter for retry attempts
func (db *DB) exponentialBackoff(attempt int) time.Duration {
	backoff := db.baseDelay * time.Duration(1<<attempt) // 100ms, 200ms, 400ms
	jitter := time.Duration(rand.Int63n(int64(backoff / 2)))
	return backoff + jitter
}
