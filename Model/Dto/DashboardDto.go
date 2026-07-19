package Dto

import "time"

type DashboardInvoice struct {
	ID          int       `json:"id"`
	InvoiceCode string    `json:"invoice_code"`
	ClientName  string    `json:"client_name"`
	ProjectName string    `json:"project_name"`
	StatusID    int       `json:"status_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type DashboardDelivery struct {
	ID         int       `json:"id"`
	OrderCode  string    `json:"order_code"`
	ClientName string    `json:"client_name"`
	StatusID   int       `json:"status_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type DashboardReceipt struct {
	ID         int       `json:"id"`
	Number     int       `json:"number"`
	ClientName string    `json:"client_name"`
	StatusID   int       `json:"status_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type DashboardWorkQueue struct {
	Invoices   []DashboardInvoice  `json:"invoices"`
	Deliveries []DashboardDelivery `json:"deliveries"`
	Receipts   []DashboardReceipt  `json:"receipts"`
}
