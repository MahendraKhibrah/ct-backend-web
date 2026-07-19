package Repository

import (
	"ct-backend/Model"
	"time"

	"gorm.io/gorm"
)

type (
	IDashboardRepository interface {
		GetInvoicesSince(since time.Time, statusIDs []int) ([]Model.Invoice, error)
		GetDeliveriesSince(since time.Time, statusID int) ([]Model.DeliveryOrder, error)
		GetReceiptsSince(since time.Time) ([]Model.Receipt, error)
	}

	DashboardRepository struct {
		DB *gorm.DB
	}
)

func DashboardRepositoryProvider(DB *gorm.DB) *DashboardRepository {
	return &DashboardRepository{DB: DB}
}

func (h *DashboardRepository) GetInvoicesSince(since time.Time, statusIDs []int) (invoices []Model.Invoice, err error) {
	err = h.DB.
		Preload("Client").
		Where("created_at >= ?", since).
		Where("invoice_status_id IN ?", statusIDs).
		Order("created_at DESC, id DESC").
		Find(&invoices).Error

	return invoices, err
}

func (h *DashboardRepository) GetDeliveriesSince(since time.Time, statusID int) (deliveries []Model.DeliveryOrder, err error) {
	err = h.DB.
		Preload("Invoice.Client").
		Where("created_at >= ?", since).
		Where("status = ?", statusID).
		Order("created_at DESC, id DESC").
		Find(&deliveries).Error

	return deliveries, err
}

func (h *DashboardRepository) GetReceiptsSince(since time.Time) (receipts []Model.Receipt, err error) {
	err = h.DB.
		Preload("Client").
		Where("created_at >= ?", since).
		Order("created_at DESC, id DESC").
		Find(&receipts).Error

	return receipts, err
}
