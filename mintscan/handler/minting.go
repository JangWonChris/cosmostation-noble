package handler

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"

	"go.uber.org/zap"
)

// GetMintingInflation returns inflation rate of the network.
func GetMintingInflation(rw http.ResponseWriter, r *http.Request) {
	resp, err := s.client.HandleResponseHeight("/minting/inflation")
	if err != nil {
		zap.L().Error("failed to get inflation information", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, resp)
	return
}
