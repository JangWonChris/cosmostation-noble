package models

type ResultAccountResponse struct {
	Balance              []Balance              `json:"balance"`
	Rewards              []Rewards              `json:"rewards"`
	Commission           []Commission           `json:"commission"`
	Delegations          []Delegations          `json:"delegations"`
	UnbondingDelegations []UnbondingDelegations `json:"unbonding_delegations"`
}

type ResultDelegations struct {
	DelegatorAddress string  `json:"delegator_address"`
	ValidatorAddress string  `json:"validator_address"`
	Moniker          string  `json:"moniker"`
	Shares           string  `json:"shares"`
	Amount           string  `json:"amount"`
	Rewards          Rewards `json:"delegator_rewards"`
}

/*
	LCD
*/

type Balance struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Rewards struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Commission struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Delegations struct {
	DelegatorAddress string  `json:"delegator_address"`
	ValidatorAddress string  `json:"validator_address"`
	Moniker          string  `json:"moniker"`
	Shares           string  `json:"shares"`
	Amount           string  `json:"amount"`
	Rewards          Rewards `json:"delegator_rewards"`
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
