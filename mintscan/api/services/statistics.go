package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/db"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/utils"

	"github.com/tendermint/tendermint/rpc/client"
	resty "gopkg.in/resty.v1"
)

// GetMarketStats returns market statistics
func GetMarketStats(config *config.Config, db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	resp, _ := resty.R().Get(config.Market.CoinGecko.Endpoint)

	var coinGeckoMarket models.CoinGeckoMarket
	err := json.Unmarshal(resp.Body(), &coinGeckoMarket)
	if err != nil {
		fmt.Printf("failed to unmarshal coingecko market data: %t\n", err)
	}

	// Query every price of an hour for 24 hours
	limit := 24
	prices, _ := db.QueryOneDayPrices(limit)

	if len(prices) <= 0 {
		log.Fatal("failed to query prices")
	}

	priceStats := make([]*models.PriceStats, 0)

	for _, market := range prices {
		tempPriceStats := &models.PriceStats{
			Price: market.Price,
			Time:  market.Time,
		}
		priceStats = append(priceStats, tempPriceStats)
	}

	resultMarket := &models.ResultMarket{
		Price:             coinGeckoMarket.MarketData.CurrentPrice.Usd,
		Currency:          prices[0].Currency,
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

	utils.Respond(w, resultMarket)
	return nil
}

// GetNetworkStats returns network statistics
func GetNetworkStats(config *config.Config, db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	var limit int

	var statsNetwork models.StatsNetwork1H
	cntStats, _ := db.Model(&statsNetwork).Count()

	switch {
	case cntStats == 1:
		return json.NewEncoder(w).Encode(&models.StatsNetwork1H{})
	case cntStats <= 24:
		limit = cntStats
	default:
		limit = 24
	}

	// Query 24 network stats
	network1HStats, err := db.QueryNetworkStats(limit)
	if err != nil {
		utils.Respond(w, models.StatsNetwork1H{})
		return nil
	}

	if len(network1HStats) <= 0 {
		log.Fatal("failed to query network stats")
	}

	// Query bonded tokens percentage change for 24 hours
	network24HStats, err := db.QueryBondedRateIn24H()
	if err != nil {
		utils.Respond(w, models.StatsNetwork1H{})
		return nil
	}

	if len(network24HStats) <= 0 {
		log.Fatal("failed to query bonded tokens stats")
	}

	// Calculate change rate of bonded tokens in 24hours
	// (LatestBondedTokens - SecondLatestBondedTokens) / SecondLatestBondedTokens
	diff := network24HStats[0].BondedTokens - network24HStats[1].BondedTokens
	changeRateIn24H := diff / network24HStats[1].BondedTokens

	bondedTokensStats := make([]*models.BondedTokensStats, 0)

	for _, network1HStat := range network1HStats {
		tempBondedTokensStats := &models.BondedTokensStats{
			BondedTokens: network1HStat.BondedTokens,
			BondedRatio:  network1HStat.BondedRatio,
			LastUpdated:  network1HStat.Time,
		}
		bondedTokensStats = append(bondedTokensStats, tempBondedTokensStats)
	}

	resultNetworkInfo := &models.NetworkInfo{
		BondendTokensPercentChange24H: changeRateIn24H,
		BondedTokensStats:             bondedTokensStats,
	}

	utils.Respond(w, resultNetworkInfo)
	return nil
}
