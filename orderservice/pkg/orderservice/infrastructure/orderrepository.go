package infrastructure

import (
	"database/sql"
	"orderservice/pkg/orderservice/model"
)

func CreateRepository(db *sql.DB) model.OrderRepository {
	return &MySqlOrderRepository{
		db: db,
	}
}

type MySqlOrderRepository struct {
	db *sql.DB
}

func (m *MySqlOrderRepository) GetAll() (map[string]model.Order, error) {
	rows, err := m.db.Query("SELECT id, order_id, quantity, cost FROM menu_item")
	if err != nil {
		return nil, nil
	}

	defer rows.Close()

	var orderId string
	var itemId string
	var quantity int
	var cost sql.NullFloat64
	orders := make(map[string]model.Order)
	for rows.Next() {
		err := rows.Scan(&itemId, &orderId, &quantity, &cost)
		if err != nil {
			return orders, err
		}

		order, ok := orders[orderId]
		if !ok {
			order = model.Order{
				Id:        orderId,
				MenuItems: []model.MenuItem{},
			}
		}

		order.MenuItems = append(order.MenuItems, model.MenuItem{
			Id:       itemId,
			Quantity: quantity,
		})
		orders[orderId] = order
	}

	return orders, nil
}

func (m *MySqlOrderRepository) GetById(id string) (*model.Order, error) {
	rows, err := m.db.Query("SELECT id, order_id, quantity, cost FROM menu_item WHERE order_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderId string
	var itemId string
	var quantity int
	var cost sql.NullFloat64
	order := model.Order{MenuItems: make([]model.MenuItem, 0)}
	for rows.Next() {
		err := rows.Scan(&itemId, &orderId, &quantity, &cost)

		if err != nil {
			return nil, err
		}
		order.MenuItems = append(order.MenuItems, model.MenuItem{
			Id:       itemId,
			Quantity: quantity,
		})
		order.Id = orderId
	}
	return &order, nil
}

func (m *MySqlOrderRepository) Add(order *model.Order) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	_, err = m.db.Exec("INSERT INTO `order` SET id = ?", order.Id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	for _, item := range order.MenuItems {
		_, err = m.db.Exec("INSERT INTO menu_item SET id = ?, order_id = ?, quantity = ?", item.Id, order.Id, item.Quantity)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (m *MySqlOrderRepository) Delete(order *model.Order) error {
	_, err := m.db.Exec("DELETE FROM `order` WHERE id = ?", order.Id)
	if err != nil {
		return err
	}

	return nil
}

func (m *MySqlOrderRepository) Update(order *model.Order) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	for _, item := range order.MenuItems {
		_, err = m.db.Exec("UPDATE menu_item SET quantity = ?  WHERE id = ?", item.Quantity, item.Id)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	return err
}
