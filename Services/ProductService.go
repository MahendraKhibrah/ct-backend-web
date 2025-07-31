package Services

import (
	"ct-backend/Model"
	"ct-backend/Model/Dto"
	"ct-backend/Repository"
	"errors"
	"github.com/gin-gonic/gin"
)

type (
	IProductService interface {
		AddProduct(name string) error
		GetAllProduct(ctx *gin.Context, request *Dto.GetProductRequest) (products []*Model.Product, err error)
		EditNameProduct(id int, name string) (err error)
	}

	ProductService struct {
		ProductRepository Repository.IProductRepository
	}
)

func ProductServiceProvider(productRepository Repository.IProductRepository) *ProductService {
	return &ProductService{
		ProductRepository: productRepository,
	}
}

func (h *ProductService) AddProduct(name string) error {
	if product, _ := h.ProductRepository.GetProductByName(name); product != nil {
		return errors.New("Produk ini sudah ada")
	}
	return h.ProductRepository.AddProduct(name)
}

func (h *ProductService) GetAllProduct(ctx *gin.Context, request *Dto.GetProductRequest) (products []*Model.Product, err error) {
	return h.ProductRepository.GetAllProduct(ctx, request)
}

func (h *ProductService) EditNameProduct(id int, name string) (err error) {
	if _, err := h.ProductRepository.GetProductById(id); err != nil {
		return err
	}
	return h.ProductRepository.EditNameProduct(id, name)
}

// implement delete product
