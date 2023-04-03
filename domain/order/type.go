package order

import (
	"ecommerce/domain/user"
)

type Service interface {
	Order(u user.User) error
	GetOrderHistories(u user.User) ([]*Order, error)
	GetAllOrders() ([]*Order, error)
	GetAllProducts() ([]*Product, error)
}

type Repository interface {
	SaveOrder(o Order) error
	GetOrderByCustomer(u user.User) ([]Order, error)
	GetAllOrders() ([]Order, error)
	GetAllProducts() ([]Product, error)
}
