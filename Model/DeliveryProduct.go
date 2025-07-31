package Model

import "time"

type DeliveryProduct struct {
	ID         int
	DeliveryID int
	SalesID    int
	Quantity   int
	CreatedAt  time.Time
	Sale       Sale `gorm:"foreignKey:SalesID"`
}

type ShortDeliveryProduct struct {
	ID       int
	Name     string
	Quantity int
	SaleID   int
	Unit     string
}
