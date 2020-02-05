package schema

import "time"

// BlockCosmoshub3 has block information
type BlockCosmoshub3 struct {
	ID            int64     `json:"id" sql:",pk"`
	Height        int64     `json:"height"`
	Proposer      string    `json:"proposer"`
	BlockHash     string    `json:"block_hash" sql:",unique"`
	ParentHash    string    `json:"parent_hash" sql:",notnull"`
	NumPrecommits int64     `json:"num_pre_commits" sql:",notnull"`
	NumTxs        int64     `json:"num_txs" sql:"default:0"`
	TotalTxs      int64     `json:"total_txs" sql:"default:0"`
	Timestamp     time.Time `json:"timestamp" sql:"default:now()"`
}
