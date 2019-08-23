package lcd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	dtypes "github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	"github.com/go-pg/pg"
	resty "gopkg.in/resty.v1"
)

// SaveProposals queries the governance proposals from REST server and store them in our database
func SaveProposals(db *pg.DB, config *config.Config) {
	resp, err := resty.R().Get(config.Node.LCDURL + "/gov/proposals")
	if err != nil {
		fmt.Printf("query /gov/proposals error - %v\n", err)
	}

	var responseWithHeight dtypes.ResponseWithHeight
	err = json.Unmarshal(resp.Body(), &responseWithHeight)
	if err != nil {
		fmt.Printf("unmarshal proposals error - %v\n", err)
	}

	proposals := make([]*dtypes.Proposal, 0)
	err = json.Unmarshal(responseWithHeight.Result, &proposals)
	if err != nil {
		fmt.Printf("unmarshal proposals error - %v\n", err)
	}

	// proposal information for our database table
	proposalInfo := make([]*dtypes.ProposalInfo, 0)
	if len(proposals) > 0 {
		for _, proposal := range proposals {
			var tallyInfo dtypes.TallyInfo
			tallyResp, _ := resty.R().Get(config.Node.LCDURL + "/gov/proposals/" + proposal.ProposalID + "/tally")
			err = json.Unmarshal(tallyResp.Body(), &tallyInfo)
			if err != nil {
				fmt.Printf("unmarshal tallyInfo error - %v\n", err)
			}

			proposalID, _ := strconv.ParseInt(proposal.ProposalID, 10, 64)

			var totalDepositAmount string
			var totalDepositDenom string
			if proposal.TotalDeposit != nil {
				totalDepositAmount = proposal.TotalDeposit[0].Amount
				totalDepositDenom = proposal.TotalDeposit[0].Denom
			}

			tempProposalInfo := &dtypes.ProposalInfo{
				ID:                 proposalID,
				Title:              proposal.ProposalContent.Value.Title,
				Description:        proposal.ProposalContent.Value.Description,
				ProposalType:       proposal.ProposalContent.Type,
				ProposalStatus:     proposal.ProposalStatus,
				Yes:                tallyInfo.Yes,
				Abstain:            tallyInfo.Abstain,
				No:                 tallyInfo.No,
				NoWithVeto:         tallyInfo.NoWithVeto,
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

	// // exist and update proposerInfo
	// if len(proposalInfo) > 0 {
	// 	var tempProposalInfo dtypes.ProposalInfo
	// 	for i := 0; i < len(proposalInfo); i++ {
	// 		// check if a validator already voted
	// 		count, _ := db.Model(&tempProposalInfo).
	// 			Where("id = ?", proposalInfo[i].ID).
	// 			Count()
	// 		if count > 0 {
	// 			// save and update proposalInfo
	// 			_, _ = db.Model(&tempProposalInfo).
	// 				Set("title = ?", proposalInfo[i].Title).
	// 				Set("description = ?", proposalInfo[i].Description).
	// 				Set("proposal_type = ?", proposalInfo[i].ProposalType).
	// 				Set("proposal_status = ?", proposalInfo[i].ProposalStatus).
	// 				Set("yes = ?", proposalInfo[i].Yes).
	// 				Set("abstain = ?", proposalInfo[i].Abstain).
	// 				Set("no = ?", proposalInfo[i].No).
	// 				Set("no_with_veto = ?", proposalInfo[i].NoWithVeto).
	// 				Set("submit_time = ?", proposalInfo[i].SubmitTime).
	// 				Set("deposit_end_time = ?", proposalInfo[i].DepositEndtime).
	// 				Set("total_deposit_amount = ?", proposalInfo[i].TotalDepositAmount).
	// 				Set("total_deposit_denom = ?", proposalInfo[i].TotalDepositDenom).
	// 				Set("voting_start_time = ?", proposalInfo[i].VotingStartTime).
	// 				Set("voting_end_time = ?", proposalInfo[i].VotingEndTime).
	// 				Where("id = ?", proposalInfo[i].ID).
	// 				Update()
	// 		} else {
	// 			err := db.Insert(&proposalInfo)
	// 			if err != nil {
	// 				fmt.Printf("error - save and update proposalInfo: %v\n", err)
	// 			}
	// 		}
	// 	}
	// }
}
