package Controller

import (
	"ct-backend/Model"
	"ct-backend/Model/Common"
	"ct-backend/Model/Dto"
	"ct-backend/Services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type (
	IInvoiceController interface {
		AddInvoice(ctx *gin.Context)
		GetAllInvoice(ctx *gin.Context)
		GetInvoiceById(ctx *gin.Context)
		LockInvoice(ctx *gin.Context)
		AddSaleToInvoice(ctx *gin.Context)
		UpdateSale(ctx *gin.Context)
		DeleteSale(ctx *gin.Context)
		UpdateFaktur(ctx *gin.Context)
		UpdateMainInformation(ctx *gin.Context)
		UpdateNote(ctx *gin.Context)
		UpdateStatus(ctx *gin.Context)
		GetAllSale(ctx *gin.Context)
		DeleteInvoice(ctx *gin.Context)
		UpdatePoFile(ctx *gin.Context)
		GetPoUrl(ctx *gin.Context)
		UpdateFakturFile(ctx *gin.Context)
		GetFakturUrl(ctx *gin.Context)
		GetPreviousSale(ctx *gin.Context)
	}

	InvoiceController struct {
		InvoiceService Services.IInvoiceService
		StorageService Services.IStorageService
	}
)

func InvoiceControllerProvider(invoiceService Services.IInvoiceService, storageService Services.IStorageService) *InvoiceController {
	return &InvoiceController{
		InvoiceService: invoiceService,
		StorageService: storageService,
	}
}

func (h *InvoiceController) AddInvoice(ctx *gin.Context) {
	var request *Dto.CreateInvoiceRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	if err := h.InvoiceService.AddInvoice(request); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *InvoiceController) GetAllInvoice(ctx *gin.Context) {
	var request *Dto.GetInvoicesRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	invoices, err := h.InvoiceService.GetAllInvoice(request, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	pagination := Common.Pagination{
		Total: ctx.GetInt("total_pages"),
		Limit: ctx.GetInt("page_size"),
		Page:  ctx.GetInt("page"),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    invoices,
		"meta":    pagination,
	})
}

func (h *InvoiceController) GetInvoiceById(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "id must be a number",
		})
		return
	}

	invoice, err := h.InvoiceService.GetInvoiceById(id)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    invoice,
	})
}

func (h *InvoiceController) LockInvoice(ctx *gin.Context) {
	var id Dto.IdRequest

	if err := ctx.ShouldBind(&id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.InvoiceService.LockInvoice(&id); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *InvoiceController) AddSaleToInvoice(ctx *gin.Context) {
	var request *Dto.AddSaleRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.InvoiceService.AddSaleToInvoice(request); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.InvoiceService.UpdateStatus(&Dto.UpdateStatusRequest{
		InvoiceId:       request.InvoiceId,
		InvoiceStatusId: 3,
	}); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *InvoiceController) UpdateSale(ctx *gin.Context) {
	var request *Dto.UpdateSaleRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.InvoiceService.UpdateSale(request); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *InvoiceController) DeleteSale(ctx *gin.Context) {
	var id Dto.IdRequest

	if err := ctx.ShouldBind(&id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.InvoiceService.DeleteSale(id); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *InvoiceController) UpdateFaktur(ctx *gin.Context) {
	var request *Dto.UpdateFakturRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.InvoiceService.UpdateFaktur(request); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *InvoiceController) UpdateMainInformation(ctx *gin.Context) {
	var request *Dto.UpdateMainInformationRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.InvoiceService.UpdateMainInformation(request); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *InvoiceController) UpdateNote(ctx *gin.Context) {
	var request *Dto.UpdateNoteRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.InvoiceService.UpdateNote(request); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *InvoiceController) UpdateStatus(ctx *gin.Context) {
	var request *Dto.UpdateStatusRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.InvoiceService.UpdateStatus(request); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *InvoiceController) GetAllSale(ctx *gin.Context) {
	invoiceId, err := strconv.Atoi(ctx.Param("invoiceId"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invoiceId must be a number",
		})
		return
	}

	sales, err := h.InvoiceService.GetAllSale(invoiceId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    sales,
	})
}

func (h *InvoiceController) DeleteInvoice(ctx *gin.Context) {
	var id Dto.IdRequest

	if err := ctx.ShouldBind(&id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.InvoiceService.DeleteInvoice(id); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *InvoiceController) UpdatePoFile(ctx *gin.Context) {
	var FileObj *Common.FileDto
	if err := ctx.ShouldBind(&FileObj); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	invoiceId, err := strconv.Atoi(ctx.Param("invoiceId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invoiceId must be a number",
		})
		return
	}

	if FileObj.Data == "" || FileObj.Data == "null" {
		err = h.InvoiceService.UpdateDocument(Dto.UpdateDocumentRequest{
			InvoiceId:    invoiceId,
			DocumentType: "po_path",
			DocumentPath: FileObj.File.Filename,
		})

		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"message": err.Error(),
			})
			return
		} else {
			FileObj.Data = FileObj.File.Filename
		}
	}

	if err = h.StorageService.UploadFile(&Model.S3ObjectRequest{
		File:   FileObj.File,
		Bucket: "cahaya-teknik",
		Key:    "invoice-po/" + FileObj.Data,
	}); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    FileObj.Data,
	})
}

func (h *InvoiceController) GetPoUrl(ctx *gin.Context) {
	var request *Dto.GetDocumentUrlRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	url, err := h.StorageService.GeneratePresignedURL(&Model.S3UrlRequest{
		Bucket: "cahaya-teknik",
		Key:    "invoice-po/" + request.Key,
	})

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    url,
	})
}

func (h *InvoiceController) UpdateFakturFile(ctx *gin.Context) {
	var FileObj *Common.FileDto
	if err := ctx.ShouldBind(&FileObj); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	invoiceId, err := strconv.Atoi(ctx.Param("invoiceId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invoiceId must be a number",
		})
		return
	}

	if FileObj.Data == "" || FileObj.Data == "null" {
		err = h.InvoiceService.UpdateDocument(Dto.UpdateDocumentRequest{
			InvoiceId:    invoiceId,
			DocumentType: "faktur",
			DocumentPath: FileObj.File.Filename,
		})

		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"message": err.Error(),
			})
			return
		} else {
			FileObj.Data = FileObj.File.Filename
		}

		invoice, err := h.InvoiceService.GetInvoiceById(invoiceId)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"message": err.Error(),
			})
			return
		}

		if invoice.InvoiceStatusId < 5 {
			err = h.InvoiceService.UpdateStatus(&Dto.UpdateStatusRequest{
				InvoiceId:       invoiceId,
				InvoiceStatusId: 5,
			})
		}
	}

	if err = h.StorageService.UploadFile(&Model.S3ObjectRequest{
		File:   FileObj.File,
		Bucket: "cahaya-teknik",
		Key:    "invoice-faktur/" + FileObj.Data,
	}); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    FileObj.Data,
	})
}

func (h *InvoiceController) GetFakturUrl(ctx *gin.Context) {
	var request *Dto.GetDocumentUrlRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	url, err := h.StorageService.GeneratePresignedURL(&Model.S3UrlRequest{
		Bucket: "cahaya-teknik",
		Key:    "invoice-faktur/" + request.Key,
	})

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    url,
	})
}

func (h *InvoiceController) GetPreviousSale(ctx *gin.Context) {
	var request Dto.GetPreviousSalesRequest
	fmt.Println("GetPreviousSale called")

	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	sales, err := h.InvoiceService.GetPreviousSale(request)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    sales,
	})
}
