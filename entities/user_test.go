package entities

import (
	"testing"
)

func TestCheckPassword(t *testing.T) {
	user := User{Password: "secret"}

	tests := []struct {
		providedPassword string
		expectError      bool
	}{
		{"secret", false},
		{"not_secret", true},
		{"secret_not", true},
		{"not_secret_not", true},
	}

	for _, test := range tests {
		err := user.CheckPassword(test.providedPassword)
		if (err != nil) != test.expectError {
			t.Errorf("user.Password = %q. CheckPassword(%q) = %v; expectError = %v",
				user.Password, test.providedPassword, err, test.expectError)
		}
	}
}

func TestCanSpendCoins(t *testing.T) {
	user := User{Balance: 10}

	tests := []struct {
		coinsToSpend uint
		expectError  bool
	}{
		{0, false},
		{10, false},
		{11, true},
	}

	for _, test := range tests {
		err := user.CanSpendCoins(test.coinsToSpend)
		if (err != nil) != test.expectError {
			t.Errorf("user.Balance = %v. CanSpendCoins(%v) = %v; expectError = %v",
				user.Balance, test.coinsToSpend, err, test.expectError)
		}
	}
}

func TestIsSameUser(t *testing.T) {
	user := User{Username: "nick"}

	tests := []struct {
		providedUsername string
		expectError      bool
	}{
		{"nick", true},
		{"Nick", false},
		{"NICK", false},
		{"ni", false},
		{"ck", false},
		{"pref_nick", false},
		{"nick_suff", false},
		{"pref_nick_suff", false},
		{"mouse", false},
		{"", false},
	}

	for _, test := range tests {
		err := user.IsSameUser(test.providedUsername)
		if (err != nil) != test.expectError {
			t.Errorf("user.Username = %q. IsSameUser(%q) = %v; expectError = %v",
				user.Username, test.providedUsername, err, test.expectError)
		}
	}
}
