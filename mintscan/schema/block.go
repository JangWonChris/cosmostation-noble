package schema

import "time"

// Block has block information.
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

// NewBlock returns new Block.
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
