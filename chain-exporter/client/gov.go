package client

import (
	"context"
	"fmt"
	"log"

	"github.com/cosmostation/cosmostation-cosmos/chain-config/custom"
	"github.com/cosmostation/mintscan-database/schema"

	// cosmos-sdk
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibccoretypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	paramstypesproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// GetGovQueryClient returns a object of queryClient
func (c *Client) GetGovQueryClient() govtypes.QueryClient {
	return govtypes.NewQueryClient(c.GRPC)
}

// GetProposals은 MBL GetProposals()를 wrap한 함수
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
		var contentI govtypes.Content
		err = custom.AppCodec.UnpackAny(proposal.Content, &contentI)
		if err != nil {
			log.Println(err)
		}

		// kind of proposals
		// distributiontypes.CommunityPoolSpendProposal
		// govtypes.TextProposal
		// ibccoretypes.ClientUpdateProposal
		// paramstypes.ParameterChangeProposal
		// upgradetypes.SoftwareUpgradeProposal
		// upgradetypes.CancelSoftwareUpgradeProposal
		var proposalType string
		switch i := contentI.(type) {
		case *govtypes.TextProposal:
			proposalType = i.ProposalType()
		case *distributiontypes.CommunityPoolSpendProposal:
			proposalType = i.ProposalType()
		case *ibccoretypes.ClientUpdateProposal:
			proposalType = i.ProposalType()
		case *paramstypesproposal.ParameterChangeProposal:
			proposalType = i.ProposalType()
		case *upgradetypes.SoftwareUpgradeProposal:
			proposalType = i.ProposalType()
		case *upgradetypes.CancelSoftwareUpgradeProposal:
			proposalType = i.ProposalType()
		default:
			log.Printf("unrecognized type : %T\n", i)
		}

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

		p := schema.Proposal{
			ID:           proposal.ProposalId,
			Title:        proposal.GetTitle(),
			Description:  contentI.GetDescription(),
			ProposalType: proposalType,
			// Description:        proposal.GetContent().GetDescription(),
			// ProposalType:       proposal.GetContent().ProposalType(),
			ProposalStatus:     proposal.Status.String(),
			Yes:                tally.Yes.String(),
			Abstain:            tally.Abstain.String(),
			No:                 tally.No.String(),
			NoWithVeto:         tally.NoWithVeto.String(),
			SubmitTime:         proposal.SubmitTime,
			DepositEndTime:     proposal.DepositEndTime,
			TotalDepositAmount: totalDepositAmount,
			TotalDepositDenom:  totalDepositDenom,
			VotingStartTime:    proposal.VotingStartTime,
			VotingEndTime:      proposal.VotingEndTime,
		}

		result = append(result, p)
	}

	return result, nil
}
