package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/db"
	errors "github.com/cosmostation/cosmostation-cosmos/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/utils"

	"github.com/gorilla/mux"
	resty "gopkg.in/resty.v1"
)

// GetProposals returns all existing proposals
func GetProposals(db *db.Database, config *config.Config, w http.ResponseWriter, r *http.Request) error {
	// Query proposals
	proposals := db.QueryProposals()

	if len(proposals) <= 0 {
		return json.NewEncoder(w).Encode(proposals)
	}

	result := make([]*models.ResultProposal, 0)

	for _, proposal := range proposals {
		// Check if the address matches any moniker in our database
		operAddr := utils.ValAddressFromAccAddress(proposal.Proposer)
		validator, _ := db.QueryValidatorByOperAddr(operAddr)

		// insert proposal data
		tempProposal := &models.ResultProposal{
			ProposalID:           proposal.ID,
			TxHash:               proposal.TxHash,
			Proposer:             proposal.Proposer,
			Moniker:              validator.Moniker,
			Title:                proposal.Title,
			Description:          proposal.Description,
			ProposalType:         proposal.ProposalType,
			ProposalStatus:       proposal.ProposalStatus,
			Yes:                  proposal.Yes,
			Abstain:              proposal.Abstain,
			No:                   proposal.No,
			NoWithVeto:           proposal.NoWithVeto,
			InitialDepositAmount: proposal.InitialDepositAmount,
			InitialDepositDenom:  proposal.InitialDepositDenom,
			TotalDepositAmount:   proposal.TotalDepositAmount,
			TotalDepositDenom:    proposal.TotalDepositDenom,
			SubmitTime:           proposal.SubmitTime,
			DepositEndtime:       proposal.DepositEndtime,
			VotingStartTime:      proposal.VotingStartTime,
			VotingEndTime:        proposal.VotingEndTime,
		}
		result = append(result, tempProposal)
	}

	utils.Respond(w, result)
	return nil
}

// GetProposal receives proposal id and returns particular proposal
func GetProposal(db *db.Database, config *config.Config, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	proposalID := vars["proposalId"]

	// Query particular proposal
	proposal := db.QueryProposal(proposalID)
	if proposal.ID == 0 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Check if the address matches any moniker in our database
	operAddr := utils.ValAddressFromAccAddress(proposal.Proposer)
	validator, _ := db.QueryValidatorByOperAddr(operAddr)

	result := &models.ResultProposal{
		ProposalID:           proposal.ID,
		TxHash:               proposal.TxHash,
		Proposer:             proposal.Proposer,
		Moniker:              validator.Moniker,
		Title:                proposal.Title,
		Description:          proposal.Description,
		ProposalType:         proposal.ProposalType,
		ProposalStatus:       proposal.ProposalStatus,
		Yes:                  proposal.Yes,
		Abstain:              proposal.Abstain,
		No:                   proposal.No,
		NoWithVeto:           proposal.NoWithVeto,
		InitialDepositAmount: proposal.InitialDepositAmount,
		InitialDepositDenom:  proposal.InitialDepositDenom,
		TotalDepositAmount:   proposal.TotalDepositAmount,
		TotalDepositDenom:    proposal.TotalDepositDenom,
		SubmitTime:           proposal.SubmitTime,
		DepositEndtime:       proposal.DepositEndtime,
		VotingStartTime:      proposal.VotingStartTime,
		VotingEndTime:        proposal.VotingEndTime,
	}

	utils.Respond(w, result)
	return nil
}

// GetVotes receives proposal id and returns voting information
func GetVotes(db *db.Database, config *config.Config, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	proposalID := vars["proposalId"]

	// Check if the proposal exists
	exists := db.QueryExistsProposal(proposalID)
	if !exists {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Query all votes
	votes := db.QueryVotes(proposalID)

	if len(votes) <= 0 {
		utils.Respond(w, &models.ResultVote{
			Tally: &models.ResultTally{},
			Votes: []*models.Votes{},
		})
		return nil
	}

	resultVotes := make([]*models.Votes, 0)

	for _, vote := range votes {
		moniker, _ := db.ConvertCosmosAddressToMoniker(vote.Voter)

		tempVoteInfo := &models.Votes{
			Voter:   vote.Voter,
			Moniker: moniker,
			Option:  vote.Option,
			TxHash:  vote.TxHash,
			Time:    vote.Time,
		}
		resultVotes = append(resultVotes, tempVoteInfo)
	}

	// Query tally information
	resp, err := resty.R().Get(config.Node.LCDURL + "/gov/proposals/" + proposalID + "/tally")

	var tally models.Tally
	err = json.Unmarshal(models.ReadRespWithHeight(resp).Result, &tally)
	if err != nil {
		fmt.Printf("failed to unmarshal tally: %t", err)
	}

	// Query vote options for the proposal
	yes, no, noWithVeto, abstain := db.QueryVoteOptions(proposalID)

	resultTally := &models.ResultTally{
		YesAmount:        tally.Yes,
		NoAmount:         tally.No,
		AbstainAmount:    tally.Abstain,
		NoWithVetoAmount: tally.NoWithVeto,
		YesNum:           yes,
		AbstainNum:       abstain,
		NoNum:            no,
		NoWithVetoNum:    noWithVeto,
	}

	result := &models.ResultVote{
		Tally: resultTally,
		Votes: resultVotes,
	}

	utils.Respond(w, result)
	return nil
}

// GetDeposits receives proposal id and returns deposit information
func GetDeposits(db *db.Database, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	proposalID := vars["proposalId"]

	// Check if the proposal exists
	exists := db.QueryExistsProposal(proposalID)
	if !exists {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	result := make([]*models.ResultDeposit, 0)

	// Query all deposits
	deposits := db.QueryDeposits(proposalID)

	if len(deposits) <= 0 {
		utils.Respond(w, result)
		return nil
	}

	for _, deposit := range deposits {
		moniker, _ := db.ConvertCosmosAddressToMoniker(deposit.Depositor)

		tempResultDeposit := &models.ResultDeposit{
			Depositor:     deposit.Depositor,
			Moniker:       moniker,
			DepositAmount: deposit.Amount,
			DepositDenom:  deposit.Denom,
			Height:        deposit.Height,
			TxHash:        deposit.TxHash,
			Time:          deposit.Time,
		}
		result = append(result, tempResultDeposit)
	}

	utils.Respond(w, result)
	return nil
}
