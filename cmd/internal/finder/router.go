package finder

import (
	"github.com/gin-gonic/gin"
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/repository/paymentstorer"
)

// Start starts the finder router
// It starts the finder router and returns an error if the builder fails
func Start(rg *gin.RouterGroup, db paymentstorer.PaymentDB) error {
	h, err := Build(db)
	if err != nil {
		return err
	}

	rg.GET("/payments/:id", h.Find)
	rg.GET("/payments/:id/events", h.FindEvents)
	return nil
}
