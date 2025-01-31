package models

import (
	"errors"
)

type User struct {
	ID      int     `json:"id"`
	Balance float64 `json:"balance"`
}

func (u *User) ValidateBalance(amount float64) error {
	if u.Balance < amount {
		return errors.New("saldo insuficiente")
	}
	return nil
}
