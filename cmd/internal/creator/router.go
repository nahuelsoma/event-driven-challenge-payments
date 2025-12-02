package creator

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/repository/paymentstorer"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/messagebroker"
)

// Start starts the creator router
// It starts the creator router and returns an error if the builder fails
func Start(rg *gin.RouterGroup, db paymentstorer.PaymentDB, c *http.Client, mbc *messagebroker.Connection, exchange, queueName string) error {
	h, err := Build(db, c, mbc, exchange, queueName)
	if err != nil {
		return err
	}

	rg.POST("/payments", h.Create)
	return nil
}
