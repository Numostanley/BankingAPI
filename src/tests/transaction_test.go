package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Numostanley/BankingAPI/internal/db"
	"github.com/Numostanley/BankingAPI/internal/handlers"
	"github.com/Numostanley/BankingAPI/internal/models"
	"github.com/Numostanley/BankingAPI/internal/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTransactionHandlers(t *testing.T) {
	testDB := createTestDB()
	defer func() {
		db.Database.DB = nil
		if sqlDB, err := testDB.DB(); err == nil {
			err := sqlDB.Close()
			if err != nil {
				return
			}
		}
	}()

	var user1 models.Account
	if err := db.Database.DB.Model(models.Account{}).First(&user1).Error; err != nil {
		log.Println("Error retrieving user1:", err)
		return
	}
	accountID := user1.ID

	// Test CreateTransactionHandler
	requestBody1 := []byte(fmt.Sprintf(`{"account_id": "%s", "reference": "ref123", "amount": 100.0}`, accountID))
	req1, err := http.NewRequest("POST", "/v1/payments/", bytes.NewBuffer(requestBody1))
	if err != nil {
		t.Fatal(err)
	}

	rr1 := httptest.NewRecorder()
	handlers.CreateTransactionHandler(rr1, req1)
	if status := rr1.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response1 map[string]interface{}
	if err := json.Unmarshal(rr1.Body.Bytes(), &response1); err != nil {
		t.Errorf("error parsing response body: %v", err)
	}
	expectedKeys := []string{"success", "message", "data", "error"}
	for _, key := range expectedKeys {
		if _, ok := response1[key]; !ok {
			t.Errorf("response missing expected key: %s", key)
		}
	}
	assert.Equal(t, "Transaction created successfully", response1["message"])

	// Test for duplicate transaction reference
	requestBody2 := []byte(fmt.Sprintf(`{"account_id": "%s", "reference": "ref123", "amount": 100.0}`, accountID))
	req2, err := http.NewRequest("POST", "/v1/payments/", bytes.NewBuffer(requestBody2))
	if err != nil {
		t.Fatal(err)
	}

	rr2 := httptest.NewRecorder()
	handlers.CreateTransactionHandler(rr2, req2)
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var response2 map[string]interface{}
	if err := json.Unmarshal(rr2.Body.Bytes(), &response2); err != nil {
		t.Errorf("error parsing response body: %v", err)
	}
	assert.Equal(t, "Error creating transaction: transaction with this reference already exists", response2["error"])

	// Test GetTransactionHandler
	req3, err := http.NewRequest("GET", "/v1/payments/ref123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr3 := httptest.NewRecorder()
	handlers.GetTransactionHandler(rr3, req3)
	if status := rr3.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response3 map[string]interface{}
	if err := json.Unmarshal(rr3.Body.Bytes(), &response3); err != nil {
		t.Errorf("error parsing response body: %v", err)
	}
	expectedKeys2 := []string{"success", "message", "data", "error"}
	for _, key := range expectedKeys2 {
		if _, ok := response3[key]; !ok {
			t.Errorf("response missing expected key: %s", key)
		}
	}

	// Test for insufficient account balance
	requestBody4 := []byte(fmt.Sprintf(`{"account_id": "%s", "reference": "ref1235", "amount": 10000.0}`, accountID))
	req4, err := http.NewRequest("POST", "/v1/payments/", bytes.NewBuffer(requestBody4))
	if err != nil {
		t.Fatal(err)
	}

	rr4 := httptest.NewRecorder()
	handlers.CreateTransactionHandler(rr4, req4)
	if status := rr4.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var response4 map[string]interface{}
	if err := json.Unmarshal(rr4.Body.Bytes(), &response4); err != nil {
		t.Errorf("error parsing response body: %v", err)
	}
	assert.Equal(t, "Insfficient account balance!!!", response4["error"])

	// Test for invalid transaction reference
	req5, err := http.NewRequest("GET", "/v1/payments/ref12398", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr5 := httptest.NewRecorder()
	handlers.GetTransactionHandler(rr5, req5)
	if status := rr5.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var response5 map[string]interface{}
	if err := json.Unmarshal(rr5.Body.Bytes(), &response5); err != nil {
		t.Errorf("error parsing response body: %v", err)
	}
	assert.Equal(t, "Error retreiving transaction: record not found", response5["error"])

}

func createTestDB() *gorm.DB {
	testDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: &db.CustomLogger{},
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	if err := testDB.AutoMigrate(
		&models.Account{},
		&models.Transaction{},
	); err != nil {
		panic(err)
	}
	log.Println("Migrations Complete")

	db.Database.DB = testDB
	utils.SeedAccount("../extras/accounts.json")
	return testDB
}
