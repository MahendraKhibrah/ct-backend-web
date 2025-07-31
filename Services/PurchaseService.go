package Services

import (
	"ct-backend/Model"
	"ct-backend/Model/Dto"
	"ct-backend/Repository"
)

type (
	IPurchaseService interface {
		AddPurchase(request *Dto.CreatePurchaseRequest) error
		GetAllPurchase() (purchases []Model.Purchase, err error)
		PayDebt(id int) (err error)
		DeletePurchase(id int) (err error)
	}

	PurchaseService struct {
		PurchaseRepository Repository.IPurchaseRepository
		ProductRepository  Repository.IProductRepository
	}
)

func PurchaseServiceProvider(purchaseRepository Repository.IPurchaseRepository, productRepository Repository.IProductRepository) *PurchaseService {
	return &PurchaseService{
		PurchaseRepository: purchaseRepository,
		ProductRepository:  productRepository,
	}
}

func (h *PurchaseService) AddPurchase(request *Dto.CreatePurchaseRequest) (err error) {

	if err = h.PurchaseRepository.AddPurchase(request); err != nil {
		return err
	}

	if err = h.ProductRepository.SumStockProduct(request.ProductId, request.Count, nil); err != nil {
		return err
	}

	return nil
}

func (h *PurchaseService) GetAllPurchase() (purchases []Model.Purchase, err error) {
	return h.PurchaseRepository.GetAllPurchase()
}

func (h *PurchaseService) PayDebt(id int) (err error) {
	return h.PurchaseRepository.ChangeIsPaidStatus(id)
}

func (h *PurchaseService) DeletePurchase(id int) (err error) {
	var (
		purchase *Model.Purchase
	)

	if purchase, err = h.PurchaseRepository.GetPurchaseById(id); err != nil {
		return err
	}

	if err = h.ProductRepository.SumStockProduct(purchase.ProductId, purchase.Count*-1, nil); err != nil {
		return err
	}

	if err = h.PurchaseRepository.DeletePurchase(id); err != nil {
		return err
	}

	return nil
}
