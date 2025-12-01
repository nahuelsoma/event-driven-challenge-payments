package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/creator"
)

// StartAPI initializes and starts the HTTP API server
func StartAPI(database interface{}, httpClient interface{}, messageBroker interface{}) error {
	// Create payment storer repository
	ps, err := creator.NewPaymentStorerRepository(database)
	if err != nil {
		return fmt.Errorf("api: failed to create payment storer repository: %w", err)
	}

	// Create wallet reserver repository
	wr, err := creator.NewWalletReserverRepository(httpClient)
	if err != nil {
		return fmt.Errorf("api: failed to create wallet reserver repository: %w", err)
	}

	// Create publisher repository
	pp, err := creator.NewPaymentPublisherRepository(messageBroker)
	if err != nil {
		return fmt.Errorf("api: failed to create publisher repository: %w", err)
	}

	// Create payment creator service
	pc, err := creator.NewPaymentCreatorService(ps, wr, pp)
	if err != nil {
		return fmt.Errorf("api: failed to create payment creator service: %w", err)
	}

	// Create handler
	h, err := creator.NewHandler(pc)
	if err != nil {
		return fmt.Errorf("api: failed to create handler: %w", err)
	}

	r := gin.New()

	apiV1 := r.Group("/api/v1")
	{
		apiV1.POST("/payments", h.Create)
	}

	r.NoRoute(func(c *gin.Context) {
		slog.WarnContext(c.Request.Context(), "Route not found",
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"client_ip", c.ClientIP(),
		)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "the requested resource was not found",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	slog.Info("Starting API server", "port", port)

	if err := r.Run("0.0.0.0:" + port); err != nil {
		return fmt.Errorf("api: failed to start server: %w", err)
	}

	return nil
}
