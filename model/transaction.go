package model

import (
	"encoding/json"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdktypestx "github.com/cosmos/cosmos-sdk/types/tx"

	"github.com/cosmostation/cosmostation-cosmos/app"
	"github.com/cosmostation/cosmostation-cosmos/custom"
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
	Code       uint32 `json:"code"` // grpc 이후 안씀
	TxHash     string `json:"tx_hash"`
	Messages   string `json:"messages"`   // grpc 이후 사용
	Fee        string `json:"fee"`        // grpc이후 안씀
	Signatures string `json:"signatures"` // 안씀
	Memo       string `json:"memo"`       // grpc 이후 안씀
	GasWanted  int64  `json:"gas_wanted"`
	GasUsed    int64  `json:"gas_used"`
	Logs       string `json:"logs"`
	RawLog     string `json:"raw_log"`
	Timestamp  string `json:"timestamp"`
}

// OldMobileParseTransaction tx를 민트스캔 응답구조에 맞게 파싱하여 리턴
func OldMobileParseTransaction(a *app.App, tx mdschema.Transaction) (result *OldMobileTransactionHistory) {
	if tx.ID == 0 {
		return
	}
	var txResp sdktypes.TxResponse
	unmarshal := custom.AppCodec.UnmarshalJSON
	unmarshal(tx.Chunk, &txResp)

	txI := txResp.GetTx()
	sdkTx, ok := txI.(*sdktypestx.Tx)
	if !ok {
		return
	}
	msgs := sdkTx.GetBody().GetMessages()
	jsonRaws := make([]json.RawMessage, len(msgs), len(msgs))
	var err error
	for i, msg := range msgs {
		jsonRaws[i], err = custom.AppCodec.MarshalJSON(msg)
		if err != nil {
			return
		}
	}
	msgsBz, err := json.Marshal(jsonRaws)
	if err != nil {
		return
	}

	feeBz, err := custom.AppCodec.MarshalJSON(sdkTx.GetAuthInfo().GetFee())
	if err != nil {
		return
	}

	temp := OldMobileTransactionHistory{
		ChainID:   a.ChainNumMap[tx.ChainInfoID],
		Height:    txResp.Height,
		Code:      txResp.Code,
		TxHash:    txResp.TxHash,
		Messages:  string(msgsBz),
		Fee:       string(feeBz),
		GasWanted: txResp.GasWanted,
		GasUsed:   txResp.GasUsed,
		Logs:      txResp.Logs.String(),
		RawLog:    txResp.RawLog,
		Memo:      sdkTx.GetBody().GetMemo(),
		Timestamp: txResp.Timestamp,
	}

	result = &temp

	return
}

// OldMobileParseTransactions []txs를 파싱하여 리턴
func OldMobileParseTransactions(a *app.App, txs []mdschema.Transaction) (results []*OldMobileTransactionHistory) {
	for i := range txs {
		if txs[i].ChainInfoID == 4 {
			results = append(results, OldMobileParseTransaction(a, txs[i]))
		}
	}
	return results
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
