package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Database struct
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

type DelegatorsDelegation struct {
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
