package schema

import "time"

// Proposal has proposal information.
type Proposal struct {
	ID                   int64     `json:"proposal_id" sql:",pk"`
	TxHash               string    `json:"tx_hash" sql:",unique"`
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
}

// Deposit has deposit information.
type Deposit struct {
	ID         int64     `json:"id" sql:",pk"`
	Height     int64     `json:"height"`
	ProposalID uint64    `json:"proposal_id"`
	Depositor  string    `json:"depositor"`
	Amount     string    `json:"amount"`
	Denom      string    `json:"denom"`
	TxHash     string    `json:"tx_hash" sql:",unique"`
	GasWanted  int64     `json:"gas_wanted"`
	GasUsed    int64     `json:"gas_used"`
	Timestamp  time.Time `json:"timestamp" sql:"default:now()"`
}

// Vote has vote information.
type Vote struct {
	ID         int64     `json:"id" sql:",pk"`
	Height     int64     `json:"height"`
	ProposalID uint64    `json:"proposal_id"`
	Voter      string    `json:"voter"`
	Option     string    `json:"option"`
	TxHash     string    `json:"tx_hash"`
	GasWanted  int64     `json:"gas_wanted"`
	GasUsed    int64     `json:"gas_used"`
	Timestamp  time.Time `json:"timestamp" sql:"default:now()"`
}

// NewProposal returns new Proposal.
func NewProposal(p Proposal) *Proposal {
	return &Proposal{
		TxHash:               p.TxHash,
		Proposer:             p.Proposer,
		Title:                p.Title,
		Description:          p.Description,
		ProposalType:         p.ProposalType,
		ProposalStatus:       p.ProposalStatus,
		Yes:                  p.Yes,
		Abstain:              p.Abstain,
		No:                   p.No,
		NoWithVeto:           p.NoWithVeto,
		InitialDepositAmount: p.InitialDepositAmount,
		InitialDepositDenom:  p.InitialDepositDenom,
		TotalDepositAmount:   p.TotalDepositAmount,
		TotalDepositDenom:    p.TotalDepositDenom,
		SubmitTime:           p.SubmitTime,
		DepositEndtime:       p.DepositEndtime,
		VotingStartTime:      p.VotingStartTime,
		VotingEndTime:        p.VotingEndTime,
	}
}

// NewDeposit returns new Deposit.
func NewDeposit(d Deposit) *Deposit {
	return &Deposit{
		Height:     d.Height,
		ProposalID: d.ProposalID,
		Depositor:  d.Depositor,
		Amount:     d.Amount,
		Denom:      d.Denom,
		TxHash:     d.TxHash,
		GasWanted:  d.GasWanted,
		GasUsed:    d.GasUsed,
		Timestamp:  d.Timestamp,
	}
}

// NewVote returns new Vote.
func NewVote(v Vote) *Vote {
	return &Vote{
		Height:     v.Height,
		ProposalID: v.ProposalID,
		Voter:      v.Voter,
		Option:     v.Option,
		TxHash:     v.TxHash,
		GasWanted:  v.GasWanted,
		GasUsed:    v.GasUsed,
		Timestamp:  v.Timestamp,
	}
}
