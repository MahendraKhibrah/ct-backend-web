package Services

import (
	"ct-backend/Model"
	"ct-backend/Model/Dto"
	"ct-backend/Repository"
	"github.com/gin-gonic/gin"
	"time"
)

type (
	IReceiptService interface {
		GetAllReceipt(ctx *gin.Context) (receipts []Model.Receipt, err error)
		GetReceiptById(id int) (receipt *Model.Receipt, err error)
		CreateReceipt(clientId int) (m *Model.Receipt, err error)
		GetInvoiceByReceiptId(id int) ([]Model.ReceiptInvoice, error)
		AddInvoiceReceipt(receiptInvoice *Dto.ReceiptInvoiceRequest) (*Model.ReceiptInvoice, error)
		LockReceipt(receiptId int) error
		DeleteReceiptInvoice(id int) error
		GetAvailableInvoices(clientId int) ([]Model.Invoice, error)
		GetClientReceipts() ([]Model.Client, error)
		PayReceipt(receiptId int) error
	}

	ReceiptService struct {
		ReceiptRepository Repository.IReceiptRepository
		InvoiceRepository Repository.IInvoiceRepository
	}
)

func ReceiptServiceProvider(receiptRepository Repository.IReceiptRepository, invoiceRepository Repository.IInvoiceRepository) *ReceiptService {
	return &ReceiptService{
		ReceiptRepository: receiptRepository,
		InvoiceRepository: invoiceRepository,
	}
}

func (h *ReceiptService) GetAllReceipt(ctx *gin.Context) (receipts []Model.Receipt, err error) {
	return h.ReceiptRepository.GetReceipts(ctx)
}

func (h *ReceiptService) GetReceiptById(id int) (receipt *Model.Receipt, err error) {
	return h.ReceiptRepository.GetReceiptById(id)
}

func (h *ReceiptService) CreateReceipt(clientId int) (m *Model.Receipt, err error) {
	lastReceipt, _ := h.ReceiptRepository.GetLast()

	receipt := &Model.Receipt{
		ClientId: clientId,
		Status:   1,
	}

	receipt.Number = 1

	if lastReceipt != nil && lastReceipt.CreatedAt.Year() == time.Now().Year() {
		receipt.Number = lastReceipt.Number + 1
	}

	return h.ReceiptRepository.CreateReceipt(receipt)
}

func (h *ReceiptService) GetInvoiceByReceiptId(id int) ([]Model.ReceiptInvoice, error) {
	return h.ReceiptRepository.GetInvoiceByReceiptId(id)
}

func (h *ReceiptService) AddInvoiceReceipt(receiptInvoice *Dto.ReceiptInvoiceRequest) (*Model.ReceiptInvoice, error) {
	err := h.InvoiceRepository.UpdateInvoiceTotalPrice(receiptInvoice.InvoiceId)
	if err != nil {
		return nil, err
	}

	return h.ReceiptRepository.AddInvoiceReceipt(receiptInvoice)
}

func (h *ReceiptService) LockReceipt(receiptId int) error {
	return h.ReceiptRepository.LockReceipt(receiptId)
}

func (h *ReceiptService) DeleteReceiptInvoice(id int) error {
	return h.ReceiptRepository.DeleteInvoiceReceipt(id)
}

func (h *ReceiptService) GetAvailableInvoices(clientId int) ([]Model.Invoice, error) {
	allInvoices, err := h.InvoiceRepository.GetAllForReceipt()
	if err != nil {
		return nil, err
	}

	var filteredInvoices []Model.Invoice
	for _, invoice := range allInvoices {
		if invoice.Client.ID == clientId {
			filteredInvoices = append(filteredInvoices, invoice)
		}
	}

	return filteredInvoices, nil
}

func (h *ReceiptService) GetClientReceipts() ([]Model.Client, error) {
	invoices, err := h.InvoiceRepository.GetAllForReceipt()

	if err != nil {
		return nil, err
	}

	var uniqueClients []Model.Client
	clientMap := make(map[int]bool)

	for _, invoice := range invoices {
		if _, exists := clientMap[invoice.Client.ID]; !exists {
			clientMap[invoice.Client.ID] = true
			uniqueClients = append(uniqueClients, invoice.Client)
		}
	}
	return uniqueClients, nil
}

func (h *ReceiptService) PayReceipt(receiptId int) error {
	return h.ReceiptRepository.PayReceipt(receiptId)
}
