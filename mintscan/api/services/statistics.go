package services

import (
	"encoding/json"
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/db"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/utils"

	"github.com/rs/zerolog/log"
	"github.com/tendermint/tendermint/rpc/client"
	resty "gopkg.in/resty.v1"
)

// GetMarketStats returns marketInfo
func GetMarketStats(config *config.Config, db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	limit := 24

	// query current price
	resp, _ := resty.R().Get(config.Market.CoinGecko.URL)

	var coinGeckoMarket models.CoinGeckoMarket
	err := json.Unmarshal(resp.Body(), &coinGeckoMarket)
	if err != nil {
		log.Info().Str(models.Service, models.LogStatistics).Str(models.Method, "GetMarketStats").Err(err).Msg("unmarshal coinGeckoMarket error")
	}

	// query price chart
	var statsCoingeckoMarket1H []models.StatsCoingeckoMarket1H
	_ = db.Model(&statsCoingeckoMarket1H).
		Order("id DESC").
		Limit(limit).
		Select()

	priceStats := make([]*models.PriceStats, 0)

	for _, market := range statsCoingeckoMarket1H {
		tempPriceStats := &models.PriceStats{
			Price: market.Price,
			Time:  market.Time,
		}
		priceStats = append(priceStats, tempPriceStats)
	}

	resultMarket := &models.ResultMarket{
		Price:             coinGeckoMarket.MarketData.CurrentPrice.Usd,
		Currency:          statsCoingeckoMarket1H[0].Currency,
		MarketCapRank:     statsCoingeckoMarket1H[0].MarketCapRank,
		PercentChange1H:   statsCoingeckoMarket1H[0].PercentChange1H,
		PercentChange24H:  statsCoingeckoMarket1H[0].PercentChange24H,
		PercentChange7D:   statsCoingeckoMarket1H[0].PercentChange7D,
		PercentChange30D:  statsCoingeckoMarket1H[0].PercentChange30D,
		TotalVolume:       statsCoingeckoMarket1H[0].TotalVolume,
		CirculatingSupply: statsCoingeckoMarket1H[0].CirculatingSupply,
		LastUpdated:       statsCoingeckoMarket1H[0].LastUpdated,
		PriceStats:        priceStats,
	}

	utils.Respond(w, resultMarket)
	return nil
}

// GetNetworkStats returns network stats
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

	// query bonded tokens chart
	var statsNetwork1H []models.StatsNetwork1H
	err := db.Model(&statsNetwork1H).
		Order("id DESC").
		Limit(limit).
		Select()
	if err != nil {
		return json.NewEncoder(w).Encode(&models.StatsNetwork1H{})
	}

	bondedTokensStats := make([]*models.BondedTokensStats, 0)
	if len(statsNetwork1H) > 0 {
		for _, networkStat := range statsNetwork1H {
			tempBondedTokensStats := &models.BondedTokensStats{
				BondedTokens: networkStat.BondedTokens,
				BondedRatio:  networkStat.BondedRatio,
				LastUpdated:  networkStat.Time,
			}
			bondedTokensStats = append(bondedTokensStats, tempBondedTokensStats)
		}
	}

	// bonded tokens percentage change in 24 hours
	var statsNetwork24H []models.StatsNetwork24H
	_ = db.Model(&statsNetwork24H).
		Order("id DESC").
		Limit(2).
		Select()

	// bonded tokens rate change in last 24 hours
	percentChange24H := float64(0)

	// TODO: cosmoshub-3 업그레이드 후 네트워크 데이터가 아예 없을경우, 1개일 경우 예외처리 하기
	if len(statsNetwork24H) > 0 {
		latestBondedTokens := statsNetwork24H[0].BondedTokens
		before24HBondedTokens := statsNetwork24H[1].BondedTokens
		diff := latestBondedTokens - before24HBondedTokens
		percentChange24H = diff / before24HBondedTokens
	}

	resultNetworkInfo := &models.NetworkInfo{
		BondendTokensPercentChange24H: percentChange24H,
		BondedTokensStats:             bondedTokensStats,
	}

	utils.Respond(w, resultNetworkInfo)
	return nil
}
