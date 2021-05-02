package exporter

import (
	"github.com/cosmostation/cosmostation-cosmos/chain-config/custom"
	"github.com/cosmostation/mintscan-backend-library/db/schema"
	"go.uber.org/zap"
)

func (ex *Exporter) SaveStatsMarket5M() {
	data, err := ex.client.GetCoinGeckoMarketData(custom.CoinGeckgoCoinID)
	if err != nil {
		zap.S().Errorf("failed to get market data: %s", err)
		return
	}

	market := &schema.StatsMarket5M{
		Price:             data.MarketData.CurrentPrice.Usd,
		Currency:          custom.Currency,
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
