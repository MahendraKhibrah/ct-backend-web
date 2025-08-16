package Repository

import (
	"ct-backend/Model"
	"ct-backend/Model/Dto"
	"ct-backend/Utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	IDeliveryRepository interface {
		GetById(id int) (delivery *Model.DeliveryOrder, err error)
		GetAll(ctx *gin.Context) (deliveries []Model.DeliveryOrder, err error)
		CreateDeliveryOrder(request *Dto.IdRequest, orderCode string) (err error)
		CreateDeliveryProduct(request *Dto.CreateDeliveryProductRequest) (err error)
		UpdateDeliveryProduct(request *Dto.UpdateDeliveryProductRequest) (err error)
		DeleteDeliveryProduct(request *Dto.DeleteDeliveryProductRequest) (err error)
		GetAllDeliveryProduct(deliveryId int) (deliveryProducts []Model.DeliveryProduct, err error)
		UpdateSender(request *Dto.UpdateSenderRequest) (err error)
		UpdateMainInformation(request *Dto.UpdateDeliveryInformationRequest) (err error)
		GetAllByInvoiceId(invoiceId int) (deliveries []Model.DeliveryOrder, err error)
		UpdateDeliveryStatus(request *Dto.UpdateDeliveryStatusRequest) (err error)
		DeleteDelivery(request *Dto.IdRequest, trx *gorm.DB) (err error)
		GetDeliveryProductById(id int) (deliveryProduct *Model.DeliveryProduct, err error)
	}

	DeliveryRepository struct {
		DB *gorm.DB
	}
)

func DeliveryRepositoryProvider(DB *gorm.DB) *DeliveryRepository {
	return &DeliveryRepository{
		DB: DB,
	}
}

func (h *DeliveryRepository) GetById(id int) (delivery *Model.DeliveryOrder, err error) {
	if err := h.DB.Preload("Invoice.Client").Where("id = ?", id).First(&delivery).Error; err != nil {
		return nil, err
	}

	return delivery, nil
}

func (h *DeliveryRepository) GetAll(ctx *gin.Context) (deliveries []Model.DeliveryOrder, err error) {
	if err := h.DB.Preload("Invoice.Client").Where("order_code ILIKE ?", "%"+ctx.Query("order_code")+"%").Order("created_at desc").Scopes(Utils.Paginate(ctx)).Find(&deliveries).Error; err != nil {
		return nil, err
	}

	return deliveries, nil
}

func (h *DeliveryRepository) CreateDeliveryOrder(request *Dto.IdRequest, orderCode string) (err error) {
	deliveryOrder := &Model.DeliveryOrder{
		InvoiceId: request.Id,
		Status:    1,
		OrderCode: orderCode,
	}

	if err := h.DB.Create(&deliveryOrder).Error; err != nil {
		return err
	}

	return nil
}

func (h *DeliveryRepository) CreateDeliveryProduct(request *Dto.CreateDeliveryProductRequest) (err error) {
	deliveryProduct := &Model.DeliveryProduct{
		DeliveryID: request.DeliveryID,
		SalesID:    request.SalesID,
		Quantity:   request.Quantity,
	}

	if err := h.DB.Create(&deliveryProduct).Error; err != nil {
		return err
	}

	return nil
}

func (h *DeliveryRepository) UpdateDeliveryProduct(request *Dto.UpdateDeliveryProductRequest) (err error) {
	deliveryProduct := &Model.DeliveryProduct{
		ID: request.ID,
	}

	if err := h.DB.Model(&deliveryProduct).Update("quantity", request.Quantity).Error; err != nil {
		return err
	}

	return nil
}

func (h *DeliveryRepository) DeleteDeliveryProduct(request *Dto.DeleteDeliveryProductRequest) (err error) {
	deliveryProduct := &Model.DeliveryProduct{
		ID: request.ID,
	}

	if err := h.DB.Delete(&deliveryProduct).Error; err != nil {
		return err
	}

	return nil
}

func (h *DeliveryRepository) GetAllDeliveryProduct(deliveryId int) (deliveryProducts []Model.DeliveryProduct, err error) {
	if err := h.DB.Preload("Sale.Product").Where("delivery_id = ?", deliveryId).Find(&deliveryProducts).Error; err != nil {
		return nil, err
	}

	return deliveryProducts, nil
}

func (h *DeliveryRepository) UpdateSender(request *Dto.UpdateSenderRequest) (err error) {
	deliveryOrder := &Model.DeliveryOrder{
		ID: request.ID,
	}

	if err := h.DB.Model(&deliveryOrder).Update("sender_id", request.SenderId).Error; err != nil {
		return err
	}

	return nil
}

func (h *DeliveryRepository) UpdateMainInformation(request *Dto.UpdateDeliveryInformationRequest) (err error) {
	deliveryOrder := &Model.DeliveryOrder{
		ID: request.DeliveryId,
	}

	if err := h.DB.Model(&deliveryOrder).Update("note", request.Note).Update("place", request.Place).Error; err != nil {
		return err
	}

	return nil
}

func (h *DeliveryRepository) GetAllByInvoiceId(invoiceId int) (deliveries []Model.DeliveryOrder, err error) {
	if err := h.DB.Where("invoice_id = ?", invoiceId).Find(&deliveries).Error; err != nil {
		return nil, err
	}

	return deliveries, nil
}

func (h *DeliveryRepository) UpdateDeliveryStatus(request *Dto.UpdateDeliveryStatusRequest) (err error) {
	deliveryOrder := &Model.DeliveryOrder{
		ID: request.DeliveryId,
	}

	if err := h.DB.Model(&deliveryOrder).Update("status", request.Status).Error; err != nil {
		return err
	}

	return nil
}

func (h *DeliveryRepository) DeleteDelivery(request *Dto.IdRequest, trx *gorm.DB) (err error) {
	db := trx
	if db == nil {
		db = h.DB
	}

	deliveryOrder := &Model.DeliveryOrder{
		ID: request.Id,
	}

	if err := db.Delete(&deliveryOrder).Error; err != nil {
		return err
	}

	if err := db.Where("delivery_id = ?", request.Id).Delete(&Model.DeliveryProduct{}).Error; err != nil {
		return err
	}

	return nil
}

func (h *DeliveryRepository) GetDeliveryProductById(id int) (deliveryProduct *Model.DeliveryProduct, err error) {
	if err := h.DB.Preload("Sale.Product").Where("id = ?", id).First(&deliveryProduct).Error; err != nil {
		return nil, err
	}

	return deliveryProduct, nil
}
