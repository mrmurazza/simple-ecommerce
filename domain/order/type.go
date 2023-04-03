package order

import (
	"ecommerce/dto/request"
)

type Service interface {
	Order(req request.CreateOrderRequest) error
	CheckoutOrder(userId, orderId int) error
	GetOrderHistories(userId int) ([]Order, error)
	GetAllOrders() ([]Order, error)
	GetAllProducts() ([]Product, error)
}

type Repository interface {
	SaveNewOrder(o Order) error
	UpdateOrder(o *Order) error
	GetByIdAndUserId(id, userId int) (*Order, error)
	GetOrderByCustomer(userId int) ([]Order, error)
	GetAllOrders() ([]Order, error)
	GetAllProducts() ([]Product, error)
	GetProductByIds(ids []int) ([]Product, error)
}
