package types

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Account struct {
	ID         int       `json:"id" db:"id"`
	FirstName  string    `json:"firstName" db:"first_name"`
	LastName   string    `json:"lastName" db:"last_name"`
	Patronymic string    `json:"patronymic" db:"patronymic"`
	Password   string    `json:"-"`
	Wallet     int64     `json:"wallet"`
	Balance    int64     `json:"balance" db:"balance"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}

func NewAccount(name, lastname, patronymic, password string) (*Account, error) {
	crPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		FirstName:  name,
		LastName:   lastname,
		Patronymic: patronymic,
		Password:   string(crPass),
		CreatedAt:  time.Now().UTC(),
	}, nil
}
