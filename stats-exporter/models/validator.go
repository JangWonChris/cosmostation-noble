package models

import (
	"time"
)

const (
	// BondedValidatorStatus is status code when a validator is live.
	BondedValidatorStatus = 2

	// UnbondingValidatorStatus is status code when a validator is not live.
	UnbondingValidatorStatus = 1

	// UnbondedValidatorStatus is status code when a validator is jailed.
	UnbondedValidatorStatus = 0
)

// Validator defines the structure for validator
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

// ValidatorDelegation defines the structure for delegations for a validator
type ValidatorDelegation struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
	Balance          string `json:"balance"`
	// Balance          Coin   `json:"balance"`
}

// SelfDelegation defines the structure for self-bonded delegation
type SelfDelegation struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
	Balance          string `json:"balance"`
	// Balance          Coin   `json:"balance"`
}
