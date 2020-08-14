package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Delegations is a struct for REST API
type Delegations struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
	Balance          string `json:"balance"`
	// Balance          Coin   `json:"balance"`
}

// Rewards is a struct for REST API
type Rewards struct {
	ValidatorAddress string `json:"validator_address"`
	Reward           []Coin `json:"reward"`
}

// UnbondingDelegations is a struct for REST API
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

// ModuleAccount is module account on chain
type ModuleAccount struct {
	Address       string    `json:"address"`
	AccountNumber uint64    `json:"account_number"`
	Coins         sdk.Coins `json:"coins"`
	Permissions   []string  `json:"permissions"`
	Name          string    `json:"name"`
}
