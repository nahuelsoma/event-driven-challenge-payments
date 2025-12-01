package finder

import (
	"github.com/gin-gonic/gin"
)

func Start(rg *gin.RouterGroup, databaseConn interface{}) error {
	pr, err := NewPaymentReaderRepository(databaseConn)
	if err != nil {
		return err
	}

	pf, err := NewPaymentFinderService(pr)
	if err != nil {
		return err
	}

	h, err := NewHandler(pf)
	if err != nil {
		return err
	}

	rg.GET("/payments/:id", h.Find)
	return nil
}
