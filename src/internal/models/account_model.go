package models

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Account struct {
	*gorm.Model
	ID          uuid.UUID       `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	FirstName   string          `json:"first_name" gorm:"type:text;not null"`
	LastName    string          `json:"last_name" gorm:"type:text;not null"`
	Email       string          `json:"email" gorm:"type:text;not null;uniqueIndex"`
	PhoneNumber string          `json:"phone_number" gorm:"type:text;not null;uniqueIndex"`
	Balance     decimal.Decimal `json:"balance" gorm:"type:decimal(10,2);"`
}

func (user *Account) CreateAccountID() {
	user.ID = uuid.New()
}

func (user *Account) CreateAccount(db *gorm.DB) error {
	err := UserValidation(user)
	if err != nil {
		return err
	}

	emailExists, _ := AccountExistsByEmail(db, user.Email)
	if emailExists {
		err = fmt.Errorf("email already exists")
		return err
	}

	phoneNumberExists, _ := AccountExistsByPhone(db, user.PhoneNumber)
	if phoneNumberExists {
		err = fmt.Errorf("phone already exists")
		return err
	}

	user.CreateAccountID()
	newUser := db.Create(&user)
	if newUser.Error != nil {
		err = fmt.Errorf("error creating user: %v", newUser.Error)
		return err
	}
	return nil
}

func (user *Account) GetFullName() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}

func AccountExistsByEmail(db *gorm.DB, emailToCheck string) (bool, error) {
	var user Account
	result := db.Where("email = ?", emailToCheck).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	} else if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func AccountExistsByPhone(db *gorm.DB, phoneToCheck string) (bool, error) {
	var user Account
	result := db.Where("phone_number = ?", phoneToCheck).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	} else if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func UserValidation(user *Account) error {
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	if user.PhoneNumber == "" {
		return fmt.Errorf("phone_number is required")
	}
	if user.FirstName == "" {
		return fmt.Errorf("first_name is required")
	}
	if user.LastName == "" {
		return fmt.Errorf("last_name is required")
	}
	return nil
}

func GetAccountByID(userID uuid.UUID, db *gorm.DB) (*Account, error) {
	user := Account{ID: userID}
	fetchedUser := db.Where("id = ?", userID).First(&user)

	if fetchedUser.Error != nil {
		return nil, fmt.Errorf("error returning user %s", fetchedUser.Error)
	}
	return &user, nil
}
