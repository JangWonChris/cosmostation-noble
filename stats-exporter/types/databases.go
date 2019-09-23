package types

import (
	"time"
)

// BlockInfo is a struct for database struct
type BlockInfo struct {
	ID        int64     `json:"id" sql:",pk"`
	BlockHash string    `json:"block_hash"`
	Height    int64     `json:"height"`
	Proposer  string    `json:"proposer"`
	TotalTxs  int64     `json:"total_txs" sql:"default:0"`
	NumTxs    int64     `json:"num_txs" sql:"default:0"`
	Time      time.Time `json:"time"`
}

// ValidatorInfo is a struct for database struct
type ValidatorInfo struct {
	ID                   int64     `sql:",pk"`
	Rank                 int64     `json:"rank"`
	Address              string    `json:"address"`
	OperatorAddress      string    `json:"operator_address" sql:",unique"`
	ConsensusPubkey      string    `json:"consensus_pubkey"`
	Proposer             string    `json:"proposer"`
	Jailed               bool      `json:"jailed" sql:"default:false,notnull"`
	Status               int       `json:"status" sql:"default:0"`
	Tokens               string    `json:"tokens"`
	DelegatorShares      string    `json:"delegator_shares"`
	Moniker              string    `json:"moniker"`
	Identity             string    `json:"identity"`
	Website              string    `json:"website"`
	Details              string    `json:"details"`
	UnbondingHeight      string    `json:"unbonding_height"`
	UnbondingTime        time.Time `json:"unbonding_time" sql:"default:null"`
	CommissionRate       string    `json:"rate"`
	CommissionMaxRate    string    `json:"max_rate"`
	CommissionChangeRate string    `json:"max_change_rate"`
	UpdateTime           time.Time `json:"update_time" sql:"default:null"`
	MinSelfDelegation    string    `json:"min_self_delegation"`
	KeybaseURL           string    `json:"keybase_url"`
}
