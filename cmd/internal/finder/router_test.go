package finder

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/repository/paymentstorer"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/database"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestStart(t *testing.T) {
	tests := []struct {
		name          string
		db            paymentstorer.PaymentDB
		expectedError bool
	}{
		{
			name:          "when database is nil it should return error",
			db:            nil,
			expectedError: true,
		},
		{
			name:          "when database is provided it should start router successfully and no error",
			db:            &database.MockDB{},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			router := gin.New()
			rg := router.Group("/api/v1")

			// Act
			err := Start(rg, tt.db)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
