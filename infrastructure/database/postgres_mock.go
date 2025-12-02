package database

import (
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of the database for testing
type MockDB struct {
	mock.Mock
}

// QueryRowContext returns a row scanner for a single row query
func (m *MockDB) QueryRowContext(ctx context.Context, query string, args ...any) RowScanner {
	callArgs := m.Called(ctx, query, args)
	return callArgs.Get(0).(RowScanner)
}

// QueryContext returns rows for a multi-row query
func (m *MockDB) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
	callArgs := m.Called(ctx, query, args)
	if callArgs.Get(0) == nil {
		return nil, callArgs.Error(1)
	}
	return callArgs.Get(0).(Rows), callArgs.Error(1)
}

// WithTransaction executes a function within a transaction
func (m *MockDB) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

// MockRowScanner is a mock implementation of RowScanner
type MockRowScanner struct {
	mock.Mock
}

// Scan scans the row into the provided destinations
func (m *MockRowScanner) Scan(dest ...any) error {
	args := m.Called(dest)
	return args.Error(0)
}

// MockRows is a mock implementation of Rows
type MockRows struct {
	mock.Mock
}

// Next advances to the next row
func (m *MockRows) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

// Scan scans the current row into the provided destinations
func (m *MockRows) Scan(dest ...any) error {
	args := m.Called(dest)
	return args.Error(0)
}

// Close closes the rows
func (m *MockRows) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Err returns any error that occurred during iteration
func (m *MockRows) Err() error {
	args := m.Called()
	return args.Error(0)
}

