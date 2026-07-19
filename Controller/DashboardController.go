package Controller

import (
	"ct-backend/Services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	IDashboardController interface {
		GetWorkQueue(ctx *gin.Context)
	}

	DashboardController struct {
		DashboardService Services.IDashboardService
	}
)

func DashboardControllerProvider(dashboardService Services.IDashboardService) *DashboardController {
	return &DashboardController{DashboardService: dashboardService}
}

func (h *DashboardController) GetWorkQueue(ctx *gin.Context) {
	workQueue, err := h.DashboardService.GetWorkQueue()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    workQueue,
	})
}
