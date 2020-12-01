package client

import (
	"context"
	"fmt"
	"log"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"

	// cosmos-sdk
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibccoretypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	paramstypesproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// GetGovQueryClient returns a object of queryClient
func (c *Client) GetGovQueryClient() govtypes.QueryClient {
	return govtypes.NewQueryClient(c.grpcClient)
}

// GetProposals returns all governance proposals
func (c *Client) GetProposals() (result []schema.Proposal, err error) {
	queryClient := c.GetGovQueryClient()
	resp, err := queryClient.Proposals(context.Background(), &govtypes.QueryProposalsRequest{})
	if err != nil {
		return []schema.Proposal{}, fmt.Errorf("failed to request gov proposals: %s", err)
	}

	if len(resp.Proposals) <= 0 {
		return []schema.Proposal{}, nil
	}

	for _, proposal := range resp.Proposals {
		totalDepositAmount := make([]string, len(proposal.TotalDeposit))
		totalDepositDenom := make([]string, len(proposal.TotalDeposit))
		for i, td := range proposal.TotalDeposit {
			totalDepositAmount[i] = td.Amount.String()
			totalDepositDenom[i] = td.Denom
		}

		request := govtypes.QueryTallyResultRequest{ProposalId: proposal.ProposalId}
		tallyResultResp, err := queryClient.TallyResult(context.Background(), &request)
		if err != nil {
			return []schema.Proposal{}, fmt.Errorf("failed to request gov proposals: %s", err)
		}
		tally := tallyResultResp.Tally

		// desp := proposal.GetContent().GetDescription()
		// log.Println("desp :", proposal)
		log.Println("proposal.GetContent() : ", proposal.GetContent())
		log.Println("proposal.Content(any) :", proposal.Content)

		var contentI govtypes.Content
		err = codec.AppCodec.UnpackAny(proposal.Content, &contentI)
		if err != nil {
			log.Println(err)
		}
		log.Println("UnpackAny :", contentI)

		switch i := contentI.(type) {
		case *govtypes.TextProposal:
			log.Println(i.Title)
			log.Println(i.Description)
		case *distributiontypes.CommunityPoolSpendProposal:
			log.Println(i.Title)
			log.Println(i.Description)
		case *ibccoretypes.ClientUpdateProposal:
			log.Println(i.Title)
			log.Println(i.Description)
		case *paramstypesproposal.ParameterChangeProposal:
			log.Println(i.Title)
			log.Println(i.Description)
		case *upgradetypes.SoftwareUpgradeProposal:
			log.Println(i.Title)
			log.Println(i.Description)
		case *upgradetypes.CancelSoftwareUpgradeProposal:
			log.Println(i.Title)
			log.Println(i.Description)
		default:
			log.Println("default")
		}
		// kind of proposals
		// distributiontypes.CommunityPoolSpendProposal
		// govtypes.TextProposal
		// ibccoretypes.ClientUpdateProposal
		// paramstypes.ParameterChangeProposal
		// upgradetypes.SoftwareUpgradeProposal
		// upgradetypes.CancelSoftwareUpgradeProposal

		p := schema.NewProposal(schema.Proposal{
			ID:           proposal.ProposalId,
			Title:        proposal.GetTitle(),
			Description:  "content description nil",
			ProposalType: "proposaltype nil",
			// Description:        proposal.GetContent().GetDescription(),
			// ProposalType:       proposal.GetContent().ProposalType(),
			ProposalStatus:     proposal.Status.String(),
			Yes:                tally.Yes.String(),
			Abstain:            tally.Abstain.String(),
			No:                 tally.No.String(),
			NoWithVeto:         tally.NoWithVeto.String(),
			SubmitTime:         proposal.SubmitTime,
			DepositEndtime:     proposal.DepositEndTime,
			TotalDepositAmount: totalDepositAmount,
			TotalDepositDenom:  totalDepositDenom,
			VotingStartTime:    proposal.VotingStartTime,
			VotingEndTime:      proposal.VotingEndTime,
		})

		result = append(result, *p)
	}

	// resp, err := c.apiClient.R().Get("/gov/proposals")
	// if err != nil {
	// 	return []schema.Proposal{}, fmt.Errorf("failed to request gov proposals: %s", err)
	// }

	// var proposals []types.Proposal
	// err = json.Unmarshal(types.ReadRespWithHeight(resp).Result, &proposals)
	// if err != nil {
	// 	return []schema.Proposal{}, fmt.Errorf("failed to unmarshal gov proposals: %s", err)
	// }

	// if len(proposals) <= 0 {
	// 	return []schema.Proposal{}, nil
	// }

	// for _, proposal := range proposals {
	// 	proposalID, _ := strconv.ParseInt(proposal.ID, 10, 64)

	// 	var totalDepositAmount string
	// 	var totalDepositDenom string
	// 	if proposal.TotalDeposit != nil {
	// 		totalDepositAmount = proposal.TotalDeposit[0].Amount
	// 		totalDepositDenom = proposal.TotalDeposit[0].Denom
	// 	}

	// 	resp, err := c.apiClient.R().Get("/gov/proposals/" + proposal.ID + "/tally")
	// 	if err != nil {
	// 		return []schema.Proposal{}, fmt.Errorf("failed to request gov tally: %s", err)
	// 	}

	// 	var tally types.Tally
	// 	err = json.Unmarshal(types.ReadRespWithHeight(resp).Result, &tally)
	// 	if err != nil {
	// 		return []schema.Proposal{}, fmt.Errorf("failed to unmarshal gov tally: %s", err)
	// 	}

	// 	p := schema.NewProposal(schema.Proposal{
	// 		ID:                 proposalID,
	// 		Title:              proposal.Content.Value.Title,
	// 		Description:        proposal.Content.Value.Description,
	// 		ProposalType:       proposal.Content.Type,
	// 		ProposalStatus:     proposal.ProposalStatus,
	// 		Yes:                tally.Yes,
	// 		Abstain:            tally.Abstain,
	// 		No:                 tally.No,
	// 		NoWithVeto:         tally.NoWithVeto,
	// 		SubmitTime:         proposal.SubmitTime,
	// 		DepositEndtime:     proposal.DepositEndTime,
	// 		TotalDepositAmount: totalDepositAmount,
	// 		TotalDepositDenom:  totalDepositDenom,
	// 		VotingStartTime:    proposal.VotingStartTime,
	// 		VotingEndTime:      proposal.VotingEndTime,
	// 	})

	// 	result = append(result, *p)
	// }

	return result, nil
}
