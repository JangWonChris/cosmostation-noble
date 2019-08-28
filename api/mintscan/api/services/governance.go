package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	errors "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models"
	dbtypes "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/types"
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
	proposalInfo := make([]*dbtypes.ProposalInfo, 0)
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
		var validatorInfo dbtypes.ValidatorInfo
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
	var proposalInfo dbtypes.ProposalInfo
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
	var validatorInfo dbtypes.ValidatorInfo
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
	// Receive proposal id
	vars := mux.Vars(r)
	proposalID := vars["proposalId"]

	// Check if proposal id exists
	var proposalInfo dbtypes.ProposalInfo
	err := db.Model(&proposalInfo).
		Where("id = ?", proposalID).
		Select()
	if err != nil {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Query all votes
	voteInfo := make([]*dbtypes.VoteInfo, 0)
	_ = db.Model(&voteInfo).
		Where("proposal_id = ?", proposalID).
		Order("id DESC").
		Select()

	// Check if votes exists
	if len(voteInfo) <= 0 {
		return json.NewEncoder(w).Encode(&models.ResultVoteInfo{
			Tally: &models.Tally{},
			Votes: []*models.Votes{},
		})
	}

	// Query count for respective votes
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

	// Votes
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

	var responseWithHeight models.ResponseWithHeight
	err = json.Unmarshal(resp.Body(), &responseWithHeight)
	if err != nil {
		fmt.Printf("unmarshal responseWithHeight error - %v\n", err)
	}

	var tallyInfo models.TallyInfo
	err = json.Unmarshal(responseWithHeight.Result, &tallyInfo)
	if err != nil {
		fmt.Printf("Proposal unmarshal error - %v\n", err)
	}

	// Tally
	tempTally := &models.Tally{
		YesAmount:        tallyInfo.Yes,
		NoAmount:         tallyInfo.No,
		AbstainAmount:    tallyInfo.Abstain,
		NoWithVetoAmount: tallyInfo.NoWithVeto,
		YesNum:           yesCnt,
		AbstainNum:       abstainCnt,
		NoNum:            noCnt,
		NoWithVetoNum:    noWithVetoCnt,
	}

	resultVoteInfo := &models.ResultVoteInfo{
		Tally: tempTally,
		Votes: votes,
	}

	utils.Respond(w, resultVoteInfo)
	return nil
}

// GetProposalDeposits receives proposal id and returns deposit information
func GetProposalDeposits(db *pg.DB, w http.ResponseWriter, r *http.Request) error {
	// Receive proposal id
	vars := mux.Vars(r)
	proposalID := vars["proposalId"]

	// Check if proposal id exists
	var proposalInfo dbtypes.ProposalInfo
	err := db.Model(&proposalInfo).
		Where("id = ?", proposalID).
		Select()
	if err != nil {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Result Response
	resultDepositInfo := make([]*models.DepositInfo, 0)

	// Query all deposit info
	depositInfo := make([]*dbtypes.DepositInfo, 0)
	_ = db.Model(&depositInfo).
		Where("proposal_id = ?", proposalID).
		Order("id DESC").
		Select()

	// Check if the deposit exists
	if len(depositInfo) <= 0 {
		return json.NewEncoder(w).Encode(resultDepositInfo)
	}

	for _, deposit := range depositInfo {
		// Convert Cosmos Address to Opeartor Address
		moniker, _ := utils.ConvertCosmosAddressToMoniker(deposit.Depositor, db)

		// Insert deposits
		tempDepositInfo := &models.DepositInfo{
			Depositor:     deposit.Depositor,
			Moniker:       moniker,
			DepositAmount: deposit.Amount,
			DepositDenom:  deposit.Denom,
			Height:        deposit.Height,
			TxHash:        deposit.TxHash,
			Time:          deposit.Time,
		}
		resultDepositInfo = append(resultDepositInfo, tempDepositInfo)
	}

	utils.Respond(w, resultDepositInfo)
	return nil
}
