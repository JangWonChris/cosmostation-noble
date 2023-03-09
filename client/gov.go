package client

import (
	"context"
	"fmt"
	"log"

	"github.com/cosmostation/cosmostation-noble/custom"
	mdschema "github.com/cosmostation/mintscan-database/schema"
	"go.uber.org/zap"

	query "github.com/cosmos/cosmos-sdk/types/query"
	// cosmos-sdk
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramstypesproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibccoretypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
)

// GetGovQueryClient returns a object of queryClient
func (c *Client) GetGovQueryClient() govtypes.QueryClient {
	return govtypes.NewQueryClient(c.GRPC)
}

func (c *Client) GetNumberofProposals() (uint64, error) {
	queryClient := c.GetGovQueryClient()
	resp, err := queryClient.Proposals(context.Background(), &govtypes.QueryProposalsRequest{
		Pagination: &query.PageRequest{
			CountTotal: true,
			Limit:      1,
		},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to request gov proposals: %s", err)
	}

	return resp.Pagination.GetTotal(), nil
}

func (c *Client) GetAllProposals() (result []mdschema.Proposal, err error) {
	keyExists := true
	var nextKey []byte
	queryClient := c.GetGovQueryClient()
	for keyExists {
		resp, err := queryClient.Proposals(context.Background(), &govtypes.QueryProposalsRequest{
			Pagination: &query.PageRequest{
				Key:   nextKey,
				Limit: 10,
			},
		})
		if err != nil {
			return []mdschema.Proposal{}, fmt.Errorf("failed to request gov proposals: %s", err)
		}

		if len(resp.Proposals) <= 0 {
			return []mdschema.Proposal{}, nil
		}

		nextKey = resp.Pagination.GetNextKey()
		keyExists = len(nextKey) > 0

		for _, proposal := range resp.Proposals {
			chunk, err := custom.AppCodec.MarshalJSON(&proposal)
			if err != nil {
				return result, fmt.Errorf("failed to marshal proposal: %s", err)
			}
			var contentI govtypes.Content
			err = custom.AppCodec.UnpackAny(proposal.Content, &contentI)
			if err != nil {
				log.Println(err)
			}

			var proposalType string
			switch i := contentI.(type) {
			case *govtypes.TextProposal:
				proposalType = i.ProposalType()
			case *distributiontypes.CommunityPoolSpendProposal:
				proposalType = i.ProposalType()
			case *paramstypesproposal.ParameterChangeProposal:
				proposalType = i.ProposalType()
			case *upgradetypes.SoftwareUpgradeProposal:
				proposalType = i.ProposalType()
			case *upgradetypes.CancelSoftwareUpgradeProposal:
				proposalType = i.ProposalType()

			// ibc
			case *ibccoretypes.ClientUpdateProposal:
				proposalType = i.ProposalType()
			case *ibccoretypes.UpgradeProposal:
				proposalType = i.ProposalType()
			default:
				proposalType = custom.GetProposalType(i)
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
				return []mdschema.Proposal{}, fmt.Errorf("failed to request gov proposals: %s", err)
			}
			tally := tallyResultResp.Tally

			p := mdschema.Proposal{
				ID:                 proposal.ProposalId,
				Title:              proposal.GetTitle(),
				Description:        contentI.GetDescription(),
				ProposalType:       proposalType,
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
				Chunk:              chunk,
			}

			result = append(result, p)
		}
	}

	return result, nil
}

// GetProposals은 MBL GetProposals()를 wrap한 함수
func (c *Client) GetProposalsByStatus(s govtypes.ProposalStatus) (result []mdschema.Proposal, err error) {
	keyExists := true
	var nextKey []byte
	queryClient := c.GetGovQueryClient()
	for keyExists {
		resp, err := queryClient.Proposals(context.Background(), &govtypes.QueryProposalsRequest{
			ProposalStatus: s,
			Pagination: &query.PageRequest{
				Limit: 1000,
			},
		})
		if err != nil {
			return []mdschema.Proposal{}, fmt.Errorf("failed to request gov proposals: %s", err)
		}

		if len(resp.Proposals) <= 0 {
			return []mdschema.Proposal{}, nil
		}

		nextKey = resp.Pagination.GetNextKey()
		keyExists = len(nextKey) > 0

		for _, proposal := range resp.Proposals {
			var contentI govtypes.Content
			err = custom.AppCodec.UnpackAny(proposal.Content, &contentI)
			if err != nil {
				log.Println(err)
			}

			var proposalType string
			switch i := contentI.(type) {
			case *govtypes.TextProposal:
				proposalType = i.ProposalType()
			case *distributiontypes.CommunityPoolSpendProposal:
				proposalType = i.ProposalType()
			case *paramstypesproposal.ParameterChangeProposal:
				proposalType = i.ProposalType()
			case *upgradetypes.SoftwareUpgradeProposal:
				proposalType = i.ProposalType()
			case *upgradetypes.CancelSoftwareUpgradeProposal:
				proposalType = i.ProposalType()

			// ibc
			case *ibccoretypes.ClientUpdateProposal:
				proposalType = i.ProposalType()
			case *ibccoretypes.UpgradeProposal:
				proposalType = i.ProposalType()
			default:
				proposalType = custom.GetProposalType(i)
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
				return []mdschema.Proposal{}, fmt.Errorf("failed to request gov proposals: %s", err)
			}
			tally := tallyResultResp.Tally

			p := mdschema.Proposal{
				ID:                 proposal.ProposalId,
				Title:              proposal.GetTitle(),
				Description:        contentI.GetDescription(),
				ProposalType:       proposalType,
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
	}

	return result, nil
}

// GetProposal은 특정 프로포절 정보를 GRPC 얻어온다.
func (c *Client) GetProposal(id uint64) (result *mdschema.Proposal, err error) {

	queryClient := c.GetGovQueryClient()
	resp, err := queryClient.Proposal(context.Background(), &govtypes.QueryProposalRequest{ProposalId: id})
	if err != nil {
		return &mdschema.Proposal{}, fmt.Errorf("failed to request gov proposals: %s", err)
	}

	chunk, err := custom.AppCodec.MarshalJSON(&resp.Proposal)
	if err != nil {
		return result, fmt.Errorf("failed to marshal proposal: %s", err)
	}

	var contentI govtypes.Content
	err = custom.AppCodec.UnpackAny(resp.Proposal.Content, &contentI)
	if err != nil {
		zap.S().Error(err)
		return &mdschema.Proposal{}, err
	}

	var proposalType string
	switch i := contentI.(type) {
	case *govtypes.TextProposal:
		proposalType = i.ProposalType()
	case *distributiontypes.CommunityPoolSpendProposal:
		proposalType = i.ProposalType()
	case *paramstypesproposal.ParameterChangeProposal:
		proposalType = i.ProposalType()
	case *upgradetypes.SoftwareUpgradeProposal:
		proposalType = i.ProposalType()
	case *upgradetypes.CancelSoftwareUpgradeProposal:
		proposalType = i.ProposalType()

	// ibc
	case *ibccoretypes.ClientUpdateProposal:
		proposalType = i.ProposalType()
	case *ibccoretypes.UpgradeProposal:
		proposalType = i.ProposalType()

	default:
		proposalType = custom.GetProposalType(i)
	}

	totalDepositAmount := make([]string, len(resp.Proposal.TotalDeposit))
	totalDepositDenom := make([]string, len(resp.Proposal.TotalDeposit))
	for i, td := range resp.Proposal.TotalDeposit {
		totalDepositAmount[i] = td.Amount.String()
		totalDepositDenom[i] = td.Denom
	}

	request := govtypes.QueryTallyResultRequest{ProposalId: resp.Proposal.ProposalId}
	tallyResultResp, err := queryClient.TallyResult(context.Background(), &request)
	if err != nil {
		return &mdschema.Proposal{}, fmt.Errorf("failed to request gov proposals: %s", err)
	}
	tally := tallyResultResp.Tally

	p := &mdschema.Proposal{
		ID:                 resp.Proposal.ProposalId,
		Title:              resp.Proposal.GetTitle(),
		Description:        contentI.GetDescription(),
		ProposalType:       proposalType,
		ProposalStatus:     resp.Proposal.Status.String(),
		Yes:                tally.Yes.String(),
		Abstain:            tally.Abstain.String(),
		No:                 tally.No.String(),
		NoWithVeto:         tally.NoWithVeto.String(),
		SubmitTime:         resp.Proposal.SubmitTime,
		DepositEndTime:     resp.Proposal.DepositEndTime,
		TotalDepositAmount: totalDepositAmount,
		TotalDepositDenom:  totalDepositDenom,
		VotingStartTime:    resp.Proposal.VotingStartTime,
		VotingEndTime:      resp.Proposal.VotingEndTime,
		Chunk:              chunk,
	}

	return p, nil
}
