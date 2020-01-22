package db

import "github.com/cosmostation/cosmostation-cosmos/mintscan/api/schema"

// QueryProposals queries proposals that are saved in database
func (db *Database) QueryProposals() []schema.ProposalInfo {
	proposals := make([]schema.ProposalInfo, 0)
	_ = db.Model(&proposals).
		Select()

	return proposals
}

// QueryProposal queries particular proposal
func (db *Database) QueryProposal(id string) schema.ProposalInfo {
	var proposal schema.ProposalInfo
	_ = db.Model(&proposal).
		Where("id = ?", id).
		Select()

	return proposal
}

// QueryExistsProposal queries if proposal exists
func (db *Database) QueryExistsProposal(id string) bool {
	var proposal schema.ProposalInfo
	exists, _ := db.Model(&proposal).
		Where("id = ?", id).
		Exists()

	return exists
}

// QueryVotes queries all votes
func (db *Database) QueryVotes(id string) []schema.VoteInfo {
	votes := make([]schema.VoteInfo, 0)
	_ = db.Model(&votes).
		Where("proposal_id = ?", id).
		Order("id DESC").
		Select()

	return votes
}

// QueryVoteOptions queries all vote options for the proposal
func (db *Database) QueryVoteOptions(id string) (int, int, int, int) {
	votes := make([]schema.VoteInfo, 0)

	yes, _ := db.Model(&votes).
		Where("proposal_id = ? AND option = ?", id, "Yes").
		Count()

	no, _ := db.Model(&votes).
		Where("proposal_id = ? AND option = ?", id, "No").
		Count()

	noWithVeto, _ := db.Model(&votes).
		Where("proposal_id = ? AND option = ?", id, "NoWithVeto").
		Count()

	abstain, _ := db.Model(&votes).
		Where("proposal_id = ? AND option = ?", id, "Abstain").
		Count()

	return yes, no, noWithVeto, abstain
}

// QueryDeposits queries all deposits
func (db *Database) QueryDeposits(id string) []schema.DepositInfo {
	deposits := make([]schema.DepositInfo, 0)
	_ = db.Model(&deposits).
		Where("proposal_id = ?", id).
		Order("id DESC").
		Select()

	return deposits
}
