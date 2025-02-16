package model

type Inventory struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinsReceived struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type CoinsSent struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
type CoinHistory struct {
	Received []CoinsReceived `json:"received"`
	Sent     []CoinsSent     `json:"sent"`
}
type Info struct {
	Coins       int          `json:"coins"`
	Inventory   []Inventory  `json:"inventory"`
	CoinHistory *CoinHistory `json:"coinHistory"`
}
