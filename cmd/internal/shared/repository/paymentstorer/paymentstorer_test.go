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
		setupMock       func(mockDB *database.MockDB, mockScanner *database.MockRowScanner)
		expectedPayment *domain.Payment
		expectedError   error
	}{
		{
			name:      "when payment exists it should return payment and no error",
			paymentID: "pay_123",
			setupMock: func(mockDB *database.MockDB, mockScanner *database.MockRowScanner) {
				mockScanner.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
					dest := args.Get(0).([]any)
					*dest[0].(*string) = "pay_123"
					*dest[1].(*string) = "key_456"
					*dest[2].(*string) = "user_789"
					*dest[3].(*float64) = 100.50
					*dest[4].(*domain.Currency) = domain.CurrencyUSD
					*dest[5].(*domain.Status) = domain.StatusPending
					*dest[6].(*time.Time) = fixedTime
					*dest[7].(*time.Time) = fixedTime
				}).Return(nil)
				mockDB.On("QueryRowContext", mock.Anything, mock.Anything, mock.Anything).Return(mockScanner)
			},
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
			name:      "when payment does not exist it should return ErrPaymentNotFound",
			paymentID: "pay_nonexistent",
			setupMock: func(mockDB *database.MockDB, mockScanner *database.MockRowScanner) {
				mockScanner.On("Scan", mock.Anything).Return(sql.ErrNoRows)
				mockDB.On("QueryRowContext", mock.Anything, mock.Anything, mock.Anything).Return(mockScanner)
			},
			expectedPayment: nil,
			expectedError:   domain.ErrPaymentNotFound,
		},
		{
			name:      "when database error occurs it should return wrapped error",
			paymentID: "pay_123",
			setupMock: func(mockDB *database.MockDB, mockScanner *database.MockRowScanner) {
				mockScanner.On("Scan", mock.Anything).Return(errors.New("connection refused"))
				mockDB.On("QueryRowContext", mock.Anything, mock.Anything, mock.Anything).Return(mockScanner)
			},
			expectedPayment: nil,
			expectedError:   errors.New("payment repository: get by id: connection refused"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDB := new(database.MockDB)
			mockScanner := new(database.MockRowScanner)
			tt.setupMock(mockDB, mockScanner)

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
		setupMock       func(mockDB *database.MockDB, mockScanner *database.MockRowScanner)
		expectedPayment *domain.Payment
		expectedError   error
	}{
		{
			name:           "when payment exists it should return payment and no error",
			idempotencyKey: "key_456",
			setupMock: func(mockDB *database.MockDB, mockScanner *database.MockRowScanner) {
				mockScanner.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
					dest := args.Get(0).([]any)
					*dest[0].(*string) = "pay_123"
					*dest[1].(*string) = "key_456"
					*dest[2].(*string) = "user_789"
					*dest[3].(*float64) = 100.50
					*dest[4].(*domain.Currency) = domain.CurrencyUSD
					*dest[5].(*domain.Status) = domain.StatusPending
					*dest[6].(*time.Time) = fixedTime
					*dest[7].(*time.Time) = fixedTime
				}).Return(nil)
				mockDB.On("QueryRowContext", mock.Anything, mock.Anything, mock.Anything).Return(mockScanner)
			},
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
			name:           "when payment does not exist it should return nil payment and no error",
			idempotencyKey: "key_nonexistent",
			setupMock: func(mockDB *database.MockDB, mockScanner *database.MockRowScanner) {
				mockScanner.On("Scan", mock.Anything).Return(sql.ErrNoRows)
				mockDB.On("QueryRowContext", mock.Anything, mock.Anything, mock.Anything).Return(mockScanner)
			},
			expectedPayment: nil,
			expectedError:   nil,
		},
		{
			name:           "when database error occurs it should return wrapped error",
			idempotencyKey: "key_456",
			setupMock: func(mockDB *database.MockDB, mockScanner *database.MockRowScanner) {
				mockScanner.On("Scan", mock.Anything).Return(errors.New("connection refused"))
				mockDB.On("QueryRowContext", mock.Anything, mock.Anything, mock.Anything).Return(mockScanner)
			},
			expectedPayment: nil,
			expectedError:   errors.New("payment repository: get by idempotency key: connection refused"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDB := new(database.MockDB)
			mockScanner := new(database.MockRowScanner)
			tt.setupMock(mockDB, mockScanner)

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
		name          string
		payment       *domain.Payment
		setupMock     func(mockDB *database.MockDB)
		expectedError error
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
			setupMock: func(mockDB *database.MockDB) {
				mockDB.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
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
			setupMock: func(mockDB *database.MockDB) {
				mockDB.On("WithTransaction", mock.Anything, mock.Anything).Return(errors.New("insert event: duplicate key"))
			},
			expectedError: errors.New("payment repository: save: insert event: duplicate key"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDB := new(database.MockDB)
			tt.setupMock(mockDB)

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
		name          string
		paymentID     string
		status        domain.Status
		gatewayRef    string
		setupMock     func(mockDB *database.MockDB)
		expectedError error
	}{
		{
			name:       "when payment exists it should update status successfully and no error",
			paymentID:  "pay_123",
			status:     domain.StatusCompleted,
			gatewayRef: "gw_ref_456",
			setupMock: func(mockDB *database.MockDB) {
				mockDB.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:       "when updating to failed status it should update successfully and no error",
			paymentID:  "pay_123",
			status:     domain.StatusFailed,
			gatewayRef: "",
			setupMock: func(mockDB *database.MockDB) {
				mockDB.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:       "when transaction fails it should return wrapped error",
			paymentID:  "pay_123",
			status:     domain.StatusCompleted,
			gatewayRef: "gw_ref_456",
			setupMock: func(mockDB *database.MockDB) {
				mockDB.On("WithTransaction", mock.Anything, mock.Anything).Return(errors.New("payment not found: pay_123"))
			},
			expectedError: errors.New("payment repository: update status: payment not found: pay_123"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDB := new(database.MockDB)
			tt.setupMock(mockDB)

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
		setupMock      func(mockDB *database.MockDB, mockRows *database.MockRows)
		expectedEvents []*domain.Event
		expectedError  error
	}{
		{
			name:      "when events exist it should return events and no error",
			paymentID: "pay_123",
			setupMock: func(mockDB *database.MockDB, mockRows *database.MockRows) {
				callCount := 0
				mockRows.On("Next").Return(true).Times(2)
				mockRows.On("Next").Return(false).Once()
				mockRows.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
					dest := args.Get(0).([]any)
					callCount++
					if callCount == 1 {
						*dest[0].(*string) = "evt_1"
						*dest[1].(*string) = "pay_123"
						*dest[2].(*int) = 1
						*dest[3].(*string) = "created"
						*dest[4].(*json.RawMessage) = json.RawMessage(`{"status":"pending"}`)
						*dest[5].(*time.Time) = fixedTime
					} else {
						*dest[0].(*string) = "evt_2"
						*dest[1].(*string) = "pay_123"
						*dest[2].(*int) = 2
						*dest[3].(*string) = "completed"
						*dest[4].(*json.RawMessage) = json.RawMessage(`{"status":"completed"}`)
						*dest[5].(*time.Time) = fixedTime
					}
				}).Return(nil).Times(2)
				mockRows.On("Close").Return(nil)
				mockRows.On("Err").Return(nil)
				mockDB.On("QueryContext", mock.Anything, mock.Anything, mock.Anything).Return(mockRows, nil)
			},
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
			name:      "when no events exist it should return empty slice and no error",
			paymentID: "pay_123",
			setupMock: func(mockDB *database.MockDB, mockRows *database.MockRows) {
				mockRows.On("Next").Return(false).Once()
				mockRows.On("Close").Return(nil)
				mockRows.On("Err").Return(nil)
				mockDB.On("QueryContext", mock.Anything, mock.Anything, mock.Anything).Return(mockRows, nil)
			},
			expectedEvents: []*domain.Event{},
			expectedError:  nil,
		},
		{
			name:      "when query fails it should return wrapped error",
			paymentID: "pay_123",
			setupMock: func(mockDB *database.MockDB, mockRows *database.MockRows) {
				mockDB.On("QueryContext", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("connection refused"))
			},
			expectedEvents: nil,
			expectedError:  errors.New("payment repository: get events by payment id: connection refused"),
		},
		{
			name:      "when scan fails it should return wrapped error",
			paymentID: "pay_123",
			setupMock: func(mockDB *database.MockDB, mockRows *database.MockRows) {
				mockRows.On("Next").Return(true).Once()
				mockRows.On("Scan", mock.Anything).Return(errors.New("scan error"))
				mockRows.On("Close").Return(nil)
				mockDB.On("QueryContext", mock.Anything, mock.Anything, mock.Anything).Return(mockRows, nil)
			},
			expectedEvents: nil,
			expectedError:  errors.New("payment repository: scan event: scan error"),
		},
		{
			name:      "when rows iteration fails it should return wrapped error",
			paymentID: "pay_123",
			setupMock: func(mockDB *database.MockDB, mockRows *database.MockRows) {
				mockRows.On("Next").Return(false).Once()
				mockRows.On("Close").Return(nil)
				mockRows.On("Err").Return(errors.New("iteration error"))
				mockDB.On("QueryContext", mock.Anything, mock.Anything, mock.Anything).Return(mockRows, nil)
			},
			expectedEvents: nil,
			expectedError:  errors.New("payment repository: iterate events: iteration error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDB := new(database.MockDB)
			mockRows := new(database.MockRows)
			tt.setupMock(mockDB, mockRows)

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
