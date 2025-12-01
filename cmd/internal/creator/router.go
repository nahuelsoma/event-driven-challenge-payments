package creator

import (
	"github.com/gin-gonic/gin"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/messagebroker"
)

func Start(rg *gin.RouterGroup, dc paymentStorerDB, c interface{}, mbc *messagebroker.Connection, queueName string) error {
	h, err := Build(dc, c, mbc, queueName)
	if err != nil {
		return err
	}

	rg.POST("/payments", h.Create)
	return nil
}
