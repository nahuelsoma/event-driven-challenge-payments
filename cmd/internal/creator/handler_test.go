package creator

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name           string
		paymentCreator PaymentCreator
		expectedError  string
	}{
		{
			name:           "when payment creator is provided it should create handler successfully and no error",
			paymentCreator: new(MockPaymentCreatorService),
			expectedError:  "",
		},
		{
			name:           "when payment creator is nil it should return error",
			paymentCreator: nil,
			expectedError:  "payment handler: payment resolver cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (Payment creator already prepared in test struct)

			// Act
			result, err := NewHandler(tt.paymentCreator)

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

func TestHandler_Create(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name               string
		idempotencyKey     string
		requestBody        interface{}
		mockPayment        *domain.Payment
		mockCreateError    error
		shouldCallCreate   bool
		expectedStatusCode int
		expectedMessage    string
	}{
		{
			name:           "when request is valid it should create payment and return 201",
			idempotencyKey: "key_123",
			requestBody: PaymentRequest{
				UserID:   "user_123",
				Amount:   100.50,
				Currency: domain.CurrencyUSD,
			},
			mockPayment: &domain.Payment{
				ID:             "pay_123",
				IdempotencyKey: "key_123",
				UserID:         "user_123",
				Amount:         100.50,
				Currency:       domain.CurrencyUSD,
				Status:         domain.StatusReserved,
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			mockCreateError:    nil,
			shouldCallCreate:   true,
			expectedStatusCode: http.StatusCreated,
			expectedMessage:    "payment created successfully",
		},
		{
			name:               "when idempotency key is missing it should return 400",
			idempotencyKey:     "",
			requestBody:        PaymentRequest{UserID: "user_123", Amount: 100.50, Currency: domain.CurrencyUSD},
			shouldCallCreate:   false,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Idempotency-Key header is required",
		},
		{
			name:               "when request body is invalid JSON it should return 400",
			idempotencyKey:     "key_123",
			requestBody:        "invalid json",
			shouldCallCreate:   false,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "invalid request body",
		},
		{
			name:               "when user ID is empty it should return 400",
			idempotencyKey:     "key_123",
			requestBody:        PaymentRequest{UserID: "", Amount: 100.50, Currency: domain.CurrencyUSD},
			shouldCallCreate:   false,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "user ID is required",
		},
		{
			name:               "when amount is zero it should return 400",
			idempotencyKey:     "key_123",
			requestBody:        PaymentRequest{UserID: "user_123", Amount: 0, Currency: domain.CurrencyUSD},
			shouldCallCreate:   false,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "amount must be greater than 0",
		},
		{
			name:               "when currency is invalid it should return 400",
			idempotencyKey:     "key_123",
			requestBody:        PaymentRequest{UserID: "user_123", Amount: 100.50, Currency: domain.Currency("INVALID")},
			shouldCallCreate:   false,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "invalid currency",
		},
		{
			name:               "when payment creator fails it should return 500",
			idempotencyKey:     "key_123",
			requestBody:        PaymentRequest{UserID: "user_123", Amount: 100.50, Currency: domain.CurrencyUSD},
			mockPayment:        nil,
			mockCreateError:    errors.New("database error"),
			shouldCallCreate:   true,
			expectedStatusCode: http.StatusInternalServerError,
			expectedMessage:    "failed to create payment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCreator := new(MockPaymentCreatorService)
			if tt.shouldCallCreate {
				mockCreator.On("Create", mock.Anything, tt.idempotencyKey, mock.Anything).Return(tt.mockPayment, tt.mockCreateError)
			}

			handler := &Handler{paymentCreator: mockCreator}

			var body []byte
			switch v := tt.requestBody.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.idempotencyKey != "" {
				req.Header.Set("Idempotency-Key", tt.idempotencyKey)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Act
			handler.Create(c)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMessage, response["message"])

			mockCreator.AssertExpectations(t)
		})
	}
}
