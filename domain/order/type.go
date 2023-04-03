package order

import (
	"ecommerce/dto/request"
)

type Service interface {
	Order(req request.CreateOrderRequest) error
	GetOrderHistories(userId int) ([]Order, error)
	GetAllOrders() ([]Order, error)
	GetAllProducts() ([]Product, error)
}

type Repository interface {
	SaveOrder(o Order) error
	GetOrderByCustomer(userId int) ([]Order, error)
	GetAllOrders() ([]Order, error)
	GetAllProducts() ([]Product, error)
	GetProductByIds(ids []int) ([]Product, error)
}
