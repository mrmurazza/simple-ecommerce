package impl

import (
	"ecommerce/domain/order"
	"ecommerce/domain/user"
	"ecommerce/dto/request"
	"errors"
	"fmt"
	"strings"
	"time"
)

type service struct {
	repo    order.Repository
	userSvc user.Service
}

func NewService(repo order.Repository, userSvc user.Service) order.Service {
	return &service{
		repo:    repo,
		userSvc: userSvc,
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

	return s.repo.SaveNewOrder(order)
}

func (s *service) CheckoutOrder(userId, orderId int) error {
	o, err := s.repo.GetByIdAndUserId(userId, orderId)
	if err != nil {
		return err
	}
	if o == nil {
		return errors.New("not found")
	}

	now := time.Now()
	o.Status = order.StatusDone
	o.UpdatedAt = &now

	err = s.repo.UpdateOrder(o)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetOrderHistories(userId int) ([]order.Order, error) {
	return s.repo.GetOrderByCustomer(userId)
}

func (s *service) GetAllOrders() ([]order.Order, error) {
	return s.repo.GetAllOrders()
}

func (s *service) GetAllProducts() ([]order.Product, error) {
	return s.repo.GetAllProducts()
}

func (s *service) RemindPendingOrder() error {
	pendingOrders, err := s.repo.GetByStatus(order.StatusPending, true)
	if err != nil {
		return nil
	}

	userIds := make([]int, 0)
	orderByUserIds := make(map[int][]*order.Order)
	for _, po := range pendingOrders {
		userIds = append(userIds, po.CustomerId)

		if _, ok := orderByUserIds[po.CustomerId]; !ok {
			orderByUserIds[po.CustomerId] = make([]*order.Order, 0)
		}
		orderByUserIds[po.CustomerId] = append(orderByUserIds[po.CustomerId], &po)
	}

	users, err := s.userSvc.GetByIds(userIds)
	if err != nil {
		return nil
	}

	for _, u := range users {
		currPOs := orderByUserIds[u.ID]

		s.sendEmailReminder(u, currPOs)

	}
	return nil
}

func (s *service) sendEmailReminder(user user.User, orders []*order.Order) error {
	for _, o := range orders {
		productList := make([]string, 0)
		for _, ou := range o.Units {
			productList = append(productList, fmt.Sprintf("%d. %s", ou.ProductId, ou.Name))
		}

		productContent := strings.Join(productList, "\n")
		link := fmt.Sprintf("%s/%s/%d", "localhost:8080/api/v1", "checkout", o.ID)
		emailContent := fmt.Sprintf(
			"Hi %s, segera selesaikan order kamu di %s. Dengan rincian product %s",
			user.Email,
			link,
			productContent)

		// Dummy replace sending email
		fmt.Println(emailContent)
	}

	return nil
}
