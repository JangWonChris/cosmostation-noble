package schema

// TransactionInfo is a struct for database table
type TransactionInfo struct {
	ID         int64  `json:"id" sql:",pk"`
	Height     int64  `json:"height"`
	TxHash     string `json:"tx_hash"  sql:",default:false,notnull"`
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

// [
//     {
//         "type": "cosmos-sdk/MsgSend",
//         "value": {
//             "amount": [
//                 {
//                     "denom": "umuon",
//                     "amount": "50000000000"
//                 }
//             ],
//             "to_address": "cosmos1tt5pxlaazxql7qpg4cfx6adk0m76u3pavzq5yn",
//             "from_address": "cosmos1y3v65vhz5f93k8uk7vnz5v7yr7ks5gdce0udvs"
//         }
//     }
// ]
