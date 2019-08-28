package types

// Coinmarketcap API
type CoinmarketcapQuotes struct {
	Status struct {
		Timestamp    string `json:"timestamp"`
		ErrorCode    int64  `json:"error_code"`
		ErrorMessage string `json:"error_message"`
		Elapsed      int64  `json:"elapsed"`
		CreditCount  int64  `json:"credit_count"`
	} `json:"status"`
	Data struct {
		Num struct {
			ID                int64    `json:"id"`
			Name              string   `json:"name"`
			Symbol            string   `json:"symbol"`
			Slug              string   `json:"slug"`
			CirculatingSupply float64  `json:"circulating_supply"`
			TotalSupply       float64  `json:"total_supply"`
			MaxSupply         int64    `json:"max_supply"`
			DateAdded         string   `json:"date_added"`
			NumMarketPairs    int64    `json:"num_market_pairs"`
			Tags              []string `json:"tags"`
			Platform          string   `json:"platform"`
			CmcRank           int64    `json:"cmc_rank"`
			LastUpdated       string   `json:"last_updated"`
			Quote             struct {
				USD struct {
					Price            float64 `json:"price"`
					Volume24H        float64 `json:"volume_24h"`
					PercentChange1H  float64 `json:"percent_change_1h"`
					PercentChange24H float64 `json:"percent_change_24h"`
					PercentChange7D  float64 `json:"percent_change_7d"`
					MarketCap        float64 `json:"market_cap"`
					LastUpdated      string  `json:"last_updated"`
				} `json:"USD"`
			} `json:"quote"`
		} `json:"3794"`
	} `json:"data"`
}
