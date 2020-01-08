package lcd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/db"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	"github.com/rs/zerolog/log"
	resty "gopkg.in/resty.v1"
)

// SaveProposals saves governance proposals in database
func SaveProposals(db *db.Database, config *config.Config) {
	resp, err := resty.R().Get(config.Node.LCDURL + "/gov/proposals")
	if err != nil {
		fmt.Printf("query /gov/proposals error - %v\n", err)
	}

	proposals := make([]*types.Proposal, 0)
	err = json.Unmarshal(types.ReadRespWithHeight(resp).Result, &proposals)
	if err != nil {
		log.Info().Str(types.Service, types.LogGovernance).Str(types.Method, "SaveProposals").Err(err).Msg("unmarshal proposals error")
	}

	// proposal information for our database table
	proposalInfo := make([]*schema.ProposalInfo, 0)
	if len(proposals) > 0 {
		for _, proposal := range proposals {
			proposalID, _ := strconv.ParseInt(proposal.ID, 10, 64)

			var totalDepositAmount string
			var totalDepositDenom string
			if proposal.TotalDeposit != nil {
				totalDepositAmount = proposal.TotalDeposit[0].Amount
				totalDepositDenom = proposal.TotalDeposit[0].Denom
			}

			tallyResp, _ := resty.R().Get(config.Node.LCDURL + "/gov/proposals/" + proposal.ID + "/tally")

			var tally types.Tally
			err = json.Unmarshal(types.ReadRespWithHeight(tallyResp).Result, &tally)
			if err != nil {
				log.Info().Str(types.Service, types.LogGovernance).Str(types.Method, "SaveProposals").Err(err).Msg("unmarshal tally error")
			}

			tempProposalInfo := &schema.ProposalInfo{
				ID:                 proposalID,
				Title:              proposal.Content.Value.Title,
				Description:        proposal.Content.Value.Description,
				ProposalType:       proposal.Content.Type,
				ProposalStatus:     proposal.ProposalStatus,
				Yes:                tally.Yes,
				Abstain:            tally.Abstain,
				No:                 tally.No,
				NoWithVeto:         tally.NoWithVeto,
				SubmitTime:         proposal.SubmitTime,
				DepositEndtime:     proposal.DepositEndTime,
				TotalDepositAmount: totalDepositAmount,
				TotalDepositDenom:  totalDepositDenom,
				VotingStartTime:    proposal.VotingStartTime,
				VotingEndTime:      proposal.VotingEndTime,
				Alerted:            false,
			}
			proposalInfo = append(proposalInfo, tempProposalInfo)
		}
	}

	// update proposerInfo
	if len(proposalInfo) > 0 {
		var tempProposalInfo schema.ProposalInfo
		for _, proposal := range proposalInfo {
			exist, _ := db.Model(&tempProposalInfo).
				Where("id = ?", proposal.ID).
				Exists()

			if exist {
				// save and update proposalInfo
				_, _ = db.Model(&tempProposalInfo).
					Set("title = ?", proposal.Title).
					Set("description = ?", proposal.Description).
					Set("proposal_type = ?", proposal.ProposalType).
					Set("proposal_status = ?", proposal.ProposalStatus).
					Set("yes = ?", proposal.Yes).
					Set("abstain = ?", proposal.Abstain).
					Set("no = ?", proposal.No).
					Set("no_with_veto = ?", proposal.NoWithVeto).
					Set("deposit_end_time = ?", proposal.DepositEndtime).
					Set("total_deposit_amount = ?", proposal.TotalDepositAmount).
					Set("total_deposit_denom = ?", proposal.TotalDepositDenom).
					Set("submit_time = ?", proposal.SubmitTime).
					Set("voting_start_time = ?", proposal.VotingStartTime).
					Set("voting_end_time = ?", proposal.VotingEndTime).
					Where("id = ?", proposal.ID).
					Update()
			} else {
				err := db.Insert(proposal)
				if err != nil {
					fmt.Printf("error - save and update proposalInfo: %v\n", err)
				}
			}
		}
	}
}
