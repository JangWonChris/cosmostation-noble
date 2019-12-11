package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

/*
	Threse structs are for REST API structs
*/

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
		} `json:"commission_rates"`
		UpdateTime time.Time `json:"update_time"`
	} `json:"commission"`
	MinSelfDelegation string `json:"min_self_delegation"`
}

type ValidatorDelegations struct {
	DelegatorAddress string  `json:"delegator_address"`
	ValidatorAddress string  `json:"validator_address"`
	Shares           sdk.Dec `json:"shares"`
}

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
