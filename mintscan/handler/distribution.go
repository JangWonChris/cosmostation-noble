package handler

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

// GetRewardsBetweenDelegatorAndValidator returns delegations rewards between a delegator and a validator.
func GetRewardsBetweenDelegatorAndValidator(rw http.ResponseWriter, r *http.Request) {
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

	resp, err := s.client.HandleResponseHeight("/distribution/delegators/" + delAddr + "/rewards/" + valAddr)
	if err != nil {
		zap.L().Error("failed to get community pool", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, resp)
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

	resp, err := s.client.HandleResponseHeight("/distribution/delegators/" + delAddr + "/withdraw_address")
	if err != nil {
		zap.L().Error("failed to get delegator withdraw address", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, resp)
	return
}

// GetCommunityPool returns current community pool.
func GetCommunityPool(rw http.ResponseWriter, r *http.Request) {
	resp, err := s.client.HandleResponseHeight("/distribution/community_pool")
	if err != nil {
		zap.L().Error("failed to get community pool", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, resp)
	return
}
