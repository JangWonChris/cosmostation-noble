package db

import (
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
)

// UpdateProposal updates the given proposal
func (db *Database) UpdateProposal(data *schema.ProposalInfo) (bool, error) {
	var proposalInfo schema.ProposalInfo
	_, err := db.Model(&proposalInfo).
		Set("title = ?", data.Title).
		Set("description = ?", data.Description).
		Set("proposal_type = ?", data.ProposalType).
		Set("proposal_status = ?", data.ProposalStatus).
		Set("yes = ?", data.Yes).
		Set("abstain = ?", data.Abstain).
		Set("no = ?", data.No).
		Set("no_with_veto = ?", data.NoWithVeto).
		Set("deposit_end_time = ?", data.DepositEndtime).
		Set("total_deposit_amount = ?", data.TotalDepositAmount).
		Set("total_deposit_denom = ?", data.TotalDepositDenom).
		Set("submit_time = ?", data.SubmitTime).
		Set("voting_start_time = ?", data.VotingStartTime).
		Set("voting_end_time = ?", data.VotingEndTime).
		Where("id = ?", data.ID).
		Update()
	if err != nil {
		return false, nil
	}
	return true, nil
}

// UpdateKeyBase updates the given validator info
func (db *Database) UpdateKeyBase(id int64, url string) (bool, error) {
	var validatorInfo schema.ValidatorInfo
	_, err := db.Model(&validatorInfo).
		Set("keybase_url = ?", url).
		Where("id = ?", id).
		Update()
	if err != nil {
		return false, err
	}
	return true, nil
}
