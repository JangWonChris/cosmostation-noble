package schema

import "time"

// TransactionIndex indexes all txs that occurred in Cosmos Network
type TransactionIndex struct {
	ID        int32     `json:"id" sql:",pk"`
	Height    int64     `json:"height" sql:",notnull"`
	TxHash    string    `json:"tx_hash" sql:",notnull,unique"`
	Timestamp time.Time `json:"timestamp" sql:"default:now()"`
	ChainID   string    `json:"chain_id" sql:",notnull"`
}
