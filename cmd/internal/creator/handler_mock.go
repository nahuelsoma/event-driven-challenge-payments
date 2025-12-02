package creator

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// MockHandler is a mock implementation of Handler
type MockHandler struct {
	mock.Mock
}

// Create handles POST /payments requests
func (m *MockHandler) Create(c *gin.Context) {
	m.Called(c)
}
