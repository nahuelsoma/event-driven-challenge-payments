package finder

import (
	"github.com/gin-gonic/gin"
)

func Start(rg *gin.RouterGroup, db paymentReaderDB) error {
	pr, err := NewPaymentReaderRepository(db)
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
