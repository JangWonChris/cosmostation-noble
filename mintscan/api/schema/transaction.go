package schema

// TransactionInfo is a struct for database table
type TransactionInfo struct {
	ID         int64  `json:"id" sql:",pk"`
	Height     int64  `json:"height"`
	TxHash     string `json:"tx_hash"  sql:",default:false,notnull"`
	GasWanted  int64  `json:"gas_wanted" sql:"default:0"`
	GasUsed    int64  `json:"gas_used" sql:"default:0"`
	Messages   string `json:"messages" sql:"default: '[]'::jsonb"`
	Fee        string `json:"fee" sql:"default: '{}'::jsonb"`
	Signatures string `json:"signatures" sql:"default: '[]'::jsonb"`
	Memo       string `json:"memo"`
	Logs       string `json:"logs" sql:"default: '[]'::jsonb"`
	Events     string `json:"events" sql:"default: '[]'::jsonb"`
	Time       string `json:"time"` // format that TxResponse returns
}
