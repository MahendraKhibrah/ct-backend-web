package Route

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitInvoice(c *gin.RouterGroup, db *gorm.DB) {
	r := InvoiceDI(db)
	m := CommonMiddlewareDI()

	c.Use(m.Authentication)

	c.GET("/token-validator", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "token ok",
		})
	})
	c.POST("/add-invoice", r.AddInvoice)
	c.GET("/get-all-invoice", r.GetAllInvoice)
	c.GET("/get-invoice/:id", r.GetInvoiceById)
	c.POST("/lock-invoice", r.LockInvoice)
	c.POST("/add-sale-to-invoice", r.AddSaleToInvoice)
	c.PUT("/update-sale", r.UpdateSale)
	c.DELETE("/delete-sale", r.DeleteSale)
	c.PUT("/update-faktur", r.UpdateFaktur)
	c.PUT("/update-main-information", r.UpdateMainInformation)
	c.PUT("/update-note", r.UpdateNote)
	c.PUT("/update-status", r.UpdateStatus)
	c.GET("/get-all-sale/:invoiceId", r.GetAllSale)
	c.DELETE("/delete-invoice", r.DeleteInvoice)
	c.POST("/update-po-file/:invoiceId", r.UpdatePoFile)
	c.POST("/get-po-url", r.GetPoUrl)
	c.POST("/update-faktur-file/:invoiceId", r.UpdateFakturFile)
	c.POST("/get-faktur-url", r.GetFakturUrl)
	c.GET("/get-previous-sale", r.GetPreviousSale)
}
