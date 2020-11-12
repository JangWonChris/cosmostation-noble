package client

import (
	"encoding/json"
	"fmt"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"
)

// --------------------
// CoinGecko
// --------------------

// GetCoinGeckoMarketData returns current market data from CoinGecko API based upon params
func (c *Client) GetCoinGeckoMarketData(id string) (model.CoinGeckoMarketData, error) {
	queryStr := "/coins/" + id + "?localization=false&tickers=false&community_data=false&developer_data=false&sparkline=false"

	resp, err := c.coinGeckoClient.R().Get(queryStr)
	if err != nil {
		return model.CoinGeckoMarketData{}, err
	}

	if resp.IsError() {
		return model.CoinGeckoMarketData{}, fmt.Errorf("failed to respond: %s", err)
	}

	var data model.CoinGeckoMarketData
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return model.CoinGeckoMarketData{}, err
	}

	return data, nil
}

// GetCoinGeckoCoinPrice returns simple coin price
func (c *Client) GetCoinGeckoCoinPrice(id string) (json.RawMessage, error) {
	queryStr := "/simple/price?ids=" + id + "&vs_currencies=usd&include_market_cap=false&include_last_updated_at=true"

	resp, err := c.coinGeckoClient.R().Get(queryStr)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to respond: %s", err)
	}

	var rawData json.RawMessage
	err = json.Unmarshal(resp.Body(), &rawData)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

// GetCoinMarketChartData returns current market chart data from CoinGecko API based upon params.
func (c *Client) GetCoinMarketChartData(id string, from string, to string) (data model.CoinGeckoMarketDataChart, err error) {
	resp, err := c.coinGeckoClient.R().Get("/coins/" + id + "/market_chart/range?id=" + id + "&vs_currency=usd&from=" + from + "&to=" + to)
	if err != nil {
		return model.CoinGeckoMarketDataChart{}, err
	}

	if resp.IsError() {
		return model.CoinGeckoMarketDataChart{}, fmt.Errorf("failed to request: %s", err)
	}

	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return model.CoinGeckoMarketDataChart{}, err
	}

	return data, nil
}
