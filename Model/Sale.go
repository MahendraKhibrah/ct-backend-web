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
	CreatedAt    time.Time
	Product      Product `gorm:"foreignKey:ProductId" json:"Product,omitempty"`
}
