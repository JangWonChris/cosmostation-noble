package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/schema"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

// [TODO] Cosmos SDK v0.38.x 부터 Proposal이 Voting 단계로 넘어가지 않았을 경우 REST API에서 삭제한다.
// Chain Exporter가 재시작을 한번이라도 했을 경우 디비에 쌓이질 않는다. REST API에만 의존할까?

// GetProposals returns all existing proposals
func GetProposals(rw http.ResponseWriter, r *http.Request) {
	proposals, err := s.db.QueryProposals()
	if err != nil {
		zap.L().Error("failed to query proposals", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	if len(proposals) <= 0 {
		zap.L().Debug("found no proposals saved in database")
		model.Respond(rw, []schema.Proposal{})
		return
	}

	result := make([]*model.ResultProposal, 0)

	for _, p := range proposals {
		// Error doesn't need to be handled since any accoount can propose proposal
		val, _ := s.db.QueryValidatorByAny(p.Proposer)

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
			DepositEndtime:       p.DepositEndtime,
			VotingStartTime:      p.VotingStartTime,
			VotingEndTime:        p.VotingEndTime,
		}

		result = append(result, proposal)
	}

	model.Respond(rw, result)
	return
}

// GetProposal receives proposal id and returns particular proposal
func GetProposal(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["proposal_id"]

	// Query particular proposal
	p, _ := s.db.QueryProposal(id)
	if p.ID == 0 {
		zap.L().Debug("this proposal does not exist", zap.String("id", id))
		errors.ErrNotExist(rw, http.StatusNotFound)
		return
	}

	// Error doesn't need to be handled since any accoount can propose proposal
	val, _ := s.db.QueryValidatorByAny(p.Proposer)

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
		DepositEndtime:       p.DepositEndtime,
		VotingStartTime:      p.VotingStartTime,
		VotingEndTime:        p.VotingEndTime,
	}

	model.Respond(rw, result)
	return
}

// GetDeposits receives proposal id and returns deposit information
func GetDeposits(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["proposal_id"]

	// Query particular proposal
	p, _ := s.db.QueryProposal(id)
	if p.ID == 0 {
		zap.L().Debug("this proposal does not exist", zap.String("id", id))
		errors.ErrNotExist(rw, http.StatusNotFound)
		return
	}

	result := make([]*model.ResultDeposit, 0)

	deposits, _ := s.db.QueryDeposits(id)
	if len(deposits) <= 0 {
		zap.L().Debug("this proposal does not have any deposit yet", zap.String("id", id))
		model.Respond(rw, result)
		return
	}

	for _, d := range deposits {
		// Error doesn't need to be handled since any accoount can propose proposal
		val, _ := s.db.QueryValidatorByAny(d.Depositor)

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

// GetVotes returns vote transactions with the given proposal id
func GetVotes(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["proposal_id"]

	// Query particular proposal
	p, _ := s.db.QueryProposal(id)
	if p.ID == 0 {
		zap.L().Debug("this proposal does not exist", zap.String("id", id))
		errors.ErrNotExist(rw, http.StatusNotFound)
		return
	}

	// Query all votes
	votes, _ := s.db.QueryVotes(id)
	if len(votes) <= 0 {
		model.Respond(rw, &model.ResultVote{
			Tally: &model.ResultTally{},
			Votes: []*model.Votes{},
		})
		return
	}

	rv := make([]*model.Votes, 0)

	for _, v := range votes {
		// Error doesn't need to be handled since any accoount can propose proposal
		val, _ := s.db.QueryValidatorByAny(v.Voter)

		vote := &model.Votes{
			Voter:   v.Voter,
			Moniker: val.Moniker,
			Option:  v.Option,
			TxHash:  v.TxHash,
			Time:    v.Timestamp,
		}

		rv = append(rv, vote)
	}

	// Query tally information.
	resp, err := s.client.HandleResponseHeight("/gov/proposals/" + id + "/tally")
	if err != nil {
		zap.L().Error("failed to get tally info", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	var tally model.Tally
	err = json.Unmarshal(resp.Result, &tally)
	if err != nil {
		zap.L().Error("failed to unmarshal tally", zap.Error(err))
		errors.ErrFailedUnmarshalJSON(rw, http.StatusInternalServerError)
		return
	}

	// Query vote options for the proposal
	yes, no, noWithVeto, abstain := s.db.QueryVoteOptions(id)

	rt := &model.ResultTally{
		YesAmount:        tally.Yes,
		NoAmount:         tally.No,
		AbstainAmount:    tally.Abstain,
		NoWithVetoAmount: tally.NoWithVeto,
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
