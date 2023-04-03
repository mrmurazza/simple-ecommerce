package order

import (
	"time"
)

type Status string

const (
	StatusPending Status = "PENDING"
	StatusDone    Status = "DONE"
)

func GetAllOrderStatus() map[Status]bool {
	return map[Status]bool{
		StatusPending: true,
		StatusDone:    true,
	}
}

type Order struct {
	ID          int
	CustomerId  int
	Status      Status
	TotalQty    int
	TotalAmount int

	CreatedAt *time.Time `gorm:"default:current_timestamp"`
	UpdatedAt *time.Time `gorm:"default:current_timestamp"`

	Units []OrderUnit
}

type OrderUnit struct {
	ID          int
	OrderId     int
	ProductId   int
	Price       int
	Qty         int
	Name        string
	Description string
	Image       string

	CreatedAt *time.Time `gorm:"default:current_timestamp"`
}

type Product struct {
	ID          int
	Name        string
	Price       int
	Description string
	Image       string

	CreatedAt *time.Time `gorm:"default:current_timestamp"`
	UpdatedAt *time.Time `gorm:"default:current_timestamp"`
}
