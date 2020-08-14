package model

import (
	"time"
)

const (
	Kava           = "kava"
	Binance        = "binancecoin"
	BNB            = "bnb"
	USDX           = "usdx"
	USDXStableCoin = "usdx-stablecoin"
)

type PriceStats struct {
	Price float64   `json:"price"`
	Time  time.Time `json:"time"`
}

type NetworkInfo struct {
	BondendTokensPercentChange24H float64              `json:"bonded_tokens_percent_change_24h"`
	BondedTokensStats             []*BondedTokensStats `json:"bonded_tokens_stats"`
}

type BondedTokensStats struct {
	BondedTokens float64   `json:"bonded_tokens"`
	BondedRatio  float64   `json:"bonded_ratio"`
	LastUpdated  time.Time `json:"last_updated"`
}
