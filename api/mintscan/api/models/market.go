package models

import (
	"time"
)

/*
	Legacy Code
*/

type MarketInfo struct {
	Price            float64       `json:"price"`
	Currency         string        `json:"currency"`
	PercentChange1H  float64       `json:"percent_change_1h"`
	PercentChange24H float64       `json:"percent_change_24h"`
	LastUpdated      time.Time     `json:"last_updated"`
	PriceStats       []*PriceStats `json:"price_stats"`
}

type PriceStats struct {
	Price float64   `json:"price"`
	Time  time.Time `json:"time"`
}

type NetworkInfo struct {
	BondendTokensPercentChange24H float64              `json:"bonded_tokens_percent_change_24h"`
	BondedTokensStats             []*BondedTokensStats `json:"bonded_tokens_stats"`
}

type BondedTokensStats struct {
	BondedTokens int64     `json:"bonded_tokens"`
	BondedRatio  float64   `json:"bonded_ratio"`
	LastUpdated  time.Time `json:"last_updated"`
}
