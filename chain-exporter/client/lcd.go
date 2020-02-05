package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"

	resty "gopkg.in/resty.v1"
)

// Tx queries for a transaction from the REST client and decodes it into a sdk.Tx
// if the transaction exists. An error is returned if the tx doesn't exist or
// decoding fails.
func (c Client) Tx(hash string) (sdk.TxResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/txs/%s", c.clientNode, hash))
	if err != nil {
		return sdk.TxResponse{}, err
	}

	defer resp.Body.Close()

	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	var tx sdk.TxResponse

	if err := c.cdc.UnmarshalJSON(bz, &tx); err != nil {
		return sdk.TxResponse{}, err
	}

	return tx, nil
}

// SaveProposals saves governance proposals in database
func (c Client) SaveProposals() {
	resp, err := resty.R().Get(c.clientNode + "/gov/proposals")
	if err != nil {
		fmt.Printf("failed to request /gov/proposals: %v \n", err)
	}

	proposals := make([]*types.Proposal, 0)
	err = json.Unmarshal(types.ReadRespWithHeight(resp).Result, &proposals)
	if err != nil {
		fmt.Printf("failed to unmarshal Proposal: %v \n", err)
	}

	// proposal information for our database table
	result := make([]*schema.Proposal, 0)

	if len(proposals) > 0 {
		for _, proposal := range proposals {
			proposalID, _ := strconv.ParseInt(proposal.ID, 10, 64)

			var totalDepositAmount string
			var totalDepositDenom string
			if proposal.TotalDeposit != nil {
				totalDepositAmount = proposal.TotalDeposit[0].Amount
				totalDepositDenom = proposal.TotalDeposit[0].Denom
			}

			tallyResp, _ := resty.R().Get(c.clientNode + "/gov/proposals/" + proposal.ID + "/tally")

			var tally types.Tally
			err = json.Unmarshal(types.ReadRespWithHeight(tallyResp).Result, &tally)
			if err != nil {
				fmt.Printf("failed to unmarshal Tally: %v \n", err)
			}

			tempProposal := &schema.Proposal{
				ID:                 proposalID,
				Title:              proposal.Content.Value.Title,
				Description:        proposal.Content.Value.Description,
				ProposalType:       proposal.Content.Type,
				ProposalStatus:     proposal.ProposalStatus,
				Yes:                tally.Yes,
				Abstain:            tally.Abstain,
				No:                 tally.No,
				NoWithVeto:         tally.NoWithVeto,
				SubmitTime:         proposal.SubmitTime,
				DepositEndtime:     proposal.DepositEndTime,
				TotalDepositAmount: totalDepositAmount,
				TotalDepositDenom:  totalDepositDenom,
				VotingStartTime:    proposal.VotingStartTime,
				VotingEndTime:      proposal.VotingEndTime,
				Alerted:            false,
			}

			result = append(result, tempProposal)
		}
	}

	if len(result) > 0 {
		for _, proposal := range result {
			exist, _ := c.db.ExistProposal(proposal.ID)

			if exist {
				result, _ := c.db.UpdateProposal(proposal)
				if !result {
					log.Printf("failed to update Proposal ID: %d", proposal.ID)
				}
			} else {
				result, _ := c.db.InsertProposal(proposal)
				if !result {
					log.Printf("failed to save Proposal ID: %d", proposal.ID)
				}
			}
		}
	}
}

// SaveBondedValidators saves bonded validators information in database
func (c Client) SaveBondedValidators() {
	resp, _ := resty.R().Get(c.clientNode + "/staking/validators?status=bonded")

	var bondedValidators []*types.Validator
	err := json.Unmarshal(types.ReadRespWithHeight(resp).Result, &bondedValidators)
	if err != nil {
		fmt.Printf("failed to request /staking/validators?status=bonded: %v \n", err)
	}

	// sort out bondedValidators by highest tokens
	sort.Slice(bondedValidators[:], func(i, j int) bool {
		tempToken1, _ := strconv.Atoi(bondedValidators[i].Tokens)
		tempToken2, _ := strconv.Atoi(bondedValidators[j].Tokens)
		return tempToken1 > tempToken2
	})

	// bondedValidator information for our database table
	validator := make([]*schema.Validator, 0)

	for i, bondedValidator := range bondedValidators {
		tempValidator := &schema.Validator{
			Rank:                 i + 1,
			OperatorAddress:      bondedValidator.OperatorAddress,
			Address:              utils.AccAddressFromOperatorAddress(bondedValidator.OperatorAddress),
			ConsensusPubkey:      bondedValidator.ConsensusPubkey,
			Proposer:             utils.ConsAddrFromConsPubkey(bondedValidator.ConsensusPubkey),
			Jailed:               bondedValidator.Jailed,
			Status:               bondedValidator.Status,
			Tokens:               bondedValidator.Tokens,
			DelegatorShares:      bondedValidator.DelegatorShares,
			Moniker:              bondedValidator.Description.Moniker,
			Identity:             bondedValidator.Description.Identity,
			Website:              bondedValidator.Description.Website,
			Details:              bondedValidator.Description.Details,
			UnbondingHeight:      bondedValidator.UnbondingHeight,
			UnbondingTime:        bondedValidator.UnbondingTime,
			CommissionRate:       bondedValidator.Commission.CommissionRates.Rate,
			CommissionMaxRate:    bondedValidator.Commission.CommissionRates.MaxRate,
			CommissionChangeRate: bondedValidator.Commission.CommissionRates.MaxChangeRate,
			MinSelfDelegation:    bondedValidator.MinSelfDelegation,
			UpdateTime:           bondedValidator.Commission.UpdateTime,
		}
		validator = append(validator, tempValidator)
	}

	if len(validator) > 0 {
		result, err := c.db.InsertOrUpdateValidators(validator)
		if !result {
			log.Printf("failed to insert or update bonded validators: %t", err)
		}
	}
}

// SaveUnbondingAndUnBondedValidators saves unbonding and unbonded validators information in database
func (c Client) SaveUnbondingAndUnBondedValidators() {
	resp, _ := resty.R().Get(c.clientNode + "/staking/validators?status=unbonding")

	var unbondingValidators []*types.Validator
	err := json.Unmarshal(types.ReadRespWithHeight(resp).Result, &unbondingValidators)
	if err != nil {
		fmt.Printf("failed to request /staking/validators?status=unbonding: %v \n", err)
	}

	// sort out bondedValidators by highest tokens
	sort.Slice(unbondingValidators[:], func(i, j int) bool {
		tempToken1, _ := strconv.Atoi(unbondingValidators[i].Tokens)
		tempToken2, _ := strconv.Atoi(unbondingValidators[j].Tokens)
		return tempToken1 > tempToken2
	})

	// validators information for our database table
	validator := make([]*schema.Validator, 0)
	if len(unbondingValidators) > 0 {
		for _, unbondingValidator := range unbondingValidators {
			tempValidator := &schema.Validator{
				OperatorAddress:      unbondingValidator.OperatorAddress,
				Address:              utils.AccAddressFromOperatorAddress(unbondingValidator.OperatorAddress),
				ConsensusPubkey:      unbondingValidator.ConsensusPubkey,
				Proposer:             utils.ConsAddrFromConsPubkey(unbondingValidator.ConsensusPubkey),
				Jailed:               unbondingValidator.Jailed,
				Status:               unbondingValidator.Status,
				Tokens:               unbondingValidator.Tokens,
				DelegatorShares:      unbondingValidator.DelegatorShares,
				Moniker:              unbondingValidator.Description.Moniker,
				Identity:             unbondingValidator.Description.Identity,
				Website:              unbondingValidator.Description.Website,
				Details:              unbondingValidator.Description.Details,
				UnbondingHeight:      unbondingValidator.UnbondingHeight,
				UnbondingTime:        unbondingValidator.UnbondingTime,
				CommissionRate:       unbondingValidator.Commission.CommissionRates.Rate,
				CommissionMaxRate:    unbondingValidator.Commission.CommissionRates.MaxRate,
				CommissionChangeRate: unbondingValidator.Commission.CommissionRates.MaxChangeRate,
				MinSelfDelegation:    unbondingValidator.MinSelfDelegation,
				UpdateTime:           unbondingValidator.Commission.UpdateTime,
			}
			validator = append(validator, tempValidator)
		}
	} else {
		// save unbonded validators after succesfully saved unbonding validators
		c.saveUnbondedValidators()
	}

	// first rank
	status := 2
	rank, _ := c.db.QueryFirstRankValidatorByStatus(status)

	for i, validator := range validator {
		validator.Rank = (rank.Rank + 1 + i)
	}

	// save and update validatorInfo
	if len(validator) > 0 {
		result, err := c.db.InsertOrUpdateValidators(validator)
		if !result {
			log.Printf("failed to insert or update unbonding validators: %t", err)
		}

		// save unbonded validators after succesfully saved unbonding validators
		c.saveUnbondedValidators()
	}
}

// saveUnbondedValidators saves unbonded validators information in database
func (c Client) saveUnbondedValidators() {
	resp, _ := resty.R().Get(c.clientNode + "/staking/validators?status=unbonded")

	var unbondedValidators []*types.Validator
	err := json.Unmarshal(types.ReadRespWithHeight(resp).Result, &unbondedValidators)
	if err != nil {
		fmt.Printf("failed to request /staking/validators?status=unbonded: %v \n", err)
	}

	// sort out bondedValidators by highest tokens
	sort.Slice(unbondedValidators[:], func(i, j int) bool {
		tempToken1, _ := strconv.Atoi(unbondedValidators[i].Tokens)
		tempToken2, _ := strconv.Atoi(unbondedValidators[j].Tokens)
		return tempToken1 > tempToken2
	})

	// validators information for our database table
	validator := make([]*schema.Validator, 0)
	if len(unbondedValidators) > 0 {
		for _, unbondedValidator := range unbondedValidators {
			tempValidator := &schema.Validator{
				OperatorAddress:      unbondedValidator.OperatorAddress,
				Address:              utils.AccAddressFromOperatorAddress(unbondedValidator.OperatorAddress),
				ConsensusPubkey:      unbondedValidator.ConsensusPubkey,
				Proposer:             utils.ConsAddrFromConsPubkey(unbondedValidator.ConsensusPubkey),
				Jailed:               unbondedValidator.Jailed,
				Status:               unbondedValidator.Status,
				Tokens:               unbondedValidator.Tokens,
				DelegatorShares:      unbondedValidator.DelegatorShares,
				Moniker:              unbondedValidator.Description.Moniker,
				Identity:             unbondedValidator.Description.Identity,
				Website:              unbondedValidator.Description.Website,
				Details:              unbondedValidator.Description.Details,
				UnbondingHeight:      unbondedValidator.UnbondingHeight,
				UnbondingTime:        unbondedValidator.UnbondingTime,
				CommissionRate:       unbondedValidator.Commission.CommissionRates.Rate,
				CommissionMaxRate:    unbondedValidator.Commission.CommissionRates.MaxRate,
				CommissionChangeRate: unbondedValidator.Commission.CommissionRates.MaxChangeRate,
				MinSelfDelegation:    unbondedValidator.MinSelfDelegation,
				UpdateTime:           unbondedValidator.Commission.UpdateTime,
			}
			validator = append(validator, tempValidator)
		}
	}

	// first rank
	status := 1
	rank, _ := c.db.QueryFirstRankValidatorByStatus(status)

	for i, validator := range validator {
		validator.Rank = (rank.Rank + 1 + i)
	}

	// save and update validator
	if len(validator) > 0 {
		result, err := c.db.InsertOrUpdateValidators(validator)
		if !result {
			log.Printf("failed to insert or update unbonded validators: %t", err)
		}
	}
}
