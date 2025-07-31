package Controller

import (
	"ct-backend/Model/Common"
	"ct-backend/Model/Dto"
	"ct-backend/Services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type (
	IReceiptController interface {
		GetReceipts(ctx *gin.Context)
		GetReceiptById(ctx *gin.Context)
		CreateReceipt(ctx *gin.Context)
		GetReceiptByInvoiceId(ctx *gin.Context)
		AddInvoiceToReceipt(ctx *gin.Context)
		LockReceipt(ctx *gin.Context)
		DeleteReceiptInvoice(ctx *gin.Context)
		GetAvailableInvoices(ctx *gin.Context)
		GetClientReceipts(ctx *gin.Context)
		PayReceipt(ctx *gin.Context)
	}

	ReceiptController struct {
		service Services.IReceiptService
	}
)

func ReceiptControllerProvider(service Services.IReceiptService) *ReceiptController {
	return &ReceiptController{service: service}
}

func (h *ReceiptController) GetReceipts(ctx *gin.Context) {
	receipts, err := h.service.GetAllReceipt(ctx)

	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})

		return
	}

	pagination := Common.Pagination{
		Total: ctx.GetInt("total_pages"),
		Limit: ctx.GetInt("page_size"),
		Page:  ctx.GetInt("page"),
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    receipts,
		"meta":    pagination,
	})
}

func (h *ReceiptController) GetReceiptById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
	}

	receipt, err := h.service.GetReceiptById(id)

	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})

		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    receipt,
	})
}

func (h *ReceiptController) CreateReceipt(ctx *gin.Context) {
	var request *Dto.IdRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	_, err := h.service.CreateReceipt(request.Id)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *ReceiptController) GetReceiptByInvoiceId(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
	}

	receipts, err := h.service.GetInvoiceByReceiptId(id)

	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})

		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    receipts,
	})
}

func (h *ReceiptController) AddInvoiceToReceipt(ctx *gin.Context) {
	var request *Dto.ReceiptInvoiceRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	_, err := h.service.AddInvoiceReceipt(request)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *ReceiptController) LockReceipt(ctx *gin.Context) {
	var request *Dto.IdRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := h.service.LockReceipt(request.Id)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *ReceiptController) DeleteReceiptInvoice(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
	}

	err = h.service.DeleteReceiptInvoice(id)

	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})

		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
	})
}

func (h *ReceiptController) GetAvailableInvoices(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	invoices, err := h.service.GetAvailableInvoices(id)

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

func (h *ReceiptController) GetClientReceipts(ctx *gin.Context) {
	clientReceipts, err := h.service.GetClientReceipts()

	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})

		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    clientReceipts,
	})
}

func (h *ReceiptController) PayReceipt(ctx *gin.Context) {
	var request *Dto.IdRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := h.service.PayReceipt(request.Id)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
