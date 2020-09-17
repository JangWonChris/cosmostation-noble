package model

// MoonPay wraps MoonPay api key
type MoonPay struct {
	APIKey string `json:"api_key"`
}

// ResultMoonPay wraps signautre
type ResultMoonPay struct {
	Signature string `json:"signature"`
}
