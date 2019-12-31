package models

// Coin is a struct for REST API
type Coin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
