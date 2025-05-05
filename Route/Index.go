package Route

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(c *gin.Engine, db *gorm.DB) {
	prefix := c.Group("/")
	InitAuth(prefix, db)
	InitProduct(prefix, db)
	InitPurchase(prefix, db)
	InitInvoice(prefix, db)
	InitClient(prefix, db)
	InitDelivery(prefix, db)
	InitUser(prefix, db)
	InitReceipt(prefix, db)
}
