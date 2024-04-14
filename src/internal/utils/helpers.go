package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Numostanley/BankingAPI/internal/db"
	"github.com/Numostanley/BankingAPI/internal/models"
	"github.com/Numostanley/BankingAPI/internal/serializers"
)

func RespondWithError(w http.ResponseWriter, code int, data serializers.ResponseSerializer) {
	if code > 499 {
		log.Println("Responding with 5XX error: ", data.Error)
	}
	RespondWithJSON(w, code, data)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		return
	}
}

func OpenFile(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func CloseFile(file *os.File) {
	if file != nil {
		err := file.Close()
		if err != nil {
			return
		}
	}
}

func SeedAccount() {
	filename := "extras/accounts.json"

	file, err := OpenFile(filename)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer CloseFile(file)

	var accountParams []models.Account

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&accountParams)
	if err != nil {
		log.Println("error decoding json: ", err)
	}

	database := db.Database.DB
	for _, account := range accountParams {
		_, err := models.GetAccountByID(account.ID, database)

		if err != nil {
			err := account.CreateAccount(database)
			if err != nil {
				return
			}
			log.Println(
				account,
			)
		}
	}
}

func MockThirdPartyAPI() (bool, error) {
	time.Sleep(time.Second * 3)

	rand.New(rand.NewSource(time.Now().UnixNano()))
	result := rand.Intn(2)

	if result == 0 {
		return true, nil
	} else {
		return false, fmt.Errorf("third-party API request failed")
	}
}
