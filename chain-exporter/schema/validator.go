package schema

import "time"

// Validator has validator information
type Validator struct {
	ID                   int64     `sql:",pk"`
	Rank                 int       `json:"rank"`
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

// PowerEventHistory has validator's power event history information
type PowerEventHistory struct {
	ID                   int64     `sql:",pk"`
	IDValidator          int       `json:"id_validator" sql:"default:0"`
	Height               int64     `json:"height"`
	Moniker              string    `json:"moniker"`
	OperatorAddress      string    `json:"operator_address"`
	Proposer             string    `json:"proposer"`
	VotingPower          float64   `json:"voting_power" sql:"default:0"`
	EventType            string    `json:"event_type" sql:"default:null"`
	NewVotingPowerAmount float64   `json:"new_voting_power_amount" sql:"new_voting_power_amount"`
	NewVotingPowerDenom  string    `json:"new_voting_power_denom" sql:"new_voting_power_denom"`
	TxHash               string    `json:"tx_hash" sql:"default:null"`
	Time                 time.Time `json:"time" sql:"default:null"`
}

// Miss has validator's range of missing blocks information
type Miss struct {
	ID           int64     `json:"id" sql:",pk"`
	Address      string    `json:"address"`
	StartHeight  int64     `json:"start_height"`
	EndHeight    int64     `json:"end_height"`
	MissingCount int64     `json:"missing_count"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Alerted      bool      `json:"alerted" sql:",default:false,notnull"`
}

// MissDetail has validator's missing blocks information
type MissDetail struct {
	ID       int64     `json:"id" sql:",pk"`
	Address  string    `json:"address"`
	Height   int64     `json:"height"`
	Proposer string    `json:"proposer_address"`
	Time     time.Time `json:"start_time"`
	Alerted  bool      `json:"alerted" sql:",default:false,notnull"`
}

// Evidence has evidence of slashing information
type Evidence struct {
	ID       int64     `json:"id" sql:",pk"`
	Proposer string    `json:"proposer"`
	Height   int64     `json:"height"`
	Hash     string    `json:"hash" sql:",unique"`
	Time     time.Time `json:"time"`
}
