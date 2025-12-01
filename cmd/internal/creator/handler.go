package creator

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"

	"github.com/gin-gonic/gin"
)

// PaymentResolver defines the interface for payment creation business logic
type PaymentCreator interface {
	Create(ctx context.Context, pr *PaymentRequest) (*domain.Payment, error)
}

// Handler handles HTTP requests for payment operations
type Handler struct {
	paymentCreator PaymentCreator
}

// NewHandler creates a new Payment controller
func NewHandler(pc PaymentCreator) (*Handler, error) {
	if pc == nil {
		return nil, errors.New("payment handler: payment resolver cannot be nil")
	}

	return &Handler{
		paymentCreator: pc,
	}, nil
}

// Create handles POST /payments requests
func (h *Handler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var pr PaymentRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&pr); err != nil {
		slog.ErrorContext(ctx, "Failed to decode payment request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   "bad request",
		})
		return
	}

	if err := pr.Validate(); err != nil {
		slog.WarnContext(ctx, "Invalid payment request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed",
			"error":   "bad request",
		})
		return
	}

	payment, err := h.paymentCreator.Create(ctx, &pr)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create payment", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create payment",
			"error":   "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "payment created successfully",
		"data":    payment,
	})
}
