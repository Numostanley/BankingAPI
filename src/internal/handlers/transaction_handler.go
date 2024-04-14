package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Numostanley/BankingAPI/internal/db"
	"github.com/Numostanley/BankingAPI/internal/models"
	"github.com/Numostanley/BankingAPI/internal/serializers"
	"github.com/Numostanley/BankingAPI/internal/utils"
)

func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	data := serializers.ResponseSerializer{
		Success: true,
		Message: "",
		Data:    struct{}{},
		Error:   "",
	}

	var newTransaction models.Transaction

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newTransaction)
	if err != nil {
		data.Error = fmt.Sprintf("Error parsing JSON: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}
	database := db.Database.DB

	account, err := models.GetAccountByID(newTransaction.AccountID, database)
	if err != nil {
		data.Error = fmt.Sprintf("Invalid account ID: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	if newTransaction.Amount.GreaterThan(account.Balance) {
		data.Error = "Insfficient account balance!!!"
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	err = newTransaction.CreateTransaction(database)

	if err != nil {
		data.Error = fmt.Sprintf("Error creating transaction: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	success, err := utils.MockThirdPartyAPI()
	if !success {
		newTransaction.Status = models.TransactionStatus.Fail
		database.Save(&newTransaction)

		data.Error = fmt.Sprintf("Error creating transaction: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	newTransaction.Status = models.TransactionStatus.Success
	database.Save(&newTransaction)

	newBalance := account.Balance.Sub(newTransaction.Amount)
	account.Balance = newBalance
	database.Save(&account)

	response := serializers.TransactionSerializer{}
	data.Data = response.GetUserResponse(&newTransaction)
	data.Message = "Transaction created successfully"

	utils.RespondWithJSON(w, 200, data)
}

func GetTransactionHandler(w http.ResponseWriter, r *http.Request) {
	data := serializers.ResponseSerializer{
		Success: true,
		Message: "",
		Data:    struct{}{},
		Error:   "",
	}

	reference := r.URL.Path[len("/v1/payments/"):]
	database := db.Database.DB

	transaction, err := models.GetTransactiontByReference(reference, database)
	if err != nil {
		data.Error = fmt.Sprintf("Error retreiving transaction: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	response := serializers.TransactionSerializer{}
	data.Data = response.GetUserResponse(transaction)
	data.Message = "Transaction retreived successfully"

	utils.RespondWithJSON(w, 200, data)
}
