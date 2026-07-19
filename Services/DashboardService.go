package Services

import (
	"ct-backend/Model/Dto"
	"ct-backend/Repository"
	"time"
)

const (
	dashboardInvoicePeriodDays   = 14
	dashboardDeliveryPeriodDays  = 14
	dashboardReceiptPeriodMonths = 1
)

var dashboardInvoiceStatusIDs = []int{1, 3, 4}

type (
	IDashboardService interface {
		GetWorkQueue() (*Dto.DashboardWorkQueue, error)
	}

	DashboardService struct {
		DashboardRepository Repository.IDashboardRepository
	}
)

func DashboardServiceProvider(dashboardRepository Repository.IDashboardRepository) *DashboardService {
	return &DashboardService{DashboardRepository: dashboardRepository}
}

func (h *DashboardService) GetWorkQueue() (*Dto.DashboardWorkQueue, error) {
	now := time.Now()
	invoiceSince := now.AddDate(0, 0, -dashboardInvoicePeriodDays)
	deliverySince := now.AddDate(0, 0, -dashboardDeliveryPeriodDays)
	receiptSince := now.AddDate(0, -dashboardReceiptPeriodMonths, 0)

	invoices, err := h.DashboardRepository.GetInvoicesSince(invoiceSince, dashboardInvoiceStatusIDs)
	if err != nil {
		return nil, err
	}

	deliveries, err := h.DashboardRepository.GetDeliveriesSince(deliverySince, 1)
	if err != nil {
		return nil, err
	}

	receipts, err := h.DashboardRepository.GetReceiptsSince(receiptSince)
	if err != nil {
		return nil, err
	}

	workQueue := &Dto.DashboardWorkQueue{
		Invoices:   make([]Dto.DashboardInvoice, 0, len(invoices)),
		Deliveries: make([]Dto.DashboardDelivery, 0, len(deliveries)),
		Receipts:   make([]Dto.DashboardReceipt, 0, len(receipts)),
	}

	for _, invoice := range invoices {
		workQueue.Invoices = append(workQueue.Invoices, Dto.DashboardInvoice{
			ID:          invoice.ID,
			InvoiceCode: invoice.InvoiceCode,
			ClientName:  invoice.Client.Name,
			ProjectName: invoice.ProjectName,
			StatusID:    invoice.InvoiceStatusId,
			CreatedAt:   invoice.CreatedAt,
		})
	}

	for _, delivery := range deliveries {
		workQueue.Deliveries = append(workQueue.Deliveries, Dto.DashboardDelivery{
			ID:         delivery.ID,
			OrderCode:  delivery.OrderCode,
			ClientName: delivery.Invoice.Client.Name,
			StatusID:   delivery.Status,
			CreatedAt:  delivery.CreatedAt,
		})
	}

	for _, receipt := range receipts {
		workQueue.Receipts = append(workQueue.Receipts, Dto.DashboardReceipt{
			ID:         receipt.ID,
			Number:     receipt.Number,
			ClientName: receipt.Client.Name,
			StatusID:   receipt.Status,
			CreatedAt:  receipt.CreatedAt,
		})
	}

	return workQueue, nil
}
