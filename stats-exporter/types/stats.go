package types

import "time"

type CoingeckoMarketStats struct {
	ID               int64     `json:"id" sql:",pk"`
	Price            float64   `json:"price"`
	Currency         string    `json:"currency"`
	PercentChange1H  float64   `json:"percent_change_1h"`
	PercentChange24H float64   `json:"percent_change_24h"`
	PercentChange7D  float64   `json:"percent_change_7d"`
	LastUpdated      time.Time `json:"last_updated"`
	Time             time.Time `json:"time"`
}

type ValidatorStats struct {
	ID                int64     `json:"id" sql:",pk"`
	Moniker           string    `json:"moniker"`
	OperatorAddress   string    `json:"operator_address"`
	Address           string    `json:"address"`
	Proposer          string    `json:"proposer"`
	SelfBonded        string    `json:"self_bonded"`
	DelegatorShares   string    `json:"delegator_shares"`
	DelegatorNum      int       `json:"delegator_num"`
	VotingPowerChange string    `json:"voting_power_change"`
	Time              time.Time `json:"time"`
}

type NetworkStats struct {
	ID              int64     `json:"id" sql:",pk"`
	BlockTime       float64   `json:"block_time"`
	BondedTokens    int64     `json:"bonded_tokens"`
	BondedRatio     float64   `json:"bonded_ratio"`
	NotBondedTokens int64     `json:"not_bonded_tokens"`
	LastUpdated     time.Time `json:"last_updated"`
}

type CoinmarketcapMarketStats struct {
	ID               int64     `json:"id" sql:",pk"`
	Price            float64   `json:"price"`
	Currency         string    `json:"currency"`
	Volume24H        float64   `json:"volumt_24h"`
	PercentChange1H  float64   `json:"percent_change_1h"`
	PercentChange24H float64   `json:"percent_change_24h"`
	PercentChange7D  float64   `json:"percent_change_7d"`
	LastUpdated      string    `json:"last_updated"`
	Time             time.Time `json:"time"`
}
