package Controller

import (
	"ct-backend/Model"
	"ct-backend/Model/Common"
	"ct-backend/Model/Dto"
	"ct-backend/Services"
	"github.com/gin-gonic/gin"
	"strconv"
)

type (
	IDeliveryController interface {
		GetById(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		CreateDeliveryOrder(ctx *gin.Context)
		CreateDeliveryProduct(ctx *gin.Context)
		UpdateDeliveryProduct(ctx *gin.Context)
		DeleteDeliveryProduct(ctx *gin.Context)
		GetAllDeliveryProduct(ctx *gin.Context)
		UpdateSender(ctx *gin.Context)
		UpdateMainInformation(ctx *gin.Context)
		GetPreviousNote(ctx *gin.Context)
		LockDeliveryOrder(ctx *gin.Context)
		GetAvailableInvoices(ctx *gin.Context)
		GetAvailableSales(ctx *gin.Context)
		AddAllAvailableSale(ctx *gin.Context)
	}

	DeliveryController struct {
		DeliveryService Services.IDeliveryService
	}
)

func DeliveryControllerProvider(deliveryService Services.IDeliveryService) *DeliveryController {
	return &DeliveryController{
		DeliveryService: deliveryService,
	}
}

func (h *DeliveryController) GetById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	delivery, err := h.DeliveryService.GetById(id)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    delivery,
	})
}

func (h *DeliveryController) GetAll(ctx *gin.Context) {
	deliveries, err := h.DeliveryService.GetAll(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	if deliveries == nil {
		deliveries = []Model.ShortDeliveryOrder{}
	}

	pagination := Common.Pagination{
		Total: ctx.GetInt("total_pages"),
		Limit: ctx.GetInt("page_size"),
		Page:  ctx.GetInt("page"),
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    deliveries,
		"meta":    pagination,
	})
}

func (h *DeliveryController) CreateDeliveryOrder(ctx *gin.Context) {
	var request *Dto.IdRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.DeliveryService.CreateDeliveryOrder(request); err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
	})
}

func (h *DeliveryController) CreateDeliveryProduct(ctx *gin.Context) {
	var request *Dto.CreateDeliveryProductRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.DeliveryService.CreateDeliveryProduct(request); err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
	})
}

func (h *DeliveryController) UpdateDeliveryProduct(ctx *gin.Context) {
	var request *Dto.UpdateDeliveryProductRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.DeliveryService.UpdateDeliveryProduct(request); err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
	})
}

func (h *DeliveryController) DeleteDeliveryProduct(ctx *gin.Context) {
	var request *Dto.DeleteDeliveryProductRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.DeliveryService.DeleteDeliveryProduct(request); err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
	})
}

func (h *DeliveryController) GetAllDeliveryProduct(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	deliveryProducts, err := h.DeliveryService.GetAllDeliveryProduct(&Dto.IdRequest{Id: id})
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    deliveryProducts,
	})
}

func (h *DeliveryController) UpdateSender(ctx *gin.Context) {
	var request *Dto.UpdateSenderRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.DeliveryService.UpdateSender(request); err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
	})
}

func (h *DeliveryController) UpdateMainInformation(ctx *gin.Context) {
	var request *Dto.UpdateDeliveryInformationRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.DeliveryService.UpdateMainInformation(request); err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
	})
}

func (h *DeliveryController) GetPreviousNote(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	note, err := h.DeliveryService.GetPreviousNote(id)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    note,
	})
}

func (h *DeliveryController) LockDeliveryOrder(ctx *gin.Context) {
	var request *Dto.LockDeliveryOrderRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.DeliveryService.LockDeliveryOrder(request); err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
	})
}

func (h *DeliveryController) GetAvailableInvoices(ctx *gin.Context) {
	invoices, err := h.DeliveryService.GetAvailableInvoices()
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    invoices,
	})
}

func (h *DeliveryController) GetAvailableSales(ctx *gin.Context) {
	invoiceId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	sales, err := h.DeliveryService.GetAvailableSales(invoiceId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    sales,
	})
}

func (h *DeliveryController) AddAllAvailableSale(ctx *gin.Context) {
	var request *Dto.LockDeliveryOrderRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	sales, err := h.DeliveryService.GetAvailableSales(request.InvoiceId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	for _, sale := range sales {
		request := &Dto.CreateDeliveryProductRequest{
			DeliveryID: request.DeliveryId,
			SalesID:    sale.ID,
			Quantity:   sale.NotSentCount,
		}

		if err := h.DeliveryService.CreateDeliveryProduct(request); err != nil {
			ctx.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	ctx.JSON(200, gin.H{
		"message": "success",
	})
}
