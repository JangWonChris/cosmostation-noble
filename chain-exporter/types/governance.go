package types

import "time"

// Proposal Information - Parsing from LCD
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

// Voting Information
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

// Deposit Information
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

// LCD struct
type Proposal struct {
	ProposalContent struct {
		Type  string `json:"type"`
		Value struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		} `json:"value"`
	} `json:"proposal_content"`
	ProposalID       string `json:"proposal_id"`
	ProposalStatus   string `json:"proposal_status"`
	FinalTallyResult struct {
		Yes        string `json:"yes"`
		Abstain    string `json:"abstain"`
		No         string `json:"no"`
		NoWithVeto string `json:"no_with_veto"`
	} `json:"final_tally_result"`
	SubmitTime     time.Time `json:"submit_time"`
	DepositEndTime time.Time `json:"deposit_end_time"`
	TotalDeposit   []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"total_deposit"`
	VotingStartTime time.Time `json:"voting_start_time"`
	VotingEndTime   time.Time `json:"voting_end_time"`
}

// type Proposal struct {
// 	Type  string `json:"type"`
// 	Value struct {
// 		ProposalID       string `json:"proposal_id"`
// 		Title            string `json:"title"`
// 		Description      string `json:"description"`
// 		ProposalType     string `json:"proposal_type"`
// 		ProposalStatus   string `json:"proposal_status"`
// 		FinalTallyResult struct {
// 			Yes        string `json:"yes"`
// 			Abstain    string `json:"abstain"`
// 			No         string `json:"no"`
// 			NoWithVeto string `json:"no_with_veto"`
// 		} `json:"final_tally_result"`
// 		SubmitTime     string `json:"submit_time"`
// 		DepositEndTime string `json:"deposit_end_time"`
// 		TotalDeposit   []struct {
// 			Denom  string `json:"denom"`
// 			Amount string `json:"amount"`
// 		} `json:"total_deposit"`
// 		VotingStartTime string `json:"voting_start_time"`
// 		VotingEndTime   string `json:"voting_end_time"`
// 	} `json:"value"`
// }

// LCD struct
type TallyInfo struct {
	Yes        string `json:"yes"`
	Abstain    string `json:"abstain"`
	No         string `json:"no"`
	NoWithVeto string `json:"no_with_veto"`
}
