package Repository

import (
	"ct-backend/Model"
	"ct-backend/Model/Dto"
	"ct-backend/Utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	IProductRepository interface {
		AddProduct(name string) error
		GetProductByName(name string) (product *Model.Product, err error)
		GetAllProduct(ctx *gin.Context, request *Dto.GetProductRequest) (products []*Model.Product, err error)
		EditNameProduct(id int, name string) (err error)
		SumStockProduct(id int, stock int, trx *gorm.DB) (err error)
		GetProductById(id int) (product *Model.Product, err error)
	}

	ProductRepository struct {
		DB *gorm.DB
	}
)

func ProductRepositoryProvider(DB *gorm.DB) *ProductRepository {
	return &ProductRepository{
		DB: DB,
	}
}

func (h *ProductRepository) AddProduct(name string) (err error) {
	product := &Model.Product{
		Name: name,
	}

	if err := h.DB.Create(&product).Error; err != nil {
		return err
	}

	return nil
}

func (h *ProductRepository) GetProductByName(name string) (product *Model.Product, err error) {
	if err := h.DB.
		Where("name = ?", name).
		First(&product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (h *ProductRepository) GetAllProduct(ctx *gin.Context, request *Dto.GetProductRequest) (products []*Model.Product, err error) {
	query := h.DB

	if request.Search != "" {
		query = query.Where("name LIKE ?", "%"+request.Search+"%")
	}

	if err := query.Scopes(Utils.Paginate(ctx)).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, err
}

func (h *ProductRepository) EditNameProduct(id int, name string) (err error) {
	if err := h.DB.
		Model(&Model.Product{}).
		Where("id = ?", id).
		Update("name", name).Error; err != nil {
		return err
	}

	return nil
}

func (h *ProductRepository) SumStockProduct(id int, stock int, trx *gorm.DB) (err error) {
	db := trx
	if db == nil {
		db = h.DB
	}

	if err := db.
		Debug().
		Model(&Model.Product{}).
		Where("id = ?", id).
		Update("stock", gorm.Expr("stock + ?", stock)).Error; err != nil {
		return err
	}

	return nil
}

func (h *ProductRepository) GetProductById(id int) (product *Model.Product, err error) {

	if err := h.DB.
		Where("id = ?", id).
		First(&product).Error; err != nil {
		return nil, err
	}

	return product, nil
}
