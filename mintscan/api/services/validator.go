package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/db"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/schema"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/utils"
	resty "gopkg.in/resty.v1"

	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/rpc/client"

	"github.com/cosmos/cosmos-sdk/codec"
)

// GetValidators returns all existing validators
func GetValidators(db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	validators := make([]schema.ValidatorInfo, 0)

	status := r.FormValue("status")

	switch status {
	case "":
		validators, _ = db.QueryValidators()
	case "active":
		validators, _ = db.QueryActiveValidators()
	case "inactive":
		validators, _ = db.QueryInActiveValidators()
	default:
		return json.NewEncoder(w).Encode(validators)
	}

	if len(validators) <= 0 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Query the latest two blocks information
	blocks := db.QueryLastestTwoBlocks()

	// Sort validators by their tokens in descending order
	sort.Slice(validators[:], func(i, j int) bool {
		tempToken1, _ := strconv.Atoi(validators[i].Tokens)
		tempToken2, _ := strconv.Atoi(validators[j].Tokens)
		return tempToken1 > tempToken2
	})

	result := make([]*models.ResultValidator, 0)

	for _, validator := range validators {
		var missBlockCount int

		// if a validator is jailed, missing block is 100
		if validator.Jailed {
			missBlockCount = 100
		} else {
			missBlockCount, _ = db.QueryMissingBlocksCount(validator.Proposer, int(blocks[1].Height), 99)
		}

		uptime := &models.Uptime{
			Address:      validator.Proposer,
			MissedBlocks: missBlockCount,
			OverBlocks:   100,
		}

		tempResultValidator := &models.ResultValidator{
			Rank:                 validator.Rank,
			OperatorAddress:      validator.OperatorAddress,
			ConsensusPubkey:      validator.ConsensusPubkey,
			Jailed:               validator.Jailed,
			Status:               validator.Status,
			Tokens:               validator.Tokens,
			DelegatorShares:      validator.DelegatorShares,
			Moniker:              validator.Moniker,
			Identity:             validator.Identity,
			Website:              validator.Website,
			Details:              validator.Details,
			UnbondingHeight:      validator.UnbondingHeight,
			UnbondingTime:        validator.UnbondingTime,
			CommissionRate:       validator.CommissionRate,
			CommissionMaxRate:    validator.CommissionMaxRate,
			CommissionChangeRate: validator.CommissionChangeRate,
			UpdateTime:           validator.UpdateTime,
			Uptime:               *uptime,
			MinSelfDelegation:    validator.MinSelfDelegation,
			KeybaseURL:           validator.KeybaseURL,
		}
		result = append(result, tempResultValidator)
	}

	utils.Respond(w, result)
	return nil
}

// GetValidator receives validator address and returns that validator
func GetValidator(db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	address := vars["address"]

	// Check if the validator address exists
	validator, _ := db.ConvertToProposer(address)
	if validator.Proposer == "" {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// query the latest block height saved in database - currently use second highest block height in database to easing client's handling
	blocks := db.QueryLastestTwoBlocks()

	var missBlockCount int

	// check if a validator is jailed, missing block is 100
	if validator.Jailed {
		missBlockCount = 100
	} else {
		missBlockCount, _ = db.QueryMissingBlocksCount(validator.Proposer, int(blocks[1].Height), 99)
	}

	tempUptime := &models.Uptime{
		Address:      validator.Proposer,
		MissedBlocks: missBlockCount,
		OverBlocks:   100,
	}

	// Query a validator's bonded information
	validatorSetInfo := db.QueryValidatorBondedInfo(validator.Proposer)

	result := &models.ResultValidatorDetail{
		Rank:                 validator.Rank,
		OperatorAddress:      validator.OperatorAddress,
		ConsensusPubkey:      validator.ConsensusPubkey,
		BondedHeight:         validatorSetInfo.Height,
		BondedTime:           validatorSetInfo.Time,
		Jailed:               validator.Jailed,
		Status:               validator.Status,
		Tokens:               validator.Tokens,
		DelegatorShares:      validator.DelegatorShares,
		Moniker:              validator.Moniker,
		Identity:             validator.Identity,
		Website:              validator.Website,
		Details:              validator.Details,
		UnbondingHeight:      validator.UnbondingHeight,
		UnbondingTime:        validator.UnbondingTime,
		CommissionRate:       validator.CommissionRate,
		CommissionMaxRate:    validator.CommissionMaxRate,
		CommissionChangeRate: validator.CommissionChangeRate,
		UpdateTime:           validator.UpdateTime,
		Uptime:               *tempUptime,
		MinSelfDelegation:    validator.MinSelfDelegation,
		KeybaseURL:           validator.KeybaseURL,
	}

	utils.Respond(w, result)
	return nil
}

// GetValidatorBlockMisses receives validator address and returns the validator's block consencutive misses
func GetValidatorBlockMisses(db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	address := vars["address"]

	// Check if the validator address exists
	validatorInfo, _ := db.ConvertToProposer(address)

	if validatorInfo.Proposer == "" {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	address = validatorInfo.Proposer

	// Query a range of missing blocks that a validator has missed
	limit := int(50)
	blocks, _ := db.QueryMissingBlocks(address, limit)

	resultMisses := make([]*models.ResultMisses, 0)
	for _, block := range blocks {
		tempResultMisses := &models.ResultMisses{
			StartHeight:  block.StartHeight,
			EndHeight:    block.EndHeight,
			MissingCount: block.MissingCount,
			StartTime:    block.StartTime,
			EndTime:      block.EndTime,
		}
		resultMisses = append(resultMisses, tempResultMisses)
	}

	utils.Respond(w, resultMisses)
	return nil
}

// GetValidatorBlockMissesDetail receives validator address and returns the validator's block misses (uptime)
// When uptime is 100%, there is no missing blocks in database therefore it returns an empty array.
// When uptime is from 1% through 99%, just return how many missing blocks are recorded in database.
// When uptime is 0%, there are two cases. Missing last 100 blocks or a validator is unbonded | unbonding state.
func GetValidatorBlockMissesDetail(db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	address := vars["address"]

	// check if the validator exists
	validatorInfo, _ := db.ConvertToProposer(address)
	if validatorInfo.Proposer == "" {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	address = validatorInfo.Proposer

	latestHeight, _ := db.QueryLatestBlockHeight()
	if latestHeight == -1 {
		fmt.Printf("failed to query latest block height")
	}

	result := make([]*models.ResultMissesDetail, 0)

	// Query validator's missing blocks for the last 100
	// Note: use second highest block height in database to ease client side's handling
	blocks, _ := db.QueryMissingBlocksInDetail(address, latestHeight, 104)
	if len(blocks) <= 0 {
		return json.NewEncoder(w).Encode(result)
	}

	for _, block := range blocks {
		tempResultMissesDetail := &models.ResultMissesDetail{
			Height: block.Height,
			Time:   block.Time,
		}
		result = append(result, tempResultMissesDetail)
	}

	utils.Respond(w, result)
	return nil
}

// GetValidatorEvents receives validator address and returns the validator's events
func GetValidatorEvents(db *db.Database, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	address := vars["address"]

	// Check if the address is validator
	validatorInfo, _ := db.ConvertToProposer(address)
	if validatorInfo.Proposer == "" {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	address = validatorInfo.Proposer

	limit := int(50) // default limit is 50
	before := int(0)
	after := int(0)
	offset := int(0)

	if len(r.URL.Query()["limit"]) > 0 {
		limit, _ = strconv.Atoi(r.URL.Query()["limit"][0])
	}

	if len(r.URL.Query()["before"]) > 0 {
		before, _ = strconv.Atoi(r.URL.Query()["before"][0])
	}

	if len(r.URL.Query()["after"]) > 0 {
		after, _ = strconv.Atoi(r.URL.Query()["after"][0])
	}

	if len(r.URL.Query()["offset"]) > 0 {
		offset, _ = strconv.Atoi(r.URL.Query()["offset"][0])
	}

	if limit > 50 {
		errors.ErrOverMaxLimit(w, http.StatusRequestedRangeNotSatisfiable)
		return nil
	}

	// Edge case
	// Some validators existed in cosmoshub-1 or cosmoshub-2 but not in cosmoshub-3 won't have any power event history
	// Return empty array for client to handle this
	validatorID, _ := db.QueryValidatorID(address)
	if validatorID == 0 {
		utils.Respond(w, []models.ResultVotingPowerHistory{})
		return nil
	}

	if validatorID == -1 {
		errors.ErrInternalServer(w, http.StatusInternalServerError)
	}

	events := make([]schema.ValidatorSetInfo, 0)

	switch {
	case before > 0:
		events, _ = db.QueryValidatorPowerEvents(validatorID, limit, before, after, offset)
	case after > 0:
		events, _ = db.QueryValidatorPowerEvents(validatorID, limit, before, after, offset)
	case offset >= 0:
		events, _ = db.QueryValidatorPowerEvents(validatorID, limit, before, after, offset)
	}

	result := make([]*models.ResultVotingPowerHistory, 0)

	for i, event := range events {
		tempResultValidatorSet := &models.ResultVotingPowerHistory{
			ID:             i + 1,
			Height:         event.Height,
			EndHeight:      1,
			EventType:      event.EventType,
			VotingPower:    event.VotingPower * 1000000,
			NewVotingPower: event.NewVotingPowerAmount * 1000000,
			TxHash:         event.TxHash,
			Timestamp:      event.Time,
		}
		result = append(result, tempResultValidatorSet)
	}

	utils.Respond(w, result)
	return nil
}

// GetRedelegations receives delegator, srcvalidator, dstvalidator address and returns redelegations information
func GetRedelegations(config *config.Config, db *db.Database, w http.ResponseWriter, r *http.Request) error {
	endpoint := "/staking/redelegations?"
	if len(r.URL.Query()["delegator"]) > 0 {
		endpoint += fmt.Sprintf("delegator=%s&", r.URL.Query()["delegator"][0])
	}

	if len(r.URL.Query()["validator_from"]) > 0 {
		endpoint += fmt.Sprintf("validator_from=%s&", r.URL.Query()["validator_from"][0])
	}

	if len(r.URL.Query()["validator_to"]) > 0 {
		endpoint += fmt.Sprintf("validator_to=%s&", r.URL.Query()["validator_to"][0])
	}

	resp, _ := resty.R().Get(config.Node.LCDEndpoint + endpoint)

	var redelegations []models.Redelegations
	err := json.Unmarshal(models.ReadRespWithHeight(resp).Result, &redelegations)
	if err != nil {
		fmt.Printf("failed to unmarshal redelegations: %t\n", err)
	}

	utils.Respond(w, redelegations)
	return nil
}

// GetValidatorDelegations receives validator address and returns all existing delegations that are delegated to the validator
func GetValidatorDelegations(codec *codec.Codec, config *config.Config, db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	operatorAddress := vars["address"]

	validator, _ := db.ConvertToProposer(operatorAddress)

	// Check if the validator address exists
	if validator.Proposer == "" {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Query all delegations of the validator
	resp, _ := resty.R().Get(config.Node.LCDEndpoint + "/staking/validators/" + validator.OperatorAddress + "/delegations")

	var delegations []*models.ValidatorDelegations
	err := json.Unmarshal(models.ReadRespWithHeight(resp).Result, &delegations)
	if err != nil {
		fmt.Printf("failed to unmarshal delegations: %t\n", err)
	}

	// Calculate the amount of uatom, which should divide validator's token by delegator_shares
	tokens, _ := strconv.ParseFloat(validator.Tokens, 64)
	delegatorShares, _ := strconv.ParseFloat(validator.DelegatorShares, 64)
	uatom := tokens / delegatorShares

	var validatorDelegations []*models.ValidatorDelegations
	if len(delegations) > 0 {
		for _, delegation := range delegations {
			shares, _ := strconv.ParseFloat(delegation.Shares.String(), 64)
			amount := fmt.Sprintf("%f", shares*uatom)

			tempValidatorDelegations := &models.ValidatorDelegations{
				DelegatorAddress: delegation.DelegatorAddress,
				ValidatorAddress: delegation.ValidatorAddress,
				Shares:           delegation.Shares,
				Amount:           amount,
			}
			validatorDelegations = append(validatorDelegations, tempValidatorDelegations)
		}
	}

	// Query delegation change rate in 24 hours by 24 rows order by descending id
	validatorStats, _ := db.QueryValidatorStats24H(validator.Proposer, 2)

	delegatorNumChange24H := int(0)
	latestDelegatorNum := int(0)

	if len(validatorStats) > 1 {
		latestDelegatorNum = validatorStats[0].DelegatorNum
		before24DelegatorNum := validatorStats[1].DelegatorNum
		delegatorNumChange24H = latestDelegatorNum - before24DelegatorNum
	}

	result := &models.ResultValidatorDelegations{
		TotalDelegatorNum:     len(delegations),
		DelegatorNumChange24H: delegatorNumChange24H,
		ValidatorDelegations:  validatorDelegations,
	}

	utils.Respond(w, result)
	return nil
}
