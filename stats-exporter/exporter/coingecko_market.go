package exporter

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/types"

	resty "gopkg.in/resty.v1"
)

// SaveCoinGeckoMarketStats1H saves coingecko market statistics every hour
func (ses *StatsExporterService) SaveCoinGeckoMarketStats1H() {
	log.Println("CoingGecko Market Stats 1H")

	// Request CoinGecko API
	var coinGeckoMarketInfo types.CoinGeckoMarketInfo
	resp, err := resty.R().SetHeader("Accepts", "application/json").Get(ses.config.Market.CoinGecko.URL)
	if err != nil {
		fmt.Print("Query CoinGecko API Request error - ", err)
	}

	err = json.Unmarshal(resp.Body(), &coinGeckoMarketInfo)
	if err != nil {
		fmt.Printf("CoinGeckoMarketInfo unmarshal error - %v\n", err)
	}

	// Insert into marketStats slice
	statsCoingeckoMarket := make([]*types.StatsCoingeckoMarket1H, 0)
	tempStatsCoingeckoMarket := &types.StatsCoingeckoMarket1H{
		Price:    coinGeckoMarketInfo.MarketData.CurrentPrice.Usd,
		Currency: types.Currency,
		Time:     time.Now(),
	}
	statsCoingeckoMarket = append(statsCoingeckoMarket, tempStatsCoingeckoMarket)

	// Save
	_, err = ses.db.Model(&statsCoingeckoMarket).Insert()
	if err != nil {
		fmt.Printf("error - save MarketStats: %v\n", err)
	}
}

// SaveCoinGeckoMarketStats24H saves coingecko market statistics 24 hours
func (ses *StatsExporterService) SaveCoinGeckoMarketStats24H() {
	log.Println("CoingGecko Market Stats 24H")

	// Request CoinGecko API
	var coinGeckoMarketInfo types.CoinGeckoMarketInfo
	resp, err := resty.R().SetHeader("Accepts", "application/json").Get(ses.config.Market.CoinGecko.URL)
	if err != nil {
		fmt.Print("Query CoinGecko API Request error - ", err)
	}

	err = json.Unmarshal(resp.Body(), &coinGeckoMarketInfo)
	if err != nil {
		fmt.Printf("CoinGeckoMarketInfo unmarshal error - %v\n", err)
	}

	// Insert into marketStats slice
	statsCoingeckoMarket := make([]*types.StatsCoingeckoMarket24H, 0)
	tempStatsCoingeckoMarket := &types.StatsCoingeckoMarket24H{
		Price:    coinGeckoMarketInfo.MarketData.CurrentPrice.Usd,
		Currency: types.Currency,
		Time:     time.Now(),
	}
	statsCoingeckoMarket = append(statsCoingeckoMarket, tempStatsCoingeckoMarket)

	// Save
	_, err = ses.db.Model(&statsCoingeckoMarket).Insert()
	if err != nil {
		fmt.Printf("error - save MarketStats: %v\n", err)
	}
}
