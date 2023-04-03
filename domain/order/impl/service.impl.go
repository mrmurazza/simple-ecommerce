package impl

import (
	"ecommerce/domain/order"
	"ecommerce/domain/user"
)

type service struct {
	repo order.Repository
}

func NewService(repo order.Repository) order.Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Order(u user.User) error {
	return nil
}

func (s *service) GetOrderHistories(u user.User) ([]order.Order, error) {
	return s.GetOrderHistories(u)
}

func (s *service) GetAllOrders() ([]order.Order, error) {
	return s.repo.GetAllOrders()
}

func (s *service) GetAllProducts() ([]order.Product, error) {
	return s.GetAllProducts()
}
