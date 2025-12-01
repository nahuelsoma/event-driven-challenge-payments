package finder

import (
	"github.com/gin-gonic/gin"
)

func Start(rg *gin.RouterGroup, db paymentReaderDB) error {
	h, err := Build(db)
	if err != nil {
		return err
	}

	rg.GET("/payments/:id", h.Find)
	rg.GET("/payments/:id/events", h.FindEvents)
	return nil
}
