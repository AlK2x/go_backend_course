package orderservice

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOrderById(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v2/order/{orderId}", nil)
	r = mux.SetURLVars(r, map[string]string{"orderId": "123"})
	err := getOrderById(w, r)
	if err != nil {
		t.Fatal(err)
	}

	response := w.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Status code is wrong. Have: %d, want: %d", response.StatusCode, http.StatusOK)
	}

	jsonStr, err := ioutil.ReadAll(response.Body)
	t.Log(jsonStr)
	_ = response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	order := Order{}
	if err = json.Unmarshal(jsonStr, &order); err != nil {
		t.Errorf("Can't parse json with error %v", err)
	}
}

func TestGetOrderList(t *testing.T) {
	w := httptest.NewRecorder()
	err := getOrderList(w, nil)
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

	list := OrderListResponse{}
	if err = json.Unmarshal(jsonString, &list); err != nil {
		t.Errorf("Can't parse json response with error %v", err)
	}
}
