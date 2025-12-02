package creator

import (
	"testing"
	"time"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

func TestPaymentRequest_Validate(t *testing.T) {
	tests := []struct {
		name          string
		request       *PaymentRequest
		expectedError string
	}{
		{
			name: "when request has all valid fields it should pass validation and no error",
			request: &PaymentRequest{
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
			},
			expectedError: "",
		},
		{
			name: "when user ID is empty it should return error with message 'user ID is required'",
			request: &PaymentRequest{
				UserID:   "",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
			},
			expectedError: "user ID is required",
		},
		{
			name: "when amount is zero it should return error with message 'amount must be greater than 0'",
			request: &PaymentRequest{
				UserID:   "user_123",
				Amount:   0,
				Currency: domain.CurrencyUSD,
			},
			expectedError: "amount must be greater than 0",
		},
		{
			name: "when amount is negative it should return error with message 'amount must be greater than 0'",
			request: &PaymentRequest{
				UserID:   "user_123",
				Amount:   -50.00,
				Currency: domain.CurrencyUSD,
			},
			expectedError: "amount must be greater than 0",
		},
		{
			name: "when currency is invalid it should return error with message 'invalid currency'",
			request: &PaymentRequest{
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.Currency("INVALID"),
			},
			expectedError: "invalid currency",
		},
		{
			name: "when currency is empty it should return error with message 'invalid currency'",
			request: &PaymentRequest{
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.Currency(""),
			},
			expectedError: "invalid currency",
		},
		{
			name: "when request has minimum valid amount it should pass validation and no error",
			request: &PaymentRequest{
				UserID:   "user_123",
				Amount:   0.01,
				Currency: domain.CurrencyEUR,
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (Request already prepared in test struct)

			// Act
			err := tt.request.Validate()

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

func TestNewPayment(t *testing.T) {
	tests := []struct {
		name           string
		idempotencyKey string
		userID         string
		amount         float64
		currency       domain.Currency
	}{
		{
			name:           "when creating payment with valid data it should return payment with correct fields",
			idempotencyKey: "key_123",
			userID:         "user_456",
			amount:         100.50,
			currency:       domain.CurrencyUSD,
		},
		{
			name:           "when creating payment with EUR currency it should return payment with EUR currency",
			idempotencyKey: "key_789",
			userID:         "user_012",
			amount:         250.00,
			currency:       domain.CurrencyEUR,
		},
		{
			name:           "when creating payment with minimum amount it should return payment with correct amount",
			idempotencyKey: "key_min",
			userID:         "user_min",
			amount:         0.01,
			currency:       domain.CurrencyUSD,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			beforeCreate := time.Now()

			// Act
			result := NewPayment(tt.idempotencyKey, tt.userID, tt.amount, tt.currency)

			// Assert
			afterCreate := time.Now()

			assert.NotNil(t, result)
			assert.NotEmpty(t, result.ID)
			assert.Equal(t, tt.idempotencyKey, result.IdempotencyKey)
			assert.Equal(t, tt.userID, result.UserID)
			assert.Equal(t, tt.amount, result.Amount)
			assert.Equal(t, tt.currency, result.Currency)
			assert.Equal(t, domain.StatusPending, result.Status)

			// Verify timestamps are set correctly
			assert.True(t, result.CreatedAt.After(beforeCreate) || result.CreatedAt.Equal(beforeCreate))
			assert.True(t, result.CreatedAt.Before(afterCreate) || result.CreatedAt.Equal(afterCreate))
			assert.True(t, result.UpdatedAt.After(beforeCreate) || result.UpdatedAt.Equal(beforeCreate))
			assert.True(t, result.UpdatedAt.Before(afterCreate) || result.UpdatedAt.Equal(afterCreate))
		})
	}
}

