package handler

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/schema"

	"go.uber.org/zap"
)

const (
	// RequiredLimit is the limit number of items that are required for clients to handle market chart graph.
	RequiredLimit = 25
)

// GetMarketStats returns market statistics
// TODO: find better and cleaner way to handle this API.
func GetMarketStats(rw http.ResponseWriter, r *http.Request) {
	resp, err := s.client.GetCoinGeckoMarketData(model.Kava)
	if err != nil {
		zap.L().Error("failed to get market data from CoinGecko", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	// query every price of an hour for 24 hours
	prices, err := s.db.QueryPrices1D(RequiredLimit)
	if err != nil {
		zap.L().Error("failed to query prices for 1 day", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	if len(prices) <= 0 {
		zap.L().Debug("failed to query prices", zap.Int("len(prices)", len(prices)))
		errors.ErrNotExist(rw, http.StatusNotFound)
		return
	}

	priceStats := make([]*model.PriceStats, 0)

	for _, market := range prices {
		tempPriceStats := &model.PriceStats{
			Price: market.Price,
			Time:  market.Timestamp,
		}
		priceStats = append(priceStats, tempPriceStats)
	}

	result := &model.ResultMarket{
		Price:             resp.MarketData.CurrentPrice.Usd,
		Currency:          model.Currency,
		MarketCapRank:     prices[0].MarketCapRank,
		PercentChange1H:   prices[0].PercentChange1H,
		PercentChange24H:  prices[0].PercentChange24H,
		PercentChange7D:   prices[0].PercentChange7D,
		PercentChange30D:  prices[0].PercentChange30D,
		TotalVolume:       prices[0].TotalVolume,
		CirculatingSupply: prices[0].CirculatingSupply,
		LastUpdated:       prices[0].LastUpdated,
		PriceStats:        priceStats,
	}

	model.Respond(rw, result)
	return
}

// GetNetworkStats returns network statistics
func GetNetworkStats(rw http.ResponseWriter, r *http.Request) {
	var limit int

	var statsNetwork schema.StatsNetwork1H
	cntStats, _ := s.db.Model(&statsNetwork).Count()

	switch {
	case cntStats == 1:
		model.Respond(rw, schema.StatsNetwork1H{})
		return
	case cntStats <= 24:
		limit = cntStats
	default:
		limit = RequiredLimit
	}

	// query 24 network stats
	network1HStats, err := s.db.QueryNetworkStats(limit)
	if err != nil {
		model.Respond(rw, schema.StatsNetwork1H{})
		return
	}

	if len(network1HStats) <= 0 {
		zap.L().Debug("failed to query network stats", zap.Int("len(network1HStats)", len(network1HStats)))
	}

	// Query bonded tokens percentage change for 24 hours
	network24HStats, err := s.db.QueryBondedRateIn1D()
	if err != nil {
		model.Respond(rw, schema.StatsNetwork1H{})
		return
	}

	if len(network24HStats) <= 0 {
		zap.L().Debug("failed to query bonded tokens stats", zap.Int("len(network1HStats)", len(network1HStats)))
	}

	// Calculate change rate of bonded tokens in 24hours
	// (LatestBondedTokens - SecondLatestBondedTokens) / SecondLatestBondedTokens
	diff := network24HStats[0].BondedTokens - network24HStats[1].BondedTokens
	changeRateIn24H := diff / network24HStats[1].BondedTokens

	bondedTokensStats := make([]*model.BondedTokensStats, 0)

	for _, network1HStat := range network1HStats {
		temp := &model.BondedTokensStats{
			BondedTokens: network1HStat.BondedTokens,
			BondedRatio:  network1HStat.BondedRatio,
			LastUpdated:  network1HStat.Timestamp,
		}
		bondedTokensStats = append(bondedTokensStats, temp)
	}

	result := &model.NetworkInfo{
		BondendTokensPercentChange24H: changeRateIn24H,
		BondedTokensStats:             bondedTokensStats,
	}

	model.Respond(rw, result)
	return
}
