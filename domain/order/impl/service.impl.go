package impl

import (
	"ecommerce/domain/order"
	"ecommerce/domain/user"
	"ecommerce/dto/request"
	"fmt"
	"time"
)

type service struct {
	repo order.Repository
}

func NewService(repo order.Repository) order.Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Order(req request.CreateOrderRequest) error {
	productIds := make([]int, 0)
	for _, detailReq := range req.Products {
		productIds = append(productIds, detailReq.ProductId)

	}

	products, err := s.repo.GetProductByIds(productIds)
	if err != nil {
		return err
	}

	existingProductMap := make(map[int]*order.Product)
	for _, p := range products {
		existingProductMap[p.ID] = &p
	}

	now := time.Now()
	totalQty := 0
	totalAmount := 0
	orderUnits := make([]order.OrderUnit, 0)
	for _, detailReq := range req.Products {
		product, ok := existingProductMap[detailReq.ProductId]
		if !ok {
			return fmt.Errorf("product id %d not found", detailReq.ProductId)
		}

		if product.Qty < detailReq.Qty {
			return fmt.Errorf("existing stock insufficient for product %d", detailReq.ProductId)
		}

		orderUnits = append(orderUnits, order.OrderUnit{
			ProductId:   product.ID,
			Name:        product.Name,
			Description: product.Description,
			Image:       product.Image,
			Price:       product.Price,
			Qty:         detailReq.Qty,
			CreatedAt:   &now,
		})
		totalQty += detailReq.Qty
		totalAmount += detailReq.Qty * product.Price
	}

	order := order.Order{
		CustomerId:  req.UserId,
		Status:      order.StatusPending,
		TotalQty:    totalQty,
		TotalAmount: totalAmount,
		CreatedAt:   &now,
		UpdatedAt:   &now,
		Units:       orderUnits,
	}

	return s.repo.SaveOrder(order)
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
