package exporter

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/types"

	resty "gopkg.in/resty.v1"
)

/*
	https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest
*/

// SaveCoinMarketCapMarketStats1H saves coinmarketcap statistics every hour
func (ses *StatsExporterService) SaveCoinMarketCapMarketStats1H() {
	log.Println("CoinMarketCap Market Stats 1H")

	// Request CoinMarketCap API
	var coinMarketCapQuotes types.CoinmarketcapQuotes
	resp, err := resty.R().
		SetQueryParam("id", ses.config.Market.CoinmarketCap.CoinID).
		SetQueryParam("convert", "USD").
		SetHeader("Accepts", "application/json").
		SetHeader("X-CMC_PRO_API_KEY", ses.config.Market.CoinmarketCap.APIKey). // API KEY [회사 계정으로 만들어야 될 필요가 있다. 요청 건수 제한]
		Get(ses.config.Market.CoinmarketCap.URL)
	if err != nil {
		fmt.Print("Query CoinMarketCap API request error - ", err)
	}

	err = json.Unmarshal(resp.Body(), &coinMarketCapQuotes)
	if err != nil {
		fmt.Printf("CoinMarketCapQuotes unmarshal error - %v\n", err)
	}

	// Insert into marketStats slice
	statsCoinmarketcapMarket := make([]*types.StatsCoinmarketcapMarket1H, 0)
	tempStatsCoinmarketcapMarket := &types.StatsCoinmarketcapMarket1H{
		Price:     coinMarketCapQuotes.Data.Num.Quote.USD.Price,
		Currency:  types.Currency,
		Volume24H: coinMarketCapQuotes.Data.Num.Quote.USD.Volume24H,
		Time:      time.Now(),
	}
	statsCoinmarketcapMarket = append(statsCoinmarketcapMarket, tempStatsCoinmarketcapMarket)

	// Save
	_, err = ses.db.Model(&statsCoinmarketcapMarket).Insert()
	if err != nil {
		fmt.Printf("error - save MarketStats: %v\n", err)
	}
}

// SaveCoinMarketCapMarketStats24H saves coinmarketcap statistics 24 hours
func (ses *StatsExporterService) SaveCoinMarketCapMarketStats24H() {
	log.Println("CoinMarketCap Market Stats 24H")

	// Request CoinMarketCap API
	var coinMarketCapQuotes types.CoinmarketcapQuotes
	resp, err := resty.R().
		SetQueryParam("id", ses.config.Market.CoinmarketCap.CoinID). // Cosmos ID
		SetQueryParam("convert", "USD").
		SetHeader("Accepts", "application/json").
		SetHeader("X-CMC_PRO_API_KEY", ses.config.Market.CoinmarketCap.APIKey). // API KEY [회사 계정으로 만들어야 될 필요가 있다. 요청 건수 제한]
		Get(ses.config.Market.CoinmarketCap.URL)
	if err != nil {
		fmt.Print("Query CoinMarketCap API request error - ", err)
	}

	err = json.Unmarshal(resp.Body(), &coinMarketCapQuotes)
	if err != nil {
		fmt.Printf("CoinMarketCapQuotes unmarshal error - %v\n", err)
	}

	// Insert into marketStats slice
	statsCoinmarketcapMarket := make([]*types.StatsCoinmarketcapMarket24H, 0)
	tempStatsCoinmarketcapMarket := &types.StatsCoinmarketcapMarket24H{
		Price:     coinMarketCapQuotes.Data.Num.Quote.USD.Price,
		Currency:  types.Currency,
		Volume24H: coinMarketCapQuotes.Data.Num.Quote.USD.Volume24H,
		Time:      time.Now(),
	}
	statsCoinmarketcapMarket = append(statsCoinmarketcapMarket, tempStatsCoinmarketcapMarket)

	// Save
	_, err = ses.db.Model(&statsCoinmarketcapMarket).Insert()
	if err != nil {
		fmt.Printf("error - save MarketStats: %v\n", err)
	}
}
