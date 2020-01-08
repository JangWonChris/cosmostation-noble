package exporter

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/types"

	resty "gopkg.in/resty.v1"
)

// SaveCoinGeckoMarketStats1H saves coingecko market stats every hour
func (ses *StatsExporterService) SaveCoinGeckoMarketStats1H() {
	var coinGeckoMarket types.CoinGeckoMarket
	resp, err := resty.R().
		SetHeader("Accepts", "application/json").
		Get(ses.config.Market.CoinGecko.URL)
	if err != nil {
		fmt.Printf("failed to request CoinGecko API: %v \n", err)
	}

	err = json.Unmarshal(resp.Body(), &coinGeckoMarket)
	if err != nil {
		fmt.Printf("failed to unmarshal CoinGeckoMarket: %v \n", err)
	}

	statsCoingeckoMarket := &schema.StatsCoingeckoMarket1H{
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

	result, _ := ses.db.InsertCoinGeckoMarket1H(*statsCoingeckoMarket)
	if result {
		log.Println("succesfully saved CoingGecko Market Stats 1H")
	}
}

// SaveCoinGeckoMarketStats24H saves coingecko market stats 24 hours
func (ses *StatsExporterService) SaveCoinGeckoMarketStats24H() {
	var coinGeckoMarket types.CoinGeckoMarket
	resp, err := resty.R().SetHeader("Accepts", "application/json").Get(ses.config.Market.CoinGecko.URL)
	if err != nil {
		fmt.Printf("failed to request CoinGecko API: %v \n", err)
	}

	err = json.Unmarshal(resp.Body(), &coinGeckoMarket)
	if err != nil {
		fmt.Printf("failed to unmarshal CoinGeckoMarket: %v \n", err)
	}

	statsCoingeckoMarket := &schema.StatsCoingeckoMarket24H{
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

	result, _ := ses.db.InsertCoinGeckoMarket24H(*statsCoingeckoMarket)
	if result {
		log.Println("succesfully saved CoingGecko Market Stats 24H")
	}
}

// SaveCoinMarketCapMarketStats1H saves coinmarketcap stats every hour
func (ses *StatsExporterService) SaveCoinMarketCapMarketStats1H() {
	var coinMarketCapQuotes types.CoinmarketcapQuotes
	resp, err := resty.R().
		SetQueryParam("id", ses.config.Market.CoinmarketCap.CoinID).
		SetQueryParam("convert", "USD").
		SetHeader("Accepts", "application/json").
		SetHeader("X-CMC_PRO_API_KEY", ses.config.Market.CoinmarketCap.APIKey). // [TODO] API KEY - 요청 건수 제한으로 인해 회사 계정으로 만들어야 될 필요
		Get(ses.config.Market.CoinmarketCap.URL)
	if err != nil {
		fmt.Printf("failed to request CoinMarketCap API: %v \n", err)
	}

	err = json.Unmarshal(resp.Body(), &coinMarketCapQuotes)
	if err != nil {
		fmt.Printf("failed to unmarshal CoinMarketCapQuotes: %v \n", err)
	}

	statsCoinmarketcapMarket := &schema.StatsCoinmarketcapMarket1H{
		Price:     coinMarketCapQuotes.Data.Num.Quote.USD.Price,
		Currency:  types.Currency,
		Volume24H: coinMarketCapQuotes.Data.Num.Quote.USD.Volume24H,
		Time:      time.Now(),
	}

	result, _ := ses.db.InsertCoinMarketCapMarket1H(*statsCoinmarketcapMarket)
	if result {
		log.Println("succesfully saved CoinMarketCap Market Stats 1H")
	}
}

// SaveCoinMarketCapMarketStats24H saves coinmarketcap stats 24 hours
func (ses *StatsExporterService) SaveCoinMarketCapMarketStats24H() {
	var coinMarketCapQuotes types.CoinmarketcapQuotes
	resp, err := resty.R().
		SetQueryParam("id", ses.config.Market.CoinmarketCap.CoinID). // Cosmos ID
		SetQueryParam("convert", "USD").
		SetHeader("Accepts", "application/json").
		SetHeader("X-CMC_PRO_API_KEY", ses.config.Market.CoinmarketCap.APIKey).
		Get(ses.config.Market.CoinmarketCap.URL)
	if err != nil {
		fmt.Printf("failed to request CoinMarketCap API: %v \n", err)
	}

	err = json.Unmarshal(resp.Body(), &coinMarketCapQuotes)
	if err != nil {
		fmt.Printf("failed to unmarshal CoinMarketCapQuotes: %v \n", err)
	}

	statsCoinmarketcapMarket := &schema.StatsCoinmarketcapMarket24H{
		Price:     coinMarketCapQuotes.Data.Num.Quote.USD.Price,
		Currency:  types.Currency,
		Volume24H: coinMarketCapQuotes.Data.Num.Quote.USD.Volume24H,
		Time:      time.Now(),
	}

	result, _ := ses.db.InsertCoinMarketCapMarket24H(*statsCoinmarketcapMarket)
	if result {
		log.Println("succesfully saved CoinMarketCap Market Stats 24H")
	}
}
