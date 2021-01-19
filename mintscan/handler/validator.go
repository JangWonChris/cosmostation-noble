package handler

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"
	"github.com/cosmostation/mintscan-backend-library/db/schema"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

// GetValidators returns all validators.
func GetValidators(rw http.ResponseWriter, r *http.Request) {
	var status string

	if len(r.URL.Query()["status"]) > 0 {
		status = r.URL.Query()["status"][0]
	}

	vals := make([]schema.Validator, 0)

	switch status {
	case model.ActiveValidator:
		vals, _ = s.db.QueryValidatorsByStatus(int(stakingtypes.Bonded))
	case model.InactiveValidator:
		unbondingVals, _ := s.db.QueryValidatorsByStatus(int(stakingtypes.Unbonding))
		unbondedVals, _ := s.db.QueryValidatorsByStatus(int(stakingtypes.Unbonded))
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

	latestDBHeight, err := s.db.QueryLatestBlockHeight()
	if err != nil {
		zap.S().Errorf("failed to query latest block height: %s", err)
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	result := make([]*model.ResultValidator, 0)

	for _, val := range vals {
		// Default is missing the last 100 blocks
		missBlockCount := model.MissingAllBlocks

		if val.Status == int(stakingtypes.Bonded) {
			blocks, err := s.db.QueryValidatorUptime(val.Proposer, latestDBHeight-1)
			if err != nil {
				zap.S().Errorf("failed to query validator's missing blocks: %s", err)
				errors.ErrInternalServer(rw, http.StatusInternalServerError)
				return
			}
			missBlockCount = len(blocks)
		}

		uptime := &model.Uptime{
			Address:      val.Proposer,
			MissedBlocks: missBlockCount,
			OverBlocks:   model.MissingAllBlocks,
		}

		val := &model.ResultValidator{
			Rank:            val.Rank,
			AccountAddress:  val.Address,
			OperatorAddress: val.OperatorAddress,
			ConsensusPubkey: val.ConsensusPubkey,
			Jailed:          val.Jailed,
			Status:          val.Status,
			Tokens:          val.Tokens,
			DelegatorShares: val.DelegatorShares,
			Moniker:         val.Moniker,
			Identity:        val.Identity,
			Website:         val.Website,
			Details:         val.Details,
			// UnbondingHeight:      val.UnbondingHeight,
			//jeonghwan 문자열로 되있는데, 숫자로 바꿔도 무관?
			UnbondingHeight:      fmt.Sprintf("%d", val.UnbondingHeight),
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

	val, err := s.db.QueryValidatorByAnyAddr(address)
	if err != nil {
		zap.S().Errorf("failed to query validator information: %s", err)
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	if val.Address == "" {
		errors.ErrNotExist(rw, http.StatusNotFound)
		return
	}

	latestDBHeight, err := s.db.QueryLatestBlockHeight()
	if err != nil {
		zap.S().Errorf("failed to query latest block height: %s", err)
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	// Default is missing the last 100 blocks
	missBlockCount := model.MissingAllBlocks

	if val.Status == int(stakingtypes.Bonded) {
		blocks, err := s.db.QueryValidatorUptime(val.Proposer, latestDBHeight-1)
		if err != nil {
			zap.S().Errorf("failed to query validator's missing blocks: %s", err)
			errors.ErrInternalServer(rw, http.StatusInternalServerError)
			return
		}
		missBlockCount = len(blocks)
	}

	uptime := &model.Uptime{
		Address:      val.Proposer,
		MissedBlocks: missBlockCount,
		OverBlocks:   model.MissingAllBlocks,
	}

	// Query a validator's bonded information
	powerEventHistory, _ := s.db.QueryValidatorBondedInfo(val.Proposer)

	result := &model.ResultValidatorDetail{
		Rank:            val.Rank,
		AccountAddress:  val.Address,
		OperatorAddress: val.OperatorAddress,
		ConsensusPubkey: val.ConsensusPubkey,
		BondedHeight:    powerEventHistory.Height,
		BondedTime:      powerEventHistory.Timestamp,
		Jailed:          val.Jailed,
		Status:          val.Status,
		Tokens:          val.Tokens,
		DelegatorShares: val.DelegatorShares,
		Moniker:         val.Moniker,
		Identity:        val.Identity,
		Website:         val.Website,
		Details:         val.Details,
		//jeonghwan 문자열로 되있는데, 숫자로 바꿔도 무관?
		UnbondingHeight:      fmt.Sprintf("%d", val.UnbondingHeight),
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

// GetValidatorUptime returns a validator's uptime, which counts a number of missing blocks for the last 100 blocks.
// When uptime is 100%: there is not a single missing block saved in database. Therfore, it returns an empty array.
// When uptime is from 1% ~ 99%: simply return a number of missing blocks.
// When uptime is 0%: Case 1. return 100 missing blocks.
// When uptime is 0%: Case 2. a validator is unbonding or unbonded state.
func GetValidatorUptime(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	val, err := s.db.QueryValidatorByAnyAddr(address)
	if err != nil {
		zap.S().Errorf("failed to query validator information: %s", err)
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	if val.Address == "" {
		errors.ErrNotExist(rw, http.StatusNotFound)
		return
	}

	latestDBHeight, err := s.db.QueryLatestBlockHeight()
	if err != nil {
		zap.S().Errorf("failed to query latest block height: %s", err)
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	var result model.ResultMissesDetail
	result.LatestHeight = latestDBHeight - 1

	// Query missing blocks for the last 100 blocks
	blocks, err := s.db.QueryValidatorUptime(val.Proposer, result.LatestHeight)
	if err != nil {
		zap.S().Errorf("failed to query validator's missing blocks: %s", err)
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	if len(blocks) <= 0 {
		result.ResultUptime = []model.ResultUptime{} // empty array
		model.Respond(rw, result)
		return
	}

	for _, block := range blocks {
		uptime := &model.ResultUptime{
			Height:    block.Height,
			Timestamp: block.Timestamp,
		}

		result.ResultUptime = append(result.ResultUptime, *uptime)
	}

	model.Respond(rw, result)
	return
}

// GetValidatorUptimeRange returns the validator's block consencutive misses.
func GetValidatorUptimeRange(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	val, err := s.db.QueryValidatorByAnyAddr(address)
	if err != nil {
		zap.L().Debug("failed to query validator info", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	if val.Address == "" {
		errors.ErrNotExist(rw, http.StatusNotFound)
		return
	}

	blocks, err := s.db.QueryValidatorUptimeRange(val.Proposer)
	if len(blocks) <= 0 {
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

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

	val, err := s.db.QueryValidatorByAnyAddr(address)
	if err != nil {
		zap.L().Error("failed to query validator info", zap.Error(err))
		errors.ErrNotExist(rw, http.StatusNotFound)
		return
	}

	queryClient := stakingtypes.NewQueryClient(s.client.GetCLIContext())
	request := stakingtypes.QueryValidatorDelegationsRequest{ValidatorAddr: val.OperatorAddress}
	res, err := queryClient.ValidatorDelegations(context.Background(), &request)
	if err != nil {
		zap.L().Error("failed to get all delegations from a validator", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	// calculate the amount of uatom, which should divide validator's token by delegator_shares
	tokens, _ := strconv.ParseFloat(val.Tokens, 64)
	delegatorShares, _ := strconv.ParseFloat(val.DelegatorShares, 64)
	uatom := tokens / delegatorShares

	var validatorDelegations []*model.ValidatorDelegations
	for _, dr := range res.DelegationResponses {
		shares, _ := strconv.ParseFloat(dr.Delegation.Shares.String(), 64)
		amount := fmt.Sprintf("%f", shares*uatom)

		temp := &model.ValidatorDelegations{
			DelegatorAddress: dr.Delegation.DelegatorAddress,
			ValidatorAddress: dr.Delegation.ValidatorAddress,
			Shares:           dr.Delegation.Shares,
			Amount:           amount,
		}
		validatorDelegations = append(validatorDelegations, temp)
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
		TotalDelegatorNum:     len(res.DelegationResponses),
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

	val, err := s.db.QueryValidatorByAnyAddr(address)
	if err != nil {
		zap.S().Errorf("failed to query validator information: %s", err)
		return
	}

	if val.Address == "" {
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
	fmt.Println(len(events))

	result := make([]*model.ResultPowerEventHistory, 0)

	for _, e := range events {
		temp := &model.ResultPowerEventHistory{
			ID:             e.ID,
			Height:         e.Height,
			MsgType:        e.MsgType,
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

	val, err := s.db.QueryValidatorByAnyAddr(address)
	if err != nil {
		zap.S().Errorf("failed to query validator information: %s", err)
		return
	}

	if val.Address == "" {
		errors.ErrNotExist(rw, http.StatusNotFound)
		return
	}

	count, _ := s.db.CountPowerEventHistoryTransactions(val.Proposer)

	result := &model.ResultVotingPowerHistoryCount{
		Moniker:         val.Moniker,
		OperatorAddress: val.OperatorAddress,
		Count:           count,
	}

	model.Respond(rw, result)
	return
}

// GetRedelegationsLegacy returns all redelegations from a validator with the given query param
func GetRedelegationsLegacy(rw http.ResponseWriter, r *http.Request) {
	var delAddr string
	if len(r.URL.Query()["delegator"]) > 0 {
		delAddr = r.URL.Query()["delegator"][0]
	}

	res, err := s.client.GRPC.GetRedelegations(r.Context(), delAddr, "", "")
	if err != nil {
		zap.L().Error("failed to get all redelegations from a validator", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, res)
	return
}

// GetRedelegations returns all redelegations from a validator with the given query param
func GetRedelegations(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	delAddr := vars["delAddr"]

	var srcValidatorAddress, dstValidatorAddress string

	if len(r.URL.Query()["src_validator_addr"]) > 0 {
		srcValidatorAddress = r.URL.Query()["src_validator_addr"][0]
	}

	if len(r.URL.Query()["dst_validator_addr"]) > 0 {
		dstValidatorAddress = r.URL.Query()["dst_validator_addr"][0]
	}
	res, err := s.client.GRPC.GetRedelegations(r.Context(), delAddr, srcValidatorAddress, dstValidatorAddress)
	if err != nil {
		zap.L().Error("failed to get all redelegations from a validator", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, res)
	return
}
