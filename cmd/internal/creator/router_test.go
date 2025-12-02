package creator

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestStart(t *testing.T) {
	tests := []struct {
		name          string
		db            interface{}
		httpClient    *http.Client
		mbConnection  interface{}
		exchange      string
		queueName     string
		expectedError bool
	}{
		{
			name:          "when database is nil it should return error",
			db:            nil,
			httpClient:    &http.Client{},
			mbConnection:  nil,
			exchange:      "payments",
			queueName:     "payments.created",
			expectedError: true,
		},
		{
			name:          "when http client is nil it should return error",
			db:            nil,
			httpClient:    nil,
			mbConnection:  nil,
			exchange:      "payments",
			queueName:     "payments.created",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			router := gin.New()
			rg := router.Group("/api/v1")

			// Act
			err := Start(rg, nil, tt.httpClient, nil, tt.exchange, tt.queueName)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

