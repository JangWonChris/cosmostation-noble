package types

import "time"

// Database struct
type BlockInfo struct {
	ID        int64     `json:"id" sql:",pk"`
	BlockHash string    `json:"block_hash"`
	Height    int64     `json:"height"`
	Proposer  string    `json:"proposer"`
	TotalTxs  int64     `json:"total_txs" sql:"default:0"`
	NumTxs    int64     `json:"num_txs" sql:"default:0"`
	Time      time.Time `json:"time"`
}
