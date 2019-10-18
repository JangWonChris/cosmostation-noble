package models

import "time"

type Votes struct {
	Voter   string    `json:"voter"`
	Moniker string    `json:"moniker" sql:"default:null"`
	Option  string    `json:"option"`
	TxHash  string    `json:"tx_hash"`
	Time    time.Time `json:"time"`
}
