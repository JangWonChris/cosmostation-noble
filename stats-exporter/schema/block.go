package schema

import "time"

// Block defines the structure for block information.
type Block struct {
	ID            int64     `json:"id" sql:",pk"`
	ChainID       string    `json:"chain_id" sql:",notnull"`
	Height        int64     `json:"height"`
	Proposer      string    `json:"proposer"`
	BlockHash     string    `json:"block_hash" sql:",unique"`
	ParentHash    string    `json:"parent_hash" sql:",notnull"`
	NumSignatures int64     `json:"num_signatures" sql:",notnull"`
	NumTxs        int64     `json:"num_txs" sql:"default:0"`
	Timestamp     time.Time `json:"timestamp" sql:"default:now()"`
}

// NewBlock returns a new Block.
// Note that TotalTxs param is removed indefinitely from block header in new tendermint version v0.33.+ (Marko).
func NewBlock(b Block) *Block {
	return &Block{
		ChainID:       b.ChainID,
		Height:        b.Height,
		Proposer:      b.Proposer,
		BlockHash:     b.BlockHash,
		ParentHash:    b.ParentHash,
		NumSignatures: b.NumSignatures,
		NumTxs:        b.NumTxs,
		Timestamp:     b.Timestamp,
	}
}
