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
	"github.com/rs/zerolog/log"
)

// GetValidators returns all existing validators
func GetValidators(db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	validatorInfo := make([]*schema.ValidatorInfo, 0)

	// check status param
	statusParam := r.FormValue("status")
	switch statusParam {
	case "":
		_ = db.Model(&validatorInfo).
			Order("id ASC").
			Select()
	case "active":
		_ = db.Model(&validatorInfo).
			Where("status = ?", 2).
			Order("id ASC").
			Select()
	case "inactive":
		_ = db.Model(&validatorInfo).
			Where("status = ? OR status = ?", 0, 1).
			Order("id ASC").
			Select()
	default:
		return json.NewEncoder(w).Encode(validatorInfo)
	}

	if len(validatorInfo) <= 0 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// query the latest block height saved in database - currently use second highest block height in database to easing client's handling
	var blockInfo []schema.BlockInfo
	_ = db.Model(&blockInfo).
		Column("height").
		Order("id DESC").
		Limit(2).
		Select()

	// Sort validatorInfo by highest tokens
	sort.Slice(validatorInfo[:], func(i, j int) bool {
		tempToken1, _ := strconv.Atoi(validatorInfo[i].Tokens)
		tempToken2, _ := strconv.Atoi(validatorInfo[j].Tokens)
		return tempToken1 > tempToken2
	})

	resultValidator := make([]*models.ResultValidator, 0)
	for _, validator := range validatorInfo {
		var missBlockCount int

		// if a validator is jailed, missing block is 100
		if validator.Jailed {
			missBlockCount = 100
		} else {
			var missDetailInfo []schema.MissDetailInfo
			missBlockCount, _ = db.Model(&missDetailInfo).
				Where("address = ? AND height BETWEEN ? AND ?", validator.Proposer, blockInfo[1].Height-int64(99), blockInfo[1].Height).
				Count()
		}

		tempUptime := &models.Uptime{
			Address:      validator.Proposer,
			MissedBlocks: int64(missBlockCount),
			OverBlocks:   100,
		}

		// Insert validator data
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
			Uptime:               *tempUptime,
			MinSelfDelegation:    validator.MinSelfDelegation,
			KeybaseURL:           validator.KeybaseURL,
		}
		resultValidator = append(resultValidator, tempResultValidator)
	}

	utils.Respond(w, resultValidator)
	return nil
}

// GetValidator receives validator address and returns that validator
func GetValidator(db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	address := vars["address"]

	// check if the input validator address exists
	validatorInfo, _ := db.ConvertToProposer(address)
	if validatorInfo.Proposer == "" {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// query the latest block height saved in database - currently use second highest block height in database to easing client's handling
	var blockInfo []schema.BlockInfo
	_ = db.Model(&blockInfo).
		Column("height").
		Order("id DESC").
		Limit(2).
		Select()

	var missBlockCount int

	// check if a validator is jailed, missing block is 100
	if validatorInfo.Jailed {
		missBlockCount = 100
	} else {
		var missDetailInfo []schema.MissDetailInfo
		missBlockCount, _ = db.Model(&missDetailInfo).
			Where("address = ? AND height BETWEEN ? AND ?", validatorInfo.Proposer, blockInfo[1].Height-int64(99), blockInfo[1].Height).
			Count()
	}

	tempUptime := &models.Uptime{
		Address:      validatorInfo.Proposer,
		MissedBlocks: int64(missBlockCount),
		OverBlocks:   100,
	}

	// validator's bonded height and its timestamp
	var validatorSetInfo schema.ValidatorSetInfo
	_ = db.Model(&validatorSetInfo).
		Where("proposer = ? AND event_type = ?", validatorInfo.Proposer, "create_validator").
		Select()

	resultValidatorDetail := &models.ResultValidatorDetail{
		Rank:                 validatorInfo.Rank,
		OperatorAddress:      validatorInfo.OperatorAddress,
		ConsensusPubkey:      validatorInfo.ConsensusPubkey,
		BondedHeight:         validatorSetInfo.Height,
		BondedTime:           validatorSetInfo.Time,
		Jailed:               validatorInfo.Jailed,
		Status:               validatorInfo.Status,
		Tokens:               validatorInfo.Tokens,
		DelegatorShares:      validatorInfo.DelegatorShares,
		Moniker:              validatorInfo.Moniker,
		Identity:             validatorInfo.Identity,
		Website:              validatorInfo.Website,
		Details:              validatorInfo.Details,
		UnbondingHeight:      validatorInfo.UnbondingHeight,
		UnbondingTime:        validatorInfo.UnbondingTime,
		CommissionRate:       validatorInfo.CommissionRate,
		CommissionMaxRate:    validatorInfo.CommissionMaxRate,
		CommissionChangeRate: validatorInfo.CommissionChangeRate,
		UpdateTime:           validatorInfo.UpdateTime,
		Uptime:               *tempUptime,
		MinSelfDelegation:    validatorInfo.MinSelfDelegation,
		KeybaseURL:           validatorInfo.KeybaseURL,
	}

	utils.Respond(w, resultValidatorDetail)
	return nil
}

// GetValidatorBlockMisses receives validator address and returns the validator's block consencutive misses
func GetValidatorBlockMisses(db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	address := vars["address"]

	// check if the input validator address exists
	validatorInfo, _ := db.ConvertToProposer(address)

	if validatorInfo.Proposer == "" {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	address = validatorInfo.Proposer

	// query a validator's missing blocks
	var missInfos []schema.MissInfo
	err := db.Model(&missInfos).
		Where("address = ?", address).
		Limit(50).
		Order("start_height DESC").
		Select()
	if err != nil {
		return err
	}

	resultMisses := make([]*models.ResultMisses, 0)
	for _, missInfo := range missInfos {
		tempResultMisses := &models.ResultMisses{
			StartHeight:  missInfo.StartHeight,
			EndHeight:    missInfo.EndHeight,
			MissingCount: missInfo.MissingCount,
			StartTime:    missInfo.StartTime,
			EndTime:      missInfo.EndTime,
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

	// query the latest block height saved in database - currently use second highest block height in database to easing client's handling
	var blockInfo []schema.BlockInfo
	_ = db.Model(&blockInfo).
		Column("height").
		Order("id DESC").
		Limit(1).
		Select()

	// query a validator's missing blocks
	var missDetailInfos []schema.MissDetailInfo
	_ = db.Model(&missDetailInfos).
		Where("address = ? AND height BETWEEN ? AND ?", validatorInfo.Proposer, blockInfo[0].Height-int64(104), blockInfo[0].Height).
		Limit(104).
		Order("height DESC").
		Select()

	resultMissesDetail := make([]*models.ResultMissesDetail, 0)
	if len(missDetailInfos) <= 0 {
		return json.NewEncoder(w).Encode(resultMissesDetail)
	}

	for _, missDetailInfo := range missDetailInfos {
		tempResultMissesDetail := &models.ResultMissesDetail{
			Height: missDetailInfo.Height,
			Time:   missDetailInfo.Time,
		}
		resultMissesDetail = append(resultMissesDetail, tempResultMissesDetail)
	}

	utils.Respond(w, resultMissesDetail)
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

	// default limit and max is 50 and offset is 0 (latest blocks)
	limit := int(50)
	offset := int(0)

	if len(r.URL.Query()["limit"]) > 0 {
		limit, _ = strconv.Atoi(r.URL.Query()["limit"][0])
	}

	if limit > 50 {
		errors.ErrOverMaxLimit(w, http.StatusRequestedRangeNotSatisfiable)
		return nil
	}

	if len(r.URL.Query()["offset"]) > 0 {
		offset, _ = strconv.Atoi(r.URL.Query()["offset"][0])
	}

	validatorID, _ := db.QueryValidatorID(address)
	if validatorID == -1 {
		fmt.Printf("failed to query the latest block height from database.")
		return nil
	}

	resultVotingPowerHistory := make([]*models.ResultVotingPowerHistory, 0)
	if validatorID != 0 {
		events, _ := db.QueryValidatorPowerEvents(validatorID, limit, offset)

		for i, event := range events {
			tempResultValidatorSet := &models.ResultVotingPowerHistory{
				ID:             i + 1,
				Height:         event.Height,
				EventType:      event.EventType,
				VotingPower:    event.VotingPower * 1000000,
				NewVotingPower: event.NewVotingPowerAmount * 1000000,
				TxHash:         event.TxHash,
				Timestamp:      event.Time,
			}
			resultVotingPowerHistory = append(resultVotingPowerHistory, tempResultValidatorSet)
		}
	}

	utils.Respond(w, resultVotingPowerHistory)
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

	resp, _ := resty.R().Get(config.Node.LCDURL + endpoint)

	var redelegations []models.Redelegations
	err := json.Unmarshal(models.ReadRespWithHeight(resp).Result, &redelegations)
	if err != nil {
		log.Info().Str(models.Service, models.LogValidator).Str(models.Method, "GetRedelegations").Err(err).Msg("unmarshal redelegations error")
	}

	utils.Respond(w, redelegations)
	return nil
}

// GetValidatorDelegations receives validator address and returns all existing delegations that are delegated to the validator
func GetValidatorDelegations(codec *codec.Codec, config *config.Config, db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	operatorAddress := vars["address"]

	validatorInfo, _ := db.ConvertToProposer(operatorAddress)

	// Check if the validator address exists
	if validatorInfo.Proposer == "" {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Query all delegations of the validator
	resp, _ := resty.R().Get(config.Node.LCDURL + "/staking/validators/" + validatorInfo.OperatorAddress + "/delegations")

	var delegations []*models.ValidatorDelegations
	err := json.Unmarshal(models.ReadRespWithHeight(resp).Result, &delegations)
	if err != nil {
		log.Info().Str(models.Service, models.LogValidator).Str(models.Method, "GetValidatorDelegations").Err(err).Msg("unmarshal delegations error")
	}

	// validator's token divide by delegator_shares equals amount of uatom
	tokens, _ := strconv.ParseFloat(validatorInfo.Tokens, 64)
	delegatorShares, _ := strconv.ParseFloat(validatorInfo.DelegatorShares, 64)
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

	// query delegation change rate in 24 hours by 24 rows order by descending id
	statsValidators24H := make([]*models.StatsValidators24H, 0)
	_ = db.Model(&statsValidators24H).
		Where("proposer = ?", validatorInfo.Proposer).
		Order("id DESC").
		Limit(2).
		Select()

	// initial variables and current delegator numbers
	delegatorNumChange24H := int(0)
	latestDelegatorNum := int(0)

	// get change delegator num in 24 hours
	if len(statsValidators24H) > 1 {
		latestDelegatorNum = statsValidators24H[0].DelegatorNum
		before24DelegatorNum := statsValidators24H[1].DelegatorNum
		delegatorNumChange24H = latestDelegatorNum - before24DelegatorNum
	}

	resultValidatorDelegations := &models.ResultValidatorDelegations{
		TotalDelegatorNum:     len(delegations),
		DelegatorNumChange24H: delegatorNumChange24H,
		ValidatorDelegations:  validatorDelegations,
	}

	utils.Respond(w, resultValidatorDelegations)
	return nil
}
