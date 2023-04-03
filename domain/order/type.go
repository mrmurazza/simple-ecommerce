package order

import (
	"ecommerce/domain/user"
	"ecommerce/dto/request"
)

type Service interface {
	Order(req request.CreateOrderRequest) error
	GetOrderHistories(u user.User) ([]Order, error)
	GetAllOrders() ([]Order, error)
	GetAllProducts() ([]Product, error)
}

type Repository interface {
	SaveOrder(o Order) error
	GetOrderByCustomer(u user.User) ([]Order, error)
	GetAllOrders() ([]Order, error)
	GetAllProducts() ([]Product, error)
	GetProductByIds(ids []int) ([]Product, error)
}
