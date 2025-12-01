package creator

import (
	"github.com/gin-gonic/gin"
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/paymentstorer"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/messagebroker"
)

func Start(rg *gin.RouterGroup, db paymentstorer.PaymentDB, c interface{}, mbc *messagebroker.Connection, exchange, queueName string) error {
	h, err := Build(db, c, mbc, exchange, queueName)
	if err != nil {
		return err
	}

	rg.POST("/payments", h.Create)
	return nil
}
