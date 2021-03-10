package transport

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"orderservice/pkg/orderservice/infrastructure"
	"orderservice/pkg/orderservice/model"
	"testing"
)

func TestGetOrderById(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v2/order/{orderId}", nil)
	r = mux.SetURLVars(r, map[string]string{"orderId": "123"})
	srv := Server{orderRepository: infrastructure.NewInMemoryRepository()}

	err := srv.GetOrderById(w, r)
	if err != nil {
		t.Fatal(err)
	}

	response := w.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Status code is wrong. Have: %d, want: %d", response.StatusCode, http.StatusOK)
	}

	jsonStr, err := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	order := model.Order{}
	if err = json.Unmarshal(jsonStr, &order); err != nil {
		t.Errorf("Can't parse json with error %v", err)
	}
}

func TestCreateOrder(t *testing.T) {
	w := httptest.NewRecorder()
	msg := model.Order{MenuItems: []model.MenuItem{{
		Id:       "abc-def-ghi",
		Quantity: 1},
	}}
	body, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/api/v2/order", bytes.NewReader(body))

	srv := Server{orderRepository: infrastructure.NewInMemoryRepository()}
	err = srv.CreateOrder(w, r)
	if err != nil {
		t.Fatal(err)
	}

	response := w.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Status code is wrong. Have: %d, want: %d", response.StatusCode, http.StatusOK)
	}

	jsonStr, err := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	order := model.Order{}
	if err = json.Unmarshal(jsonStr, &order); err != nil {
		t.Errorf("Can't parse json with error %v", err)
	}
}

func TestGetOrderList(t *testing.T) {
	w := httptest.NewRecorder()
	srv := Server{orderRepository: infrastructure.NewInMemoryRepository()}
	err := srv.GetOrderList(w, nil)
	if err != nil {
		t.Fatal(err)
	}

	response := w.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Status code is wrong. Have: %d, want: %d", response.StatusCode, http.StatusOK)
	}

	jsonString, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	list := orderListResponse{}
	if err = json.Unmarshal(jsonString, &list); err != nil {
		t.Errorf("Can't parse json response with error %v", err)
	}
}
