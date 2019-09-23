package types

import "time"

const (
	Currency = "USD"
)

/*
	Every 1 hour
*/

type StatsCoingeckoMarket1H struct {
	ID       int64     `json:"id" sql:",pk"`
	Price    float64   `json:"price"`
	Currency string    `json:"currency"`
	Time     time.Time `json:"time"`
}

type StatsCoinmarketcapMarket1H struct {
	ID        int64     `json:"id" sql:",pk"`
	Price     float64   `json:"price"`
	Currency  string    `json:"currency"`
	Volume24H float64   `json:"volumt_24h"`
	Time      time.Time `json:"time"`
}

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

type StatsCoingeckoMarket24H struct {
	ID       int64     `json:"id" sql:",pk"`
	Price    float64   `json:"price"`
	Currency string    `json:"currency"`
	Time     time.Time `json:"time"`
}

type StatsCoinmarketcapMarket24H struct {
	ID        int64     `json:"id" sql:",pk"`
	Price     float64   `json:"price"`
	Currency  string    `json:"currency"`
	Volume24H float64   `json:"volumt_24h"`
	Time      time.Time `json:"time"`
}

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
