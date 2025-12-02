package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPayment_Validate(t *testing.T) {
	tests := []struct {
		name          string
		payment       *Payment
		expectedError string
	}{
		{
			name: "when payment has all required fields it should pass validation and no error",
			payment: &Payment{
				ID:     "pay_123",
				UserID: "user_456",
				Amount: 100.50,
			},
			expectedError: "",
		},
		{
			name: "when payment ID is empty it should return error with message 'payment ID is required'",
			payment: &Payment{
				ID:     "",
				UserID: "user_456",
				Amount: 100.50,
			},
			expectedError: "payment ID is required",
		},
		{
			name: "when user ID is empty it should return error with message 'user ID is required'",
			payment: &Payment{
				ID:     "pay_123",
				UserID: "",
				Amount: 100.50,
			},
			expectedError: "user ID is required",
		},
		{
			name: "when amount is zero it should return error with message 'amount must be greater than 0'",
			payment: &Payment{
				ID:     "pay_123",
				UserID: "user_456",
				Amount: 0,
			},
			expectedError: "amount must be greater than 0",
		},
		{
			name: "when amount is negative it should return error with message 'amount must be greater than 0'",
			payment: &Payment{
				ID:     "pay_123",
				UserID: "user_456",
				Amount: -50.00,
			},
			expectedError: "amount must be greater than 0",
		},
		{
			name: "when payment has minimum valid amount it should pass validation and no error",
			payment: &Payment{
				ID:     "pay_123",
				UserID: "user_456",
				Amount: 0.01,
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (Payment already prepared in test struct)

			// Act
			err := tt.payment.Validate()

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPayment_UpdateStatus(t *testing.T) {
	tests := []struct {
		name           string
		initialStatus  Status
		newStatus      Status
		expectedStatus Status
	}{
		{
			name:           "when updating status from pending to completed it should update status and timestamp",
			initialStatus:  StatusPending,
			newStatus:      StatusCompleted,
			expectedStatus: StatusCompleted,
		},
		{
			name:           "when updating status from pending to failed it should update status and timestamp",
			initialStatus:  StatusPending,
			newStatus:      StatusFailed,
			expectedStatus: StatusFailed,
		},
		{
			name:           "when updating status from pending to reserved it should update status and timestamp",
			initialStatus:  StatusPending,
			newStatus:      StatusReserved,
			expectedStatus: StatusReserved,
		},
		{
			name:           "when updating status to the same value it should update timestamp",
			initialStatus:  StatusCompleted,
			newStatus:      StatusCompleted,
			expectedStatus: StatusCompleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			initialTime := time.Now().Add(-1 * time.Hour)
			payment := &Payment{
				ID:        "pay_123",
				UserID:    "user_456",
				Amount:    100.50,
				Status:    tt.initialStatus,
				UpdatedAt: initialTime,
			}

			// Act
			beforeUpdate := time.Now()
			payment.UpdateStatus(tt.newStatus)
			afterUpdate := time.Now()

			// Assert
			assert.Equal(t, tt.expectedStatus, payment.Status)
			assert.True(t, payment.UpdatedAt.After(initialTime) || payment.UpdatedAt.Equal(initialTime))
			assert.True(t, payment.UpdatedAt.After(beforeUpdate) || payment.UpdatedAt.Equal(beforeUpdate))
			assert.True(t, payment.UpdatedAt.Before(afterUpdate) || payment.UpdatedAt.Equal(afterUpdate))
		})
	}
}

func TestPayment_Parse(t *testing.T) {
	tests := []struct {
		name            string
		body            []byte
		expectedPayment *Payment
		expectedError   bool
	}{
		{
			name: "when body contains valid JSON it should parse payment successfully and no error",
			body: []byte(`{"id":"pay_123","idempotency_key":"key_456","user_id":"user_789","amount":100.50,"currency":"USD","status":"pending"}`),
			expectedPayment: &Payment{
				ID:             "pay_123",
				IdempotencyKey: "key_456",
				UserID:         "user_789",
				Amount:         100.50,
				Currency:       CurrencyUSD,
				Status:         StatusPending,
			},
			expectedError: false,
		},
		{
			name: "when body contains partial JSON it should parse available fields and no error",
			body: []byte(`{"id":"pay_123","user_id":"user_789","amount":50.00}`),
			expectedPayment: &Payment{
				ID:     "pay_123",
				UserID: "user_789",
				Amount: 50.00,
			},
			expectedError: false,
		},
		{
			name:            "when body contains invalid JSON it should return error",
			body:            []byte(`{invalid json}`),
			expectedPayment: nil,
			expectedError:   true,
		},
		{
			name: "when body is empty JSON object it should parse empty payment and no error",
			body: []byte(`{}`),
			expectedPayment: &Payment{
				ID:     "",
				UserID: "",
				Amount: 0,
			},
			expectedError: false,
		},
		{
			name:            "when body is empty it should return error",
			body:            []byte(``),
			expectedPayment: nil,
			expectedError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			payment := &Payment{}

			// Act
			err := payment.Parse(tt.body)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPayment.ID, payment.ID)
				assert.Equal(t, tt.expectedPayment.IdempotencyKey, payment.IdempotencyKey)
				assert.Equal(t, tt.expectedPayment.UserID, payment.UserID)
				assert.Equal(t, tt.expectedPayment.Amount, payment.Amount)
				assert.Equal(t, tt.expectedPayment.Currency, payment.Currency)
				assert.Equal(t, tt.expectedPayment.Status, payment.Status)
			}
		})
	}
}

func TestPayment_Marshal(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name          string
		payment       *Payment
		expectedError bool
	}{
		{
			name: "when payment has all fields it should marshal successfully and no error",
			payment: &Payment{
				ID:             "pay_123",
				IdempotencyKey: "key_456",
				UserID:         "user_789",
				Amount:         100.50,
				Currency:       CurrencyUSD,
				Status:         StatusPending,
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			expectedError: false,
		},
		{
			name: "when payment has minimal fields it should marshal successfully and no error",
			payment: &Payment{
				ID:     "pay_123",
				UserID: "user_789",
				Amount: 50.00,
			},
			expectedError: false,
		},
		{
			name:          "when payment is empty it should marshal successfully and no error",
			payment:       &Payment{},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (Payment already prepared in test struct)

			// Act
			result, err := tt.payment.Marshal()

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)

				// Verify round-trip: unmarshal and compare
				parsedPayment := &Payment{}
				parseErr := parsedPayment.Parse(result)
				assert.NoError(t, parseErr)
				assert.Equal(t, tt.payment.ID, parsedPayment.ID)
				assert.Equal(t, tt.payment.IdempotencyKey, parsedPayment.IdempotencyKey)
				assert.Equal(t, tt.payment.UserID, parsedPayment.UserID)
				assert.Equal(t, tt.payment.Amount, parsedPayment.Amount)
				assert.Equal(t, tt.payment.Currency, parsedPayment.Currency)
				assert.Equal(t, tt.payment.Status, parsedPayment.Status)
			}
		})
	}
}

