package services

import (
	"encoding/json"
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/db"
	errors "github.com/cosmostation/cosmostation-cosmos/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/schema"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/libs/bech32"
	resty "gopkg.in/resty.v1"

	"github.com/rs/zerolog/log"
)

// GetProposals returns all existing proposals
func GetProposals(db *db.Database, config *config.Config, w http.ResponseWriter, r *http.Request) error {
	proposalInfo := make([]*schema.ProposalInfo, 0)
	_ = db.Model(&proposalInfo).Select()

	// check if any proposal exists
	if len(proposalInfo) <= 0 {
		return json.NewEncoder(w).Encode(proposalInfo)
	}

	resultProposal := make([]*models.ResultProposal, 0)

	for _, proposal := range proposalInfo {
		// convert to validator operator address
		_, decoded, _ := bech32.DecodeAndConvert(proposal.Proposer)
		cosmosOperAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixValAddr, decoded)

		// check if the address matches any moniker in our database
		var validatorInfo schema.ValidatorInfo
		_ = db.Model(&validatorInfo).
			Column("moniker").
			Where("operator_address = ?", cosmosOperAddress).
			Limit(1).
			Select()

		// insert proposal data
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
func GetProposal(db *db.Database, config *config.Config, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	proposalID := vars["proposalId"]

	// query particular proposal
	var proposalInfo schema.ProposalInfo
	err := db.Model(&proposalInfo).
		Where("id = ?", proposalID).
		Select()
	if err != nil {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// convert to validator operator address
	_, decoded, _ := bech32.DecodeAndConvert(proposalInfo.Proposer)
	cosmosOperAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixValAddr, decoded)

	// check if the address matches any moniker in our DB
	var validatorInfo schema.ValidatorInfo
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

// GetVotes receives proposal id and returns voting information
func GetVotes(db *db.Database, config *config.Config, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	proposalID := vars["proposalId"]

	// check if proposal id exists
	var proposalInfo schema.ProposalInfo
	err := db.Model(&proposalInfo).
		Where("id = ?", proposalID).
		Select()
	if err != nil {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// query all votes
	voteInfo := make([]*schema.VoteInfo, 0)
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
		moniker, _ := db.ConvertCosmosAddressToMoniker(vote.Voter)

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

	var tally models.Tally
	err = json.Unmarshal(models.ReadRespWithHeight(resp).Result, &tally)
	if err != nil {
		log.Info().Str(models.Service, models.LogGovernance).Str(models.Method, "GetVotes").Err(err).Msg("unmarshal tally error")
	}

	tempResultTally := &models.ResultTally{
		YesAmount:        tally.Yes,
		NoAmount:         tally.No,
		AbstainAmount:    tally.Abstain,
		NoWithVetoAmount: tally.NoWithVeto,
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

// GetDeposits receives proposal id and returns deposit information
func GetDeposits(db *db.Database, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	proposalID := vars["proposalId"]

	// check if proposal id exists
	var proposalInfo schema.ProposalInfo
	err := db.Model(&proposalInfo).
		Where("id = ?", proposalID).
		Select()
	if err != nil {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	resultDepositInfo := make([]*models.ResultDeposit, 0)

	// query all deposit info
	depositInfo := make([]*schema.DepositInfo, 0)
	_ = db.Model(&depositInfo).
		Where("proposal_id = ?", proposalID).
		Order("id DESC").
		Select()

	// check if the deposit exists
	if len(depositInfo) <= 0 {
		return json.NewEncoder(w).Encode(resultDepositInfo)
	}

	for _, deposit := range depositInfo {
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
		resultDepositInfo = append(resultDepositInfo, tempResultDeposit)
	}

	utils.Respond(w, resultDepositInfo)
	return nil
}
