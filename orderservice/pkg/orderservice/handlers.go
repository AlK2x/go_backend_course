package orderservice

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Order struct {
	MenuItems []OrderItem `json:"menuItems"`
}

type OrderItem struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type OrderResponse struct {
	Order
	OrderedAtTimestamp string `json:"orderedAtTimestamp"`
	Cost               int    `json:"cost"`
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	s := router.PathPrefix("/api/v2").Subrouter()

	s.Name("GetOrderById").
		Methods(http.MethodGet).
		Path("/order/{orderId}").
		Handler(Logger(createHandlerFunc(GetOrderById)))

	s.Name("GetOrderList").
		Methods(http.MethodGet).
		Path("/orders").
		Handler(Logger(createHandlerFunc(GetOrderList)))

	s.Name("GetOrderList").
		Methods(http.MethodGet).
		Path("/").
		Handler(Logger(createHandlerFunc(GetOrderList)))

	return router
}

func GetOrderById(w http.ResponseWriter, r *http.Request) error {
	orderId := mux.Vars(r)["orderID"]
	err := jsonResponse(w, OrderResponse{
		Order: Order{
			MenuItems: []OrderItem{{
				Id:       orderId,
				Quantity: 1,
			}},
		},
		OrderedAtTimestamp: time.Now().String(),
		Cost:               999,
	})
	if err != nil {
		return err
	}
	return nil
}

func GetOrderList(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(http.StatusOK)
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
