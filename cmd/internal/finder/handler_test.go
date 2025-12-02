package finder

import (
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
		paymentFinder  PaymentFinder
		expectedError  string
	}{
		{
			name:           "when payment finder is provided it should create handler successfully and no error",
			paymentFinder:  new(MockPaymentFinder),
			expectedError:  "",
		},
		{
			name:           "when payment finder is nil it should return error",
			paymentFinder:  nil,
			expectedError:  "payment handler: payment finder cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (Payment finder already prepared in test struct)

			// Act
			result, err := NewHandler(tt.paymentFinder)

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

func TestHandler_Find(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name               string
		paymentID          string
		mockPayment        *domain.Payment
		mockFindError      error
		shouldCallFind     bool
		expectedStatusCode int
		expectedMessage    string
	}{
		{
			name:      "when payment ID is valid and payment is found it should return 200",
			paymentID: "pay_123",
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
			mockFindError:      nil,
			shouldCallFind:     true,
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "payment found successfully",
		},
		{
			name:               "when payment ID is empty it should return 400",
			paymentID:          "",
			shouldCallFind:     false,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "payment ID is required",
		},
		{
			name:               "when payment is not found it should return 404",
			paymentID:          "pay_not_found",
			mockPayment:        nil,
			mockFindError:      domain.ErrPaymentNotFound,
			shouldCallFind:     true,
			expectedStatusCode: http.StatusNotFound,
			expectedMessage:    "payment not found",
		},
		{
			name:               "when find fails with internal error it should return 500",
			paymentID:          "pay_error",
			mockPayment:        nil,
			mockFindError:      errors.New("database error"),
			shouldCallFind:     true,
			expectedStatusCode: http.StatusInternalServerError,
			expectedMessage:    "failed to find payment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockFinder := new(MockPaymentFinder)
			if tt.shouldCallFind {
				filter := &PaymentFilter{PaymentID: tt.paymentID}
				mockFinder.On("Find", mock.Anything, filter).Return(tt.mockPayment, tt.mockFindError)
			}

			handler := &Handler{paymentFinder: mockFinder}

			req := httptest.NewRequest(http.MethodGet, "/payments/"+tt.paymentID, nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{gin.Param{Key: "id", Value: tt.paymentID}}

			// Act
			handler.Find(c)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMessage, response["message"])

			if tt.shouldCallFind {
				mockFinder.AssertExpectations(t)
			}
		})
	}
}

func TestHandler_FindEvents(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name               string
		paymentID          string
		mockEvents         []*domain.Event
		mockFindEventsError error
		shouldCallFindEvents bool
		expectedStatusCode   int
		expectedMessage      string
	}{
		{
			name:      "when payment ID is valid and events are found it should return 200",
			paymentID: "pay_123",
			mockEvents: []*domain.Event{
				{
					ID:        "event_1",
					PaymentID: "pay_123",
					Sequence:  1,
					EventType: "created",
					Payload:   json.RawMessage(`{"status":"pending"}`),
					CreatedAt: fixedTime,
				},
				{
					ID:        "event_2",
					PaymentID: "pay_123",
					Sequence:  2,
					EventType: "status_changed",
					Payload:   json.RawMessage(`{"status":"reserved"}`),
					CreatedAt: fixedTime,
				},
			},
			mockFindEventsError:  nil,
			shouldCallFindEvents: true,
			expectedStatusCode:   http.StatusOK,
			expectedMessage:      "events found successfully",
		},
		{
			name:                 "when payment ID is empty it should return 400",
			paymentID:            "",
			shouldCallFindEvents: false,
			expectedStatusCode:   http.StatusBadRequest,
			expectedMessage:      "payment ID is required",
		},
		{
			name:                 "when events list is empty it should return 200 with empty array",
			paymentID:            "pay_no_events",
			mockEvents:           []*domain.Event{},
			mockFindEventsError:  nil,
			shouldCallFindEvents: true,
			expectedStatusCode:   http.StatusOK,
			expectedMessage:      "events found successfully",
		},
		{
			name:                 "when find events fails with internal error it should return 500",
			paymentID:            "pay_error",
			mockEvents:            nil,
			mockFindEventsError:  errors.New("database error"),
			shouldCallFindEvents: true,
			expectedStatusCode:   http.StatusInternalServerError,
			expectedMessage:      "failed to find events",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockFinder := new(MockPaymentFinder)
			if tt.shouldCallFindEvents {
				mockFinder.On("FindEvents", mock.Anything, tt.paymentID).Return(tt.mockEvents, tt.mockFindEventsError)
			}

			handler := &Handler{paymentFinder: mockFinder}

			req := httptest.NewRequest(http.MethodGet, "/payments/"+tt.paymentID+"/events", nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{gin.Param{Key: "id", Value: tt.paymentID}}

			// Act
			handler.FindEvents(c)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMessage, response["message"])

			if tt.shouldCallFindEvents {
				mockFinder.AssertExpectations(t)
			}
		})
	}
}

