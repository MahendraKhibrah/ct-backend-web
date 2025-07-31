package Controller

import (
	"ct-backend/Model/Common"
	"ct-backend/Model/Dto"
	"ct-backend/Services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	IProductController interface {
		AddProduct(ctx *gin.Context)
		GetAllProduct(ctx *gin.Context)
		EditNameProduct(ctx *gin.Context)
	}

	ProductController struct {
		ProductService Services.IProductService
	}
)

func ProductControllerProvider(productService Services.IProductService) *ProductController {
	return &ProductController{
		ProductService: productService,
	}
}

func (h *ProductController) AddProduct(ctx *gin.Context) {
	var request *Dto.AddProductRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.ProductService.AddProduct(request.Name); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *ProductController) GetAllProduct(ctx *gin.Context) {
	var request *Dto.GetProductRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	products, err := h.ProductService.GetAllProduct(ctx, request)
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
		"data":    products,
		"meta":    pagination,
	})
}

func (h *ProductController) EditNameProduct(ctx *gin.Context) {
	var request *Dto.EditProductRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.ProductService.EditNameProduct(request.Id, request.Name); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
