package creator

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/messagebroker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewPaymentPublisherRepository(t *testing.T) {
	tests := []struct {
		name           string
		messageBroker  MessageBroker
		expectedError  string
		expectedResult bool
	}{
		{
			name:           "when message broker is provided it should create repository successfully and no error",
			messageBroker:  new(messagebroker.MockPublisher),
			expectedError:  "",
			expectedResult: true,
		},
		{
			name:           "when message broker is nil it should return error with message 'payment publisher: message broker cannot be nil'",
			messageBroker:  nil,
			expectedError:  "payment publisher: message broker cannot be nil",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (Message broker already prepared in test struct)

			// Act
			result, err := NewPaymentPublisherRepository(tt.messageBroker)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, result.messageBroker)
			}
		})
	}
}

func TestPaymentPublisherRepository_Publish(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name          string
		payment       *domain.Payment
		setupMock     func(mockBroker *messagebroker.MockPublisher)
		expectedError error
	}{
		{
			name: "when payment is valid and broker publishes successfully it should return no error",
			payment: &domain.Payment{
				ID:             "pay_123",
				IdempotencyKey: "key_456",
				UserID:         "user_789",
				Amount:         100.50,
				Currency:       domain.CurrencyUSD,
				Status:         domain.StatusPending,
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			setupMock: func(mockBroker *messagebroker.MockPublisher) {
				mockBroker.On("Publish", mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "when broker fails to publish it should return wrapped error",
			payment: &domain.Payment{
				ID:             "pay_123",
				IdempotencyKey: "key_456",
				UserID:         "user_789",
				Amount:         100.50,
				Currency:       domain.CurrencyUSD,
				Status:         domain.StatusPending,
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			setupMock: func(mockBroker *messagebroker.MockPublisher) {
				mockBroker.On("Publish", mock.Anything).Return(errors.New("connection refused"))
			},
			expectedError: errors.New("publisher: failed to publish payment: connection refused"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockBroker := new(messagebroker.MockPublisher)
			tt.setupMock(mockBroker)

			repo := &PaymentPublisherRepository{messageBroker: mockBroker}

			// Act
			err := repo.Publish(context.Background(), tt.payment)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockBroker.AssertExpectations(t)
		})
	}
}
