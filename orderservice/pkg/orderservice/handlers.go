package orderservice

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type order struct {
	MenuItems []orderItem `json:"menuItems"`
}

type orderItem struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type orderResponse struct {
	order
	OrderedAtTimestamp string `json:"orderedAtTimestamp"`
	Cost               int    `json:"cost"`
}

type orderListResponse struct {
	Id        string      `json:"id"`
	MenuItems []orderItem `json:"menuItems"`
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	s := router.PathPrefix("/api/v2").Subrouter()

	s.Name("getOrderById").
		Methods(http.MethodGet).
		Path("/order/{orderId}").
		Handler(logger(createHandlerFunc(getOrderById)))

	s.Name("getOrderList").
		Methods(http.MethodGet).
		Path("/orders").
		Handler(logger(createHandlerFunc(getOrderList)))

	return router
}

func getOrderById(w http.ResponseWriter, r *http.Request) error {
	orderId, ok := mux.Vars(r)["orderId"]
	if !ok {
		badRequestResponse(w, "OrderId not found")
		return nil
	}

	err := jsonResponse(w, orderResponse{
		order: order{
			MenuItems: []orderItem{{
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

func getOrderList(w http.ResponseWriter, _ *http.Request) error {
	list := orderListResponse{
		Id: "d290f1ee-6c56-4b01-90e6-d701748f0851",
		MenuItems: []orderItem{{
			Id:       "f290d1ce-6c234-4b31-90e6-d701748f0851",
			Quantity: 1,
		}},
	}
	err := jsonResponse(w, list)
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
