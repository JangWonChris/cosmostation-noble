package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Database struct
type ValidatorInfo struct {
	ID                   int64     `sql:",pk"`
	Rank                 int       `json:"rank"`
	OperatorAddress      string    `json:"operator_address" sql:",unique"`
	CosmosAddress        string    `json:"cosmos_address"`
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

type ValidatorDelegationsInfo struct {
	ID                  int64     `sql:",pk"`
	OperatorAddress     string    `json:"operator_address" sql:",unique"`
	CosmosAddress       string    `json:"cosmos_address"`
	TotalShares         float64   `json:"total_shares"`
	SelfDelegatedShares float64   `json:"self_delegated_shares"`
	OthersShares        float64   `json:"others"`
	DelegatorNum        int       `json:"delegator_num"`
	Time                time.Time `json:"time" sql:"default:null"`
}

// Database struct
type ValidatorSetInfo struct {
	ID                   int64     `sql:",pk"`
	IDValidator          int       `json:"id_validator" sql:"default:0"`
	Height               int64     `json:"height"`
	Proposer             string    `json:"proposer"`
	VotingPower          float64   `json:"voting_power" sql:"default:0"`
	EventType            string    `json:"event_type" sql:"default:null"`
	NewVotingPowerAmount float64   `json:"new_voting_power_amount" sql:"new_voting_power_amount"`
	NewVotingPowerDenom  string    `json:"new_voting_power_denom" sql:"new_voting_power_denom"`
	TxHash               string    `json:"tx_hash" sql:"default:null"`
	Time                 time.Time `json:"time" sql:"default:null"`
}

// LCD struct
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

type ValidatorDelegations struct {
	DelegatorAddress string  `json:"delegator_address"`
	ValidatorAddress string  `json:"validator_address"`
	Shares           sdk.Dec `json:"shares"`
}
