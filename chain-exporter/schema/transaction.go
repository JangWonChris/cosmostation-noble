package schema

// Transaction has tx information.
type Transaction struct {
	ID         int64  `json:"id" sql:",pk"`
	ChainID    string `json:"chain_id" sql:",notnull"`
	Height     int64  `json:"height"`
	Code       uint32 `json:"code" sql:"default:0"`
	TxHash     string `json:"tx_hash"  sql:",default:false,notnull,unique"`
	GasWanted  int64  `json:"gas_wanted" sql:"default:0"`
	GasUsed    int64  `json:"gas_used" sql:"default:0"`
	Messages   string `json:"messages" sql:"type:jsonb, default: '[]'::jsonb"`
	Fee        string `json:"fee" sql:"type:jsonb, default: '{}'::jsonb"`
	Signatures string `json:"signatures" sql:"type:jsonb, default: '[]'::jsonb"`
	Memo       string `json:"memo"`
	Logs       string `json:"logs" sql:"type:jsonb, default: '[]'::jsonb"`
	Timestamp  string `json:"timestamp" sql:"default:now()"` // format that TxResponse returns
}

// NewTransaction returns a new Transaction.
func NewTransaction(t Transaction) *Transaction {
	return &Transaction{
		ChainID:    t.ChainID,
		Height:     t.Height,
		Code:       t.Code,
		TxHash:     t.TxHash,
		GasWanted:  t.GasWanted,
		GasUsed:    t.GasUsed,
		Messages:   t.Messages,
		Fee:        t.Fee,
		Signatures: t.Signatures,
		Memo:       t.Memo,
		Logs:       t.Logs,
		Timestamp:  t.Timestamp,
	}
}
