package types

import (
	"encoding/json"
	"time"
)

const (
	MultiMsg = "multi_msg"
)

// GeneralTx is general tx struct that is unmarshallable for any tx_msg type
type GeneralTx struct {
	Height string `json:"height"`
	TxHash string `json:"txhash"`
	RawLog string `json:"raw_log"`
	Logs   []struct {
		MsgIndex uint8  `json:"msg_index"`
		Success  bool   `json:"success"`
		Log      string `json:"log"`
	} `json:"logs"`
	GasWanted string `json:"gas_wanted"`
	GasUsed   string `json:"gas_used"`
	Events    []struct {
		Type       string `json:"type"`
		Attributes []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
	} `json:"events"`
	Tx struct {
		Type  string `json:"type"`
		Value struct {
			Msg []struct {
				Type  string          `json:"type"`
				Value json.RawMessage `json:"value"`
			} `json:"msg"`
			Fee struct {
				Amount []struct {
					Denom  string `json:"denom"`
					Amount string `json:"amount"`
				} `json:"amount"`
				Gas string `json:"gas"`
			} `json:"fee"`
			Signatures json.RawMessage `json:"signatures"`
			Memo       string          `json:"memo"`
		} `json:"value"`
	} `json:"tx"`
	Timestamp time.Time `json:"timestamp"`
}
