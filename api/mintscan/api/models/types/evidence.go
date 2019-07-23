package types

import (
	"time"
)

type EvidenceInfo struct {
	ID      int64     `json:"id" sql:",pk"`
	Address string    `json:"address"`
	Height  int64     `json:"height"`
	Hash    string    `json:"hash"`
	Time    time.Time `json:"time"`
}
