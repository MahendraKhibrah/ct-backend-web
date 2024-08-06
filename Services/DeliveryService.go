package Services

import (
	"ct-backend/Model"
	"ct-backend/Model/Dto"
	"ct-backend/Repository"
	"strconv"
)

type (
	IDeliveryService interface {
		GetById(id int) (delivery *Model.DeliveryOrder, err error)
		GetAll() (deliveries []Model.ShortDeliveryOrder, err error)
		CreateDeliveryOrder(request *Dto.IdRequest) (err error)
		CreateDeliveryProduct(request *Dto.CreateDeliveryProductRequest) (err error)
		UpdateDeliveryProduct(request *Dto.UpdateDeliveryProductRequest) (err error)
		DeleteDeliveryProduct(request *Dto.DeleteDeliveryProductRequest) (err error)
		GetAllDeliveryProduct(request *Dto.IdRequest) (deliveryProducts []Model.ShortDeliveryProduct, err error)
		UpdateSender(request *Dto.UpdateSenderRequest) (err error)
		UpdateMainInformation(request *Dto.UpdateDeliveryInformationRequest) (err error)
		GetPreviousNote(id int) (note string, err error)
		LockDeliveryOrder(request *Dto.LockDeliveryOrderRequest) (err error)
	}

	DeliveryService struct {
		Repo        Repository.IDeliveryRepository
		InvoiceRepo Repository.IInvoiceRepository
	}
)

func DeliveryServiceProvider(repo Repository.IDeliveryRepository, invoiceRepo Repository.IInvoiceRepository) *DeliveryService {
	return &DeliveryService{
		Repo:        repo,
		InvoiceRepo: invoiceRepo,
	}
}

func (h *DeliveryService) GetById(id int) (delivery *Model.DeliveryOrder, err error) {
	return h.Repo.GetById(id)
}

func (h *DeliveryService) GetAll() (deliveries []Model.ShortDeliveryOrder, err error) {
	rawData, err := h.Repo.GetAll()
	if err != nil {
		return nil, err
	}

	for _, data := range rawData {
		deliveries = append(deliveries, Model.ShortDeliveryOrder{
			ID:         data.ID,
			OrderCode:  data.OrderCode,
			ClientName: data.Invoice.Client.Name,
			CreatedAt:  data.CreatedAt,
			Status:     data.GetStatusName(),
		})
	}

	return deliveries, nil
}

func (h *DeliveryService) CreateDeliveryOrder(request *Dto.IdRequest) (err error) {
	invoice, err := h.InvoiceRepo.GetById(request.Id)
	if err != nil {
		return err
	}

	deliveries, err := h.Repo.GetAllByInvoiceId(request.Id)
	if err != nil {
		return err
	}

	orderCode := generateOrderCode(invoice.InvoiceCode, len(deliveries))

	return h.Repo.CreateDeliveryOrder(request, orderCode)
}

func (h *DeliveryService) CreateDeliveryProduct(request *Dto.CreateDeliveryProductRequest) (err error) {
	if err = h.Repo.CreateDeliveryProduct(request); err != nil {
		return err
	}

	if err = h.InvoiceRepo.UpdateNotSentSale(&Dto.UpdateNotSentSaleRequest{SaleId: request.SalesID, Count: request.Quantity}); err != nil {
		return err
	}

	return nil
}

func (h *DeliveryService) UpdateDeliveryProduct(request *Dto.UpdateDeliveryProductRequest) (err error) {
	if err = h.Repo.UpdateDeliveryProduct(request); err != nil {
		return err
	}

	if err = h.InvoiceRepo.UpdateNotSentSale(&Dto.UpdateNotSentSaleRequest{SaleId: request.ID, Count: request.Quantity - request.CurrentQuantity}); err != nil {
		return err
	}

	return nil
}

func (h *DeliveryService) DeleteDeliveryProduct(request *Dto.DeleteDeliveryProductRequest) (err error) {
	if err = h.Repo.DeleteDeliveryProduct(request); err != nil {
		return err
	}

	if err = h.InvoiceRepo.UpdateNotSentSale(&Dto.UpdateNotSentSaleRequest{SaleId: request.SaleId, Count: request.Quantity * -1}); err != nil {
		return err
	}

	return nil
}

func (h *DeliveryService) GetAllDeliveryProduct(request *Dto.IdRequest) (deliveryProducts []Model.ShortDeliveryProduct, err error) {
	rawData, err := h.Repo.GetAllDeliveryProduct(request.Id)
	if err != nil {
		return nil, err
	}

	for _, data := range rawData {
		deliveryProducts = append(deliveryProducts, Model.ShortDeliveryProduct{
			ID:       data.ID,
			Name:     data.Sale.Product.Name,
			Quantity: data.Quantity,
		})
	}

	return deliveryProducts, nil
}

func (h *DeliveryService) UpdateSender(request *Dto.UpdateSenderRequest) (err error) {
	return h.Repo.UpdateSender(request)
}

func (h *DeliveryService) UpdateMainInformation(request *Dto.UpdateDeliveryInformationRequest) (err error) {
	return h.Repo.UpdateMainInformation(request)
}

func (h *DeliveryService) GetPreviousNote(id int) (note string, err error) {
	currentDelivery, err := h.Repo.GetById(id)
	if err != nil {
		return "", err
	}

	invoiceId := currentDelivery.InvoiceId
	deliveries, err := h.Repo.GetAllByInvoiceId(invoiceId)
	if err != nil {
		return "", err
	}

	if len(deliveries) < 2 {
		return "", nil
	}

	for i, delivery := range deliveries {
		if delivery.ID == id {
			if i == 0 {
				return "", nil
			}

			return deliveries[i-1].Note, nil
		}
	}

	return "", nil
}

func (h *DeliveryService) LockDeliveryOrder(request *Dto.LockDeliveryOrderRequest) (err error) {
	invoiceStatus := 2

	sales, err := h.InvoiceRepo.GetSalesByInvoiceId(request.InvoiceId)
	if err != nil {
		return err
	}

	invoice, err := h.InvoiceRepo.GetById(request.InvoiceId)
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
		invoiceStatus++
	}

	if invoice.PoPath != "-" && invoice.PoPath != "" {
		invoiceStatus++
	}

	if invoice.IsTaxable {
		if invoice.PoCode != "" && invoice.PoCode != "-" {
			invoiceStatus++
		}
	} else {
		invoiceStatus++
	}

	if err = h.InvoiceRepo.UpdateStatus(&Dto.UpdateStatusRequest{InvoiceId: request.InvoiceId, InvoiceStatusId: invoiceStatus}); err != nil {
		return err
	}

	if err = h.Repo.UpdateDeliveryStatus(&Dto.UpdateDeliveryStatusRequest{DeliveryId: request.DeliveryId, Status: 2}); err != nil {
		return err
	}

	return nil
}

func generateOrderCode(invoiceCode string, deliveriesCount int) (val string) {
	ascii := 61 + deliveriesCount

	val = invoiceCode + "/" + strconv.Itoa(ascii)
	return val
}
