package processor

import (
	"context"
	"errors"
	"testing"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/repository/paymentstorer"
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/repository/walletclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewPaymentProcessorService(t *testing.T) {
	tests := []struct {
		name              string
		paymentResolver   PaymentResolver
		walletResolver    WalletResolver
		gatewayProcessor  GatewayProcessor
		expectedError     string
	}{
		{
			name:              "when all dependencies are provided it should create service successfully and no error",
			paymentResolver:   new(paymentstorer.MockPaymentRepository),
			walletResolver:    new(walletclient.MockWalletClient),
			gatewayProcessor:  new(MockGatewayProcessor),
			expectedError:     "",
		},
		{
			name:              "when payment resolver is nil it should return error",
			paymentResolver:   nil,
			walletResolver:    new(walletclient.MockWalletClient),
			gatewayProcessor:  new(MockGatewayProcessor),
			expectedError:     "payment processor: resolver cannot be nil",
		},
		{
			name:              "when wallet resolver is nil it should return error",
			paymentResolver:   new(paymentstorer.MockPaymentRepository),
			walletResolver:    nil,
			gatewayProcessor:  new(MockGatewayProcessor),
			expectedError:     "payment processor: wallet resolver cannot be nil",
		},
		{
			name:              "when gateway processor is nil it should return error",
			paymentResolver:   new(paymentstorer.MockPaymentRepository),
			walletResolver:    new(walletclient.MockWalletClient),
			gatewayProcessor:  nil,
			expectedError:     "payment processor: gateway processor cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (Dependencies already prepared in test struct)

			// Act
			result, err := NewPaymentProcessorService(tt.paymentResolver, tt.walletResolver, tt.gatewayProcessor)

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

func TestPaymentProcessorService_Process(t *testing.T) {
	tests := []struct {
		name                    string
		payment                 *domain.Payment
		mockExistingPayment     *domain.Payment
		mockGetError            error
		mockGatewayRef          string
		mockGatewayError        error
		mockReleaseError        error
		mockConfirmError        error
		mockUpdateStatusError   error
		shouldCallGateway       bool
		shouldCallRelease       bool
		shouldCallConfirm       bool
		shouldCallUpdateStatus  bool
		expectedError           error
	}{
		{
			name: "when payment is already completed it should skip processing and return no error",
			payment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusReserved,
			},
			mockExistingPayment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusCompleted,
			},
			mockGetError:          nil,
			shouldCallGateway:     false,
			shouldCallRelease:     false,
			shouldCallConfirm:     false,
			shouldCallUpdateStatus: false,
			expectedError:         nil,
		},
		{
			name: "when payment is already failed it should skip processing and return no error",
			payment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusReserved,
			},
			mockExistingPayment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusFailed,
			},
			mockGetError:          nil,
			shouldCallGateway:     false,
			shouldCallRelease:     false,
			shouldCallConfirm:     false,
			shouldCallUpdateStatus: false,
			expectedError:         nil,
		},
		{
			name: "when payment has unexpected status it should skip processing and return no error",
			payment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusReserved,
			},
			mockExistingPayment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusPending,
			},
			mockGetError:          nil,
			shouldCallGateway:     false,
			shouldCallRelease:     false,
			shouldCallConfirm:     false,
			shouldCallUpdateStatus: false,
			expectedError:         nil,
		},
		{
			name: "when payment processing succeeds it should return no error",
			payment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusReserved,
			},
			mockExistingPayment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusReserved,
			},
			mockGetError:          nil,
			mockGatewayRef:        "gw_ref_123",
			mockGatewayError:      nil,
			mockConfirmError:      nil,
			mockUpdateStatusError: nil,
			shouldCallGateway:     true,
			shouldCallRelease:     false,
			shouldCallConfirm:     true,
			shouldCallUpdateStatus: true,
			expectedError:         nil,
		},
		{
			name: "when get payment fails it should return wrapped error",
			payment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusReserved,
			},
			mockExistingPayment:   nil,
			mockGetError:          errors.New("database error"),
			shouldCallGateway:     false,
			shouldCallRelease:     false,
			shouldCallConfirm:     false,
			shouldCallUpdateStatus: false,
			expectedError:         errors.New("payment processor: failed to get payment: database error"),
		},
		{
			name: "when gateway processing fails it should release funds and update status to failed and return no error",
			payment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusReserved,
			},
			mockExistingPayment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusReserved,
			},
			mockGetError:          nil,
			mockGatewayError:      errors.New("gateway error"),
			mockReleaseError:      nil,
			mockUpdateStatusError: nil,
			shouldCallGateway:     true,
			shouldCallRelease:     true,
			shouldCallConfirm:     false,
			shouldCallUpdateStatus: true,
			expectedError:         nil,
		},
		{
			name: "when gateway processing fails and release funds fails it should return wrapped error",
			payment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusReserved,
			},
			mockExistingPayment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusReserved,
			},
			mockGetError:          nil,
			mockGatewayError:      errors.New("gateway error"),
			mockReleaseError:      errors.New("release failed"),
			shouldCallGateway:     true,
			shouldCallRelease:     true,
			shouldCallConfirm:     false,
			shouldCallUpdateStatus: false,
			expectedError:         errors.New("payment processor: failed to release funds: release failed"),
		},
		{
			name: "when gateway processing fails and update status fails it should return wrapped error",
			payment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusReserved,
			},
			mockExistingPayment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusReserved,
			},
			mockGetError:          nil,
			mockGatewayError:      errors.New("gateway error"),
			mockReleaseError:      nil,
			mockUpdateStatusError: errors.New("update failed"),
			shouldCallGateway:     true,
			shouldCallRelease:     true,
			shouldCallConfirm:     false,
			shouldCallUpdateStatus: true,
			expectedError:         errors.New("payment processor: failed to update status to failed: update failed"),
		},
		{
			name: "when gateway succeeds but confirm funds fails it should return wrapped error",
			payment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusReserved,
			},
			mockExistingPayment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusReserved,
			},
			mockGetError:          nil,
			mockGatewayRef:        "gw_ref_123",
			mockGatewayError:      nil,
			mockConfirmError:      errors.New("confirm failed"),
			shouldCallGateway:     true,
			shouldCallRelease:     false,
			shouldCallConfirm:     true,
			shouldCallUpdateStatus: false,
			expectedError:         errors.New("payment processor: failed to confirm funds: confirm failed"),
		},
		{
			name: "when gateway succeeds and confirm succeeds but update status fails it should return wrapped error",
			payment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusReserved,
			},
			mockExistingPayment: &domain.Payment{
				ID:       "pay_123",
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
				Status:   domain.StatusReserved,
			},
			mockGetError:          nil,
			mockGatewayRef:        "gw_ref_123",
			mockGatewayError:      nil,
			mockConfirmError:      nil,
			mockUpdateStatusError: errors.New("update failed"),
			shouldCallGateway:     true,
			shouldCallRelease:     false,
			shouldCallConfirm:     true,
			shouldCallUpdateStatus: true,
			expectedError:         errors.New("payment processor: failed to update status to completed: update failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockPaymentResolver := new(paymentstorer.MockPaymentRepository)
			mockWalletResolver := new(walletclient.MockWalletClient)
			mockGatewayProcessor := new(MockGatewayProcessor)

			mockPaymentResolver.On("GetByID", mock.Anything, tt.payment.ID).Return(tt.mockExistingPayment, tt.mockGetError)

			if tt.shouldCallGateway {
				mockGatewayProcessor.On("Process", mock.Anything, tt.payment.ID, tt.payment.Amount).Return(tt.mockGatewayRef, tt.mockGatewayError)
			}

			if tt.shouldCallRelease {
				mockWalletResolver.On("Release", mock.Anything, tt.payment.UserID, tt.payment.Amount, tt.payment.ID).Return(tt.mockReleaseError)
			}

			if tt.shouldCallConfirm {
				mockWalletResolver.On("Confirm", mock.Anything, tt.payment.UserID, tt.payment.Amount, tt.payment.ID).Return(tt.mockConfirmError)
			}

			if tt.shouldCallUpdateStatus {
				if tt.mockGatewayError != nil {
					mockPaymentResolver.On("UpdateStatus", mock.Anything, tt.payment.ID, domain.StatusFailed, "").Return(tt.mockUpdateStatusError)
				} else {
					mockPaymentResolver.On("UpdateStatus", mock.Anything, tt.payment.ID, domain.StatusCompleted, tt.mockGatewayRef).Return(tt.mockUpdateStatusError)
				}
			}

			service := &PaymentProcessorService{
				paymentResolver:  mockPaymentResolver,
				walletResolver:   mockWalletResolver,
				gatewayProcessor: mockGatewayProcessor,
			}

			// Act
			err := service.Process(context.Background(), tt.payment)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockPaymentResolver.AssertExpectations(t)
			mockWalletResolver.AssertExpectations(t)
			mockGatewayProcessor.AssertExpectations(t)
		})
	}
}

