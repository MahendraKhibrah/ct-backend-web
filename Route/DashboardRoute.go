package Route

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitDashboard(c *gin.RouterGroup, db *gorm.DB) {
	dashboard := DashboardDI(db)
	middleware := CommonMiddlewareDI()

	c.Use(middleware.Authentication)
	c.GET("/dashboard/work-queue", dashboard.GetWorkQueue)
}
