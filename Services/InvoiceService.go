package Services

import (
	"ct-backend/Model"
	"ct-backend/Model/Dto"
	"ct-backend/Repository"
	"ct-backend/Utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

type (
	IInvoiceService interface {
		AddInvoice(request *Dto.CreateInvoiceRequest) error
		GetAllInvoice(request *Dto.GetInvoicesRequest, ctx *gin.Context) ([]Model.ShortInvoice, error)
		GetInvoiceById(id int) (Model.Invoice, error)
		LockInvoice(request *Dto.IdRequest) error
		AddSaleToInvoice(request *Dto.AddSaleRequest) error
		UpdateSale(request *Dto.UpdateSaleRequest) error
		DeleteSale(request Dto.IdRequest) error
		GetAllSale(invoiceId int) ([]Model.Sale, error)
		UpdateFaktur(request *Dto.UpdateFakturRequest) error
		UpdateMainInformation(request *Dto.UpdateMainInformationRequest) error
		UpdateNote(request *Dto.UpdateNoteRequest) error
		UpdateStatus(request *Dto.UpdateStatusRequest) error
		DeleteInvoice(request Dto.IdRequest) error
		UpdateDocument(request Dto.UpdateDocumentRequest) error
		GetPreviousSale(request Dto.GetPreviousSalesRequest) (Model.Sale, error)
	}

	InvoiceService struct {
		InvoiceRepository  Repository.IInvoiceRepository
		ProductRepository  Repository.IProductRepository
		DeliveryRepository Repository.IDeliveryRepository
		DB                 *gorm.DB
	}
)

func InvoiceServiceProvider(
	invoiceRepository Repository.IInvoiceRepository,
	ProductRepository Repository.IProductRepository,
	DeliveryRepository Repository.IDeliveryRepository,
	DB *gorm.DB) *InvoiceService {
	return &InvoiceService{
		InvoiceRepository:  invoiceRepository,
		ProductRepository:  ProductRepository,
		DeliveryRepository: DeliveryRepository,
		DB:                 DB}
}

func (h *InvoiceService) AddInvoice(request *Dto.CreateInvoiceRequest) error {
	invoice, err := h.InvoiceRepository.GetLast(request.IsTaxable)

	log.Println("ini invoice terakhir", invoice)

	if err != nil {
		return err
	}

	request.InvoiceCode, err = createInvoiceCode(invoice, request.IsTaxable)
	if invoice != nil {
		request.Seller = invoice.Seller
		request.Platform = invoice.Platform
		request.PaymentMethod = invoice.PaymentMethod
		request.PlatformDescription = invoice.PlatformDescription
		request.PlatformNumber = invoice.PlatformNumber
	}

	if err != nil {
		return err
	}
	return h.InvoiceRepository.Create(request)
}

func (h *InvoiceService) GetAllInvoice(request *Dto.GetInvoicesRequest, ctx *gin.Context) ([]Model.ShortInvoice, error) {
	invoices, err := h.InvoiceRepository.GetAll(request, ctx)
	if err != nil {
		return nil, err
	}

	var shortInvoices []Model.ShortInvoice
	for _, invoice := range invoices {

		seenIDs := make(map[int]struct{})
		finalSales := make([]Model.Sale, 0)
		for _, data := range invoice.Sales {
			if data.NotSentCount > 0 {
				invoice.InvoiceStatusId = 3
			}

			if _, exists := seenIDs[data.ProductId]; exists {
				continue
			}

			seenIDs[data.ProductId] = struct{}{}

			finalSales = append(finalSales, data)
		}

		shortInvoices = append(shortInvoices, Model.ShortInvoice{
			ID:          invoice.ID,
			InvoiceCode: invoice.InvoiceCode,
			ClientName:  invoice.Client.Name,
			ProjectName: invoice.ProjectName,
			TotalItems:  len(finalSales),
			CreatedAt:   invoice.CreatedAt,
			Status:      invoice.GetStatusName(),
			StatusId:    invoice.InvoiceStatusId,
		})
	}

	if request.MinQuantity != "" && request.MaxQuantity != "" {
		minQuantity, err := strconv.Atoi(request.MinQuantity)
		if err != nil {
			return nil, err
		}

		maxQuantity, err := strconv.Atoi(request.MaxQuantity)
		if err != nil {
			return nil, err
		}

		if minQuantity >= 0 && maxQuantity >= 0 {
			var filteredInvoices []Model.ShortInvoice
			for _, invoice := range shortInvoices {
				if invoice.TotalItems >= minQuantity && invoice.TotalItems <= maxQuantity {
					filteredInvoices = append(filteredInvoices, invoice)
				}
			}
			shortInvoices = filteredInvoices
		}
	}

	return shortInvoices, nil
}

func (h *InvoiceService) GetInvoiceById(id int) (Model.Invoice, error) {
	var (
		invoice Model.Invoice
		err     error
	)

	if invoice, err = h.InvoiceRepository.GetById(id); err != nil {
		return Model.Invoice{}, err
	}

	if invoice.InvoiceStatusId == 4 {
		for _, sale := range invoice.Sales {
			if sale.NotSentCount > 0 {
				invoice.InvoiceStatusId = 3
				break
			}
		}
	}

	return invoice, nil
}

func (h *InvoiceService) LockInvoice(request *Dto.IdRequest) error {
	err := h.InvoiceRepository.UpdateInvoiceTotalPrice(request.Id)
	if err != nil {
		return err
	}

	return h.InvoiceRepository.UpdateStatus(&Dto.UpdateStatusRequest{InvoiceId: request.Id, InvoiceStatusId: 3})
}

func (h *InvoiceService) AddSaleToInvoice(request *Dto.AddSaleRequest) error {
	err := h.ProductRepository.SumStockProduct(request.ProductId, request.Count*-1, nil)
	if err != nil {
		return err
	}

	return h.InvoiceRepository.AddSale(request)
}

func (h *InvoiceService) UpdateSale(request *Dto.UpdateSaleRequest) error {
	sale, err := h.InvoiceRepository.GetSale(request.Id)
	if err != nil {
		return err
	}

	request.NotSentCount = 0
	if sale.NotSentCount > 0 {
		count := request.Count - request.CurrentCount
		request.NotSentCount = sale.NotSentCount + count
		err := h.ProductRepository.SumStockProduct(request.ProductId, count*-1, nil)
		if err != nil {
			return err
		}
	}

	// if product change than update both their stock
	if sale.ProductId != request.ProductId {
		if err := h.ProductRepository.SumStockProduct(sale.ProductId, request.Count, nil); err != nil {
			return err
		}

		if err := h.ProductRepository.SumStockProduct(request.ProductId, request.Count*-1, nil); err != nil {
			return err
		}
	}

	return h.InvoiceRepository.UpdateSale(request)
}

func (h *InvoiceService) DeleteSale(request Dto.IdRequest) error {
	var (
		err error
	)

	trx := h.DB.Begin()

	if trx.Error != nil {
		return trx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			trx.Rollback()
			err = fmt.Errorf("panic occurred: %v", r)
		} else if err != nil {
			trx.Rollback()
		} else {
			trx.Commit()
		}
	}()

	invoiceStatus := 3

	tempSale, err := h.InvoiceRepository.GetSale(request.Id)
	if err != nil {
		return err
	}

	err = h.InvoiceRepository.DeleteSale(request)
	if err != nil {
		return err
	}

	sales, err := h.InvoiceRepository.GetSalesByInvoiceId(tempSale.InvoiceId)
	if err != nil {
		return err
	}

	invoice, err := h.InvoiceRepository.GetById(tempSale.InvoiceId)
	if err != nil {
		return err
	}

	notSentEmpty := true
	for _, sale := range sales {
		if sale.NotSentCount > 0 {
			notSentEmpty = false
			break
		}
	}

	if notSentEmpty {
		invoiceStatus = 4
		if !invoice.IsTaxable {
			invoiceStatus++
		}
	}

	if err = h.InvoiceRepository.UpdateStatus(&Dto.UpdateStatusRequest{InvoiceId: tempSale.InvoiceId, InvoiceStatusId: invoiceStatus}); err != nil {
		return err
	}

	return nil
}

func (h *InvoiceService) UpdateFaktur(request *Dto.UpdateFakturRequest) error {
	return h.InvoiceRepository.UpdateFaktur(request)
}

func (h *InvoiceService) UpdateMainInformation(request *Dto.UpdateMainInformationRequest) error {
	return h.InvoiceRepository.UpdateMainInformation(request)
}

func (h *InvoiceService) UpdateNote(request *Dto.UpdateNoteRequest) error {
	return h.InvoiceRepository.UpdateNote(request)
}

func (h *InvoiceService) UpdateStatus(request *Dto.UpdateStatusRequest) error {
	return h.InvoiceRepository.UpdateStatus(request)
}

func (h *InvoiceService) DeleteInvoice(request Dto.IdRequest) (err error) {
	trx := h.DB.Begin()
	if trx.Error != nil {
		return trx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			trx.Rollback()
			err = fmt.Errorf("panic occurred: %v", r)
		} else if err != nil {
			trx.Rollback()
		} else {
			trx.Commit()
		}
	}()

	invoice, err := h.InvoiceRepository.GetById(request.Id)
	if err != nil {
		return err
	}
	if invoice.InvoiceStatusId > 5 {
		return fmt.Errorf("cannot delete invoice with status %d", invoice.InvoiceStatusId)
	}

	products, err := h.InvoiceRepository.GetAllSale(request.Id)
	if err != nil {
		return err
	}

	for _, product := range products {
		if e := h.ProductRepository.SumStockProduct(product.ProductId, product.Quantity, trx); e != nil {
			return e
		}
	}

	if err := h.InvoiceRepository.Delete(request, trx); err != nil {
		return err
	}

	if err := h.DeliveryRepository.DeleteDelivery(&request, trx); err != nil {
		return err
	}

	return nil
}

func (h *InvoiceService) GetAllSale(invoiceId int) ([]Model.Sale, error) {
	return h.InvoiceRepository.GetAllSale(invoiceId)
}

func createInvoiceCode(invoice *Model.Invoice, isTaxable bool) (val string, err error) {
	month := Utils.MonthToRoman(int(time.Now().Month()))
	year := time.Now().Year()
	order := 1
	companyCode := "CCT"

	if !isTaxable {
		companyCode = "SAM"
	}

	if invoice != nil {
		if invoice.CreatedAt.Year() == year {
			parts := strings.Split(invoice.InvoiceCode, "/")
			order, err = strconv.Atoi(parts[0])

			if err != nil {
				fmt.Println("Error converting to integer:", err)
				return "", err
			}

			order++
		}
	}

	return fmt.Sprintf("%d/%s/%s/SBY/%d", order, month, companyCode, year-2000), nil
}

func (h *InvoiceService) UpdateDocument(request Dto.UpdateDocumentRequest) error {
	return h.InvoiceRepository.UpdateDocument(request)
}

func (h *InvoiceService) GetPreviousSale(request Dto.GetPreviousSalesRequest) (Model.Sale, error) {
	return h.InvoiceRepository.GetPreviousSale(request.ProductID, request.ClientID)
}
