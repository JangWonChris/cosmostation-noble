package handler

import (
	"context"
	"net/http"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"

	"go.uber.org/zap"
)

// GetMintingInflation returns inflation rate of the network.
func GetMintingInflation(rw http.ResponseWriter, r *http.Request) {

	queryClient := minttypes.NewQueryClient(s.client.GetCliContext())
	res, err := queryClient.Inflation(context.Background(), &minttypes.QueryInflationRequest{})
	if err != nil {
		zap.L().Error("failed to get inflation information", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, res)
	return
}
