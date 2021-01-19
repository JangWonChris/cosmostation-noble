package model

import "time"

const (
	// YES is one of the voting options that agrees to the proposal.
	YES = "Yes"

	// NO is one of the voting options that disagree with the proposal.
	NO = "No"

	// NOWITHVETO is one of the voting options that strongly disagree with the proposal.
	NOWITHVETO = "NoWithVeto"

	// ABSTAIN is the one of the voting options that gives up his/her voting right.
	ABSTAIN = "Abstain"
)

// Proposal defines the structure for a proposal information.
// type Proposal struct {
// 	Content struct {
// 		Type  string `json:"type"`
// 		Value struct {
// 			Title       string `json:"title"`
// 			Description string `json:"description"`
// 		} `json:"value"`
// 	} `json:"content"`
// 	ID               string `json:"id"`
// 	ProposalStatus   string `json:"proposal_status"`
// 	FinalTallyResult struct {
// 		Yes        string `json:"yes"`
// 		Abstain    string `json:"abstain"`
// 		No         string `json:"no"`
// 		NoWithVeto string `json:"no_with_veto"`
// 	} `json:"final_tally_result"`
// 	SubmitTime     time.Time `json:"submit_time"`
// 	DepositEndTime time.Time `json:"deposit_end_time"`
// 	TotalDeposit   []struct {
// 		Denom  string `json:"denom"`
// 		Amount string `json:"amount"`
// 	} `json:"total_deposit"`
// 	VotingStartTime time.Time `json:"voting_start_time"`
// 	VotingEndTime   time.Time `json:"voting_end_time"`
// }

// Votes defines the structure for proposal votes.
type Votes struct {
	Voter   string    `json:"voter"`
	Moniker string    `json:"moniker" sql:"default:null"`
	Option  string    `json:"option"`
	TxHash  string    `json:"tx_hash"`
	Time    time.Time `json:"time"`
}

// Tally defines the structure for a proposal's tally information.
// type Tally struct {
// 	Yes        string `json:"yes"`
// 	Abstain    string `json:"abstain"`
// 	No         string `json:"no"`
// 	NoWithVeto string `json:"no_with_veto"`
// }
