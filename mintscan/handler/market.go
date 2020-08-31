package handler

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"
	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

// GetCoinPrice returns coin price.
func GetCoinPrice(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	coinID := vars["id"]

	switch coinID {
	case model.Cosmos:
		coinID = model.Cosmos
	default:
		coinID = model.Cosmos
	}

	result, err := s.client.GetCoinGeckoCoinPrice(coinID)
	if err != nil {
		zap.L().Error("failed to query validator by proposer", zap.Error(err))
		return
	}

	model.Respond(rw, result)
	return
}
