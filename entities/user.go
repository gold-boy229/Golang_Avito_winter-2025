package entities

import (
	"fmt"
)

type User struct {
	Id       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique" json:"username"`
	Password string `json:"password"`
	Balance  uint   `gorm:"default:1000" json:"balance"`
}

func (user *User) CheckPassword(providedPassword string) (err error) {
	err = nil
	if user.Password != providedPassword {
		err = fmt.Errorf("invalid credentials")
	}
	return
}

func (user *User) CanSpendCoins(amount uint) (err error) {
	err = nil
	if user.Balance < amount {
		err = fmt.Errorf("not enough coins. Have: %d, want to spend: %d", user.Balance, amount)
	}
	return
}

func (user *User) IsAnotherUser(AnotherUser *User) (err error) {
	err = nil
	if user.Username == AnotherUser.Username {
		err = fmt.Errorf("it's the same person")
	}
	return
}
