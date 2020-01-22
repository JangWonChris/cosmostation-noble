package models

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Uptime struct {
	Address      string `json:"address"`
	MissedBlocks int    `json:"missed_blocks"`
	OverBlocks   int64  `json:"over_blocks"`
}

// Validator is a struct for REST API
type Validator struct {
	OperatorAddress string `json:"operator_address"`
	ConsensusPubkey string `json:"consensus_pubkey"`
	Jailed          bool   `json:"jailed"`
	Status          int    `json:"status"`
	Tokens          string `json:"tokens"`
	DelegatorShares string `json:"delegator_shares"`
	Description     struct {
		Moniker  string `json:"moniker"`
		Identity string `json:"identity"`
		Website  string `json:"website"`
		Details  string `json:"details"`
	} `json:"description"`
	UnbondingHeight string    `json:"unbonding_height"`
	UnbondingTime   time.Time `json:"unbonding_time"`
	Commission      struct {
		CommissionRates struct {
			Rate          string `json:"rate"`
			MaxRate       string `json:"max_rate"`
			MaxChangeRate string `json:"max_change_rate"`
		}
		UpdateTime time.Time `json:"update_time"`
	} `json:"commission"`
	MinSelfDelegation string `json:"min_self_delegation"`
}

// ValidatorDelegations is a struct for REST API
type ValidatorDelegations struct {
	DelegatorAddress string  `json:"delegator_address"`
	ValidatorAddress string  `json:"validator_address"`
	Shares           sdk.Dec `json:"shares"`
	Amount           string  `json:"amount"`
}

// Redelegations is a struct for REST API
type Redelegations struct {
	DelegatorAddress    string `json:"delegator_address"`
	ValidatorSrcAddress string `json:"validator_src_address"`
	ValidatorDstAddress string `json:"validator_dst_address"`
	Entries             []struct {
		CreationHeight int       `json:"creation_height"`
		CompletionTime time.Time `json:"completion_time"`
		InitialBalance string    `json:"initial_balance"`
		SharesDst      string    `json:"shares_dst"`
		Balance        string    `json:"balance"`
	} `json:"entries"`
}
