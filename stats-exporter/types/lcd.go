package types

// Pool is a struct for LCD Pool API
type Pool struct {
	NotBondedTokens string `json:"not_bonded_tokens"`
	BondedTokens    string `json:"bonded_tokens"`
}

// Inflation is a struct for LCD Inflation API
type Inflation struct {
	Height string `json:"height"`
	Result string `json:"result"`
}

// DelegatorDelegation is a struct for LCD DelegatorDelegation API
type DelegatorDelegation struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
	Balance          string `json:"balance"`
}

// ValidatorDelegation is a struct for LCD ValidatorDelegation API
type ValidatorDelegation struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
	Balance          string `json:"balance"`
}
