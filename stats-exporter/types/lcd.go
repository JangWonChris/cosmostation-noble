package types

type Pool struct {
	NotBondedTokens string `json:"not_bonded_tokens"`
	BondedTokens    string `json:"bonded_tokens"`
}

type Inflation struct {
	Height string `json:"height"`
	Result string `json:"height"`
}
