package model

type Order struct {
	Id        string
	MenuItems []MenuItem `json:"menuItems"`
}

type MenuItem struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type OrderRepository interface {
	GetAll() (map[string]Order, error)
	GetById(id string) (*Order, error)
	Add(order *Order) error
	Delete(order *Order) error
	Update(order *Order) error
}
