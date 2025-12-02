package finder

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/repository/paymentstorer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewPaymentFinderService(t *testing.T) {
	tests := []struct {
		name          string
		paymentReader PaymentReader
		expectedError string
	}{
		{
			name:          "when payment reader is provided it should create service successfully and no error",
			paymentReader: new(paymentstorer.MockPaymentRepository),
			expectedError: "",
		},
		{
			name:          "when payment reader is nil it should return error with message 'payment finder: reader cannot be nil'",
			paymentReader: nil,
			expectedError: "payment finder: reader cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (Payment reader already prepared in test struct)

			// Act
			result, err := NewPaymentFinderService(tt.paymentReader)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, result.paymentReader)
			}
		})
	}
}

func TestPaymentFinderService_Find(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name            string
		filter          *PaymentFilter
		mockPayment     *domain.Payment
		mockError       error
		expectedPayment *domain.Payment
		expectedError   error
	}{
		{
			name: "when payment exists it should return payment and no error",
			filter: &PaymentFilter{
				PaymentID: "pay_123",
			},
			mockPayment: &domain.Payment{
				ID:             "pay_123",
				IdempotencyKey: "key_456",
				UserID:         "user_789",
				Amount:         100.50,
				Currency:       domain.CurrencyUSD,
				Status:         domain.StatusReserved,
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			mockError: nil,
			expectedPayment: &domain.Payment{
				ID:             "pay_123",
				IdempotencyKey: "key_456",
				UserID:         "user_789",
				Amount:         100.50,
				Currency:       domain.CurrencyUSD,
				Status:         domain.StatusReserved,
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			expectedError: nil,
		},
		{
			name: "when payment reader returns error it should return wrapped error",
			filter: &PaymentFilter{
				PaymentID: "pay_nonexistent",
			},
			mockPayment:     nil,
			mockError:       domain.ErrPaymentNotFound,
			expectedPayment: nil,
			expectedError:   errors.New("payment finder: get payment: payment not found"),
		},
		{
			name: "when database error occurs it should return wrapped error",
			filter: &PaymentFilter{
				PaymentID: "pay_123",
			},
			mockPayment:     nil,
			mockError:       errors.New("connection refused"),
			expectedPayment: nil,
			expectedError:   errors.New("payment finder: get payment: connection refused"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockReader := new(paymentstorer.MockPaymentRepository)
			mockReader.On("GetByID", mock.Anything, tt.filter.PaymentID).Return(tt.mockPayment, tt.mockError)

			service := &PaymentFinderService{paymentReader: mockReader}

			// Act
			result, err := service.Find(context.Background(), tt.filter)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedPayment.ID, result.ID)
				assert.Equal(t, tt.expectedPayment.IdempotencyKey, result.IdempotencyKey)
				assert.Equal(t, tt.expectedPayment.UserID, result.UserID)
				assert.Equal(t, tt.expectedPayment.Amount, result.Amount)
				assert.Equal(t, tt.expectedPayment.Currency, result.Currency)
				assert.Equal(t, tt.expectedPayment.Status, result.Status)
			}

			mockReader.AssertExpectations(t)
		})
	}
}

func TestPaymentFinderService_FindEvents(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name           string
		paymentID      string
		mockEvents     []*domain.Event
		mockError      error
		expectedEvents []*domain.Event
		expectedError  error
	}{
		{
			name:      "when events exist it should return events and no error",
			paymentID: "pay_123",
			mockEvents: []*domain.Event{
				{
					ID:        "evt_1",
					PaymentID: "pay_123",
					Sequence:  1,
					EventType: "created",
					Payload:   json.RawMessage(`{"status":"pending"}`),
					CreatedAt: fixedTime,
				},
				{
					ID:        "evt_2",
					PaymentID: "pay_123",
					Sequence:  2,
					EventType: "reserved",
					Payload:   json.RawMessage(`{"status":"reserved"}`),
					CreatedAt: fixedTime,
				},
			},
			mockError: nil,
			expectedEvents: []*domain.Event{
				{
					ID:        "evt_1",
					PaymentID: "pay_123",
					Sequence:  1,
					EventType: "created",
					Payload:   json.RawMessage(`{"status":"pending"}`),
					CreatedAt: fixedTime,
				},
				{
					ID:        "evt_2",
					PaymentID: "pay_123",
					Sequence:  2,
					EventType: "reserved",
					Payload:   json.RawMessage(`{"status":"reserved"}`),
					CreatedAt: fixedTime,
				},
			},
			expectedError: nil,
		},
		{
			name:           "when no events exist it should return empty slice and no error",
			paymentID:      "pay_123",
			mockEvents:     []*domain.Event{},
			mockError:      nil,
			expectedEvents: []*domain.Event{},
			expectedError:  nil,
		},
		{
			name:           "when payment reader returns error it should return wrapped error",
			paymentID:      "pay_123",
			mockEvents:     nil,
			mockError:      errors.New("database error"),
			expectedEvents: nil,
			expectedError:  errors.New("payment finder: get events: database error"),
		},
		{
			name:           "when connection error occurs it should return wrapped error",
			paymentID:      "pay_123",
			mockEvents:     nil,
			mockError:      errors.New("connection refused"),
			expectedEvents: nil,
			expectedError:  errors.New("payment finder: get events: connection refused"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockReader := new(paymentstorer.MockPaymentRepository)
			mockReader.On("GetEventsByPaymentID", mock.Anything, tt.paymentID).Return(tt.mockEvents, tt.mockError)

			service := &PaymentFinderService{paymentReader: mockReader}

			// Act
			result, err := service.FindEvents(context.Background(), tt.paymentID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, len(tt.expectedEvents))
				for i, expectedEvent := range tt.expectedEvents {
					assert.Equal(t, expectedEvent.ID, result[i].ID)
					assert.Equal(t, expectedEvent.PaymentID, result[i].PaymentID)
					assert.Equal(t, expectedEvent.Sequence, result[i].Sequence)
					assert.Equal(t, expectedEvent.EventType, result[i].EventType)
				}
			}

			mockReader.AssertExpectations(t)
		})
	}
}
