package creator

import (
	"github.com/gin-gonic/gin"
)

func Start(rg *gin.RouterGroup, databaseConn paymentStorerDB, client interface{}, messageQueueConn interface{}) error {
	ps, err := NewPaymentStorerRepository(databaseConn)
	if err != nil {
		return err
	}

	wr, err := NewWalletReserverRepository(client)
	if err != nil {
		return err
	}

	pp, err := NewPaymentPublisherRepository(messageQueueConn)
	if err != nil {
		return err
	}

	pc, err := NewPaymentCreatorService(ps, wr, pp)
	if err != nil {
		return err
	}

	h, err := NewHandler(pc)
	if err != nil {
		return err
	}

	rg.POST("/payments", h.Create)
	return nil
}
