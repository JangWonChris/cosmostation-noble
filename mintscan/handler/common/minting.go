package common

import (
	"context"
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"
	"go.uber.org/zap"
)

// GetMintingInflation returns inflation rate of the network.
func GetMintingInflation(rw http.ResponseWriter, r *http.Request) {
	// 접속하는 유저마다 요청하는 inflation은 노드를 멈추게 만든다.
	model.Respond(rw, struct{}{})
	return

	res, err := s.Client.GRPC.GetInflation(context.Background())
	if err != nil {
		zap.L().Error("failed to get inflation information", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, res)
	return
}
