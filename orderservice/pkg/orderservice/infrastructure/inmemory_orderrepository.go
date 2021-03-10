package infrastructure

import (
	"orderservice/pkg/orderservice/model"
)

func NewInMemoryRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{orders: map[string]model.Order{}}
}

type InMemoryOrderRepository struct {
	orders map[string]model.Order
}

func (i InMemoryOrderRepository) GetAll() (map[string]model.Order, error) {
	return i.orders, nil
}

func (i InMemoryOrderRepository) GetById(id string) (*model.Order, error) {
	order, ok := i.orders[id]
	if !ok {
		return nil, nil
	}
	return &order, nil
}

func (i InMemoryOrderRepository) Add(order *model.Order) error {
	i.orders[order.Id] = *order
	return nil
}

func (i InMemoryOrderRepository) Delete(order *model.Order) error {
	delete(i.orders, order.Id)
	return nil
}

func (i InMemoryOrderRepository) Update(order *model.Order) error {
	i.orders[order.Id] = *order
	return nil
}
