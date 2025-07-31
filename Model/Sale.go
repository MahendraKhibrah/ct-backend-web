package Model

import "time"

type Sale struct {
	ID           int
	InvoiceId    int
	ProductId    int
	Quantity     int
	Price        int
	SendStatus   bool
	NotSentCount int
	Unit         string
	CreatedAt    time.Time
	Product      Product `gorm:"foreignKey:ProductId" json:"Product,omitempty"`
	Invoice      Invoice `gorm:"foreignKey:InvoiceId" json:"Invoice,omitempty"`
}
