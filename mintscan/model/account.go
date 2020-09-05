package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Delegations defines the structure for delegations.
type Delegations struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
	Balance          string `json:"balance"`
	// Balance          Coin   `json:"balance"` // for next update
}

// Rewards defines the structure for rewards.
type Rewards struct {
	ValidatorAddress string `json:"validator_address"`
	Reward           []Coin `json:"reward"`
}

// UnbondingDelegations defines the structure for unbonding delegations.
type UnbondingDelegations struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Moniker          string `json:"moniker"`
	Entries          []struct {
		CreationHeight string `json:"creation_height"`
		CompletionTime string `json:"completion_time"`
		InitialBalance string `json:"initial_balance"`
		Balance        string `json:"balance"`
	} `json:"entries"`
}

// ModuleAccount defines the structure for module account information.
type ModuleAccount struct {
	Address       string    `json:"address"`
	AccountNumber uint64    `json:"account_number"`
	Coins         sdk.Coins `json:"coins"`
	Permissions   []string  `json:"permissions"`
	Name          string    `json:"name"`
}
