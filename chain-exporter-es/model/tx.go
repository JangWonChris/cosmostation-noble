package model

import (
	"encoding/json"
	"time"
)

type ElasticsearchTxInfo struct {
	Hash   string          `json:"hash"`
	Height int64           `json:"height"`
	Time   time.Time       `json:"time"`
	Tx     json.RawMessage `json:"tx"`
	Result *TxResultInfo   `json:"result"`
}

type TxResultInfo struct {
	GasWanted int64           `json:"gas_wanted"`
	GasUsed   int64           `json:"gas_used"`
	Log       json.RawMessage `json:"log"`
	Tags      json.RawMessage `json:"tags"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}