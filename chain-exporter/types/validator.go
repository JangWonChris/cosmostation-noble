package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// DefaultQueryValidatorsPage is the default page number for querying validators via querier.
	DefaultQueryValidatorsPage = 1

	// DefaultQueryValidatorsPerPage is the default per page number for querying validators via querier.
	DefaultQueryValidatorsPerPage = 200

	// BondedValidatorStatus is status code when a validator is live.
	BondedValidatorStatus = 2

	// UnbondingValidatorStatus is status code when a validator is not live.
	UnbondingValidatorStatus = 1

	// UnbondedValidatorStatus is status code when a validator is jailed.
	UnbondedValidatorStatus = 0
)

// Validator defines the structure for validator information.
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

// ValidatorDelegations defines the structure for validator's delegations.
type ValidatorDelegations struct {
	DelegatorAddress string  `json:"delegator_address"`
	ValidatorAddress string  `json:"validator_address"`
	Shares           sdk.Dec `json:"shares"`
}
