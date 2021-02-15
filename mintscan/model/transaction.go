package model

import (
	"encoding/json"
	"fmt"

	"github.com/cosmostation/cosmostation-cosmos/chain-config/custom"
	codec "github.com/cosmostation/cosmostation-cosmos/chain-config/custom"
	"github.com/cosmostation/mintscan-backend-library/db/schema"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdktypesTx "github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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
		var txb sdktypesTx.TxBody
		var any codectypes.Any
		var jrs []json.RawMessage
		_, _, _ = txb, any, jrs
		err = json.Unmarshal([]byte(tx.Messages), &jrs)
		if err != nil {
			fmt.Println("json unm", err)
			return
		}
		for i := range jrs {
			var m sdktypes.Msg
			err = codec.AppCodec.UnmarshalJSON(jrs[i], &any)
			if err != nil {
				fmt.Println("proto unm", err)
			}
			// b, err := json.Marshal(any)
			// if err != nil {
			// 	fmt.Println("marshal err", err)
			// }
			fmt.Println("nalgut", any)
			fmt.Println("typeurl :", any.TypeUrl)
			fmt.Println("anystring :", any.String())
			// err = codectypes.UnpackInterfaces(any, custom.EncodingConfig.InterfaceRegistry.Resolve(any.TypeUrl))
			p, err := custom.EncodingConfig.InterfaceRegistry.Resolve(any.TypeUrl)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("p : ", p)
			a, err := codectypes.NewAnyWithValue(p)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("a cached :", a.GetCachedValue())
			fmt.Println("unpack getcached value:", any.GetCachedValue())
			switch any.GetCachedValue().(type) {
			case *banktypes.MsgSend:
				fmt.Println("msgsend")
			default:
				fmt.Printf("type %T\n", m)
			}
			fmt.Println("nalgut", any)
		}
		return
	}
	return

	// msgs := make([]Message, 0)
	// 	err = json.Unmarshal([]byte(tx.Messages), &msgs)
	// 	if err != nil {
	// 		return []ResultTx{}, err
	// 	}

	// 	var fee *Fee
	// 	err = json.Unmarshal([]byte(tx.Fee), &fee)
	// 	if err != nil {
	// 		return []ResultTx{}, err
	// 	}

	// 	var logs []Log
	// 	err = json.Unmarshal([]byte(tx.Logs), &logs)
	// 	if err != nil {
	// 		return []ResultTx{}, err
	// 	}

	// 	tx := &ResultTx{
	// 		ID:        tx.ID,
	// 		Height:    tx.Height,
	// 		TxHash:    tx.TxHash,
	// 		Logs:      logs,
	// 		GasWanted: tx.GasWanted,
	// 		GasUsed:   tx.GasUsed,
	// 		Msgs:      msgs,
	// 		Fee:       fee,
	// 		Memo:      tx.Memo,
	// 		Timestamp: tx.Timestamp,
	// 	}

	// 	result = append(result, *tx)
	// }

	return result, nil
}
