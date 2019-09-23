package types

import (
	"encoding/json"
	"time"
)

// Info를 붙여야할까? 이름 다시 생각해보기
// CoinGeckoMarket is a struct for CoinGecko API
type CoinGeckoMarket struct {
	ID                  string          `json:"id"`
	Symbol              string          `json:"symbol"`
	Name                string          `json:"name"`
	BlockTimeInMinutes  int             `json:"block_time_in_minutes"`
	Categories          json.RawMessage `json:"categories"`
	Localization        json.RawMessage `json:"localization"`
	Description         json.RawMessage `json:"description"`
	Links               json.RawMessage `json:"links"`
	Image               json.RawMessage `json:"image"`
	CountryOrigin       string          `json:"country_origin"`
	GenesisDate         json.RawMessage `json:"genesis_date"`
	IcoData             json.RawMessage `json:"ico_data"`
	MarketCapRank       int             `json:"market_cap_rank"`
	CoingeckoRank       int             `json:"coingecko_rank"`
	CoingeckoScore      float64         `json:"coingecko_score"`
	DeveloperScore      float64         `json:"developer_score"`
	CommunityScore      float64         `json:"community_score"`
	LiquidityScore      float64         `json:"liquidity_score"`
	PublicInterestScore float64         `json:"public_interest_score"`
	MarketData          struct {
		CurrentPrice struct {
			Aed float64 `json:"aed"`
			Ars float64 `json:"ars"`
			Aud float64 `json:"aud"`
			Bch float64 `json:"bch"`
			Bdt float64 `json:"bdt"`
			Bhd float64 `json:"bhd"`
			Bmd float64 `json:"bmd"`
			Bnb float64 `json:"bnb"`
			Brl float64 `json:"brl"`
			Btc float64 `json:"btc"`
			Cad float64 `json:"cad"`
			Chf float64 `json:"chf"`
			Clp float64 `json:"clp"`
			Cny float64 `json:"cny"`
			Czk float64 `json:"czk"`
			Dkk float64 `json:"dkk"`
			Eos float64 `json:"eos"`
			Eth float64 `json:"eth"`
			Eur float64 `json:"eur"`
			Gbp float64 `json:"gbp"`
			Hkd float64 `json:"hkd"`
			Huf float64 `json:"huf"`
			Idr float64 `json:"idr"`
			Ils float64 `json:"ils"`
			Inr float64 `json:"inr"`
			Jpy float64 `json:"jpy"`
			Krw float64 `json:"krw"`
			Kwd float64 `json:"kwd"`
			Lkr float64 `json:"lkr"`
			Ltc float64 `json:"ltc"`
			Mmk float64 `json:"mmk"`
			Mxn float64 `json:"mxn"`
			Myr float64 `json:"myr"`
			Nok float64 `json:"nok"`
			Nzd float64 `json:"nzd"`
			Php float64 `json:"php"`
			Pkr float64 `json:"pkr"`
			Pln float64 `json:"pln"`
			Rub float64 `json:"rub"`
			Sar float64 `json:"sar"`
			Sek float64 `json:"sek"`
			Sgd float64 `json:"sgd"`
			Thb float64 `json:"thb"`
			Try float64 `json:"try"`
			Twd float64 `json:"twd"`
			Usd float64 `json:"usd"`
			Vef float64 `json:"vef"`
			Vnd float64 `json:"vnd"`
			Xag float64 `json:"xag"`
			Xau float64 `json:"xau"`
			Xdr float64 `json:"xdr"`
			Xlm float64 `json:"xlm"`
			Xrp float64 `json:"xrp"`
			Zar float64 `json:"zar"`
		} `json:"current_price"`
		Roi                                json.RawMessage `json:"roi"`
		Ath                                json.RawMessage `json:"ath"`
		AthChangePercentage                json.RawMessage `json:"ath_change_percentage"`
		AthDate                            json.RawMessage `json:"ath_date"`
		MarketCap                          json.RawMessage `json:"market_cap"`
		MarketCapRank                      json.RawMessage `json:"market_cap_rank"`
		TotalVolume                        json.RawMessage `json:"total_volume"`
		High24H                            json.RawMessage `json:"high_24h"`
		Low24H                             json.RawMessage `json:"low_24h"`
		PriceChange24H                     json.RawMessage `json:"price_change_24h"`
		PriceChangePercentage24H           json.RawMessage `json:"price_change_percentage_24h"`
		PriceChangePercentage7D            json.RawMessage `json:"price_change_percentage_7d"`
		PriceChangePercentage14D           json.RawMessage `json:"price_change_percentage_14d"`
		PriceChangePercentage30D           json.RawMessage `json:"price_change_percentage_30d"`
		PriceChangePercentage60D           json.RawMessage `json:"price_change_percentage_60d"`
		PriceChangePercentage200D          json.RawMessage `json:"price_change_percentage_200d"`
		PriceChangePercentage1Y            json.RawMessage `json:"price_change_percentage_1y"`
		MarketCapChange24H                 json.RawMessage `json:"market_cap_change_24h"`
		MarketCapChangePercentage24H       json.RawMessage `json:"market_cap_change_percentage_24h"`
		PriceChange24HInCurrency           json.RawMessage `json:"price_change_24h_in_currency"`
		PriceChangePercentage1HInCurrency  json.RawMessage `json:"price_change_percentage_1h_in_currency"`
		PriceChangePercentage24HInCurrency json.RawMessage `json:"price_change_percentage_24h_in_currency"`
		PriceChangePercentage7DInCurrency  json.RawMessage `json:"price_change_percentage_7d_in_currency"`
		PriceChangePercentage14DInCurrency json.RawMessage `json:"price_change_percentage_14d_in_currency"`
		PriceChangePercentage30DInCurrency json.RawMessage `json:"price_change_percentage_30d_in_currency"`
		PriceChangePercentage60DInCurrency struct {
		} `json:"price_change_percentage_60d_in_currency"`
		PriceChangePercentage200DInCurrency struct {
		} `json:"price_change_percentage_200d_in_currency"`
		PriceChangePercentage1YInCurrency struct {
		} `json:"price_change_percentage_1y_in_currency"`
		MarketCapChange24HInCurrency           json.RawMessage `json:"market_cap_change_24h_in_currency"`
		MarketCapChangePercentage24HInCurrency json.RawMessage `json:"market_cap_change_percentage_24h_in_currency"`
		TotalSupply                            json.RawMessage `json:"total_supply"`
		CirculatingSupply                      json.RawMessage `json:"circulating_supply"`
		LastUpdated                            json.RawMessage `json:"last_updated"`
	} `json:"market_data"`
	CommunityData       json.RawMessage `json:"community_data"`
	DeveloperData       json.RawMessage `json:"developer_data"`
	PublicInterestStats json.RawMessage `json:"public_interest_stats"`
	StatusUpdates       json.RawMessage `json:"status_updates"`
	LastUpdated         time.Time       `json:"last_updated"`
	Tickers             json.RawMessage `json:"tickers"`
}
