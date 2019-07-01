package services

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/stats"
	u "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/utils"

	"github.com/go-pg/pg"
	"github.com/tendermint/tendermint/rpc/client"
	resty "gopkg.in/resty.v1"
)

var (
	CoinGeckoAPIURL = "https://api.coingecko.com/api/v3/coins/cosmos"
)

// GetMarketInfo returns marketInfo
func GetMarketInfo(RPCClient *client.HTTP, DB *pg.DB, Config *config.Config, w http.ResponseWriter, r *http.Request) error {
	// How many data to show in a chart
	limit := 24

	// Query all market stats
	var marketInfo stats.CoingeckoMarketStats
	err := DB.Model(&marketInfo).
		Order("id DESC").
		Limit(1).
		Select()
	if err != nil {
		return json.NewEncoder(w).Encode(&models.MarketInfo{})
	}

	// Query LCD - current price
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, _ := resty.R().Get(CoinGeckoAPIURL)

	var coinGeckoMarketInfo models.CoinGeckoMarketInfo
	err = json.Unmarshal(resp.Body(), &coinGeckoMarketInfo)
	if err != nil {
		fmt.Printf("MarketInfo unmarshal error - %v\n", err)
	}

	// Query price chart
	var marketStats []stats.CoingeckoMarketStats
	_ = DB.Model(&marketStats).
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
		Price:            coinGeckoMarketInfo.MarketData.CurrentPrice.Usd,
		Currency:         marketInfo.Currency,
		PercentChange1H:  marketInfo.PercentChange1H,
		PercentChange24H: marketInfo.PercentChange24H,
		LastUpdated:      marketInfo.LastUpdated,
		PriceStats:       priceStats,
	}

	u.Respond(w, resultMarketInfo)
	return nil
}

// GetNetworkStats returns network stats
func GetNetworkStats(RPCClient *client.HTTP, DB *pg.DB, Config *config.Config, w http.ResponseWriter, r *http.Request) error {
	// How many data to show in a chart
	limit := 24

	// Query bonded tokens chart
	var networkStats []stats.NetworkStats
	err := DB.Model(&networkStats).
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
	_ = DB.Model(&bondedTokensLatest).
		Order("id DESC").
		Limit(1).
		Select()

	// Bonded Tokens 24H before
	var bondedTokensBefore24H stats.NetworkStats
	_ = DB.Model(&bondedTokensBefore24H).
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

	u.Respond(w, resultNetworkInfo)
	return nil
}
