package models

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ResultValidator struct {
	Rank                 int       `json:"rank"`
	OperatorAddress      string    `json:"operator_address"`
	ConsensusPubkey      string    `json:"consensus_pubkey"`
	Jailed               bool      `json:"jailed"`
	Status               int       `json:"status"`
	Tokens               string    `json:"tokens"`
	DelegatorShares      string    `json:"delegator_shares"`
	Moniker              string    `json:"moniker"`
	Identity             string    `json:"identity"`
	Website              string    `json:"website"`
	Details              string    `json:"details"`
	UnbondingHeight      string    `json:"unbonding_height"`
	UnbondingTime        time.Time `json:"unbonding_time"`
	CommissionRate       string    `json:"rate"`
	CommissionMaxRate    string    `json:"max_rate"`
	CommissionChangeRate string    `json:"max_change_rate"`
	UpdateTime           time.Time `json:"update_time"`
	Uptime               Uptime    `json:"uptime"`
	MinSelfDelegation    string    `json:"min_self_delegation"`
	KeybaseURL           string    `json:"keybase_url"`
}

type ResultValidatorDetail struct {
	Rank                 int       `json:"rank"`
	OperatorAddress      string    `json:"operator_address"`
	ConsensusPubkey      string    `json:"consensus_pubkey"`
	BondedHeight         int64     `json:"bonded_height"`
	BondedTime           time.Time `json:"bonded_time"`
	Jailed               bool      `json:"jailed"`
	Status               int       `json:"status"`
	Tokens               string    `json:"tokens"`
	DelegatorShares      string    `json:"delegator_shares"`
	Moniker              string    `json:"moniker"`
	Identity             string    `json:"identity"`
	Website              string    `json:"website"`
	Details              string    `json:"details"`
	UnbondingHeight      string    `json:"unbonding_height"`
	UnbondingTime        time.Time `json:"unbonding_time"`
	CommissionRate       string    `json:"rate"`
	CommissionMaxRate    string    `json:"max_rate"`
	CommissionChangeRate string    `json:"max_change_rate"`
	UpdateTime           time.Time `json:"update_time"`
	Uptime               Uptime    `json:"uptime"`
	MinSelfDelegation    string    `json:"min_self_delegation"`
	KeybaseURL           string    `json:"keybase_url"`
}

type Uptime struct {
	Address      string `json:"address"`
	MissedBlocks int64  `json:"missed_blocks"`
	OverBlocks   int64  `json:"over_blocks"`
}

type ResultMisses struct {
	StartHeight  int64     `json:"start_height"`
	EndHeight    int64     `json:"end_height"`
	MissingCount int64     `json:"missing_count"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
}
type ResultMissesDetail struct {
	Height int64     `json:"height"`
	Time   time.Time `json:"time"`
}

type ResultVotingPowerHistory struct {
	Height         int64     `json:"height"`
	EventType      string    `json:"event_type"`
	VotingPower    float64   `json:"voting_power"`
	NewVotingPower float64   `json:"new_voting_power"`
	TxHash         string    `json:"tx_hash"`
	Timestamp      time.Time `json:"timestamp"`
}

type ResultValidatorDelegations struct {
	TotalDelegatorNum     int                     `json:"total_delegator_num"`
	DelegatorNumChange24H int                     `json:"delegator_num_change_24h"`
	ValidatorDelegations  []*ValidatorDelegations `json:"delegations"`
}

type Validator struct {
	OperatorAddress string  `json:"operator_address"`
	ConsensusPubkey string  `json:"consensus_pubkey"`
	Jailed          bool    `json:"jailed"`
	Status          int     `json:"status"`
	Tokens          sdk.Dec `json:"tokens"`
	DelegatorShares sdk.Dec `json:"delegator_shares"`
	Description     struct {
		Moniker  string `json:"moniker"`
		Identity string `json:"identity"`
		Website  string `json:"website"`
		Details  string `json:"details"`
	} `json:"description"`
	UnbondingHeight string    `json:"unbonding_height"`
	UnbondingTime   time.Time `json:"unbonding_time"`
	Commission      struct {
		Rate          sdk.Dec   `json:"rate"`
		MaxRate       sdk.Dec   `json:"max_rate"`
		MaxChangeRate sdk.Dec   `json:"max_change_rate"`
		UpdateTime    time.Time `json:"update_time"`
	} `json:"commission"`
	MinSelfDelegation string `json:"min_self_delegation"`
}

type ValidatorDelegations struct {
	DelegatorAddress string  `json:"delegator_address"`
	ValidatorAddress string  `json:"validator_address"`
	Shares           sdk.Dec `json:"shares"`
	Amount           string  `json:"amount"`
}

type Redelegations struct {
	DelegatorAddress    string `json:"delegator_address"`
	ValidatorSrcAddress string `json:"validator_src_address"`
	ValidatorDstAddress string `json:"validator_dst_address"`
	Entries             []struct {
		CreationHeight string    `json:"creation_height"`
		CompletionTime time.Time `json:"completion_time"`
		InitialBalance string    `json:"initial_balance"`
		SharesDst      string    `json:"shares_dst"`
	} `json:"entries"`
}
