package lcd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/databases"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	"github.com/rs/zerolog/log"
	resty "gopkg.in/resty.v1"
)

// SaveProposals saves governance proposals in database
func SaveProposals(db *databases.Database, config *config.Config) {
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
		for i := 0; i < len(proposalInfo); i++ {
			// check if a validator already voted
			count, _ := db.Model(&tempProposalInfo).
				Where("id = ?", proposalInfo[i].ID).
				Count()

			if count > 0 {
				// save and update proposalInfo
				_, _ = db.Model(&tempProposalInfo).
					Set("title = ?", proposalInfo[i].Title).
					Set("description = ?", proposalInfo[i].Description).
					Set("proposal_type = ?", proposalInfo[i].ProposalType).
					Set("proposal_status = ?", proposalInfo[i].ProposalStatus).
					Set("yes = ?", proposalInfo[i].Yes).
					Set("abstain = ?", proposalInfo[i].Abstain).
					Set("no = ?", proposalInfo[i].No).
					Set("no_with_veto = ?", proposalInfo[i].NoWithVeto).
					Set("deposit_end_time = ?", proposalInfo[i].DepositEndtime).
					Set("total_deposit_amount = ?", proposalInfo[i].TotalDepositAmount).
					Set("total_deposit_denom = ?", proposalInfo[i].TotalDepositDenom).
					Set("submit_time = ?", proposalInfo[i].SubmitTime).
					Set("voting_start_time = ?", proposalInfo[i].VotingStartTime).
					Set("voting_end_time = ?", proposalInfo[i].VotingEndTime).
					Where("id = ?", proposalInfo[i].ID).
					Update()
			} else {
				err := db.Insert(&proposalInfo)
				if err != nil {
					fmt.Printf("error - save and update proposalInfo: %v\n", err)
				}
			}
		}
	}
}
