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

type ValidatorStats struct {
	ID                  int64     `json:"id" sql:",pk"`
	Moniker             string    `json:"moniker"`
	OperatorAddress     string    `json:"operator_address"`
	CosmosAddress       string    `json:"cosmos_address"`
	ProposerAddress     string    `json:"proposer_address"`
	SelfBonded1H        string    `json:"self_bonded_1h"`
	DelegatorShares1H   string    `json:"delegator_shares_1h"`
	DelegatorNum1H      int       `json:"delegator_num_1h"`
	VotingPowerChange1H string    `json:"voting_power_change_1h"`
	Time                time.Time `json:"time"`
}

type NetworkStats struct {
	ID                int64     `json:"id" sql:",pk"`
	BlockTime1H       float64   `json:"block_time_1H"`
	BondedTokens1H    int64     `json:"bonded_tokens_1h"`
	BondedRatio1H     float64   `json:"bonded_ratio_1h"`
	NotBondedTokens1H int64     `json:"not_bonded_tokens_1h"`
	LastUpdated       time.Time `json:"last_updated"`
}
