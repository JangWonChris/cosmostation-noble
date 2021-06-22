package model

import (
	"encoding/json"

	// sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmostation/cosmostation-cosmos/app"
	mdschema "github.com/cosmostation/mintscan-database/schema"
)

// TxData defines the structure for transction data list.
// type TxData struct {
// 	Txs []json.RawMessage `json:"txs"`
// }

// TxList defines the structure for transaction list.
type TxList struct {
	TxHash []string `json:"tx_list"`
}

// Message defines the structure for transaction message.
// type Message struct {
// 	Type  string          `json:"type"`
// 	Value json.RawMessage `json:"value"`
// }

// Fee defines the structure for transaction fee.
// type Fee struct {
// 	Gas    string `json:"gas,omitempty"`
// 	Amount []struct {
// 		Amount string `json:"amount,omitempty"`
// 		Denom  string `json:"denom,omitempty"`
// 	} `json:"amount,omitempty"`
// }

// Event defines the structure for transaction event.
// type Event struct {
// 	Type       string `json:"type"`
// 	Attributes []struct {
// 		Key   string `json:"key"`
// 		Value string `json:"value"`
// 	} `json:"attributes"`
// }

// Log defines the structure for transaction log.
// type Log struct {
// 	MsgIndex int     `json:"msg_index"`
// 	Log      string  `json:"log"`
// 	Events   []Event `json:"events"`
// }

// 모바일 히스토리 용
type OldMobileTransactionHistory struct {
	ID         int64  `json:"id"`
	ChainID    string `json:"chain_id"`
	Height     int64  `json:"height"`
	Code       uint32 `json:"code"`
	TxHash     string `json:"tx_hash"`
	Messages   string `json:"messages"`
	Fee        string `json:"fee"`
	Signatures string `json:"signatures"`
	Memo       string `json:"memo"`
	GasWanted  int64  `json:"gas_wanted"`
	GasUsed    int64  `json:"gas_used"`
	Logs       string `json:"logs"`
	RawLog     string `json:"raw_log"`
	Timestamp  string `json:"timestamp"`
}

// ParseTransaction receives single transaction from database and return it after unmarshal them.
func OldMobileParseTransaction(a *app.App, tx mdschema.Transaction) (result *OldMobileTransactionHistory) {
	if tx.ID == 0 {
		return
	}
	// var txResp sdktypes.TxResponse

	// custom.AppCodec.UnmarshalJSON(tx.Chunk, &txResp)

	result.ChainID = a.ChainNumMap[tx.ChainInfoID]
	result.Code = tx.Code
	result.TxHash = tx.Hash
	result.Timestamp = tx.Timestamp.String()
	// result.Height = txResp.Height

	return
}

// ParseTransaction receives single transaction from database and return it after unmarshal them.
func ParseTransaction(a *app.App, tx mdschema.Transaction) (result *ResultTx) {
	if tx.ID == 0 {
		return
	}
	var jsonRaws json.RawMessage

	jsonRaws = tx.Chunk

	header := ResultTxHeader{
		ID:        tx.ID,
		ChainID:   a.ChainNumMap[tx.ChainInfoID],
		BlockID:   tx.BlockID,
		Timestamp: tx.Timestamp.String(),
	}

	result = &ResultTx{
		ResultTxHeader: header,
		Data:           jsonRaws,
	}

	return result
}

// ParseTransactions receives result transactions from database and return them after unmarshal them.
func ParseTransactions(a *app.App, txs []mdschema.Transaction) (results []*ResultTx) {
	for i := range txs {
		results = append(results, ParseTransaction(a, txs[i]))
	}
	return results
}
