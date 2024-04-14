package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Numostanley/BankingAPI/internal/handlers"
)

func TestCreateTransactionHandler(t *testing.T) {
	requestBody := []byte(`{"account_id": "123", "reference": "ref123", "amount": 100.0}`)
	req, err := http.NewRequest("POST", "/v1/payments/", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handlers.CreateTransactionHandler(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("error parsing response body: %v", err)
	}

	expectedKeys := []string{"success", "message", "data", "error"}
	for _, key := range expectedKeys {
		if _, ok := response[key]; !ok {
			t.Errorf("response missing expected key: %s", key)
		}
	}
}

func TestGetTransactionHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/payments/ref123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handlers.GetTransactionHandler(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("error parsing response body: %v", err)
	}

	expectedKeys := []string{"success", "message", "data", "error"}
	for _, key := range expectedKeys {
		if _, ok := response[key]; !ok {
			t.Errorf("response missing expected key: %s", key)
		}
	}
}
