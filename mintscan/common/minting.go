package common

import (
	"context"
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/app"
	"github.com/cosmostation/cosmostation-cosmos/errors"
	"github.com/cosmostation/cosmostation-cosmos/model"
	"go.uber.org/zap"
)

// GetMintingInflation returns inflation rate of the network.
func GetMintingInflation(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		res, err := a.Client.GRPC.GetInflation(context.Background())
		if err != nil {
			zap.L().Error("failed to get inflation information", zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			return
		}

		model.Respond(rw, res)
		return
	}
}
