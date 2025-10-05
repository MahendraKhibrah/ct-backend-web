package Repository

import (
	"ct-backend/Model"
	"ct-backend/Model/Dto"
	"ct-backend/Utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	IReceiptRepository interface {
		GetReceipts(ctx *gin.Context) ([]Model.Receipt, error)
		GetReceiptById(id int) (*Model.Receipt, error)
		CreateReceipt(receipt *Model.Receipt) (*Model.Receipt, error)
		GetLast() (*Model.Receipt, error)
		GetInvoiceByReceiptId(id int) ([]Model.ReceiptInvoice, error)
		AddInvoiceReceipt(receiptInvoice *Dto.ReceiptInvoiceRequest) (*Model.ReceiptInvoice, error)
		DeleteInvoiceReceipt(id int) error
		LockReceipt(receiptId int) error
		PayReceipt(receiptId int) error
	}

	ReceiptRepository struct {
		DB *gorm.DB
	}
)

func ReceiptRepositoryProvider(DB *gorm.DB) *ReceiptRepository {
	return &ReceiptRepository{
		DB: DB,
	}
}

func (h *ReceiptRepository) GetReceipts(ctx *gin.Context) ([]Model.Receipt, error) {
	var receipts []Model.Receipt
	if err := h.DB.Preload("Client").Order("created_at desc").Scopes(Utils.Paginate(ctx)).Find(&receipts).Error; err != nil {
		return nil, err
	}

	return receipts, nil
}

func (h *ReceiptRepository) GetReceiptById(id int) (*Model.Receipt, error) {
	var receipt *Model.Receipt
	if err := h.DB.Preload("Client").Where("id = ?", id).First(&receipt).Error; err != nil {
		return receipt, err
	}

	return receipt, nil
}

func (h *ReceiptRepository) CreateReceipt(receipt *Model.Receipt) (*Model.Receipt, error) {
	if err := h.DB.Create(&receipt).Error; err != nil {
		return receipt, err
	}

	return receipt, nil
}

func (h *ReceiptRepository) GetLast() (*Model.Receipt, error) {
	var receipt *Model.Receipt
	if err := h.DB.Last(&receipt).Error; err != nil {
		return receipt, err
	}

	return receipt, nil
}

func (h *ReceiptRepository) GetInvoiceByReceiptId(id int) ([]Model.ReceiptInvoice, error) {
	var receiptInvoices []Model.ReceiptInvoice
	if err := h.DB.Preload("Invoice").Where("receipt_id = ?", id).Find(&receiptInvoices).Error; err != nil {
		return nil, err
	}

	return receiptInvoices, nil
}

func (h *ReceiptRepository) AddInvoiceReceipt(receiptInvoice *Dto.ReceiptInvoiceRequest) (*Model.ReceiptInvoice, error) {
	receiptInvoiceModel := &Model.ReceiptInvoice{
		ReceiptId: receiptInvoice.ReceiptId,
		InvoiceId: receiptInvoice.InvoiceId,
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&receiptInvoiceModel).Error; err != nil {
			return err
		}

		if err := tx.Model(&Model.Invoice{}).Where("id = ?", receiptInvoice.InvoiceId).Update("invoice_status_id", 6).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return receiptInvoiceModel, nil
}

func (h *ReceiptRepository) DeleteInvoiceReceipt(id int) error {
	var receiptInvoices *Model.ReceiptInvoice

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&receiptInvoices).Error; err != nil {
			return err
		}

		if err := tx.Delete(&receiptInvoices).Error; err != nil {
			return err
		}

		if err := tx.Model(&Model.Invoice{}).Where("id = ?", receiptInvoices.InvoiceId).Update("invoice_status_id", 5).Error; err != nil {
			return err
		}

		if err := tx.Model(&Model.Receipt{}).Where("id = ?", receiptInvoices.ReceiptId).Update("status", 1).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (h *ReceiptRepository) LockReceipt(receiptId int) error {
	if err := h.DB.Model(&Model.Receipt{}).Where("id = ?", receiptId).Update("status", 2).Error; err != nil {
		return err
	}

	return nil
}

func (h *ReceiptRepository) PayReceipt(receiptId int) error {
	if err := h.DB.Model(&Model.Receipt{}).Where("id = ?", receiptId).Update("status", 3).Error; err != nil {
		return err
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var receiptInvoices []Model.ReceiptInvoice
		if err := tx.Preload("Invoice").Where("receipt_id = ?", receiptId).Find(&receiptInvoices).Error; err != nil {
			return err
		}

		for _, receiptInvoice := range receiptInvoices {
			if err := tx.Model(&Model.Invoice{}).Where("id = ?", receiptInvoice.InvoiceId).Update("invoice_status_id", 8).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(&Model.Receipt{}).Where("id = ?", receiptId).Update("status", 3).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
