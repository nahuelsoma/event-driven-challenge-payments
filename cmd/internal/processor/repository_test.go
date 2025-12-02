package processor

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGatewayProcessorRepository(t *testing.T) {
	tests := []struct {
		name          string
		client        *http.Client
		expectedError string
	}{
		{
			name:          "when http client is provided it should create repository successfully and no error",
			client:        &http.Client{},
			expectedError: "",
		},
		{
			name:          "when http client is nil it should return error with message 'gateway processor: client cannot be nil'",
			client:        nil,
			expectedError: "gateway processor: client cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (HTTP client already prepared in test struct)

			// Act
			result, err := NewGatewayProcessorRepository(tt.client)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, result.client)
			}
		})
	}
}

func TestGatewayProcessorRepository_Process(t *testing.T) {
	tests := []struct {
		name              string
		paymentID         string
		amount            float64
		expectedError     error
		expectedPattern   string
		shouldMatchPattern bool
	}{
		{
			name:              "when payment ID and amount are valid it should return gateway reference and no error",
			paymentID:         "pay_123",
			amount:            100.50,
			expectedError:     nil,
			expectedPattern:   `^gw_[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`,
			shouldMatchPattern: true,
		},
		{
			name:              "when payment ID is empty it should return gateway reference and no error",
			paymentID:         "",
			amount:            50.00,
			expectedError:     nil,
			expectedPattern:   `^gw_[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`,
			shouldMatchPattern: true,
		},
		{
			name:              "when amount is zero it should return gateway reference and no error",
			paymentID:         "pay_456",
			amount:            0.00,
			expectedError:     nil,
			expectedPattern:   `^gw_[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`,
			shouldMatchPattern: true,
		},
		{
			name:              "when amount is negative it should return gateway reference and no error",
			paymentID:         "pay_789",
			amount:            -10.00,
			expectedError:     nil,
			expectedPattern:   `^gw_[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`,
			shouldMatchPattern: true,
		},
		{
			name:              "when amount is very large it should return gateway reference and no error",
			paymentID:         "pay_large",
			amount:            999999.99,
			expectedError:     nil,
			expectedPattern:   `^gw_[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`,
			shouldMatchPattern: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			repo := &GatewayProcessorRepository{
				client: &http.Client{},
			}

			// Act
			result, err := repo.Process(context.Background(), tt.paymentID, tt.amount)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
				if tt.shouldMatchPattern {
					matched, matchErr := regexp.MatchString(tt.expectedPattern, result)
					assert.NoError(t, matchErr)
					assert.True(t, matched, "gateway reference should match pattern %s, got %s", tt.expectedPattern, result)
				}
			}
		})
	}
}

