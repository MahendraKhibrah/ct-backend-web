package Repository

import (
	"ct-backend/Model"
	"ct-backend/Model/Dto"
	"ct-backend/Utils"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
)

type (
	IInvoiceRepository interface {
		GetAll(request *Dto.GetInvoicesRequest, ctx *gin.Context) ([]Model.Invoice, error)
		GetById(id int) (Model.Invoice, error)
		GetLast(isTaxable bool) (*Model.Invoice, error)
		Create(request *Dto.CreateInvoiceRequest) (err error)
		Delete(request Dto.IdRequest, trx *gorm.DB) (err error)
		UpdateDocument(request Dto.UpdateDocumentRequest) (err error)
		AddSale(request *Dto.AddSaleRequest) (err error)
		GetAllSale(invoiceId int) ([]Model.Sale, error)
		UpdateSale(request *Dto.UpdateSaleRequest) (err error)
		DeleteSale(request Dto.IdRequest) (err error)
		UpdateFaktur(request *Dto.UpdateFakturRequest) (err error)
		UpdateMainInformation(request *Dto.UpdateMainInformationRequest) (err error)
		UpdateNote(request *Dto.UpdateNoteRequest) (err error)
		UpdateStatus(request *Dto.UpdateStatusRequest) (err error)
		UpdateNotSentSale(request *Dto.UpdateNotSentSaleRequest) (err error)
		GetSalesByInvoiceId(invoiceId int) ([]Model.Sale, error)
		GetAllForDelivery() ([]Model.Invoice, error)
		GetAllForReceipt() ([]Model.Invoice, error)
		UpdateInvoiceTotalPrice(invoiceID int) error
		GetSale(saleId int) (Model.Sale, error)
		GetPreviousSale(productID string, clientID string) (Model.Sale, error)
	}

	InvoiceRepository struct {
		DB *gorm.DB
	}
)

func InvoiceRepositoryProvider(DB *gorm.DB) *InvoiceRepository {
	return &InvoiceRepository{
		DB: DB,
	}
}

func (h *InvoiceRepository) GetAll(request *Dto.GetInvoicesRequest, ctx *gin.Context) (invoices []Model.Invoice, err error) {
	query := h.DB.Debug().Joins("JOIN clients ON clients.id = invoices.client_id").Preload("Client").Preload("Sales")

	if request.Search != "" {
		query = query.Where("clients.name ILIKE ? OR invoices.project_name ILIKE ?", "%"+request.Search+"%", "%"+request.Search+"%")
	}

	if request.Company != "" {
		query = query.Where("invoice_code ILIKE ?", "%"+request.Company+"%")
	}

	if err := query.Order("created_at DESC").Scopes(Utils.Paginate(ctx)).Find(&invoices).Error; err != nil {
		return nil, err
	}

	return invoices, nil
}

func (h *InvoiceRepository) GetLast(isTaxable bool) (invoice *Model.Invoice, err error) {
	if err := h.DB.Where("is_taxable = ?", isTaxable).Last(&invoice).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return invoice, nil
}

func (h *InvoiceRepository) GetById(id int) (invoice Model.Invoice, err error) {
	if err := h.DB.Preload("Client").Preload("Sales").Where("id = ?", id).First(&invoice).Error; err != nil {
		return Model.Invoice{}, err
	}

	return invoice, nil
}

func (h *InvoiceRepository) Create(request *Dto.CreateInvoiceRequest) (err error) {
	invoice := &Model.Invoice{
		ClientId:            request.ClientId,
		InvoiceCode:         request.InvoiceCode,
		Seller:              request.Seller,
		Platform:            request.Platform,
		PaymentMethod:       request.PaymentMethod,
		PlatformDescription: request.PlatformDescription,
		PlatformNumber:      request.PlatformNumber,
		IsTaxable:           request.IsTaxable,
		InvoiceStatusId:     1,
	}

	if err := h.DB.Create(&invoice).Error; err != nil {
		return err
	}

	return nil
}

func (h *InvoiceRepository) Delete(request Dto.IdRequest, trx *gorm.DB) (err error) {
	db := trx
	if db == nil {
		db = h.DB
	}

	if err := db.Where("id = ?", request.Id).Delete(&Model.Invoice{}).Error; err != nil {
		return err
	}

	if err := db.Where("invoice_id = ?", request.Id).Delete(&Model.Sale{}).Error; err != nil {
		return err
	}

	return nil
}

func (h *InvoiceRepository) UpdateDocument(request Dto.UpdateDocumentRequest) (err error) {
	documentType := "po_path"

	if request.DocumentType == "faktur" {
		documentType = "facture_path"
	}

	if err := h.DB.
		Model(&Model.Invoice{}).
		Where("id = ?", request.InvoiceId).
		Update(documentType, request.DocumentPath).Error; err != nil {
		return err
	}

	return nil
}

func (h *InvoiceRepository) AddSale(request *Dto.AddSaleRequest) (err error) {
	sale := &Model.Sale{
		InvoiceId:    request.InvoiceId,
		ProductId:    request.ProductId,
		Quantity:     request.Count,
		Price:        request.Price,
		NotSentCount: request.Count,
		Unit:         request.Unit,
		SendStatus:   false,
	}

	if err := h.DB.Create(&sale).Error; err != nil {
		return err
	}

	return nil
}

func (h *InvoiceRepository) UpdateSale(request *Dto.UpdateSaleRequest) (err error) {
	if err := h.DB.
		Model(&Model.Sale{}).
		Where("id = ?", request.Id).
		Update("quantity", request.Count).
		Update("not_sent_count", request.NotSentCount).
		Update("unit", request.Unit).
		Update("product_id", request.ProductId).
		Update("price", request.Price).Error; err != nil {
		return err
	}

	return nil
}

func (h *InvoiceRepository) DeleteSale(request Dto.IdRequest) (err error) {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var sale Model.Sale
	if err := tx.First(&sale, request.Id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&Model.Product{}).Where("id = ?", sale.ProductId).Update("stock", gorm.Expr("stock + ?", sale.Quantity)).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&sale).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (h *InvoiceRepository) UpdateFaktur(request *Dto.UpdateFakturRequest) (err error) {
	if err := h.DB.
		Model(&Model.Invoice{}).
		Where("id = ?", request.InvoiceId).
		Update("discount", request.Discount).
		Update("payment_term", request.PaymentTerm).Error; err != nil {
		return err
	}

	return nil
}

func (h *InvoiceRepository) UpdateMainInformation(request *Dto.UpdateMainInformationRequest) (err error) {
	if err := h.DB.
		Model(&Model.Invoice{}).
		Where("id = ?", request.InvoiceId).
		Update("note", request.Note).
		Update("po_code", request.PoCode).
		Update("seller", request.Seller).
		Update("platform", request.Platform).
		Update("payment_method", request.PaymentMethod).
		Update("platform_description", request.PlatformDescription).
		Update("project_name", request.Project).
		Update("date", request.Date).
		Update("client_id", request.ClientID).
		Update("platform_number", request.PlatformNumber).Error; err != nil {
		return err
	}

	return nil
}

func (h *InvoiceRepository) UpdateNote(request *Dto.UpdateNoteRequest) (err error) {
	if err := h.DB.
		Model(&Model.Invoice{}).
		Where("id = ?", request.InvoiceId).
		Update("note", request.Note).Error; err != nil {
		return err
	}

	return nil
}

func (h *InvoiceRepository) UpdateStatus(request *Dto.UpdateStatusRequest) (err error) {
	if err := h.DB.
		Model(&Model.Invoice{}).
		Where("id = ?", request.InvoiceId).
		Update("invoice_status_id", request.InvoiceStatusId).Error; err != nil {
		return err
	}

	return nil
}

func (h *InvoiceRepository) GetAllSale(invoiceId int) (sales []Model.Sale, err error) {
	if err := h.DB.Preload("Product").Where("invoice_id = ?", invoiceId).Find(&sales).Error; err != nil {
		return nil, err
	}

	return sales, nil
}

func (h *InvoiceRepository) UpdateNotSentSale(request *Dto.UpdateNotSentSaleRequest) (err error) {
	if err := h.DB.
		Model(&Model.Sale{}).
		Where("id = ?", request.SaleId).
		Update("not_sent_count", gorm.Expr("not_sent_count - ?", request.Count)).Error; err != nil {
		return err
	}

	return nil
}

func (h *InvoiceRepository) GetSalesByInvoiceId(invoiceId int) (sales []Model.Sale, err error) {
	if err := h.DB.Preload("Product").Where("invoice_id = ?", invoiceId).Find(&sales).Error; err != nil {
		return nil, err
	}

	return sales, nil
}

func (h *InvoiceRepository) GetAllForDelivery() (invoices []Model.Invoice, err error) {
	if err := h.DB.Preload("Client").
		Joins("JOIN sales ON sales.invoice_id = invoices.id").
		Where("sales.not_sent_count > 0").
		Find(&invoices).Error; err != nil {
		return nil, err
	}

	return invoices, nil
}

func (h *InvoiceRepository) GetAllForReceipt() (invoices []Model.Invoice, err error) {
	if err := h.DB.Preload("Client").Where("invoice_status_id > ? AND invoice_status_id < ?", 3, 6).Find(&invoices).Error; err != nil {
		return nil, err
	}

	return invoices, nil
}

func (h *InvoiceRepository) GetSale(saleId int) (sale Model.Sale, err error) {
	if err := h.DB.Preload("Product").Where("id = ?", saleId).First(&sale).Error; err != nil {
		return Model.Sale{}, err
	}

	return sale, nil
}

func (h *InvoiceRepository) UpdateInvoiceTotalPrice(invoiceID int) error {
	var total float64
	var invoice Model.Invoice

	err := h.DB.Model(&Model.Sale{}).
		Where("invoice_id = ?", invoiceID).
		Select("SUM(price * quantity) as total").
		Scan(&total).Error
	if err != nil {
		return err
	}

	invoice, err = h.GetById(invoiceID)
	if err != nil {
		return err
	}

	if invoice.IsTaxable {
		total = math.Round(total * 1.11)
	}

	err = h.DB.Model(&Model.Invoice{}).
		Where("id = ?", invoiceID).
		Update("total_price", total).Error
	if err != nil {
		return err
	}

	return nil
}

func (h *InvoiceRepository) GetPreviousSale(productID string, clientID string) (sale Model.Sale, err error) {
	if err := h.DB.
		Joins("JOIN invoices ON invoices.id = sales.invoice_id").
		Where("invoices.client_id = ?", clientID).
		Where("sales.product_id = ?", productID).
		Where("sales.price > ?", 1).
		Order("sales.created_at DESC").
		First(&sale).Error; err != nil {
		return Model.Sale{}, err
	}

	return sale, nil
}
