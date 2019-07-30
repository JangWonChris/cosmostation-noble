package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// REST API struct
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
		Rate          string    `json:"rate"`
		MaxRate       string    `json:"max_rate"`
		MaxChangeRate string    `json:"max_change_rate"`
		UpdateTime    time.Time `json:"update_time"`
	} `json:"commission"`
	MinSelfDelegation string `json:"min_self_delegation"`
}

type ValidatorDelegationsInfo struct {
	ID                  int64     `sql:",pk"`
	OperatorAddress     string    `json:"operator_address" sql:",unique"`
	Address             string    `json:"cosmos_address"`
	TotalShares         float64   `json:"total_shares"`
	SelfDelegatedShares float64   `json:"self_delegated_shares"`
	OthersShares        float64   `json:"others"`
	DelegatorNum        int       `json:"delegator_num"`
	Time                time.Time `json:"time" sql:"default:null"`
}

type ValidatorDelegations struct {
	DelegatorAddress string  `json:"delegator_address"`
	ValidatorAddress string  `json:"validator_address"`
	Shares           sdk.Dec `json:"shares"`
}

// KeyBase struct
type KeyBase struct {
	Status struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
	} `json:"status"`
	Them []struct {
		ID       string `json:"id"`
		Pictures struct {
			Primary struct {
				URL    string `json:"url"`
				Source string `json:"source"`
			} `json:"primary"`
		} `json:"pictures"`
	} `json:"them"`
}
