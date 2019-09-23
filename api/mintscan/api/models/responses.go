package models

import (
	"time"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/types"
)

// ResultBlock is a struct for block result response
type ResultBlock struct {
	Height          int64     `json:"height"`
	Proposer        string    `json:"proposer"`
	OperatorAddress string    `json:"operator_address"`
	Moniker         string    `json:"moniker"`
	BlockHash       string    `json:"block_hash"`
	Identity        string    `json:"identity"`
	NumTxs          int64     `json:"num_txs"`
	TxData          TxData    `json:"tx_data"`
	Time            time.Time `json:"time"`
}

// ResultBlocksByOperatorAddress is a struct for block result response
type ResultBlocksByOperatorAddress struct {
	Height                 int64     `json:"height"`
	Proposer               string    `json:"proposer"`
	OperatorAddress        string    `json:"operator_address"`
	Moniker                string    `json:"moniker"`
	BlockHash              string    `json:"block_hash"`
	Identity               string    `json:"identity"`
	NumTxs                 int64     `json:"num_txs"`
	TotalNumProposerBlocks int       `json:"total_num_proposer_blocks"`
	TxData                 TxData    `json:"tx_data"`
	Time                   time.Time `json:"time"`
}

// ResultBlockByHeight is a struct for block result response
type ResultBlockByHeight struct {
	Height            int64                       `json:"height"`
	Time              time.Time                   `json:"time"`
	NumTxs            int64                       `json:"num_txs"`
	Proposer          ResultBlockByHeightProposer `json:"proposer"`
	MissingValidators []struct {
		Address         []string `json:"address"`
		OperatorAddress string   `json:"operator_address"`
		Moniker         string   `json:"moniker"`
		VotingPower     float64  `json:"voting_power"`
	} `json:"missing_validators"`
	Evidence []struct {
		Address string `json:"address"`
	} `json:"evidence"`
	TxData TxData `json:"tx_data"`
}

// ResultBlockByHeightProposer is a struct for block result response
type ResultBlockByHeightProposer struct {
	Address         string  `json:"address"`
	OperatorAddress string  `json:"operator_address"`
	Moniker         string  `json:"moniker"`
	VotingPower     float64 `json:"voting_power"`
}

// ResultDelegations is a struct for delegations result response
type ResultDelegations struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Moniker          string `json:"moniker"`
	Shares           string `json:"shares"`
	Balance          string `json:"balance"`
	Amount           string `json:"amount"`
	Rewards          []Coin `json:"delegator_rewards"`
}

// ResultProposal is a struct for proposal result response
type ResultProposal struct {
	ProposalID           int64  `json:"proposal_id"`
	TxHash               string `json:"tx_hash"`
	Proposer             string `json:"proposer" sql:"default:null"`
	Moniker              string `json:"moniker" sql:"default:null"`
	Title                string `json:"title"`
	Description          string `json:"description"`
	ProposalType         string `json:"proposal_type"`
	ProposalStatus       string `json:"proposal_status"`
	Yes                  string `json:"yes"`
	Abstain              string `json:"abstain"`
	No                   string `json:"no"`
	NoWithVeto           string `json:"no_with_veto"`
	InitialDepositAmount string `json:"initial_deposit_amount" sql:"default:null"`
	InitialDepositDenom  string `json:"initial_deposit_denom" sql:"default:null"`
	TotalDepositAmount   string `json:"total_deposit_amount"`
	TotalDepositDenom    string `json:"total_deposit_denom"`
	SubmitTime           string `json:"submit_time"`
	DepositEndtime       string `json:"deposit_end_time" sql:"deposit_end_time"`
	VotingStartTime      string `json:"voting_start_time"`
	VotingEndTime        string `json:"voting_end_time"`
}

// ResultInflation is a struct for inflation result response
type ResultInflation struct {
	Inflation float64 `json:"inflation"`
}

// ResultVote is a struct for vote information result response
type ResultVote struct {
	Tally *ResultTally `json:"tally"`
	Votes []*Votes     `json:"votes"`
}

// ResultProposalDetail is a struct for deposit detail information result response
type ResultProposalDetail struct {
	ProposalID         int64             `json:"proposal_id"`
	TotalVotesNum      int               `json:"total_votes_num"`
	TotalDepositAmount float64           `json:"total_deposit_amount"`
	ResultVoteInfo     ResultVote        `json:"vote_info"`
	DepositInfo        types.DepositInfo `json:"deposit_info"`
}

// ResultStatus is a struct for status result response
type ResultStatus struct {
	ChainID                string    `json:"chain_id"`
	BlockHeight            int64     `json:"block_height"`
	BlockTime              float64   `json:"block_time"`
	TotalTxsNum            int64     `json:"total_txs_num"`
	TotalValidatorNum      int       `json:"total_validator_num"`
	UnjailedValidatorNum   int       `json:"unjailed_validator_num"`
	JailedValidatorNum     int       `json:"jailed_validator_num"`
	TotalSupplyTokens      float64   `json:"total_supply_tokens"`
	TotalCirculatingTokens float64   `json:"total_circulating_tokens"`
	BondedTokens           float64   `json:"bonded_tokens"`
	NotBondedTokens        float64   `json:"not_bonded_tokens"`
	Time                   time.Time `json:"time"`
}

// ResultTransactionInfo is a struct for tx result response
type ResultTransactionInfo struct {
	Height int64     `json:"height"`
	TxHash string    `json:"tx_hash"`
	Time   time.Time `json:"time"`
}

// ResultValidator is a struct for validator result response
type ResultValidator struct {
	Rank                 int       `json:"rank"`
	OperatorAddress      string    `json:"operator_address"`
	ConsensusPubkey      string    `json:"consensus_pubkey"`
	Jailed               bool      `json:"jailed"`
	Status               int       `json:"status"`
	Tokens               string    `json:"tokens"`
	DelegatorShares      string    `json:"delegator_shares"`
	Moniker              string    `json:"moniker"`
	Identity             string    `json:"identity"`
	Website              string    `json:"website"`
	Details              string    `json:"details"`
	UnbondingHeight      string    `json:"unbonding_height"`
	UnbondingTime        time.Time `json:"unbonding_time"`
	CommissionRate       string    `json:"rate"`
	CommissionMaxRate    string    `json:"max_rate"`
	CommissionChangeRate string    `json:"max_change_rate"`
	UpdateTime           time.Time `json:"update_time"`
	Uptime               Uptime    `json:"uptime"`
	MinSelfDelegation    string    `json:"min_self_delegation"`
	KeybaseURL           string    `json:"keybase_url"`
}

// ResultValidatorDetail is a struct for validator detail result response
type ResultValidatorDetail struct {
	Rank                 int       `json:"rank"`
	OperatorAddress      string    `json:"operator_address"`
	ConsensusPubkey      string    `json:"consensus_pubkey"`
	BondedHeight         int64     `json:"bonded_height"`
	BondedTime           time.Time `json:"bonded_time"`
	Jailed               bool      `json:"jailed"`
	Status               int       `json:"status"`
	Tokens               string    `json:"tokens"`
	DelegatorShares      string    `json:"delegator_shares"`
	Moniker              string    `json:"moniker"`
	Identity             string    `json:"identity"`
	Website              string    `json:"website"`
	Details              string    `json:"details"`
	UnbondingHeight      string    `json:"unbonding_height"`
	UnbondingTime        time.Time `json:"unbonding_time"`
	CommissionRate       string    `json:"rate"`
	CommissionMaxRate    string    `json:"max_rate"`
	CommissionChangeRate string    `json:"max_change_rate"`
	UpdateTime           time.Time `json:"update_time"`
	Uptime               Uptime    `json:"uptime"`
	MinSelfDelegation    string    `json:"min_self_delegation"`
	KeybaseURL           string    `json:"keybase_url"`
}

// ResultMisses is a struct for validator miss blocks result response
type ResultMisses struct {
	StartHeight  int64     `json:"start_height"`
	EndHeight    int64     `json:"end_height"`
	MissingCount int64     `json:"missing_count"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
}

// ResultMissesDetail is a struct for validator miss block detail result response
type ResultMissesDetail struct {
	Height int64     `json:"height"`
	Time   time.Time `json:"time"`
}

// ResultVotingPowerHistory is a struct for validator voting power history result response
type ResultVotingPowerHistory struct {
	Height         int64     `json:"height"`
	EventType      string    `json:"event_type"`
	VotingPower    float64   `json:"voting_power"`
	NewVotingPower float64   `json:"new_voting_power"`
	TxHash         string    `json:"tx_hash"`
	Timestamp      time.Time `json:"timestamp"`
}

// ResultValidatorDelegations is a struct for validator delegations result response
type ResultValidatorDelegations struct {
	TotalDelegatorNum     int                           `json:"total_delegator_num"`
	DelegatorNumChange24H int                           `json:"delegator_num_change_24h"`
	ValidatorDelegations  []*types.ValidatorDelegations `json:"delegations"`
}

// ResultRewards is a struct for rewards result response
type ResultRewards struct {
	Rewards []Rewards `json:"rewards"`
	Total   []Coin    `json:"total"`
}

// ResultTally is a struct for tally result response
type ResultTally struct {
	YesAmount        string `json:"yes_amount"`
	AbstainAmount    string `json:"abstain_amount"`
	NoAmount         string `json:"no_amount"`
	NoWithVetoAmount string `json:"no_with_veto_amount"`
	YesNum           int    `json:"yes_num"`
	AbstainNum       int    `json:"abstain_num"`
	NoNum            int    `json:"no_num"`
	NoWithVetoNum    int    `json:"no_with_veto_num"`
}

// ResultDeposit is a struct for deposit result response
type ResultDeposit struct {
	Depositor     string    `json:"depositor"`
	Moniker       string    `json:"moniker" sql:"default:null"`
	DepositAmount int64     `json:"deposit_amount"`
	DepositDenom  string    `json:"deposit_denom"`
	Height        int64     `json:"height"`
	TxHash        string    `json:"tx_hash"`
	Time          time.Time `json:"time"`
}

// ResultMarket is a struct for market result response
type ResultMarket struct {
	Price             float64       `json:"price"`
	Currency          string        `json:"currency"`
	MarketCapRank     uint8         `json:"market_cap_rank"`
	PercentChange1H   float64       `json:"percent_change_1h"`
	PercentChange24H  float64       `json:"percent_change_24h"`
	PercentChange7D   float64       `json:"percent_change_7d"`
	PercentChange30D  float64       `json:"percent_change_30d"`
	TotalVolume       uint64        `json:"total_volume"`
	CirculatingSupply float64       `json:"circulating_supply"`
	LastUpdated       time.Time     `json:"last_updated"`
	PriceStats        []*PriceStats `json:"price_stats"`
}
