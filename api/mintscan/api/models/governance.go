package models

import "time"

type ResultProposal struct {
	ProposalID           int64  `json:"proposal_id"`
	TxHash               string `json:"tx_hash"`
	Proposer             string `json:"proposer" sql:"default:null"`
	Moniker              string `json:"moniker" sql:"default:null"`
	Title                string `json:"title"`
	Description          string `json:"description"`
	ProposalType         string `json:"proposal_type"`
	ProposalStatus       string `json:"proposal_status"`
	Yes                  string `json:"yes"`
	Abstain              string `json:"abstain"`
	No                   string `json:"no"`
	NoWithVeto           string `json:"no_with_veto"`
	InitialDepositAmount string `json:"initial_deposit_amount" sql:"default:null"`
	InitialDepositDenom  string `json:"initial_deposit_denom" sql:"default:null"`
	TotalDepositAmount   string `json:"total_deposit_amount"`
	TotalDepositDenom    string `json:"total_deposit_denom"`
	SubmitTime           string `json:"submit_time"`
	DepositEndtime       string `json:"deposit_end_time" sql:"deposit_end_time"`
	VotingStartTime      string `json:"voting_start_time"`
	VotingEndTime        string `json:"voting_end_time"`
}

type (
	ResultVoteInfo struct {
		Tally *Tally   `json:"tally"`
		Votes []*Votes `json:"votes"`
	}

	Tally struct {
		YesAmount        string `json:"yes_amount"`
		AbstainAmount    string `json:"abstain_amount"`
		NoAmount         string `json:"no_amount"`
		NoWithVetoAmount string `json:"no_with_veto_amount"`
		YesNum           int    `json:"yes_num"`
		AbstainNum       int    `json:"abstain_num"`
		NoNum            int    `json:"no_num"`
		NoWithVetoNum    int    `json:"no_with_veto_num"`
	}

	Votes struct {
		Voter   string    `json:"voter"`
		Moniker string    `json:"moniker" sql:"default:null"`
		Option  string    `json:"option"`
		TxHash  string    `json:"tx_hash"`
		Time    time.Time `json:"time"`
	}

	TallyInfo struct {
		Yes        string `json:"yes"`
		Abstain    string `json:"abstain"`
		No         string `json:"no"`
		NoWithVeto string `json:"no_with_veto"`
	}
)

type (
	DepositInfo struct {
		Depositor     string    `json:"depositor"`
		Moniker       string    `json:"moniker" sql:"default:null"`
		DepositAmount int64     `json:"deposit_amount"`
		DepositDenom  string    `json:"deposit_denom"`
		Height        int64     `json:"height"`
		TxHash        string    `json:"tx_hash"`
		Time          time.Time `json:"time"`
	}
)

type ResultProposalDetail struct {
	ProposalID         int64          `json:"proposal_id"`
	TotalVotesNum      int            `json:"total_votes_num"`
	TotalDepositAmount float64        `json:"total_deposit_amount"`
	ResultVoteInfo     ResultVoteInfo `json:"vote_info"`
	DepositInfo        DepositInfo    `json:"deposit_info"`
}
