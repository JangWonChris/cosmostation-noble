package model

import (
	"encoding/json"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/schema"
)

type TxData struct {
	Txs []string `json:"txs"`
}

type TxList struct {
	TxHash []string `json:"tx_list"`
}

type Message struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

type Fee struct {
	Gas    string `json:"gas,omitempty"`
	Amount []struct {
		Amount string `json:"amount,omitempty"`
		Denom  string `json:"denom,omitempty"`
	} `json:"amount,omitempty"`
}

type Event struct {
	Type       string `json:"type"`
	Attributes []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"attributes"`
}

type Log struct {
	MsgIndex int     `json:"msg_index"`
	Log      string  `json:"log"`
	Events   []Event `json:"events"`
}

// ParseTransaction receives single transaction from database and return it after unmarshal them.
func ParseTransaction(tx schema.Transaction) (result *ResultTx, err error) {
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
			Timestamp: tx.Timestamp,
		}

		result = append(result, *tx)
	}

	return result, nil
}
