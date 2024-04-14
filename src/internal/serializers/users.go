package serializers

import (
	"time"

	"github.com/Numostanley/BankingAPI/internal/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type UserDetailSerializer struct {
	ID          uuid.UUID       `json:"id"`
	FullName    string          `json:"full_name"`
	Email       string          `json:"email"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	PhoneNumber string          `json:"phone_number"`
	Balance     decimal.Decimal `json:"balance"`
}

func (userRes *UserDetailSerializer) GetUserResponse(user *models.Account) UserDetailSerializer {
	response := UserDetailSerializer{
		ID:          user.ID,
		FullName:    user.GetFullName(),
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		PhoneNumber: user.PhoneNumber,
		Balance:     user.Balance,
	}

	return response
}
