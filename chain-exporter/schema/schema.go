package schema

import (
	"time"
)

// BlockInfo is a struct for database table
type BlockInfo struct {
	ID        int64     `json:"id" sql:",pk"`
	BlockHash string    `json:"block_hash" sql:",unique"`
	Height    int64     `json:"height"`
	Proposer  string    `json:"proposer"`
	TotalTxs  int64     `json:"total_txs" sql:"default:0"`
	NumTxs    int64     `json:"num_txs" sql:"default:0"`
	Time      time.Time `json:"time"`
}

// MissInfo is a struct for database table
type MissInfo struct {
	ID           int64     `json:"id" sql:",pk"`
	Address      string    `json:"address"`
	StartHeight  int64     `json:"start_height"`
	EndHeight    int64     `json:"end_height"`
	MissingCount int64     `json:"missing_count"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Alerted      bool      `json:"alerted" sql:",default:false,notnull"`
}

// MissDetailInfo is a struct for database table
type MissDetailInfo struct {
	ID       int64     `json:"id" sql:",pk"`
	Address  string    `json:"address"`
	Height   int64     `json:"height"`
	Proposer string    `json:"proposer_address"`
	Time     time.Time `json:"start_time"`
	Alerted  bool      `json:"alerted" sql:",default:false,notnull"`
}

// EvidenceInfo is a struct for database table
type EvidenceInfo struct {
	ID       int64     `json:"id" sql:",pk"`
	Proposer string    `json:"proposer"`
	Height   int64     `json:"height"`
	Hash     string    `json:"hash"`
	Time     time.Time `json:"time"`
}

// TransactionInfo is a struct for database table
type TransactionInfo struct {
	ID      int64     `json:"id" sql:",pk"`
	Height  int64     `json:"height"`
	TxHash  string    `json:"tx_hash"`
	MsgType string    `json:"msg_type"`
	Memo    string    `json:"memo"`
	Time    time.Time `json:"time"`
}

// ProposalInfo is a struct for database table
type ProposalInfo struct {
	ID                   int64     `json:"proposal_id" sql:",pk"`
	TxHash               string    `json:"tx_hash"`
	Proposer             string    `json:"proposer" sql:"default:null"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	ProposalType         string    `json:"proposal_type"`
	ProposalStatus       string    `json:"proposal_status"`
	Yes                  string    `json:"yes"`
	Abstain              string    `json:"abstain"`
	No                   string    `json:"no"`
	NoWithVeto           string    `json:"no_with_veto"`
	InitialDepositAmount string    `json:"initial_deposit_amount" sql:"default:null"`
	InitialDepositDenom  string    `json:"initial_deposit_denom" sql:"default:null"`
	TotalDepositAmount   string    `json:"total_deposit_amount"`
	TotalDepositDenom    string    `json:"total_deposit_denom"`
	SubmitTime           time.Time `json:"submit_time"`
	DepositEndtime       time.Time `json:"deposit_end_time" sql:"deposit_end_time"`
	VotingStartTime      time.Time `json:"voting_start_time"`
	VotingEndTime        time.Time `json:"voting_end_time"`
	Alerted              bool      `sql:"default:false,notnull" json:"alerted"`
}

// VoteInfo is a struct for database table
type VoteInfo struct {
	ID         int64     `json:"id" sql:",pk"`
	Height     int64     `json:"height"`
	ProposalID int64     `json:"proposal_id"`
	Voter      string    `json:"voter"`
	Option     string    `json:"option"`
	TxHash     string    `json:"tx_hash"`
	GasWanted  int64     `json:"gas_wanted"`
	GasUsed    int64     `json:"gas_used"`
	Time       time.Time `json:"time"`
}

// DepositInfo is a struct for database table
type DepositInfo struct {
	ID         int64     `json:"id" sql:",pk"`
	Height     int64     `json:"height"`
	ProposalID int64     `json:"proposal_id"`
	Depositor  string    `json:"depositor"`
	Amount     string    `json:"amount"`
	Denom      string    `json:"denom"`
	TxHash     string    `json:"tx_hash"`
	GasWanted  int64     `json:"gas_wanted"`
	GasUsed    int64     `json:"gas_used"`
	Time       time.Time `json:"time"`
}

// ValidatorInfo is a struct for database table
type ValidatorInfo struct {
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

// ValidatorSetInfo is a struct for database table
type ValidatorSetInfo struct {
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
