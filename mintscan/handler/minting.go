package handler

import (
	"net/http"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	clienttypes "github.com/cosmostation/cosmostation-cosmos/mintscan/client/types"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"

	"go.uber.org/zap"
)

// GetMintingInflation returns inflation rate of the network.
func GetMintingInflation(rw http.ResponseWriter, r *http.Request) {
	resp, err := s.client.RequestWithRestServer(clienttypes.PrefixMint + "/inflation")
	if err != nil {
		zap.L().Error("failed to get inflation information", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}
	var inflationResponse minttypes.QueryInflationResponse
	if err = s.client.GetCliContext().JSONMarshaler.UnmarshalJSON(resp, &inflationResponse); err != nil {
		zap.L().Error("failed to unmarshal data", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, inflationResponse)
	return
}
