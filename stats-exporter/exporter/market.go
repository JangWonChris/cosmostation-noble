package exporter

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/types"

	resty "gopkg.in/resty.v1"
)

func (ses *StatsExporterService) SaveCoinGeckoMarketStats() {
	log.Println("CoingGecko Market Stats")

	// Request CoinGecko API
	resp, err := resty.R().
		SetHeader("Accepts", "application/json").
		Get(ses.config.CoinGeckoURL)
	if err != nil {
		fmt.Print("CoinGecko API Request error - ", err)
	}

	// Unmarshal the response
	var coinGeckoMarketInfo types.CoinGeckoMarketInfo
	err = json.Unmarshal(resp.Body(), &coinGeckoMarketInfo)
	if err != nil {
		fmt.Printf("CoinGeckoMarketInfo unmarshal error - %v\n", err)
	}

	// Insert into marketStats slice
	marketStats := make([]*types.CoingeckoMarketStats, 0)
	tempMarketStats := &types.CoingeckoMarketStats{
		Price:            coinGeckoMarketInfo.MarketData.CurrentPrice.Usd,
		Currency:         "USD",
		PercentChange1H:  coinGeckoMarketInfo.MarketData.PriceChangePercentage1HInCurrency.Usd,
		PercentChange24H: coinGeckoMarketInfo.MarketData.PriceChangePercentage24HInCurrency.Usd,
		PercentChange7D:  coinGeckoMarketInfo.MarketData.PriceChangePercentage7DInCurrency.Usd,
		LastUpdated:      coinGeckoMarketInfo.MarketData.LastUpdated,
		Time:             time.Now(),
	}
	marketStats = append(marketStats, tempMarketStats)

	// Save
	_, err = ses.db.Model(&marketStats).Insert()
	if err != nil {
		fmt.Printf("error - save MarketStats: %v\n", err)
	}
}

func (ses *StatsExporterService) SaveCoinMarketCapMarketStats() {
	log.Println("CoinMarketCap Market Stats")

	// Request CoinMarketCap API
	resp, err := resty.R().
		SetQueryParam("id", ses.config.Coinmarketcap.CoinID). // Cosmos ID
		SetQueryParam("convert", "USD").
		SetHeader("Accepts", "application/json").
		SetHeader("X-CMC_PRO_API_KEY", ses.config.Coinmarketcap.APIKey). // API KEY [회사 계정으로 만들어야 될 필요가 있다. 요청 건수 제한]
		Get(ses.config.CoinmarketcapURL)
	if err != nil {
		fmt.Print("CoinMarketCap API Request error - ", err)
	}

	// Unmarshal the response
	var coinMarketCapQuotes types.CoinmarketcapQuotes
	err = json.Unmarshal(resp.Body(), &coinMarketCapQuotes)
	if err != nil {
		fmt.Printf("CoinMarketCapQuotes unmarshal error - %v\n", err)
	}

	// Insert into marketStats slice
	marketStats := make([]*types.CoinmarketcapMarketStats, 0)
	tempMarketStats := &types.CoinmarketcapMarketStats{
		Price:            coinMarketCapQuotes.Data.Num3794.Quote.USD.Price,
		Currency:         "USD",
		Volume24H:        coinMarketCapQuotes.Data.Num3794.Quote.USD.Volume24H,
		PercentChange1H:  coinMarketCapQuotes.Data.Num3794.Quote.USD.PercentChange1H,
		PercentChange24H: coinMarketCapQuotes.Data.Num3794.Quote.USD.PercentChange24H,
		PercentChange7D:  coinMarketCapQuotes.Data.Num3794.Quote.USD.PercentChange7D,
		LastUpdated:      coinMarketCapQuotes.Data.Num3794.Quote.USD.LastUpdated,
		Time:             time.Now(),
	}
	marketStats = append(marketStats, tempMarketStats)

	// Save
	_, err = ses.db.Model(&marketStats).Insert()
	if err != nil {
		fmt.Printf("error - save MarketStats: %v\n", err)
	}
}
