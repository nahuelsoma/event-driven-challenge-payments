package paymentstorer

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewStorer(t *testing.T) {
	tests := []struct {
		name           string
		db             PaymentDB
		expectedError  string
		expectedResult bool
	}{
		{
			name:           "when database is provided it should create repository successfully and no error",
			db:             new(database.MockDB),
			expectedError:  "",
			expectedResult: true,
		},
		{
			name:           "when database is nil it should return error with message 'payment repository: database cannot be nil'",
			db:             nil,
			expectedError:  "payment repository: database cannot be nil",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (Database already prepared in test struct)

			// Act
			result, err := NewStorer(tt.db)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, result.db)
			}
		})
	}
}

func TestPaymentRepository_GetByID(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name            string
		paymentID       string
		mockPayment     *domain.Payment
		mockScanError   error
		expectedPayment *domain.Payment
		expectedError   error
	}{
		{
			name:      "when payment exists it should return payment and no error",
			paymentID: "pay_123",
			mockPayment: &domain.Payment{
				ID:             "pay_123",
				IdempotencyKey: "key_456",
				UserID:         "user_789",
				Amount:         100.50,
				Currency:       domain.CurrencyUSD,
				Status:         domain.StatusPending,
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			mockScanError: nil,
			expectedPayment: &domain.Payment{
				ID:             "pay_123",
				IdempotencyKey: "key_456",
				UserID:         "user_789",
				Amount:         100.50,
				Currency:       domain.CurrencyUSD,
				Status:         domain.StatusPending,
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			expectedError: nil,
		},
		{
			name:          "when payment does not exist it should return ErrPaymentNotFound",
			paymentID:     "pay_nonexistent",
			mockPayment:   nil,
			mockScanError: sql.ErrNoRows,
			expectedPayment: nil,
			expectedError:   domain.ErrPaymentNotFound,
		},
		{
			name:          "when database error occurs it should return wrapped error",
			paymentID:     "pay_123",
			mockPayment:   nil,
			mockScanError: errors.New("connection refused"),
			expectedPayment: nil,
			expectedError:   errors.New("payment repository: get by id: connection refused"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDB := new(database.MockDB)
			mockScanner := new(database.MockRowScanner)

			if tt.mockPayment != nil {
				mockScanner.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
					dest := args.Get(0).([]any)
					*dest[0].(*string) = tt.mockPayment.ID
					*dest[1].(*string) = tt.mockPayment.IdempotencyKey
					*dest[2].(*string) = tt.mockPayment.UserID
					*dest[3].(*float64) = tt.mockPayment.Amount
					*dest[4].(*domain.Currency) = tt.mockPayment.Currency
					*dest[5].(*domain.Status) = tt.mockPayment.Status
					*dest[6].(*time.Time) = tt.mockPayment.CreatedAt
					*dest[7].(*time.Time) = tt.mockPayment.UpdatedAt
				}).Return(tt.mockScanError)
			} else {
				mockScanner.On("Scan", mock.Anything).Return(tt.mockScanError)
			}

			mockDB.On("QueryRowContext", mock.Anything, mock.Anything, mock.Anything).Return(mockScanner)

			repo := &PaymentRepository{db: mockDB}

			// Act
			result, err := repo.GetByID(context.Background(), tt.paymentID)

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

			mockDB.AssertExpectations(t)
			mockScanner.AssertExpectations(t)
		})
	}
}

func TestPaymentRepository_GetByIDempotencyKey(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name            string
		idempotencyKey  string
		mockPayment     *domain.Payment
		mockScanError   error
		expectedPayment *domain.Payment
		expectedError   error
	}{
		{
			name:           "when payment exists it should return payment and no error",
			idempotencyKey: "key_456",
			mockPayment: &domain.Payment{
				ID:             "pay_123",
				IdempotencyKey: "key_456",
				UserID:         "user_789",
				Amount:         100.50,
				Currency:       domain.CurrencyUSD,
				Status:         domain.StatusPending,
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			mockScanError: nil,
			expectedPayment: &domain.Payment{
				ID:             "pay_123",
				IdempotencyKey: "key_456",
				UserID:         "user_789",
				Amount:         100.50,
				Currency:       domain.CurrencyUSD,
				Status:         domain.StatusPending,
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			expectedError: nil,
		},
		{
			name:            "when payment does not exist it should return nil payment and no error",
			idempotencyKey:  "key_nonexistent",
			mockPayment:     nil,
			mockScanError:   sql.ErrNoRows,
			expectedPayment: nil,
			expectedError:   nil,
		},
		{
			name:            "when database error occurs it should return wrapped error",
			idempotencyKey:  "key_456",
			mockPayment:     nil,
			mockScanError:   errors.New("connection refused"),
			expectedPayment: nil,
			expectedError:   errors.New("payment repository: get by idempotency key: connection refused"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDB := new(database.MockDB)
			mockScanner := new(database.MockRowScanner)

			if tt.mockPayment != nil {
				mockScanner.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
					dest := args.Get(0).([]any)
					*dest[0].(*string) = tt.mockPayment.ID
					*dest[1].(*string) = tt.mockPayment.IdempotencyKey
					*dest[2].(*string) = tt.mockPayment.UserID
					*dest[3].(*float64) = tt.mockPayment.Amount
					*dest[4].(*domain.Currency) = tt.mockPayment.Currency
					*dest[5].(*domain.Status) = tt.mockPayment.Status
					*dest[6].(*time.Time) = tt.mockPayment.CreatedAt
					*dest[7].(*time.Time) = tt.mockPayment.UpdatedAt
				}).Return(tt.mockScanError)
			} else {
				mockScanner.On("Scan", mock.Anything).Return(tt.mockScanError)
			}

			mockDB.On("QueryRowContext", mock.Anything, mock.Anything, mock.Anything).Return(mockScanner)

			repo := &PaymentRepository{db: mockDB}

			// Act
			result, err := repo.GetByIDempotencyKey(context.Background(), tt.idempotencyKey)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else if tt.expectedPayment == nil {
				assert.NoError(t, err)
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

			mockDB.AssertExpectations(t)
			mockScanner.AssertExpectations(t)
		})
	}
}

func TestPaymentRepository_Save(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name                string
		payment             *domain.Payment
		mockTransactionError error
		expectedError       error
	}{
		{
			name: "when payment is valid it should save successfully and no error",
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
			mockTransactionError: nil,
			expectedError:        nil,
		},
		{
			name: "when transaction fails it should return wrapped error",
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
			mockTransactionError: errors.New("insert event: duplicate key"),
			expectedError:        errors.New("payment repository: save: insert event: duplicate key"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDB := new(database.MockDB)
			mockDB.On("WithTransaction", mock.Anything, mock.Anything).Return(tt.mockTransactionError)

			repo := &PaymentRepository{db: mockDB}

			// Act
			err := repo.Save(context.Background(), tt.payment)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestPaymentRepository_UpdateStatus(t *testing.T) {
	tests := []struct {
		name                string
		paymentID           string
		status              domain.Status
		gatewayRef          string
		mockTransactionError error
		expectedError       error
	}{
		{
			name:                "when payment exists it should update status successfully and no error",
			paymentID:           "pay_123",
			status:              domain.StatusCompleted,
			gatewayRef:          "gw_ref_456",
			mockTransactionError: nil,
			expectedError:       nil,
		},
		{
			name:                "when updating to failed status it should update successfully and no error",
			paymentID:           "pay_123",
			status:              domain.StatusFailed,
			gatewayRef:          "",
			mockTransactionError: nil,
			expectedError:       nil,
		},
		{
			name:                "when transaction fails it should return wrapped error",
			paymentID:           "pay_123",
			status:              domain.StatusCompleted,
			gatewayRef:          "gw_ref_456",
			mockTransactionError: errors.New("payment not found: pay_123"),
			expectedError:       errors.New("payment repository: update status: payment not found: pay_123"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDB := new(database.MockDB)
			mockDB.On("WithTransaction", mock.Anything, mock.Anything).Return(tt.mockTransactionError)

			repo := &PaymentRepository{db: mockDB}

			// Act
			err := repo.UpdateStatus(context.Background(), tt.paymentID, tt.status, tt.gatewayRef)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestPaymentRepository_GetEventsByPaymentID(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name           string
		paymentID      string
		mockEvents     []*domain.Event
		mockQueryError error
		mockScanError  error
		mockRowsError  error
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
					EventType: "completed",
					Payload:   json.RawMessage(`{"status":"completed"}`),
					CreatedAt: fixedTime,
				},
			},
			mockQueryError: nil,
			mockScanError:  nil,
			mockRowsError:  nil,
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
					EventType: "completed",
					Payload:   json.RawMessage(`{"status":"completed"}`),
					CreatedAt: fixedTime,
				},
			},
			expectedError: nil,
		},
		{
			name:           "when no events exist it should return empty slice and no error",
			paymentID:      "pay_123",
			mockEvents:     []*domain.Event{},
			mockQueryError: nil,
			mockScanError:  nil,
			mockRowsError:  nil,
			expectedEvents: []*domain.Event{},
			expectedError:  nil,
		},
		{
			name:           "when query fails it should return wrapped error",
			paymentID:      "pay_123",
			mockEvents:     nil,
			mockQueryError: errors.New("connection refused"),
			mockScanError:  nil,
			mockRowsError:  nil,
			expectedEvents: nil,
			expectedError:  errors.New("payment repository: get events by payment id: connection refused"),
		},
		{
			name:           "when scan fails it should return wrapped error",
			paymentID:      "pay_123",
			mockEvents:     nil,
			mockQueryError: nil,
			mockScanError:  errors.New("scan error"),
			mockRowsError:  nil,
			expectedEvents: nil,
			expectedError:  errors.New("payment repository: scan event: scan error"),
		},
		{
			name:           "when rows iteration fails it should return wrapped error",
			paymentID:      "pay_123",
			mockEvents:     []*domain.Event{},
			mockQueryError: nil,
			mockScanError:  nil,
			mockRowsError:  errors.New("iteration error"),
			expectedEvents: nil,
			expectedError:  errors.New("payment repository: iterate events: iteration error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDB := new(database.MockDB)
			mockRows := new(database.MockRows)

			if tt.mockQueryError != nil {
				mockDB.On("QueryContext", mock.Anything, mock.Anything, mock.Anything).Return(nil, tt.mockQueryError)
			} else {
				eventCount := len(tt.mockEvents)
				if tt.mockScanError != nil {
					mockRows.On("Next").Return(true).Once()
					mockRows.On("Scan", mock.Anything).Return(tt.mockScanError)
					mockRows.On("Close").Return(nil)
				} else if tt.mockRowsError != nil {
					mockRows.On("Next").Return(false).Once()
					mockRows.On("Close").Return(nil)
					mockRows.On("Err").Return(tt.mockRowsError)
				} else {
					if eventCount > 0 {
						mockRows.On("Next").Return(true).Times(eventCount)
						mockRows.On("Next").Return(false).Once()
						scanCallCount := 0
						mockRows.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
							dest := args.Get(0).([]any)
							event := tt.mockEvents[scanCallCount]
							*dest[0].(*string) = event.ID
							*dest[1].(*string) = event.PaymentID
							*dest[2].(*int) = event.Sequence
							*dest[3].(*string) = event.EventType
							*dest[4].(*json.RawMessage) = event.Payload
							*dest[5].(*time.Time) = event.CreatedAt
							scanCallCount++
						}).Return(nil).Times(eventCount)
					} else {
						mockRows.On("Next").Return(false).Once()
					}
					mockRows.On("Close").Return(nil)
					mockRows.On("Err").Return(nil)
				}
				mockDB.On("QueryContext", mock.Anything, mock.Anything, mock.Anything).Return(mockRows, nil)
			}

			repo := &PaymentRepository{db: mockDB}

			// Act
			result, err := repo.GetEventsByPaymentID(context.Background(), tt.paymentID)

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

			mockDB.AssertExpectations(t)
		})
	}
}
