package creator

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/repository/paymentstorer"
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/repository/walletclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewPaymentCreatorService(t *testing.T) {
	tests := []struct {
		name             string
		paymentStorer    PaymentStorer
		walletReserver   WalletReserver
		paymentPublisher PaymentPublisher
		expectedError    string
	}{
		{
			name:             "when all dependencies are provided it should create service successfully and no error",
			paymentStorer:    new(paymentstorer.MockPaymentRepository),
			walletReserver:   new(walletclient.MockWalletClient),
			paymentPublisher: new(MockPaymentPublisherRepository),
			expectedError:    "",
		},
		{
			name:             "when payment storer is nil it should return error",
			paymentStorer:    nil,
			walletReserver:   new(walletclient.MockWalletClient),
			paymentPublisher: new(MockPaymentPublisherRepository),
			expectedError:    "payment creator: storer cannot be nil",
		},
		{
			name:             "when wallet reserver is nil it should return error",
			paymentStorer:    new(paymentstorer.MockPaymentRepository),
			walletReserver:   nil,
			paymentPublisher: new(MockPaymentPublisherRepository),
			expectedError:    "payment creator: wallet reserver cannot be nil",
		},
		{
			name:             "when payment publisher is nil it should return error",
			paymentStorer:    new(paymentstorer.MockPaymentRepository),
			walletReserver:   new(walletclient.MockWalletClient),
			paymentPublisher: nil,
			expectedError:    "payment creator: publisher cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (Dependencies already prepared in test struct)

			// Act
			result, err := NewPaymentCreatorService(tt.paymentStorer, tt.walletReserver, tt.paymentPublisher)

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

func TestPaymentCreatorService_Create(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name           string
		idempotencyKey string
		request        *PaymentRequest
		setupMocks     func(storer *paymentstorer.MockPaymentRepository, reserver *walletclient.MockWalletClient, publisher *MockPaymentPublisherRepository)
		expectedError  error
		expectPayment  bool
	}{
		{
			name:           "when payment already exists it should return existing payment and no error",
			idempotencyKey: "key_existing",
			request: &PaymentRequest{
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
			},
			setupMocks: func(storer *paymentstorer.MockPaymentRepository, reserver *walletclient.MockWalletClient, publisher *MockPaymentPublisherRepository) {
				existingPayment := &domain.Payment{
					ID:             "pay_existing",
					IdempotencyKey: "key_existing",
					UserID:         "user_123",
					Amount:         100.50,
					Currency:       domain.CurrencyUSD,
					Status:         domain.StatusReserved,
					CreatedAt:      fixedTime,
					UpdatedAt:      fixedTime,
				}
				storer.On("GetByIDempotencyKey", mock.Anything, "key_existing").Return(existingPayment, nil)
			},
			expectedError: nil,
			expectPayment: true,
		},
		{
			name:           "when new payment is created successfully it should return payment and no error",
			idempotencyKey: "key_new",
			request: &PaymentRequest{
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
			},
			setupMocks: func(storer *paymentstorer.MockPaymentRepository, reserver *walletclient.MockWalletClient, publisher *MockPaymentPublisherRepository) {
				storer.On("GetByIDempotencyKey", mock.Anything, "key_new").Return(nil, nil)
				storer.On("Save", mock.Anything, mock.Anything).Return(nil)
				reserver.On("Reserve", mock.Anything, "user_123", 100.50, mock.Anything).Return(nil)
				storer.On("UpdateStatus", mock.Anything, mock.Anything, domain.StatusReserved, "").Return(nil)
				publisher.On("Publish", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
			expectPayment: true,
		},
		{
			name:           "when get by idempotency key fails it should return wrapped error",
			idempotencyKey: "key_error",
			request: &PaymentRequest{
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
			},
			setupMocks: func(storer *paymentstorer.MockPaymentRepository, reserver *walletclient.MockWalletClient, publisher *MockPaymentPublisherRepository) {
				storer.On("GetByIDempotencyKey", mock.Anything, "key_error").Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("payment creator: get by idempotency key: database error"),
			expectPayment: false,
		},
		{
			name:           "when save payment fails it should return wrapped error",
			idempotencyKey: "key_save_error",
			request: &PaymentRequest{
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
			},
			setupMocks: func(storer *paymentstorer.MockPaymentRepository, reserver *walletclient.MockWalletClient, publisher *MockPaymentPublisherRepository) {
				storer.On("GetByIDempotencyKey", mock.Anything, "key_save_error").Return(nil, nil)
				storer.On("Save", mock.Anything, mock.Anything).Return(errors.New("save failed"))
			},
			expectedError: errors.New("payment creator: save payment: save failed"),
			expectPayment: false,
		},
		{
			name:           "when reserve funds fails it should update status to failed and return wrapped error",
			idempotencyKey: "key_reserve_error",
			request: &PaymentRequest{
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
			},
			setupMocks: func(storer *paymentstorer.MockPaymentRepository, reserver *walletclient.MockWalletClient, publisher *MockPaymentPublisherRepository) {
				storer.On("GetByIDempotencyKey", mock.Anything, "key_reserve_error").Return(nil, nil)
				storer.On("Save", mock.Anything, mock.Anything).Return(nil)
				reserver.On("Reserve", mock.Anything, "user_123", 100.50, mock.Anything).Return(errors.New("insufficient funds"))
				storer.On("UpdateStatus", mock.Anything, mock.Anything, domain.StatusFailed, "").Return(nil)
			},
			expectedError: errors.New("payment creator: reserve funds: insufficient funds"),
			expectPayment: false,
		},
		{
			name:           "when reserve funds fails and update status fails it should return update status error",
			idempotencyKey: "key_reserve_update_error",
			request: &PaymentRequest{
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
			},
			setupMocks: func(storer *paymentstorer.MockPaymentRepository, reserver *walletclient.MockWalletClient, publisher *MockPaymentPublisherRepository) {
				storer.On("GetByIDempotencyKey", mock.Anything, "key_reserve_update_error").Return(nil, nil)
				storer.On("Save", mock.Anything, mock.Anything).Return(nil)
				reserver.On("Reserve", mock.Anything, "user_123", 100.50, mock.Anything).Return(errors.New("insufficient funds"))
				storer.On("UpdateStatus", mock.Anything, mock.Anything, domain.StatusFailed, "").Return(errors.New("update failed"))
			},
			expectedError: errors.New("payment creator: update status to failed: update failed"),
			expectPayment: false,
		},
		{
			name:           "when update status to reserved fails it should return wrapped error",
			idempotencyKey: "key_reserved_error",
			request: &PaymentRequest{
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
			},
			setupMocks: func(storer *paymentstorer.MockPaymentRepository, reserver *walletclient.MockWalletClient, publisher *MockPaymentPublisherRepository) {
				storer.On("GetByIDempotencyKey", mock.Anything, "key_reserved_error").Return(nil, nil)
				storer.On("Save", mock.Anything, mock.Anything).Return(nil)
				reserver.On("Reserve", mock.Anything, "user_123", 100.50, mock.Anything).Return(nil)
				storer.On("UpdateStatus", mock.Anything, mock.Anything, domain.StatusReserved, "").Return(errors.New("update reserved failed"))
			},
			expectedError: errors.New("payment creator: update status to reserved: update reserved failed"),
			expectPayment: false,
		},
		{
			name:           "when publish fails it should return wrapped error",
			idempotencyKey: "key_publish_error",
			request: &PaymentRequest{
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
			},
			setupMocks: func(storer *paymentstorer.MockPaymentRepository, reserver *walletclient.MockWalletClient, publisher *MockPaymentPublisherRepository) {
				storer.On("GetByIDempotencyKey", mock.Anything, "key_publish_error").Return(nil, nil)
				storer.On("Save", mock.Anything, mock.Anything).Return(nil)
				reserver.On("Reserve", mock.Anything, "user_123", 100.50, mock.Anything).Return(nil)
				storer.On("UpdateStatus", mock.Anything, mock.Anything, domain.StatusReserved, "").Return(nil)
				publisher.On("Publish", mock.Anything, mock.Anything).Return(errors.New("publish failed"))
			},
			expectedError: errors.New("payment creator: publish payment: publish failed"),
			expectPayment: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockStorer := new(paymentstorer.MockPaymentRepository)
			mockReserver := new(walletclient.MockWalletClient)
			mockPublisher := new(MockPaymentPublisherRepository)
			tt.setupMocks(mockStorer, mockReserver, mockPublisher)

			service := &PaymentCreatorService{
				paymentStorer:    mockStorer,
				walletReserver:   mockReserver,
				paymentPublisher: mockPublisher,
			}

			// Act
			result, err := service.Create(context.Background(), tt.idempotencyKey, tt.request)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.expectPayment {
					assert.NotNil(t, result)
				}
			}

			mockStorer.AssertExpectations(t)
			mockReserver.AssertExpectations(t)
			mockPublisher.AssertExpectations(t)
		})
	}
}
