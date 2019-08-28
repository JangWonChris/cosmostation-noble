package models

// API 나눈 뒤 지워도 되는 부분
// type ResultAccountResponse struct {
// 	Balance              []Coin                 `json:"balance"`
// 	Rewards              []Coin                 `json:"rewards"`
// 	Commission           []Coin                 `json:"commission"`
// 	Delegations          []Delegations          `json:"delegations"`
// 	UnbondingDelegations []UnbondingDelegations `json:"unbonding_delegations"`
// }

/*
	LCD
*/

type Delegations struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
	Balance          string `json:"balance"`
}

type ResultRewards struct {
	Rewards []Rewards `json:"rewards"`
	Total   []Coin    `json:"total"`
}

type Rewards struct {
	ValidatorAddress string `json:"validator_address"`
	Reward           []Coin `json:"reward"`
}

type ResultDelegations struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Moniker          string `json:"moniker"`
	Shares           string `json:"shares"`
	Balance          string `json:"balance"`
	Amount           string `json:"amount"`
	Rewards          []Coin `json:"delegator_rewards"`
}

type UnbondingDelegations struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Moniker          string `json:"moniker"`
	Entries          []struct {
		CreationHeight string `json:"creation_height"`
		CompletionTime string `json:"completion_time"`
		InitialBalance string `json:"initial_balance"`
		Balance        string `json:"balance"`
	} `json:"entries"`
}
