package orderservice

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOrderById(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/v2/order/123", nil)
	err := GetOrderById(w, r)
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
