package model

// Coin defines the structure for coin.
type Coin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
