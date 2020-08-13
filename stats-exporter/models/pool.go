package models

// Pool defines the structure for amount of tokens in staking pool
type Pool struct {
	NotBondedTokens string `json:"not_bonded_tokens"`
	BondedTokens    string `json:"bonded_tokens"`
}
