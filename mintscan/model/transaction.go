package model

import (
	"encoding/json"

	"github.com/cosmostation/mintscan-backend-library/db/schema"
)

// TxData defines the structure for transction data list.
type TxData struct {
	Txs []string `json:"txs"`
}

// TxList defines the structure for transaction list.
type TxList struct {
	TxHash []string `json:"tx_list"`
}

// Message defines the structure for transaction message.
type Message struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

// Fee defines the structure for transaction fee.
type Fee struct {
	Gas    string `json:"gas,omitempty"`
	Amount []struct {
		Amount string `json:"amount,omitempty"`
		Denom  string `json:"denom,omitempty"`
	} `json:"amount,omitempty"`
}

// Event defines the structure for transaction event.
type Event struct {
	Type       string `json:"type"`
	Attributes []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"attributes"`
}

// Log defines the structure for transaction log.
type Log struct {
	MsgIndex int     `json:"msg_index"`
	Log      string  `json:"log"`
	Events   []Event `json:"events"`
}

// ParseTransaction receives single transaction from database and return it after unmarshal them.
func ParseTransaction(tx schema.Transaction) (result *ResultTx, err error) {

	// msgsBz, err := mintscancodec.AppCodec.UnmarshalJSON([]byte(tx.Messages), tx.GetBody())
	jsonRaws := make([]json.RawMessage, 0)
	if err := json.Unmarshal([]byte(tx.Messages), &jsonRaws); err != nil {
		return &ResultTx{}, err
	}

	// var sdkmsgs []sdktypes.Msg
	// for _, raw := range jsonRaws {
	// 	err = mintscancodec.AppCodec.UnmarshalJSON(raw, sdkmsgs)
	// 	err = codectypes.NewAnyWithValue(raw)
	// }

	// fmt.Println(sdkmsgs)

	// // fmt.Println(txBody.GetMessages())
	// for _, msg := range sdkmsgs {
	// 	_ = msg
	// 	// fmt.Println("msg typeurl:", msg.TypeUrl)
	// 	// a, ok := msg.GetCachedValue().(sdktypes.Msg)
	// 	// if !ok {
	// 	// 	fmt.Println("not sdktypes.Msg")
	// 	// }
	// 	// fmt.Println(a.String())
	// 	// switch aType := a.(type) {
	// 	// case *stakingtypes.MsgCreateValidator:
	// 	// 	fmt.Println("aType:", aType)
	// 	// 	fmt.Println("moniker:", aType.Description.Moniker)
	// 	// }
	// }

	var jsonraw json.RawMessage
	json.Unmarshal([]byte(tx.Messages), &jsonraw)

	// fmt.Println(txBody.GetMemo())

	msgs := make([]Message, 0)
	err = json.Unmarshal([]byte(tx.Messages), &msgs)
	if err != nil {
		return &ResultTx{}, err
	}

	var fee *Fee
	err = json.Unmarshal([]byte(tx.Fee), &fee)
	if err != nil {
		return &ResultTx{}, err
	}

	var logs []Log
	err = json.Unmarshal([]byte(tx.Logs), &logs)
	if err != nil {
		return &ResultTx{}, err
	}

	result = &ResultTx{
		ID:        tx.ID,
		Height:    tx.Height,
		TxHash:    tx.TxHash,
		Logs:      logs,
		GasWanted: tx.GasWanted,
		GasUsed:   tx.GasUsed,
		Msgs:      msgs,
		Fee:       fee,
		Memo:      tx.Memo,
		Timestamp: tx.Timestamp,
	}

	return result, nil
}

// ParseTransactions receives result transactions from database and return them after unmarshal them.
func ParseTransactions(txs []schema.Transaction) (result []ResultTx, err error) {

	for _, tx := range txs {

		msgs := make([]Message, 0)
		err = json.Unmarshal([]byte(tx.Messages), &msgs)
		if err != nil {
			return []ResultTx{}, err
		}

		var fee *Fee
		err = json.Unmarshal([]byte(tx.Fee), &fee)
		if err != nil {
			return []ResultTx{}, err
		}

		var logs []Log
		err = json.Unmarshal([]byte(tx.Logs), &logs)
		if err != nil {
			return []ResultTx{}, err
		}

		tx := &ResultTx{
			ID:        tx.ID,
			Height:    tx.Height,
			TxHash:    tx.TxHash,
			Logs:      logs,
			GasWanted: tx.GasWanted,
			GasUsed:   tx.GasUsed,
			Msgs:      msgs,
			Fee:       fee,
			Memo:      tx.Memo,
			Timestamp: tx.Timestamp,
		}

		result = append(result, *tx)
	}

	return result, nil
}
