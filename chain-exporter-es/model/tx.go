package model

import (
	"encoding/json"
	"time"
)

type ElasticsearchTxInfo struct {
	Height int64           `json:"height"`
	Hash   string          `json:"hash"`
	RawLog string `json:"raw_log"`
	Logs      json.RawMessage `json:"logs"`
	GasWanted int64           `json:"gas_wanted"`
	GasUsed   int64           `json:"gas_used"`
	Events json.RawMessage `json:"events"`
	Tx     json.RawMessage `json:"tx"`
	Timestamp   time.Time       `json:"timestamp"`
}


type Event struct {
	Type string `json:"type"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}