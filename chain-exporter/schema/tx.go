package schema

import "time"

// TxIndex indexes all txs that occurred in Cosmos Network
type TxIndex struct {
	ID        int32     `json:"id" sql:",pk"`
	Height    int64     `json:"height" sql:",notnull"`
	ChainID   string    `json:"chain_id" sql:",notnull"`
	TxHash    string    `json:"tx_hash" sql:",notnull,unique"`
	Timestamp time.Time `json:"timestamp" sql:"default:now()"`
}

// TxCosmoshub3 has tx information
type TxCosmoshub3 struct {
	ID         int64  `json:"id" sql:",pk"`
	Height     int64  `json:"height"`
	TxHash     string `json:"tx_hash"  sql:",default:false,notnull,unique"`
	GasWanted  int64  `json:"gas_wanted" sql:"default:0"`
	GasUsed    int64  `json:"gas_used" sql:"default:0"`
	Messages   string `json:"messages" sql:"type:jsonb, default: '[]'::jsonb"`
	Fee        string `json:"fee" sql:"type:jsonb, default: '{}'::jsonb"`
	Signatures string `json:"signatures" sql:"type:jsonb, default: '[]'::jsonb"`
	Memo       string `json:"memo"`
	Logs       string `json:"logs" sql:"type:jsonb, default: '[]'::jsonb"`
	Events     string `json:"events" sql:"type:jsonb, default: '[]'::jsonb"`
	Time       string `json:"time"` // format that TxResponse returns
}
