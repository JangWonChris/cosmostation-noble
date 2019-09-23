package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/stats"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/types"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/utils"

	"github.com/go-pg/pg"
	"github.com/tendermint/tendermint/rpc/client"
	resty "gopkg.in/resty.v1"
)

// GetMarketInfo returns marketInfo
func GetMarketInfo(config *config.Config, db *pg.DB, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	// How many data to show in a chart
	limit := 24

	// query all market stats
	var marketInfo stats.CoingeckoMarketStats
	err := db.Model(&marketInfo).
		Order("id DESC").
		Limit(1).
		Select()
	if err != nil {
		return json.NewEncoder(w).Encode(&models.MarketInfo{})
	}

	// query current price
	resp, _ := resty.R().Get(config.Market.CoinGecko.URL)

	var coinGeckoMarket types.CoinGeckoMarket
	err = json.Unmarshal(resp.Body(), &coinGeckoMarket)
	if err != nil {
		fmt.Printf("MarketInfo unmarshal error - %v\n", err)
	}

	// Query price chart
	var marketStats []stats.CoingeckoMarketStats
	_ = db.Model(&marketStats).
		Order("id DESC").
		Limit(limit).
		Select()

	priceStats := make([]*models.PriceStats, 0)
	for _, market := range marketStats {
		tempPriceStats := &models.PriceStats{
			Price: market.Price,
			Time:  market.Time,
		}
		priceStats = append(priceStats, tempPriceStats)
	}

	resultMarketInfo := &models.MarketInfo{
		Price:            coinGeckoMarket.MarketData.CurrentPrice.Usd,
		Currency:         marketInfo.Currency,
		PercentChange1H:  marketInfo.PercentChange1H,
		PercentChange24H: marketInfo.PercentChange24H,
		LastUpdated:      marketInfo.LastUpdated,
		PriceStats:       priceStats,
	}

	utils.Respond(w, resultMarketInfo)
	return nil
}

// GetNetworkStats returns network stats
func GetNetworkStats(config *config.Config, db *pg.DB, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	// How many data to show in a chart
	limit := 24

	// Query bonded tokens chart
	var networkStats []stats.NetworkStats
	err := db.Model(&networkStats).
		Order("id DESC").
		Limit(limit).
		Select()
	if err != nil {
		return json.NewEncoder(w).Encode(&models.NetworkInfo{})
	}

	bondedTokensStats := make([]*models.BondedTokensStats, 0)
	for _, networkStat := range networkStats {
		tempBondedTokensStats := &models.BondedTokensStats{
			BondedTokens: networkStat.BondedTokens1H,
			BondedRatio:  networkStat.BondedRatio1H,
			LastUpdated:  networkStat.LastUpdated,
		}
		bondedTokensStats = append(bondedTokensStats, tempBondedTokensStats)
	}

	// Latest bonded tokens that is saved in DB
	var bondedTokensLatest stats.NetworkStats
	_ = db.Model(&bondedTokensLatest).
		Order("id DESC").
		Limit(1).
		Select()

	// Bonded Tokens 24H before
	var bondedTokensBefore24H stats.NetworkStats
	_ = db.Model(&bondedTokensBefore24H).
		Where("id = ?", bondedTokensLatest.ID-23).
		Order("id DESC").
		Limit(1).
		Select()

	// Bonded tokens rate change in last 24 hours
	latestBondedTokens := bondedTokensLatest.BondedTokens1H
	before24HBondedTokens := bondedTokensBefore24H.BondedTokens1H
	diff := latestBondedTokens - before24HBondedTokens
	percentChange24H := float64(1)
	if diff > 0 {
		percentChange24H = (float64(diff) / float64(latestBondedTokens)) * 100
	} else {
		percentChange24H = (float64(diff) / float64(before24HBondedTokens)) * 100
	}

	resultNetworkInfo := &models.NetworkInfo{
		BondendTokensPercentChange24H: percentChange24H,
		BondedTokensStats:             bondedTokensStats,
	}

	utils.Respond(w, resultNetworkInfo)
	return nil
}
