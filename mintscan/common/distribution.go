package common

import (
	"context"
	"net/http"

	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmostation/cosmostation-cosmos/app"
	"github.com/cosmostation/cosmostation-cosmos/errors"
	"github.com/cosmostation/cosmostation-cosmos/model"
	ltypes "github.com/cosmostation/mintscan-backend-library/types"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

// GetDelegatorWithdrawalAddress returns delegator's rewards withdrawal address.
func GetDelegatorWithdrawalAddress(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		delAddr := vars["delAddr"]

		err := ltypes.VerifyBech32AccAddr(delAddr)
		if err != nil {
			zap.L().Debug("failed to validate delegator address", zap.Error(err))
			errors.ErrInvalidParam(rw, http.StatusBadRequest, "delegator address is invalid")
			return
		}

		queryClient := distributiontypes.NewQueryClient(a.Client.GRPC)
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
}

// GetCommunityPool returns current community pool.
func GetCommunityPool(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		queryClient := distributiontypes.NewQueryClient(a.Client.GRPC)
		res, err := queryClient.CommunityPool(context.Background(), &distributiontypes.QueryCommunityPoolRequest{})
		if err != nil {
			zap.L().Error("failed to get community pool", zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			return
		}

		model.Respond(rw, res)
		return
	}
}
