package handler

import (
	"context"
	"encoding/json"
	"net/http"

	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	clienttypes "github.com/cosmostation/cosmostation-cosmos/mintscan/client/types"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

// GetTotalRewardsFromDelegator returns delegations rewards between a delegator and a validator.
func GetTotalRewardsFromDelegator(rw http.ResponseWriter, r *http.Request) {
	// 사용 안하는 함수
	vars := mux.Vars(r)
	delAddr := vars["delAddr"]

	err := model.VerifyBech32AccAddr(delAddr)
	if err != nil {
		zap.L().Debug("failed to validate delegator address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "delegator address is invalid")
		return
	}

	resp, err := s.client.RequestWithRestServer(clienttypes.PrefixDistribution + "/delegators/" + delAddr + "/rewards")
	if err != nil {
		zap.L().Error("failed to get delegator rewards", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	var dwar distributiontypes.QueryDelegatorTotalRewardsResponse
	// var drr distributiontypes.QueryDelegationRewardsResponse
	if err = json.Unmarshal(resp, &dwar); err != nil {
		// if err = s.client.GetCliContext().JSONMarshaler.UnmarshalJSON(resp, &dwar); err != nil {
		zap.L().Error("failed to get unmarshal given response", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}
	model.Respond(rw, dwar)
	return
}

// GetRewardsBetweenDelegatorAndValidator returns delegations rewards between a delegator and a validator.
func GetRewardsBetweenDelegatorAndValidator(rw http.ResponseWriter, r *http.Request) {
	// 사용 안함
	vars := mux.Vars(r)
	delAddr := vars["delAddr"]
	valAddr := vars["valAddr"]

	err := model.VerifyBech32AccAddr(delAddr)
	if err != nil {
		zap.L().Debug("failed to validate delegator address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "delegator address is invalid")
		return
	}

	err = model.VerifyBech32ValAddr(valAddr)
	if err != nil {
		zap.L().Debug("failed to validate validator address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "validator address is invalid")
		return
	}

	resp, err := s.client.RequestWithRestServer(clienttypes.PrefixDistribution + "/delegators/" + delAddr + "/rewards/" + valAddr)
	if err != nil {
		zap.L().Error("failed to get delegator rewards", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	// var dwar distributiontypes.QueryDelegatorTotalRewardsResponse
	var drr distributiontypes.QueryDelegationRewardsResponse
	if err = s.client.GetCliContext().JSONMarshaler.UnmarshalJSON(resp, &drr); err != nil {
		zap.L().Error("failed to get unmarshal given response", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}
	model.Respond(rw, drr)
	return
}

// GetDelegatorWithdrawalAddress returns delegator's rewards withdrawal address.
func GetDelegatorWithdrawalAddress(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	delAddr := vars["delAddr"]

	err := model.VerifyBech32AccAddr(delAddr)
	if err != nil {
		zap.L().Debug("failed to validate delegator address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "delegator address is invalid")
		return
	}

	queryClient := distributiontypes.NewQueryClient(s.client.GetCliContext())
	request := distributiontypes.QueryDelegatorWithdrawAddressRequest{DelegatorAddress: delAddr}
	res, err := queryClient.DelegatorWithdrawAddress(context.Background(), &request)
	if err != nil {
		zap.L().Error("failed to get delegator withdraw address", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, res)
	return
}

// GetCommunityPool returns current community pool.
func GetCommunityPool(rw http.ResponseWriter, r *http.Request) {
	queryClient := distributiontypes.NewQueryClient(s.client.GetCliContext())
	res, err := queryClient.CommunityPool(context.Background(), &distributiontypes.QueryCommunityPoolRequest{})
	if err != nil {
		zap.L().Error("failed to get community pool", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, res)
	return
}
