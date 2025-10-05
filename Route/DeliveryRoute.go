package Route

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitDelivery(c *gin.RouterGroup, db *gorm.DB) {
	r := DeliveryDI(db)
	m := CommonMiddlewareDI()

	c.Use(m.Authentication)

	c.GET("/get-delivery/:id", r.GetById)
	c.GET("/get-all-delivery", r.GetAll)
	c.POST("/create-delivery-order", r.CreateDeliveryOrder)
	c.POST("/create-delivery-product", r.CreateDeliveryProduct)
	c.PUT("/update-delivery-product", r.UpdateDeliveryProduct)
	c.DELETE("/delete-delivery-product", r.DeleteDeliveryProduct)
	c.GET("/get-all-delivery-product/:id", r.GetAllDeliveryProduct)
	c.PUT("/update-sender", r.UpdateSender)
	c.PUT("/update-delivery-information", r.UpdateMainInformation)
	c.GET("/get-previous-note/:id", r.GetPreviousNote)
	c.PUT("/lock-delivery-order", r.LockDeliveryOrder)
	c.GET("/delivery/get-available-invoices", r.GetAvailableInvoices)
	c.GET("/delivery/get-available-sales/:id", r.GetAvailableSales)
	c.POST("/delivery/add-all-available-sale", r.AddAllAvailableSale)
}
