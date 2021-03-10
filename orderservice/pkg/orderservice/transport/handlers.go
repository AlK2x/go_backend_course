package transport

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"orderservice/pkg/orderservice/model"
	"time"
)

type orderResponse struct {
	model.Order
	OrderedAtTimestamp string `json:"orderedAtTimestamp"`
	Cost               int    `json:"cost"`
}

type orderListResponse struct {
	Id        string           `json:"id"`
	MenuItems []model.MenuItem `json:"menuItems"`
}

type createOrderResponse struct {
	Id string `json:"id"`
}

type Server struct {
	orderRepository model.OrderRepository
}

func NewServer(repo model.OrderRepository) *Server {
	return &Server{
		orderRepository: repo,
	}
}

func NewRouter(srv *Server) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	s := router.PathPrefix("/api/v2").Subrouter()

	s.Name("GetOrderById").
		Methods(http.MethodGet).
		Path("/order/{orderId}").
		Handler(createHandlerFunc(srv.GetOrderById))

	s.Name("GetOrderList").
		Methods(http.MethodGet).
		Path("/orders").
		Handler(createHandlerFunc(srv.GetOrderList))

	s.Name("CreateOrder").
		Methods(http.MethodPost).
		Path("/order").
		Handler(createHandlerFunc(srv.CreateOrder))

	s.Name("DeleteOrder").
		Methods(http.MethodDelete).
		Path("/order/{orderId}").
		Handler(createHandlerFunc(srv.DeleteOrder))

	s.Name("UpdateOrder").
		Methods(http.MethodPut).
		Path("/order/{orderId}").
		Handler(createHandlerFunc(srv.UpdateOrder))

	return logger(router)
}

func (s *Server) CreateOrder(w http.ResponseWriter, r *http.Request) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	var order model.Order
	err = json.Unmarshal(b, &order)
	if err != nil {
		return err
	}

	for _, item := range order.MenuItems {
		if item.Quantity <= 0 {
			badRequestResponse(w, "Incorrect quantity")
		}
	}

	id := uuid.New()
	order.Id = id.String()
	err = s.orderRepository.Add(&order)
	if err != nil {
		return err
	}
	err = jsonResponse(w, createOrderResponse{
		Id: id.String(),
	})

	return err
}

func (s *Server) GetOrderById(w http.ResponseWriter, r *http.Request) error {
	orderId, ok := mux.Vars(r)["orderId"]
	if !ok {
		badRequestResponse(w, "OrderId not found")
		return nil
	}

	order, err := s.orderRepository.GetById(orderId)
	if err != nil {
		return err
	}
	if order == nil {
		return jsonResponse(w, nil)
	}

	err = jsonResponse(w, orderResponse{
		Order:              *order,
		OrderedAtTimestamp: time.Now().String(),
		Cost:               999,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) GetOrderList(w http.ResponseWriter, _ *http.Request) error {
	orderMap, err := s.orderRepository.GetAll()
	if err != nil {
		return err
	}

	var response []orderListResponse
	for _, item := range orderMap {
		response = append(response, orderListResponse{
			Id:        item.Id,
			MenuItems: item.MenuItems,
		})
	}

	err = jsonResponse(w, response)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) DeleteOrder(w http.ResponseWriter, r *http.Request) error {
	orderId, ok := mux.Vars(r)["orderId"]
	if !ok {
		badRequestResponse(w, "OrderId not found")
		return nil
	}

	order, err := s.orderRepository.GetById(orderId)
	if err != nil {
		return err
	}
	if order == nil {
		return nil
	}

	err = s.orderRepository.Delete(order)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) UpdateOrder(w http.ResponseWriter, r *http.Request) error {
	orderId, ok := mux.Vars(r)["orderId"]
	if !ok {
		badRequestResponse(w, "OrderId not found")
		return nil
	}
	order, err := s.orderRepository.GetById(orderId)
	if err != nil {
		return err
	}
	if order == nil {
		return nil
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	var o model.Order
	err = json.Unmarshal(body, &o)
	if err != nil {
		return err
	}

	order.MenuItems = o.MenuItems

	err = s.orderRepository.Update(order)
	if err != nil {
		return err
	}
	return nil
}

func jsonResponse(w http.ResponseWriter, r interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8\"")
	resp, err := json.Marshal(r)
	if err != nil {
		return err
	}
	_, err = w.Write(resp)
	return err
}

func badRequestResponse(w http.ResponseWriter, err string) {
	http.Error(w, err, http.StatusBadRequest)
}

func errorResponse(w http.ResponseWriter, err string) {
	http.Error(w, err, http.StatusInternalServerError)
}

func createHandlerFunc(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			errorResponse(w, err.Error())
		}
	}
}
