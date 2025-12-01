package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/creator"
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/finder"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/database"
)

// StartAPI initializes and starts the HTTP API server
func StartAPI(database *database.DB, httpClient interface{}, messageBroker interface{}) error {
	r := gin.New()

	apiV1 := r.Group("/api/v1")

	// Each vertical owns its internal wiring
	if err := creator.Start(apiV1, database, httpClient, messageBroker); err != nil {
		return fmt.Errorf("api: failed to start creator vertical: %w", err)
	}

	if err := finder.Start(apiV1, database); err != nil {
		return fmt.Errorf("api: failed to start finder vertical: %w", err)
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
