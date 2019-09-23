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

	// request CoinGecko API
	var coinGeckoMarket types.CoinGeckoMarket
	resp, err := resty.R().SetHeader("Accepts", "application/json").Get(ses.config.Market.CoinGecko.URL)
	if err != nil {
		fmt.Print("query CoinGecko API request error - ", err)
	}

	err = json.Unmarshal(resp.Body(), &coinGeckoMarket)
	if err != nil {
		fmt.Printf("unmarshal coinGeckoMarket error - %v\n", err)
	}

	// insert into marketStats slice
	statsCoingeckoMarket := make([]*types.StatsCoingeckoMarket1H, 0)
	tempStatsCoingeckoMarket := &types.StatsCoingeckoMarket1H{
		Price:             coinGeckoMarket.MarketData.CurrentPrice.Usd,
		Currency:          types.Currency,
		MarketCapRank:     coinGeckoMarket.MarketCapRank,
		PercentChange1H:   coinGeckoMarket.MarketData.PriceChangePercentage1HInCurrency.Usd,
		PercentChange24H:  coinGeckoMarket.MarketData.PriceChangePercentage24HInCurrency.Usd,
		PercentChange7D:   coinGeckoMarket.MarketData.PriceChangePercentage7DInCurrency.Usd,
		PercentChange30D:  coinGeckoMarket.MarketData.PriceChangePercentage30DInCurrency.Usd,
		TotalVolume:       coinGeckoMarket.MarketData.TotalVolume.Usd,
		CirculatingSupply: coinGeckoMarket.MarketData.CirculatingSupply,
		LastUpdated:       coinGeckoMarket.LastUpdated,
		Time:              time.Now(),
	}
	statsCoingeckoMarket = append(statsCoingeckoMarket, tempStatsCoingeckoMarket)

	_, err = ses.db.Model(&statsCoingeckoMarket).Insert()
	if err != nil {
		fmt.Printf("error - save MarketStats1H: %v\n", err)
	}
}

// SaveCoinGeckoMarketStats24H saves coingecko market statistics 24 hours
func (ses *StatsExporterService) SaveCoinGeckoMarketStats24H() {
	log.Println("CoingGecko Market Stats 24H")

	// request CoinGecko API
	var coinGeckoMarket types.CoinGeckoMarket
	resp, err := resty.R().SetHeader("Accepts", "application/json").Get(ses.config.Market.CoinGecko.URL)
	if err != nil {
		fmt.Print("query CoinGecko API request error - ", err)
	}

	err = json.Unmarshal(resp.Body(), &coinGeckoMarket)
	if err != nil {
		fmt.Printf("unmarshal coinGeckoMarket error - %v\n", err)
	}

	// insert into marketStats slice
	statsCoingeckoMarket := make([]*types.StatsCoingeckoMarket24H, 0)
	tempStatsCoingeckoMarket := &types.StatsCoingeckoMarket24H{
		Price:             coinGeckoMarket.MarketData.CurrentPrice.Usd,
		Currency:          types.Currency,
		MarketCapRank:     coinGeckoMarket.MarketCapRank,
		PercentChange1H:   coinGeckoMarket.MarketData.PriceChangePercentage1HInCurrency.Usd,
		PercentChange24H:  coinGeckoMarket.MarketData.PriceChangePercentage24HInCurrency.Usd,
		PercentChange7D:   coinGeckoMarket.MarketData.PriceChangePercentage7DInCurrency.Usd,
		PercentChange30D:  coinGeckoMarket.MarketData.PriceChangePercentage30DInCurrency.Usd,
		TotalVolume:       coinGeckoMarket.MarketData.TotalVolume.Usd,
		CirculatingSupply: coinGeckoMarket.MarketData.CirculatingSupply,
		LastUpdated:       coinGeckoMarket.LastUpdated,
		Time:              time.Now(),
	}
	statsCoingeckoMarket = append(statsCoingeckoMarket, tempStatsCoingeckoMarket)

	_, err = ses.db.Model(&statsCoingeckoMarket).Insert()
	if err != nil {
		fmt.Printf("error - save MarketStats24H: %v\n", err)
	}
}
