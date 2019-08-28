package types

type Pool struct {
	NotBondedTokens string `json:"not_bonded_tokens"`
	BondedTokens    string `json:"bonded_tokens"`
}

type Inflation struct {
	Height string `json:"height"`
	Result string `json:"result"`
}

type DelegatorDelegation struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
	Balance          string `json:"balance"`
}

type ValidatorDelegation struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
	Balance          string `json:"balance"`
}
