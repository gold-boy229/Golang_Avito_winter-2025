package model

type InfoResponse struct {
	Coins       uint            `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity uint   `json:"quantity"`
}

type CoinHistory struct {
	Received []ReceivedOperation `json:"received"`
	Sent     []SentOperation     `json:"sent"`
}

type ReceivedOperation struct {
	FromUser string `json:"fromUser"`
	Amount   uint   `json:"amount"`
}

type SentOperation struct {
	ToUser string `json:"toUser"`
	Amount uint   `json:"amount"`
}

////

type AuthResponse struct {
	Token string `json:"token"`
}

////

type ErrorResponse struct {
	Errors string `json:"errors"`
}
