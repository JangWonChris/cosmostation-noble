package lcd

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/db"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	resty "gopkg.in/resty.v1"
)

// SaveProposals saves governance proposals in database
func SaveProposals(db *db.Database, config *config.Config) {
	resp, err := resty.R().Get(config.Node.LCDEndpoint + "/gov/proposals")
	if err != nil {
		fmt.Printf("failed to request /gov/proposals: %v \n", err)
	}

	proposals := make([]*types.Proposal, 0)
	err = json.Unmarshal(types.ReadRespWithHeight(resp).Result, &proposals)
	if err != nil {
		fmt.Printf("failed to unmarshal Proposal: %v \n", err)
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

			tallyResp, _ := resty.R().Get(config.Node.LCDEndpoint + "/gov/proposals/" + proposal.ID + "/tally")

			var tally types.Tally
			err = json.Unmarshal(types.ReadRespWithHeight(tallyResp).Result, &tally)
			if err != nil {
				fmt.Printf("failed to unmarshal Tally: %v \n", err)
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

	if len(proposalInfo) > 0 {
		for _, proposal := range proposalInfo {
			exist, _ := db.QueryExistProposal(proposal.ID)

			if exist {
				result, _ := db.UpdateProposal(proposal)
				if !result {
					log.Printf("failed to update Proposal ID: %d", proposal.ID)
				}
			} else {
				result, _ := db.InsertProposal(proposal)
				if !result {
					log.Printf("failed to save Proposal ID: %d", proposal.ID)
				}
			}
		}
	}
}
