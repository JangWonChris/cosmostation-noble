package schema

import "time"

// TODO: BlockCosmoshub3, TransactionCosmoshub3, Transaction Index

// Blockinfo has block information
type BlockInfo struct {
	ID            int64     `json:"id" sql:",pk"`
	Height        int64     `json:"height"`
	BlockHash     string    `json:"block_hash" sql:",unique"`
	ParentHash    string    `json:"parent_hash" sql:",notnull"`
	Proposer      string    `json:"proposer"`
	NumPrecommits int64     `json:"num_pre_commits" sql:",notnull"`
	NumTxs        int64     `json:"num_txs" sql:"default:0"`
	TotalTxs      int64     `json:"total_txs" sql:"default:0"`
	Time          time.Time `json:"time"`
}
