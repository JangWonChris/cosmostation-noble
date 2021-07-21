package common

import (
	"context"
	"net/http"
	"strconv"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmostation/cosmostation-cosmos/app"
	"github.com/cosmostation/cosmostation-cosmos/errors"
	"github.com/cosmostation/cosmostation-cosmos/model"
	mdschema "github.com/cosmostation/mintscan-database/schema"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

// [TODO] Cosmos SDK v0.38.x 부터 Proposal이 Voting 단계로 넘어가지 않았을 경우 REST API에서 삭제한다.
// Chain Exporter가 재시작을 한번이라도 했을 경우 디비에 쌓이질 않는다. REST API에만 의존할까?

// GetProposals returns all existing proposals
func GetProposals(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		proposals, err := a.DB.GetProposals()
		if err != nil {
			zap.L().Error("failed to query proposals", zap.Error(err))
			errors.ErrInternalServer(rw, http.StatusInternalServerError)
			return
		}

		if len(proposals) <= 0 {
			zap.L().Debug("found no proposals saved in database")
			model.Respond(rw, []mdschema.Proposal{})
			return
		}

		result := make([]*model.ResultProposal, 0)

		for _, p := range proposals {
			val, err := a.DB.GetValidatorByAnyAddr(p.Proposer)
			if err != nil {
				zap.S().Errorf("failed to query validator information: %s", err)
				return
			}

			proposal := &model.ResultProposal{
				ProposalID:           p.ID,
				TxHash:               p.TxHash,
				Proposer:             p.Proposer,
				Moniker:              val.Moniker,
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
				DepositEndtime:       p.DepositEndTime,
				VotingStartTime:      p.VotingStartTime,
				VotingEndTime:        p.VotingEndTime,
			}

			result = append(result, proposal)
		}

		model.Respond(rw, result)
		return
	}
}

// GetProposal receives proposal id and returns particular proposal
func GetProposal(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["proposal_id"]
		new_id, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			zap.S().Errorf("failed to convert proposal parameter(string to int): %s", err)
			errors.ErrInternalServer(rw, http.StatusInternalServerError)
			return
		}
		// Query particular proposal
		p, _ := a.DB.GetProposal(new_id)
		if p.ID == 0 {
			zap.L().Debug("this proposal does not exist", zap.String("id", id))
			errors.ErrNotExist(rw, http.StatusNotFound)
			return
		}

		// Error doesn't need to be handled since any accoount can propose proposal
		val, err := a.DB.GetValidatorByAnyAddr(p.Proposer)
		if err != nil {
			zap.S().Errorf("failed to query validator information: %s", err)
			return
		}

		result := &model.ResultProposal{
			ProposalID:           p.ID,
			TxHash:               p.TxHash,
			Proposer:             p.Proposer,
			Moniker:              val.Moniker,
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
			DepositEndtime:       p.DepositEndTime,
			VotingStartTime:      p.VotingStartTime,
			VotingEndTime:        p.VotingEndTime,
		}

		model.Respond(rw, result)
		return
	}
}

// GetDeposits receives proposal id and returns deposit information
func GetDeposits(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["proposal_id"]
		new_id, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			zap.S().Errorf("failed to convert proposal parameter(string to int): %s", err)
			errors.ErrInternalServer(rw, http.StatusInternalServerError)
			return
		}
		// Query particular proposal
		p, _ := a.DB.GetProposal(new_id)
		if p.ID == 0 {
			zap.L().Debug("this proposal does not exist", zap.String("id", id))
			errors.ErrNotExist(rw, http.StatusNotFound)
			return
		}

		result := make([]*model.ResultDeposit, 0)

		deposits, _ := a.DB.GetDeposits(new_id)
		if len(deposits) <= 0 {
			zap.L().Debug("this proposal does not have any deposit yet", zap.String("id", id))
			model.Respond(rw, result)
			return
		}

		for _, d := range deposits {
			val, err := a.DB.GetValidatorByAnyAddr(d.Depositor)
			if err != nil {
				zap.S().Errorf("failed to query validator information: %s", err)
				return
			}

			deposit := &model.ResultDeposit{
				Depositor:     d.Depositor,
				Moniker:       val.Moniker,
				DepositAmount: d.Amount,
				DepositDenom:  d.Denom,
				Height:        d.Height,
				TxHash:        d.TxHash,
				Timestamp:     d.Timestamp,
			}

			result = append(result, deposit)
		}

		model.Respond(rw, result)
		return
	}
}

// GetVotes returns vote transactions with the given proposal id
func GetVotes(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["proposal_id"]
		new_id, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			zap.S().Errorf("failed to convert proposal parameter(string to int): %s", err)
			errors.ErrInternalServer(rw, http.StatusInternalServerError)
			return
		}
		// Query particular proposal
		p, _ := a.DB.GetProposal(new_id)
		if p.ID == 0 {
			zap.L().Debug("this proposal does not exist", zap.String("id", id))
			errors.ErrNotExist(rw, http.StatusNotFound)
			return
		}

		// Query all votes
		votes, _ := a.DB.GetVotes(new_id)
		if len(votes) <= 0 {
			model.Respond(rw, &model.ResultVote{
				Tally: &model.ResultTally{},
				Votes: []*model.Votes{},
			})
			return
		}

		rv := make([]*model.Votes, 0)

		for _, v := range votes {
			// val, err := a.DB.GetValidatorByAnyAddr(v.Voter)
			// if err != nil {
			// 	zap.S().Errorf("failed to query validator information: %s", err)
			// 	return
			// }

			vote := &model.Votes{
				Voter: v.Voter,
				// Moniker: val.Moniker,
				Moniker: "",
				Option:  v.Option,
				TxHash:  v.TxHash,
				Time:    v.Timestamp,
			}

			rv = append(rv, vote)
		}

		queryClient := govtypes.NewQueryClient(a.Client.GRPC)
		proposalID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			zap.L().Error("failed to convert proposal id ", zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		}
		request := govtypes.QueryTallyResultRequest{ProposalId: proposalID}
		res, err := queryClient.TallyResult(context.Background(), &request)
		if err != nil {
			zap.L().Error("failed to get tally info", zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			return
		}

		// Query vote options for the proposal
		yes, no, noWithVeto, abstain, err := a.DB.QueryVoteOptions(p.ID)
		if err != nil {
			zap.L().Error("failed to get count of tally option", zap.Error(err))
			errors.ErrNotExist(rw, http.StatusNotFound)
			return
		}

		rt := &model.ResultTally{
			YesAmount:        res.Tally.Yes.ToDec().String(),
			NoAmount:         res.Tally.No.String(),
			AbstainAmount:    res.Tally.Abstain.String(),
			NoWithVetoAmount: res.Tally.NoWithVeto.String(),
			YesNum:           yes,
			AbstainNum:       abstain,
			NoNum:            no,
			NoWithVetoNum:    noWithVeto,
		}

		result := &model.ResultVote{
			Tally: rt,
			Votes: rv,
		}

		model.Respond(rw, result)
		return
	}
}
