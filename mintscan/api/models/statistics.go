package models

import "time"

const (
	// Currency
	Currency = "USD"
)

/*
	Every 1 hour
*/

// StatsCoingeckoMarket1H is a struct for market statistics from Coingecko
type StatsCoingeckoMarket1H struct {
	ID                int64     `json:"id" sql:",pk"`
	Price             float64   `json:"price"`
	Currency          string    `json:"currency"`
	MarketCapRank     uint8     `json:"market_cap_rank"`
	PercentChange1H   float64   `json:"percent_change_1h"`
	PercentChange24H  float64   `json:"percent_change_24h"`
	PercentChange7D   float64   `json:"percent_change_7d"`
	PercentChange30D  float64   `json:"percent_change_30d"`
	TotalVolume       uint64    `json:"total_volume"`
	CirculatingSupply float64   `json:"circulating_supply"`
	LastUpdated       time.Time `json:"last_updated"`
	Time              time.Time `json:"time"`
}

// StatsCoinmarketcapMarket1H is a struct for market statistics from CMK
type StatsCoinmarketcapMarket1H struct {
	ID        int64     `json:"id" sql:",pk"`
	Price     float64   `json:"price"`
	Currency  string    `json:"currency"`
	Volume24H float64   `json:"volumt_24h"`
	Time      time.Time `json:"time"`
}

// StatsValidators1H is a struct for validators statistics
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
	Time             time.Time `json:"time"`
}

// StatsNetwork1H is a struct for network statistics
type StatsNetwork1H struct {
	ID              int64     `json:"id" sql:",pk"`
	BlockTime       float64   `json:"block_time"`
	BondedTokens    float64   `json:"bonded_tokens"`
	TotalSupply     float64   `json:"total_supply"`
	NotBondedTokens float64   `json:"not_bonded_tokens"`
	BondedRatio     float64   `json:"bonded_ratio"`
	InflationRatio  float64   `json:"inflation_ratio"`
	TotalTxsNum     int64     `json:"total_txs_num"`
	Time            time.Time `json:"last_updated"`
}

/*
	Every 24 hours
*/

// StatsCoingeckoMarket24H is a struct for market statistics from Coingecko
type StatsCoingeckoMarket24H struct {
	ID                int64     `json:"id" sql:",pk"`
	Price             float64   `json:"price"`
	Currency          string    `json:"currency"`
	MarketCapRank     uint8     `json:"market_cap_rank"`
	PercentChange1H   float64   `json:"percent_change_1h"`
	PercentChange24H  float64   `json:"percent_change_24h"`
	PercentChange7D   float64   `json:"percent_change_7d"`
	PercentChange30D  float64   `json:"percent_change_30d"`
	TotalVolume       uint64    `json:"total_volume"`
	CirculatingSupply float64   `json:"circulating_supply"`
	LastUpdated       time.Time `json:"last_updated"`
	Time              time.Time `json:"time"`
}

// StatsCoinmarketcapMarket24H is a struct for market statistics from CMK
type StatsCoinmarketcapMarket24H struct {
	ID        int64     `json:"id" sql:",pk"`
	Price     float64   `json:"price"`
	Currency  string    `json:"currency"`
	Volume24H float64   `json:"volumt_24h"`
	Time      time.Time `json:"time"`
}

// StatsValidators24H is a struct for validators statistics
type StatsValidators24H struct {
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
	Time             time.Time `json:"time"`
}

// StatsNetwork24H is a struct for network statistics
type StatsNetwork24H struct {
	ID              int64     `json:"id" sql:",pk"`
	BlockTime       float64   `json:"block_time"`
	BondedTokens    float64   `json:"bonded_tokens"`
	NotBondedTokens float64   `json:"not_bonded_tokens"`
	TotalSupply     float64   `json:"total_supply"`
	BondedRatio     float64   `json:"bonded_ratio"`
	InflationRatio  float64   `json:"inflation_ratio"`
	TotalTxsNum     int64     `json:"total_txs_num"`
	Time            time.Time `json:"last_updated"`
}
