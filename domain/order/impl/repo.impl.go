package impl

import (
	"ecommerce/domain/order"
	"ecommerce/domain/user"

	"github.com/jinzhu/gorm"
)

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) order.Repository {
	return &repo{
		db: db,
	}
}

func (r *repo) SaveOrder(o order.Order) error {
	res := r.db.Create(&o)

	err := res.Error
	if err != nil {
		return err
	}

	for _, unit := range o.Units {
		unit.OrderId = o.ID
	}

	res = r.db.Create(o.Units)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) GetOrderByCustomer(u user.User) ([]order.Order, error) {
	var orders []order.Order

	err := r.db.Model(order.Order{}).
		Where("customer_id = ?", u.ID).
		Order("created_at DESC").
		Find(&orders).
		Error
	if err != nil {
		return nil, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	orderUnits, err := r.getOrderUnitsByOrders(orders)
	if err != nil {
		return nil, err
	}

	orders = r.assignOrderUnits(orders, orderUnits)
	return orders, nil
}

func (r *repo) getOrderUnitsByOrders(orders []order.Order) ([]order.OrderUnit, error) {
	orderIds := make([]int, 0)
	for _, o := range orders {
		orderIds = append(orderIds, o.ID)
	}

	var orderUnits []order.OrderUnit
	err := r.db.Model(order.OrderUnit{}).
		Where("order_id in ?", orderIds).
		Order("created_at DESC").
		Find(&orders).
		Error
	if err != nil {
		return nil, err
	}

	return orderUnits, nil
}

func (r *repo) assignOrderUnits(orders []order.Order, orderUnits []order.OrderUnit) []order.Order {
	orderUnitsByOrder := make(map[int][]order.OrderUnit)
	for _, ou := range orderUnits {
		if _, ok := orderUnitsByOrder[ou.OrderId]; !ok {
			orderUnitsByOrder[ou.OrderId] = make([]order.OrderUnit, 0)
		}

		orderUnitsByOrder[ou.OrderId] = append(orderUnitsByOrder[ou.OrderId], ou)
	}

	for _, o := range orders {
		orderUnits := orderUnitsByOrder[o.ID]
		o.Units = orderUnits
	}

	return orders
}

func (r *repo) GetAllOrders() ([]order.Order, error) {
	var orders []order.Order

	err := r.db.Model(order.Order{}).
		Order("created_at DESC").
		Find(&orders).
		Error
	if err != nil {
		return nil, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return orders, nil
}

func (r *repo) GetAllProducts() ([]order.Product, error) {
	var products []order.Product

	err := r.db.Model(order.Product{}).
		Order("created_at DESC").
		Find(&products).
		Error
	if err != nil {
		return nil, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return products, nil
}
