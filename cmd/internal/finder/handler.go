package finder

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"

	"github.com/gin-gonic/gin"
)

// PaymentFinder defines the interface for payment finding business logic
type PaymentFinder interface {
	Find(ctx context.Context, filter *PaymentFilter) (*domain.Payment, error)
	FindEvents(ctx context.Context, paymentID string) ([]*domain.Event, error)
}

// Handler handles HTTP requests for payment operations
type Handler struct {
	paymentFinder PaymentFinder
}

// NewHandler creates a new Payment controller
func NewHandler(pf PaymentFinder) (*Handler, error) {
	if pf == nil {
		return nil, errors.New("payment handler: payment finder cannot be nil")
	}

	return &Handler{
		paymentFinder: pf,
	}, nil
}

// Find handles GET /payments/:id requests
func (h *Handler) Find(c *gin.Context) {
	ctx := c.Request.Context()

	paymentID := c.Param("id")
	filter := &PaymentFilter{PaymentID: paymentID}

	if err := filter.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"error":   "bad request",
		})
		return
	}

	payment, err := h.paymentFinder.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "payment not found",
				"error":   "not found",
			})
			return
		}

		slog.ErrorContext(ctx, "Failed to find payment", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to find payment",
			"error":   "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "payment found successfully",
		"data":    payment,
	})
}

// FindEvents handles GET /payments/:id/events requests
func (h *Handler) FindEvents(c *gin.Context) {
	ctx := c.Request.Context()

	paymentID := c.Param("id")
	filter := &PaymentFilter{PaymentID: paymentID}

	if err := filter.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"error":   "bad request",
		})
		return
	}

	events, err := h.paymentFinder.FindEvents(ctx, paymentID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to find events", "error", err, "payment_id", paymentID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to find events",
			"error":   "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "events found successfully",
		"data":    events,
	})
}
