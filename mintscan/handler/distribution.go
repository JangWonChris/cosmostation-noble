package handler

import (
	"context"
	"net/http"

	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

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
