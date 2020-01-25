package schema

import "time"

// ProposalInfo has proposal information
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

// VoteInfo has vote information
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

// DepositInfo has deposit information
type DepositInfo struct {
	ID         int64     `json:"id" sql:",pk"`
	Height     int64     `json:"height"`
	ProposalID int64     `json:"proposal_id"`
	Depositor  string    `json:"depositor"`
	Amount     int64     `json:"amount"`
	Denom      string    `json:"denom"`
	TxHash     string    `json:"tx_hash"`
	GasWanted  int64     `json:"gas_wanted"`
	GasUsed    int64     `json:"gas_used"`
	Time       time.Time `json:"time"`
}
