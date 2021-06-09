package custom

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/errors"
	"github.com/cosmostation/cosmostation-cosmos/model"
	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

// GetSimpleCoinPrice returns coin price.
func GetSimpleCoinPrice(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	coinID := vars["id"]

	switch coinID {
	case model.Cosmos:
		coinID = model.Cosmos
	default:
		coinID = model.Cosmos
	}

	result, err := s.Client.GetCoinGeckoSimpleCoinPrice(coinID)
	if err != nil {
		zap.L().Error("failed to query validator by proposer", zap.Error(err))
		return
	}

	model.Respond(rw, result)
	return
}

// GetCoinMarketChartData returns market chart data using CoinGecko market chart API.
func GetCoinMarketChartData(rw http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()["id"]) <= 0 {
		errors.ErrRequiredParam(rw, http.StatusBadRequest, "id param is required")
		return
	}

	id := r.URL.Query()["id"][0]

	// Current time and its minus 24 hours
	to := time.Now().UTC()
	from := to.AddDate(0, 0, -1)

	marketChartData, err := s.Client.GetCoinMarketChartData(id, fmt.Sprintf("%d", from.Unix()), fmt.Sprintf("%d", to.Unix()))
	if err != nil {
		zap.S().Errorf("failed to request coin market chart data: %s", err)
		return
	}

	model.Respond(rw, marketChartData)
	return
}
