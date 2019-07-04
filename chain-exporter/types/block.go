package types

import (
	"time"
)

type BlockInfo struct {
	ID        int64     `json:"id" sql:",pk"`
	BlockHash string    `json:"block_hash"`
	Height    int64     `json:"height"`
	Proposer  string    `json:"proposer"`
	TotalTxs  int64     `json:"total_txs" sql:"default:0"`
	NumTxs    int64     `json:"num_txs" sql:"default:0"`
	Time      time.Time `json:"time"`
}

type EvidenceInfo struct {
	ID      int64     `json:"id" sql:",pk"`
	Address string    `json:"address"`
	Height  int64     `json:"height"`
	Hash    string    `json:"hash"`
	Time    time.Time `json:"time"`
}

type MissInfo struct {
	ID           int64     `json:"id" sql:",pk"`
	Address      string    `json:"address"`
	StartHeight  int64     `json:"start_height"`
	EndHeight    int64     `json:"end_height"`
	MissingCount int64     `json:"missing_count"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Alerted      bool      `json:"alerted" sql:",default:false,notnull"`
}

type MissDetailInfo struct {
	ID       int64     `json:"id" sql:",pk"`
	Address  string    `json:"address"`
	Height   int64     `json:"height"`
	Proposer string    `json:"proposer_address"`
	Time     time.Time `json:"start_time"`
	Alerted  bool      `json:"alerted" sql:",default:false,notnull"`
}

type TransactionInfo struct {
	ID      int64     `json:"id" sql:",pk"`
	Height  int64     `json:"height"`
	TxHash  string    `json:"tx_hash"`
	MsgType string    `json:"msg_type"`
	Time    time.Time `json:"time"`
}
