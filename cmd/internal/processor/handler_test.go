package processor

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name             string
		paymentProcessor PaymentProcessor
		expectedError    string
	}{
		{
			name:             "when payment processor is provided it should create handler successfully and no error",
			paymentProcessor: new(MockPaymentProcessorService),
			expectedError:    "",
		},
		{
			name:             "when payment processor is nil it should return error",
			paymentProcessor: nil,
			expectedError:    "processor handler: payment processor cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (Payment processor already prepared in test struct)

			// Act
			result, err := NewHandler(tt.paymentProcessor)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestHandler_HandleMessage(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name              string
		messageBody       []byte
		mockProcessError  error
		shouldCallProcess bool
		expectedError     error
	}{
		{
			name: "when message is valid and processing succeeds it should return no error",
			messageBody: func() []byte {
				payment := &domain.Payment{
					ID:             "pay_123",
					IdempotencyKey: "key_456",
					UserID:         "user_789",
					Amount:         100.50,
					Currency:       domain.CurrencyUSD,
					Status:         domain.StatusReserved,
					CreatedAt:      fixedTime,
					UpdatedAt:      fixedTime,
				}
				body, _ := json.Marshal(payment)
				return body
			}(),
			mockProcessError:  nil,
			shouldCallProcess: true,
			expectedError:     nil,
		},
		{
			name:              "when message body is invalid JSON it should return parse error",
			messageBody:       []byte("invalid json"),
			shouldCallProcess: false,
			expectedError:     assert.AnError,
		},
		{
			name: "when payment validation fails it should return validation error",
			messageBody: func() []byte {
				payment := &domain.Payment{
					ID:             "",
					IdempotencyKey: "key_456",
					UserID:         "user_789",
					Amount:         100.50,
					Currency:       domain.CurrencyUSD,
					Status:         domain.StatusReserved,
					CreatedAt:      fixedTime,
					UpdatedAt:      fixedTime,
				}
				body, _ := json.Marshal(payment)
				return body
			}(),
			shouldCallProcess: false,
			expectedError:     errors.New("payment ID is required"),
		},
		{
			name: "when payment has empty user ID it should return validation error",
			messageBody: func() []byte {
				payment := &domain.Payment{
					ID:             "pay_123",
					IdempotencyKey: "key_456",
					UserID:         "",
					Amount:         100.50,
					Currency:       domain.CurrencyUSD,
					Status:         domain.StatusReserved,
					CreatedAt:      fixedTime,
					UpdatedAt:      fixedTime,
				}
				body, _ := json.Marshal(payment)
				return body
			}(),
			shouldCallProcess: false,
			expectedError:     errors.New("user ID is required"),
		},
		{
			name: "when payment has zero amount it should return validation error",
			messageBody: func() []byte {
				payment := &domain.Payment{
					ID:             "pay_123",
					IdempotencyKey: "key_456",
					UserID:         "user_789",
					Amount:         0,
					Currency:       domain.CurrencyUSD,
					Status:         domain.StatusReserved,
					CreatedAt:      fixedTime,
					UpdatedAt:      fixedTime,
				}
				body, _ := json.Marshal(payment)
				return body
			}(),
			shouldCallProcess: false,
			expectedError:     errors.New("amount must be greater than 0"),
		},
		{
			name: "when processing fails it should return processing error",
			messageBody: func() []byte {
				payment := &domain.Payment{
					ID:             "pay_123",
					IdempotencyKey: "key_456",
					UserID:         "user_789",
					Amount:         100.50,
					Currency:       domain.CurrencyUSD,
					Status:         domain.StatusReserved,
					CreatedAt:      fixedTime,
					UpdatedAt:      fixedTime,
				}
				body, _ := json.Marshal(payment)
				return body
			}(),
			mockProcessError:  errors.New("processing failed"),
			shouldCallProcess: true,
			expectedError:     errors.New("processing failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockProcessor := new(MockPaymentProcessorService)
			if tt.shouldCallProcess {
				mockProcessor.On("Process", mock.Anything, mock.AnythingOfType("*domain.Payment")).Return(tt.mockProcessError)
			}

			handler := &Handler{paymentProcessor: mockProcessor}

			// Act
			err := handler.HandleMessage(tt.messageBody)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedError != assert.AnError {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			if tt.shouldCallProcess {
				mockProcessor.AssertExpectations(t)
			}
		})
	}
}
