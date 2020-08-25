package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/schema"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

// GetValidators returns all validators.
func GetValidators(rw http.ResponseWriter, r *http.Request) {
	var status string

	if len(r.URL.Query()["status"]) > 0 {
		status = r.URL.Query()["status"][0]
	}

	vals := make([]*schema.Validator, 0)

	switch status {
	case model.ActiveValidator:
		vals, _ = s.db.QueryValidatorsByStatus(model.BondedValidatorStatus)
	case model.InactiveValidator:
		unbondingVals, _ := s.db.QueryValidatorsByStatus(model.UnbondingValidatorStatus)
		unbondedVals, _ := s.db.QueryValidatorsByStatus(model.UnbondedValidatorStatus)
		vals = append(vals, unbondingVals...)
		vals = append(vals, unbondedVals...)
	default:
		vals, _ = s.db.QueryValidators()
	}

	if len(vals) <= 0 {
		zap.L().Debug("there are no validators exist in database")
		errors.ErrNotExist(rw, http.StatusNotFound)
		return
	}

	// Sort validators by their tokens in descending order.
	sort.Slice(vals[:], func(i, j int) bool {
		tk1, _ := strconv.Atoi(vals[i].Tokens)
		tk2, _ := strconv.Atoi(vals[j].Tokens)
		return tk1 > tk2
	})

	result := make([]*model.ResultValidator, 0)

	for _, val := range vals {
		var missBlockCount int

		if val.Status == model.BondedValidatorStatus {
			blocks, _ := s.db.QueryLastestTwoBlocks()
			missBlockCount, _ = s.db.CountMissingBlocks(val.Proposer, int(blocks[1].Height), 99)
		} else {
			missBlockCount = model.MissingAllBlocks
		}

		uptime := &model.Uptime{
			Address:      val.Proposer,
			MissedBlocks: missBlockCount,
			OverBlocks:   model.MissingAllBlocks,
		}

		val := &model.ResultValidator{
			Rank:                 val.Rank,
			AccountAddress:       val.Address,
			OperatorAddress:      val.OperatorAddress,
			ConsensusPubkey:      val.ConsensusPubkey,
			Jailed:               val.Jailed,
			Status:               val.Status,
			Tokens:               val.Tokens,
			DelegatorShares:      val.DelegatorShares,
			Moniker:              val.Moniker,
			Identity:             val.Identity,
			Website:              val.Website,
			Details:              val.Details,
			UnbondingHeight:      val.UnbondingHeight,
			UnbondingTime:        val.UnbondingTime,
			CommissionRate:       val.CommissionRate,
			CommissionMaxRate:    val.CommissionMaxRate,
			CommissionChangeRate: val.CommissionChangeRate,
			UpdateTime:           val.UpdateTime,
			Uptime:               uptime,
			MinSelfDelegation:    val.MinSelfDelegation,
			KeybaseURL:           val.KeybaseURL,
		}
		result = append(result, val)
	}

	model.Respond(rw, result)
	return
}

// GetValidator returns a validator information.
func GetValidator(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	val, err := s.db.QueryValidatorByAny(address)
	if err != nil {
		zap.L().Debug("failed to query validator informatino", zap.Error(err))
		errors.ErrNotExist(rw, http.StatusNotFound)
		return
	}

	var missBlockCount int
	if val.Status == model.BondedValidatorStatus {
		blocks, err := s.db.QueryLastestTwoBlocks()
		if err != nil {
			zap.L().Error("failed to get latest two blocks", zap.Error(err))
		}
		missBlockCount, _ = s.db.CountMissingBlocks(val.Proposer, int(blocks[1].Height), 99)
	} else {
		missBlockCount = model.MissingAllBlocks
	}

	uptime := &model.Uptime{
		Address:      val.Proposer,
		MissedBlocks: missBlockCount,
		OverBlocks:   model.MissingAllBlocks,
	}

	// query a validator's bonded information
	powerEventHistory, _ := s.db.QueryValidatorBondedInfo(val.Proposer)

	result := &model.ResultValidatorDetail{
		Rank:                 val.Rank,
		AccountAddress:       val.Address,
		OperatorAddress:      val.OperatorAddress,
		ConsensusPubkey:      val.ConsensusPubkey,
		BondedHeight:         powerEventHistory.Height,
		BondedTime:           powerEventHistory.Timestamp,
		Jailed:               val.Jailed,
		Status:               val.Status,
		Tokens:               val.Tokens,
		DelegatorShares:      val.DelegatorShares,
		Moniker:              val.Moniker,
		Identity:             val.Identity,
		Website:              val.Website,
		Details:              val.Details,
		UnbondingHeight:      val.UnbondingHeight,
		UnbondingTime:        val.UnbondingTime,
		CommissionRate:       val.CommissionRate,
		CommissionMaxRate:    val.CommissionMaxRate,
		CommissionChangeRate: val.CommissionChangeRate,
		UpdateTime:           val.UpdateTime,
		Uptime:               uptime,
		MinSelfDelegation:    val.MinSelfDelegation,
		KeybaseURL:           val.KeybaseURL,
	}

	model.Respond(rw, result)
	return
}

// GetValidatorUptime returns a validator's missing blocks in detail
// Uptime is 100%: there is no missing blocks in database therefore it returns an empty array.
// Uptime is from 1% through 99%: just return how many missing blocks are recorded in database.
// Uptime is 0%: there are two cases. Missing last 100 blocks or a validator is unbonded | unbonding state.
func GetValidatorUptime(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	val, err := s.db.QueryValidatorByAny(address)
	if err != nil {
		zap.L().Debug("failed to query validator info", zap.Error(err))
		errors.ErrNotExist(rw, http.StatusNotFound)
		return
	}

	latestHeight, _ := s.db.QueryLatestBlockHeight()
	if latestHeight == -1 {
		zap.L().Debug("failed to query latest block height", zap.Int64("latestHeight", latestHeight))
	}

	result := make([]*model.ResultMissesDetail, 0)

	// Query validator's missing blocks for the last 100
	// Note: use second highest block height in database to ease client side's handling
	blocks, _ := s.db.QueryMissingBlocksDetail(val.Proposer, latestHeight, 104)
	if len(blocks) <= 0 {
		model.Respond(rw, result)
		return
	}

	for _, block := range blocks {
		missDetail := &model.ResultMissesDetail{
			Height:    block.Height,
			Timestamp: block.Timestamp,
		}
		result = append(result, missDetail)
	}

	model.Respond(rw, result)
	return
}

// GetValidatorUptimeRange returns the validator's block consencutive misses.
func GetValidatorUptimeRange(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	val, err := s.db.QueryValidatorByAny(address)
	if err != nil {
		zap.L().Debug("failed to query validator info", zap.Error(err))
		errors.ErrNotExist(rw, http.StatusNotFound)
		return
	}

	// TODO: move to somewhere else?
	limit := int(50)

	blocks, _ := s.db.QueryMissingBlocks(val.Proposer, limit)

	result := make([]*model.ResultMisses, 0)

	for _, b := range blocks {
		miss := &model.ResultMisses{
			StartHeight:  b.StartHeight,
			EndHeight:    b.EndHeight,
			MissingCount: b.MissingCount,
			StartTime:    b.StartTime,
			EndTime:      b.EndTime,
		}
		result = append(result, miss)
	}

	model.Respond(rw, result)
	return
}

// GetValidatorDelegations receives validator address and returns all existing delegations that are delegated to the validator.
func GetValidatorDelegations(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	val, err := s.db.QueryValidatorByAny(address)
	if err != nil {
		zap.L().Error("failed to query validator info", zap.Error(err))
		errors.ErrNotExist(rw, http.StatusNotFound)
		return
	}

	// get all delegations from a validator
	resp, err := s.client.HandleResponseHeight("/staking/validators/" + val.OperatorAddress + "/delegations")
	if err != nil {
		zap.L().Error("failed to get delegations from a validator", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	var delegations []*model.ValidatorDelegations
	err = json.Unmarshal(resp.Result, &delegations)
	if err != nil {
		zap.L().Error("failed to unmarshal delegations", zap.Error(err))
		errors.ErrFailedUnmarshalJSON(rw, http.StatusInternalServerError)
		return
	}

	// calculate the amount of uatom, which should divide validator's token by delegator_shares
	tokens, _ := strconv.ParseFloat(val.Tokens, 64)
	delegatorShares, _ := strconv.ParseFloat(val.DelegatorShares, 64)
	uatom := tokens / delegatorShares

	var validatorDelegations []*model.ValidatorDelegations
	if len(delegations) > 0 {
		for _, delegation := range delegations {
			shares, _ := strconv.ParseFloat(delegation.Shares.String(), 64)
			amount := fmt.Sprintf("%f", shares*uatom)

			temp := &model.ValidatorDelegations{
				DelegatorAddress: delegation.DelegatorAddress,
				ValidatorAddress: delegation.ValidatorAddress,
				Shares:           delegation.Shares,
				Amount:           amount,
			}
			validatorDelegations = append(validatorDelegations, temp)
		}
	}

	// query delegation change rate in 24 hours by 24 rows order by descending id
	validatorStats, _ := s.db.QueryValidatorStats1D(val.Proposer, 2)

	delegatorNumChange24H := int(0)
	latestDelegatorNum := int(0)

	if len(validatorStats) > 1 {
		latestDelegatorNum = validatorStats[0].DelegatorNum
		before24DelegatorNum := validatorStats[1].DelegatorNum
		delegatorNumChange24H = latestDelegatorNum - before24DelegatorNum
	}

	result := &model.ResultValidatorDelegations{
		TotalDelegatorNum:     len(delegations),
		DelegatorNumChange24H: delegatorNumChange24H,
		ValidatorDelegations:  validatorDelegations,
	}

	model.Respond(rw, result)
	return
}

// GetValidatorPowerHistoryEvents receives validator address and returns the validator's events.
func GetValidatorPowerHistoryEvents(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	before, after, limit, err := model.ParseHTTPArgsWithBeforeAfterLimit(r, model.DefaultBefore, model.DefaultAfter, model.DefaultPowerEventHistoryLimit)
	if err != nil {
		zap.S().Debug("failed to parse HTTP args ", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "request is invalid")
		return
	}

	if limit > 50 {
		zap.S().Debug("failed to query with this limit ", zap.Int("request limit", limit))
		errors.ErrOverMaxLimit(rw, http.StatusUnauthorized)
		return
	}

	val, err := s.db.QueryValidatorByAny(address)
	if err != nil {
		zap.L().Error("failed to query validator info", zap.Error(err))
		errors.ErrNotExist(rw, http.StatusNotFound)
		return
	}

	// Note that saome validators existed in cosmoshub-1 or cosmoshub-2, but not in cosmoshub-3
	// They won't have any power event history, so return empty array for client to handle this
	validatorID, _ := s.db.QueryValidatorByID(val.Proposer)
	if validatorID == 0 {
		model.Respond(rw, []model.ResultPowerEventHistory{})
		return
	}

	if validatorID == -1 {
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	events, err := s.db.QueryValidatorVotingPowerEventHistory(validatorID, before, after, limit)
	if err != nil {
		zap.L().Error("failed to query power event history", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	result := make([]*model.ResultPowerEventHistory, 0)

	for _, e := range events {
		temp := &model.ResultPowerEventHistory{
			ID:             e.ID,
			Height:         e.Height,
			EventType:      e.EventType,
			VotingPower:    e.VotingPower * 1000000,
			NewVotingPower: e.NewVotingPowerAmount * 1000000,
			TxHash:         e.TxHash,
			Timestamp:      e.Timestamp,
		}
		result = append(result, temp)
	}

	model.Respond(rw, result)
	return
}

// GetValidatorEventsTotalCount receives validator address and total count of power event history.
func GetValidatorEventsTotalCount(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	val, err := s.db.QueryValidatorByAny(address)
	if err != nil {
		zap.L().Error("failed to query validator info", zap.Error(err))
		errors.ErrNotExist(rw, http.StatusNotFound)
		return
	}

	count, _ := s.db.CountValidatorPowerEvents(val.Proposer)

	result := &model.ResultVotingPowerHistoryCount{
		Moniker:         val.Moniker,
		OperatorAddress: val.OperatorAddress,
		Count:           count,
	}

	model.Respond(rw, result)
	return
}

// GetRedelegations returns all redelegations from a validator with the given query param
func GetRedelegations(rw http.ResponseWriter, r *http.Request) {
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

	// get all redelegations from a validator
	resp, err := s.client.HandleResponseHeight(endpoint)
	if err != nil {
		zap.L().Error("failed to get all redelegations from a validator", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, resp)
	return
}
