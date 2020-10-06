package handler

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"

	"go.uber.org/zap"
)

const (
	// requiredLimit is the limit number of items that are required for clients to handle market chart graph.
	requiredLimit = 25
)

// GetMarketStats returns market statistics
// TODO: find better and cleaner way to handle this API.
func GetMarketStats(rw http.ResponseWriter, r *http.Request) {
	currentPrice, err := s.db.QueryPriceFromMarketStat5M()
	if err != nil {
		zap.S().Errorf("failed to query current price from stat market 5m: %s", err)
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	result := &model.ResultMarket{
		Price:             currentPrice.Price,
		Currency:          currentPrice.Currency,
		MarketCapRank:     currentPrice.MarketCapRank,
		PercentChange1H:   currentPrice.PercentChange1H,
		PercentChange24H:  currentPrice.PercentChange24H,
		PercentChange7D:   currentPrice.PercentChange7D,
		PercentChange30D:  currentPrice.PercentChange30D,
		TotalVolume:       currentPrice.TotalVolume,
		CirculatingSupply: currentPrice.CirculatingSupply,
		LastUpdated:       currentPrice.LastUpdated,
	}

	model.Respond(rw, result)
	return
}

// GetNetworkStats returns network statistics
func GetNetworkStats(rw http.ResponseWriter, r *http.Request) {
	// Count network statistics to see if enough data is available to query.
	networkStatsNum, err := s.db.CountNetworkStats1H()
	if err != nil {
		zap.S().Errorf("failed to count network stats: %s", err)
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	if networkStatsNum < requiredLimit {
		zap.S().Debug("network stats num is less than required limit")
		errors.ErrNoDataAvailable(rw, http.StatusInternalServerError)
		return
	}

	network1HStats, err := s.db.QueryNetworkStats1H(requiredLimit)
	if err != nil {
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	bondedTokensStats := make([]*model.BondedTokensStats, 0)

	for _, ns := range network1HStats {
		stats := &model.BondedTokensStats{
			BondedTokens: ns.BondedTokens,
			BondedRatio:  ns.BondedRatio,
			LastUpdated:  ns.Timestamp,
		}

		bondedTokensStats = append(bondedTokensStats, stats)
	}

	// Query two latest network stats from 1D table to calculate bonded change rate within 24 hours.
	network1Dstats, err := s.db.QueryNetworkStats1D(2)
	if err != nil {
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	if len(network1Dstats) < 2 {
		zap.S().Debug("network stats data from 1 day table needs at least two")
		errors.ErrNoDataAvailable(rw, http.StatusInternalServerError)
		return
	}

	// Calculate change rate of bonded tokens in 24hours
	// (LatestBondedTokens - SecondLatestBondedTokens) / SecondLatestBondedTokens
	changeRateIn24H := (network1Dstats[0].BondedTokens - network1Dstats[1].BondedTokens) / network1Dstats[1].BondedTokens

	result := &model.NetworkInfo{
		BondendTokensPercentChange24H: changeRateIn24H,
		BondedTokensStats:             bondedTokensStats,
	}

	model.Respond(rw, result)
	return
}
