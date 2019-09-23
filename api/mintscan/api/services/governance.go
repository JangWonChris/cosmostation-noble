package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	errors "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/types"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/libs/bech32"
	resty "gopkg.in/resty.v1"
)

// GetProposals returns all existing proposals
func GetProposals(db *pg.DB, config *config.Config, w http.ResponseWriter, r *http.Request) error {
	// Query all proposals
	proposalInfo := make([]*types.ProposalInfo, 0)
	_ = db.Model(&proposalInfo).Select()

	// Check if any proposal exists
	if len(proposalInfo) <= 0 {
		return json.NewEncoder(w).Encode(proposalInfo)
	}

	resultProposal := make([]*models.ResultProposal, 0)
	for _, proposal := range proposalInfo {
		// Convert Cosmos Address to Opeartor Address
		_, decoded, _ := bech32.DecodeAndConvert(proposal.Proposer)
		cosmosOperAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixValAddr, decoded)

		// Check if the address matches any moniker in our DB
		var validatorInfo types.ValidatorInfo
		_ = db.Model(&validatorInfo).
			Column("moniker").
			Where("operator_address = ?", cosmosOperAddress).
			Limit(1).
			Select()

		// Insert proposal data
		tempProposal := &models.ResultProposal{
			ProposalID:           proposal.ID,
			TxHash:               proposal.TxHash,
			Proposer:             proposal.Proposer,
			Moniker:              validatorInfo.Moniker,
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
		resultProposal = append(resultProposal, tempProposal)
	}

	utils.Respond(w, resultProposal)
	return nil
}

// GetProposal receives proposal id and returns particular proposal
func GetProposal(db *pg.DB, config *config.Config, w http.ResponseWriter, r *http.Request) error {
	// Receive proposal id
	vars := mux.Vars(r)
	proposalID := vars["proposalId"]

	// Query particular proposal
	var proposalInfo types.ProposalInfo
	err := db.Model(&proposalInfo).
		Where("id = ?", proposalID).
		Select()
	if err != nil {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Convert Cosmos Address to Opeartor Address
	_, decoded, _ := bech32.DecodeAndConvert(proposalInfo.Proposer)
	cosmosOperAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixValAddr, decoded)

	// Check if the address matches any moniker in our DB
	var validatorInfo types.ValidatorInfo
	_ = db.Model(&validatorInfo).
		Column("moniker").
		Where("operator_address = ?", cosmosOperAddress).
		Limit(1).
		Select()

	resultProposal := &models.ResultProposal{
		ProposalID:           proposalInfo.ID,
		TxHash:               proposalInfo.TxHash,
		Proposer:             proposalInfo.Proposer,
		Moniker:              validatorInfo.Moniker,
		Title:                proposalInfo.Title,
		Description:          proposalInfo.Description,
		ProposalType:         proposalInfo.ProposalType,
		ProposalStatus:       proposalInfo.ProposalStatus,
		Yes:                  proposalInfo.Yes,
		Abstain:              proposalInfo.Abstain,
		No:                   proposalInfo.No,
		NoWithVeto:           proposalInfo.NoWithVeto,
		InitialDepositAmount: proposalInfo.InitialDepositAmount,
		InitialDepositDenom:  proposalInfo.InitialDepositDenom,
		TotalDepositAmount:   proposalInfo.TotalDepositAmount,
		TotalDepositDenom:    proposalInfo.TotalDepositDenom,
		SubmitTime:           proposalInfo.SubmitTime,
		DepositEndtime:       proposalInfo.DepositEndtime,
		VotingStartTime:      proposalInfo.VotingStartTime,
		VotingEndTime:        proposalInfo.VotingEndTime,
	}

	utils.Respond(w, resultProposal)
	return nil
}

// GetProposalVotes receives proposal id and returns voting information
func GetProposalVotes(db *pg.DB, config *config.Config, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	proposalID := vars["proposalId"]

	// check if proposal id exists
	var proposalInfo types.ProposalInfo
	err := db.Model(&proposalInfo).
		Where("id = ?", proposalID).
		Select()
	if err != nil {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// query all votes
	voteInfo := make([]*types.VoteInfo, 0)
	_ = db.Model(&voteInfo).
		Where("proposal_id = ?", proposalID).
		Order("id DESC").
		Select()

	// check if votes exists
	if len(voteInfo) <= 0 {
		return json.NewEncoder(w).Encode(&models.ResultVote{
			Tally: &models.ResultTally{},
			Votes: []*models.Votes{},
		})
	}

	// query count for respective votes
	yesCnt, _ := db.Model(&voteInfo).
		Where("proposal_id = ? AND option = ?", proposalID, "Yes").
		Count()
	abstainCnt, _ := db.Model(&voteInfo).
		Where("proposal_id = ? AND option = ?", proposalID, "Abstain").
		Count()
	noCnt, _ := db.Model(&voteInfo).
		Where("proposal_id = ? AND option = ?", proposalID, "No").
		Count()
	noWithVetoCnt, _ := db.Model(&voteInfo).
		Where("proposal_id = ? AND option = ?", proposalID, "NoWithVeto").
		Count()

	// votes
	votes := make([]*models.Votes, 0)
	for _, vote := range voteInfo {
		moniker, _ := utils.ConvertCosmosAddressToMoniker(vote.Voter, db)

		tempVoteInfo := &models.Votes{
			Voter:   vote.Voter,
			Moniker: moniker,
			Option:  vote.Option,
			TxHash:  vote.TxHash,
			Time:    vote.Time,
		}
		votes = append(votes, tempVoteInfo)
	}

	// query tally information
	resp, err := resty.R().Get(config.Node.LCDURL + "/gov/proposals/" + proposalID + "/tally")

	var responseWithHeight types.ResponseWithHeight
	err = json.Unmarshal(resp.Body(), &responseWithHeight)
	if err != nil {
		fmt.Printf("unmarshal responseWithHeight error - %v\n", err)
	}

	var tallyInfo types.Tally
	err = json.Unmarshal(responseWithHeight.Result, &tallyInfo)
	if err != nil {
		fmt.Printf("Proposal unmarshal error - %v\n", err)
	}

	tempResultTally := &models.ResultTally{
		YesAmount:        tallyInfo.Yes,
		NoAmount:         tallyInfo.No,
		AbstainAmount:    tallyInfo.Abstain,
		NoWithVetoAmount: tallyInfo.NoWithVeto,
		YesNum:           yesCnt,
		AbstainNum:       abstainCnt,
		NoNum:            noCnt,
		NoWithVetoNum:    noWithVetoCnt,
	}

	resultVote := &models.ResultVote{
		Tally: tempResultTally,
		Votes: votes,
	}

	utils.Respond(w, resultVote)
	return nil
}

// GetProposalDeposits receives proposal id and returns deposit information
func GetProposalDeposits(db *pg.DB, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	proposalID := vars["proposalId"]

	// check if proposal id exists
	var proposalInfo types.ProposalInfo
	err := db.Model(&proposalInfo).
		Where("id = ?", proposalID).
		Select()
	if err != nil {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	resultDepositInfo := make([]*models.ResultDeposit, 0)

	// query all deposit info
	depositInfo := make([]*types.DepositInfo, 0)
	_ = db.Model(&depositInfo).
		Where("proposal_id = ?", proposalID).
		Order("id DESC").
		Select()

	// check if the deposit exists
	if len(depositInfo) <= 0 {
		return json.NewEncoder(w).Encode(resultDepositInfo)
	}

	// deposits
	for _, deposit := range depositInfo {
		moniker, _ := utils.ConvertCosmosAddressToMoniker(deposit.Depositor, db)

		tempResultDeposit := &models.ResultDeposit{
			Depositor:     deposit.Depositor,
			Moniker:       moniker,
			DepositAmount: deposit.Amount,
			DepositDenom:  deposit.Denom,
			Height:        deposit.Height,
			TxHash:        deposit.TxHash,
			Time:          deposit.Time,
		}
		resultDepositInfo = append(resultDepositInfo, tempResultDeposit)
	}

	utils.Respond(w, resultDepositInfo)
	return nil
}
