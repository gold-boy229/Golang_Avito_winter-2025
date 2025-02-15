package entities

import (
	"fmt"
)

type User struct {
	Id       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Balance  uint   `json:"balance" gorm:"default:1000"`
}

func (user *User) CheckPassword(providedPassword string) error {
	if user.Password != providedPassword {
		return fmt.Errorf("invalid credentials")
	}
	return nil
}

func (user *User) CanSpendCoins(amount uint) error {
	if user.Balance < amount {
		return fmt.Errorf("not enough coins. Have: %d, want to spend: %d", user.Balance, amount)
	}
	return nil
}

func (user *User) IsSameUser(anotherUsername string) error {
	if user.Username == anotherUsername {
		return fmt.Errorf("it's the same person")
	}
	return nil
}
