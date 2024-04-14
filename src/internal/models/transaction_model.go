package models

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var TransactionStatus = struct {
	Success string
	Fail    string
	Pending string
}{
	Success: "successful",
	Fail:    "failed",
	Pending: "pending",
}

type Transaction struct {
	*gorm.Model
	AccountID uuid.UUID       `json:"account_id" gorm:"type:uuid;not null"`
	Reference string          `json:"reference" gorm:"type:text;not null;uniqueIndex"`
	Amount    decimal.Decimal `json:"amount" gorm:"type:decimal(10,2); not null"`
	Status    string          `json:"status" gorm:"type:text;not null"`
}

func (transaction *Transaction) CreateTransaction(db *gorm.DB) error {
	err := TransactionValidation(transaction)
	if err != nil {
		return err
	}

	referenceExists, _ := TransactionReferenceExists(db, transaction.Reference)
	if referenceExists {
		err = fmt.Errorf("transaction with this reference already exists")
		return err
	}

	transaction.Status = TransactionStatus.Pending

	newTransaction := db.Create(&transaction)
	if newTransaction.Error != nil {
		err := fmt.Errorf("error creating transaction: %v", newTransaction.Error)
		return err
	}
	return nil
}

func TransactionReferenceExists(db *gorm.DB, reference string) (bool, error) {
	var transaction Transaction
	result := db.Where("reference = ?", reference).First(&transaction)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	} else if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func TransactionValidation(transaction *Transaction) error {
	if transaction.AccountID == uuid.Nil {
		return fmt.Errorf("account_id is required")
	}
	if transaction.Reference == "" {
		return fmt.Errorf("referrence is required")
	}
	if transaction.Amount == decimal.NewFromFloat(0.0) {
		return fmt.Errorf("valid amount is required")
	}
	return nil
}

func GetTransactiontByReference(ref string, db *gorm.DB) (*Transaction, error) {
	transaction := Transaction{Reference: ref}
	fetchedTransaction := db.Where("reference = ?", ref).First(&transaction)

	if fetchedTransaction.Error != nil {
		return nil, fmt.Errorf("%s", fetchedTransaction.Error)
	}
	return &transaction, nil
}
