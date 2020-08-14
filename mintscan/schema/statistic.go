package schema

import "time"

// StatsMarket defines the structure for market statistics
type StatsMarket struct {
	ID                int64     `json:"id" sql:",pk"`
	Price             float64   `json:"price"`
	Currency          string    `json:"currency"`
	MarketCapRank     uint8     `json:"market_cap_rank"`
	CoinGeckoRank     uint8     `json:"coingecko_rank"`
	PercentChange1H   float64   `json:"percent_change_1h"`
	PercentChange24H  float64   `json:"percent_change_24h"`
	PercentChange7D   float64   `json:"percent_change_7d"`
	PercentChange30D  float64   `json:"percent_change_30d"`
	TotalVolume       float64   `json:"total_volume"`
	CirculatingSupply float64   `json:"circulating_supply"`
	LastUpdated       time.Time `json:"last_updated"`
	Timestamp         time.Time `json:"timestamp" sql:"default:now()"`
}

/*
	Every hour
*/

// StatsMarket1H defines the structure for market statistics
type StatsMarket1H struct {
	ID                int64     `json:"id" sql:",pk"`
	Price             float64   `json:"price"`
	Currency          string    `json:"currency"`
	MarketCapRank     uint8     `json:"market_cap_rank"`
	CoinGeckoRank     uint8     `json:"coingecko_rank"`
	PercentChange1H   float64   `json:"percent_change_1h"`
	PercentChange24H  float64   `json:"percent_change_24h"`
	PercentChange7D   float64   `json:"percent_change_7d"`
	PercentChange30D  float64   `json:"percent_change_30d"`
	TotalVolume       float64   `json:"total_volume"`
	CirculatingSupply float64   `json:"circulating_supply"`
	LastUpdated       time.Time `json:"last_updated"`
	Timestamp         time.Time `json:"timestamp" sql:"default:now()"`
}

// StatsValidators1H defines the structure for validators statistics
type StatsValidators1H struct {
	ID               int64     `json:"id" sql:",pk"`
	Moniker          string    `json:"moniker"`
	Address          string    `json:"address"`
	OperatorAddress  string    `json:"operator_address"`
	ConsensusPubkey  string    `json:"consensus_pubkey"`
	Proposer         string    `json:"proposer"`
	TotalDelegations float64   `json:"total_delegations"`
	SelfBonded       float64   `json:"self_bonded"`
	Others           float64   `json:"others"`
	DelegatorNum     int       `json:"delegator_num"`
	Timestamp        time.Time `json:"timestamp" sql:"default:now()"`
}

// StatsNetwork1H defines the structure for network statistics
type StatsNetwork1H struct {
	ID              int64     `json:"id" sql:",pk"`
	BlockTime       float64   `json:"block_time"`
	BondedTokens    float64   `json:"bonded_tokens"`
	TotalSupply     float64   `json:"total_supply"`
	NotBondedTokens float64   `json:"not_bonded_tokens"`
	BondedRatio     float64   `json:"bonded_ratio"`
	InflationRatio  float64   `json:"inflation_ratio"`
	TotalTxsNum     int       `json:"total_txs_num"`
	Timestamp       time.Time `json:"timestamp" sql:"default:now()"`
}

/*
	Every 24 hours
*/

// StatsMarket1D defines the structure for market statistics
type StatsMarket1D struct {
	ID                int64     `json:"id" sql:",pk"`
	Price             float64   `json:"price"`
	Currency          string    `json:"currency"`
	MarketCapRank     uint8     `json:"market_cap_rank"`
	CoinGeckoRank     uint8     `json:"coingecko_rank"`
	PercentChange1H   float64   `json:"percent_change_1h"`
	PercentChange24H  float64   `json:"percent_change_24h"`
	PercentChange7D   float64   `json:"percent_change_7d"`
	PercentChange30D  float64   `json:"percent_change_30d"`
	TotalVolume       float64   `json:"total_volume"`
	CirculatingSupply float64   `json:"circulating_supply"`
	LastUpdated       time.Time `json:"last_updated"`
	Timestamp         time.Time `json:"timestamp" sql:"default:now()"`
}

// StatsValidators1D defines the structure for validators statistics
type StatsValidators1D struct {
	ID               int64     `json:"id" sql:",pk"`
	Moniker          string    `json:"moniker"`
	Address          string    `json:"address"`
	OperatorAddress  string    `json:"operator_address"`
	ConsensusPubkey  string    `json:"consensus_pubkey"`
	Proposer         string    `json:"proposer"`
	TotalDelegations float64   `json:"total_delegations"`
	SelfBonded       float64   `json:"self_bonded"`
	Others           float64   `json:"others"`
	DelegatorNum     int       `json:"delegator_num"`
	Timestamp        time.Time `json:"timestamp" sql:"default:now()"`
}

// StatsNetwork1D defines the structure for network statistics
type StatsNetwork1D struct {
	ID              int64     `json:"id" sql:",pk"`
	BlockTime       float64   `json:"block_time"`
	BondedTokens    float64   `json:"bonded_tokens"`
	NotBondedTokens float64   `json:"not_bonded_tokens"`
	TotalSupply     float64   `json:"total_supply"`
	BondedRatio     float64   `json:"bonded_ratio"`
	InflationRatio  float64   `json:"inflation_ratio"`
	TotalTxsNum     int       `json:"total_txs_num"`
	Timestamp       time.Time `json:"timestamp" sql:"default:now()"`
}
