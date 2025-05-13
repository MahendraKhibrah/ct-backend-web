package Dto

type CreateInvoiceRequest struct {
	ClientId            int    `json:"client_id" binding:"required"`
	IsTaxable           bool   `json:"is_taxable"`
	InvoiceCode         string `json:"invoice_code"`
	Seller              string `json:"seller"`
	Platform            string `json:"platform"`
	PaymentMethod       string `json:"payment_method"`
	PlatformDescription string `json:"platform_description"`
	PlatformNumber      string `json:"platform_number"`
}

type UpdateDocumentRequest struct {
	InvoiceId    int    `json:"invoice_id" binding:"required"`
	DocumentType string `json:"document_type" binding:"required"`
	DocumentPath string `json:"document_path" binding:"required"`
}

type AddSaleRequest struct {
	Id        int `json:"id"`
	InvoiceId int `json:"invoice_id" binding:"required"`
	ProductId int `json:"product_id" binding:"required"`
	Count     int `json:"count" binding:"required"`
	Price     int `json:"price" binding:"required"`
}

type UpdateSaleRequest struct {
	Id           int `json:"id"`
	ProductId    int `json:"product_id" binding:"required"`
	CurrentCount int `json:"current_count" binding:"required"`
	Count        int `json:"count" binding:"required"`
	Price        int `json:"price" binding:"required"`
}

type UpdateFakturRequest struct {
	InvoiceId   int `json:"invoice_id" binding:"required"`
	Discount    int `json:"discount"`
	PaymentTerm int `json:"payment_term" binding:"required"`
}

type UpdateMainInformationRequest struct {
	InvoiceId           int    `json:"invoice_id" binding:"required"`
	PoCode              string `json:"po_code" binding:"required"`
	Note                string `json:"note" binding:"required"`
	Seller              string `json:"seller" binding:"required"`
	Platform            string `json:"platform" binding:"required"`
	PaymentMethod       string `json:"payment_method" binding:"required"`
	PlatformDescription string `json:"platform_description" binding:"required"`
	PlatformNumber      string `json:"platform_number" binding:"required"`
	Project             string `json:"project" binding:"required"`
	Date                string `json:"date"`
}

type UpdateNoteRequest struct {
	InvoiceId int    `json:"invoice_id" binding:"required"`
	Note      string `json:"note" binding:"required"`
}

type UpdateStatusRequest struct {
	InvoiceId       int `json:"invoice_id" binding:"required"`
	InvoiceStatusId int `json:"invoice_status_id" binding:"required"`
}

type UpdateNotSentSaleRequest struct {
	SaleId int `json:"sale_id" binding:"required"`
	Count  int `json:"count" binding:"required"`
}

type GetDocumentUrlRequest struct {
	Key string `json:"key" binding:"required"`
}

type GetInvoicesRequest struct {
	Search      string `form:"search"`
	MaxQuantity string `form:"max_quantity"`
	MinQuantity string `form:"min_quantity"`
	Company     string `form:"company"`
}
