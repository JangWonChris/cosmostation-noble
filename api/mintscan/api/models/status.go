package models

import "time"

type Pool struct {
	NotBondedTokens string `json:"not_bonded_tokens"`
	BondedTokens    string `json:"bonded_tokens"`
}

type ResultStatus struct {
	ChainID                string    `json:"chain_id"`
	BlockHeight            int64     `json:"block_height"`
	BlockTime              float64   `json:"block_time"`
	TotalTxsNum            int64     `json:"total_txs_num"`
	TotalValidatorNum      int       `json:"total_validator_num"`
	UnjailedValidatorNum   int       `json:"unjailed_validator_num"`
	JailedValidatorNum     int       `json:"jailed_validator_num"`
	TotalCirculatingTokens float64   `json:"total_circulating_tokens"`
	BondedTokens           float64   `json:"bonded_tokens"`
	NotBondedTokens        float64   `json:"not_bonded_tokens"`
	Time                   time.Time `json:"time"`
}
