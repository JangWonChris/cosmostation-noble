package exporter

import (
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/schema"

	"go.uber.org/zap"
)

const (
	// CoinID is identification that needs when requesting CoingGecko Market API.
	CoinID = "cosmos"

	// Currency is the currency parsed from CoinGecko Market API.
	Currency = "usd"
)

// SaveStatsMarket5M saves market statistics every 5 minutes.
func (ex *Exporter) SaveStatsMarket5M() {
	data, err := ex.client.CoinMarketData(CoinID)
	if err != nil {
		zap.S().Errorf("failed to get market data: %s", err)
		return
	}

	market := &schema.StatsMarket5M{
		Price:             data.MarketData.CurrentPrice.Usd,
		Currency:          Currency,
		MarketCapRank:     data.MarketCapRank,
		CoinGeckoRank:     data.CoingeckoRank,
		PercentChange1H:   data.MarketData.PriceChangePercentage1HInCurrency.Usd,
		PercentChange24H:  data.MarketData.PriceChangePercentage24HInCurrency.Usd,
		PercentChange7D:   data.MarketData.PriceChangePercentage7DInCurrency.Usd,
		PercentChange30D:  data.MarketData.PriceChangePercentage30DInCurrency.Usd,
		TotalVolume:       data.MarketData.TotalVolume.Usd,
		CirculatingSupply: data.MarketData.CirculatingSupply,
		LastUpdated:       data.LastUpdated,
	}

	err = ex.db.InsertMarket5M(market)
	if err != nil {
		zap.S().Errorf("failed to save market data: %s", err)
		return
	}

	zap.S().Info("successfully saved StatsMarket")
	return
}

// SaveStatsMarket1H saves market statistics every hour.
func (ex *Exporter) SaveStatsMarket1H() {
	resp, err := ex.client.CoinMarketData(CoinID)
	if err != nil {
		zap.S().Errorf("failed to get market data: %s", err)
		return
	}

	market := &schema.StatsMarket1H{
		Price:             resp.MarketData.CurrentPrice.Usd,
		Currency:          Currency,
		MarketCapRank:     resp.MarketCapRank,
		CoinGeckoRank:     resp.CoingeckoRank,
		PercentChange1H:   resp.MarketData.PriceChangePercentage1HInCurrency.Usd,
		PercentChange24H:  resp.MarketData.PriceChangePercentage24HInCurrency.Usd,
		PercentChange7D:   resp.MarketData.PriceChangePercentage7DInCurrency.Usd,
		PercentChange30D:  resp.MarketData.PriceChangePercentage30DInCurrency.Usd,
		TotalVolume:       resp.MarketData.TotalVolume.Usd,
		CirculatingSupply: resp.MarketData.CirculatingSupply,
		LastUpdated:       resp.LastUpdated,
	}

	err = ex.db.InsertMarket1H(market)
	if err != nil {
		zap.S().Errorf("failed to save market data: %s", err)
		return
	}

	zap.S().Info("successfully saved StatsMarket1H")
	return
}

// SaveStatsMarket1D saves market statistics every day.
func (ex *Exporter) SaveStatsMarket1D() {
	resp, err := ex.client.CoinMarketData(CoinID)
	if err != nil {
		zap.S().Errorf("failed to get market data: %s", err)
		return
	}

	market := &schema.StatsMarket1D{
		Price:             resp.MarketData.CurrentPrice.Usd,
		Currency:          Currency,
		MarketCapRank:     resp.MarketCapRank,
		CoinGeckoRank:     resp.CoingeckoRank,
		PercentChange1H:   resp.MarketData.PriceChangePercentage1HInCurrency.Usd,
		PercentChange24H:  resp.MarketData.PriceChangePercentage24HInCurrency.Usd,
		PercentChange7D:   resp.MarketData.PriceChangePercentage7DInCurrency.Usd,
		PercentChange30D:  resp.MarketData.PriceChangePercentage30DInCurrency.Usd,
		TotalVolume:       resp.MarketData.TotalVolume.Usd,
		CirculatingSupply: resp.MarketData.CirculatingSupply,
		LastUpdated:       resp.LastUpdated,
	}

	err = ex.db.InsertMarket1D(market)
	if err != nil {
		zap.S().Errorf("failed to save market data", err)
		return
	}

	zap.S().Info("successfully saved StatsMarket1D")
	return
}
