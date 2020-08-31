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
	case model.Kava:
		coinID = model.Kava
	case model.BNB:
		coinID = model.Binance
	case model.Binance:
		coinID = model.Binance
	case model.USDX:
		coinID = model.USDXStableCoin
	case model.USDXStableCoin:
		coinID = model.USDXStableCoin

	default:
		coinID = model.Kava
	}

	result, err := s.client.GetCoinGeckoCoinPrice(coinID)
	if err != nil {
		zap.L().Error("failed to query validator by proposer", zap.Error(err))
		return
	}

	model.Respond(rw, result)
	return
}
