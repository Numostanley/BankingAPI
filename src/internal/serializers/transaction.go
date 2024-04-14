package serializers

import (
	"time"

	"github.com/Numostanley/BankingAPI/internal/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionSerializer struct {
	AccountID uuid.UUID       `json:"account_id"`
	Reference string          `json:"reference"`
	Amount    decimal.Decimal `json:"amount"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func (res *TransactionSerializer) GetUserResponse(transaction *models.Transaction) TransactionSerializer {
	response := TransactionSerializer{
		AccountID: transaction.AccountID,
		Reference: transaction.Reference,
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		CreatedAt: transaction.CreatedAt,
		UpdatedAt: transaction.UpdatedAt,
	}

	return response
}
