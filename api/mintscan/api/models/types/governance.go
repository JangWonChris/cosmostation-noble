package types

import "time"

// Proposal is a struct for REST API
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

// Tally is a struct for REST API
type Tally struct {
	Yes        string `json:"yes"`
	Abstain    string `json:"abstain"`
	No         string `json:"no"`
	NoWithVeto string `json:"no_with_veto"`
}
