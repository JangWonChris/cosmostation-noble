package types

import (
	"encoding/json"
	"fmt"
	"time"

	resty "gopkg.in/resty.v1"
)

// ResponseWithHeight is a wrapper for returned values from REST API calls
type ResponseWithHeight struct {
	Height string          `json:"height"`
	Result json.RawMessage `json:"result"`
}

// ReadRespWithHeight reads response with height that has been changed in REST APIs from v0.36.0
func ReadRespWithHeight(resp *resty.Response) ResponseWithHeight {
	var responseWithHeight ResponseWithHeight
	err := json.Unmarshal(resp.Body(), &responseWithHeight)
	if err != nil {
		fmt.Printf("failed to unmarshal ResponseWithHeight: %v \n", err)
	}
	return responseWithHeight
}

// Pool describes Pool REST API
type Pool struct {
	NotBondedTokens string `json:"not_bonded_tokens"`
	BondedTokens    string `json:"bonded_tokens"`
}

// Inflation describes Inflation REST API
type Inflation struct {
	Height string `json:"height"`
	Result string `json:"result"`
}

// DelegatorDelegation describes DelegatorDelegation REST API
type DelegatorDelegation struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
	Balance          string `json:"balance"`
}

// ValidatorDelegation describes ValidatorDelegation REST API
type ValidatorDelegation struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
	Balance          string `json:"balance"`
}

// Validator describes Validator REST API
type Validator struct {
	OperatorAddress string `json:"operator_address"`
	ConsensusPubkey string `json:"consensus_pubkey"`
	Jailed          bool   `json:"jailed"`
	Status          int    `json:"status"`
	Tokens          string `json:"tokens"`
	DelegatorShares string `json:"delegator_shares"`
	Description     struct {
		Moniker  string `json:"moniker"`
		Identity string `json:"identity"`
		Website  string `json:"website"`
		Details  string `json:"details"`
	} `json:"description"`
	UnbondingHeight string    `json:"unbonding_height"`
	UnbondingTime   time.Time `json:"unbonding_time"`
	Commission      struct {
		CommissionRates struct {
			Rate          string `json:"rate"`
			MaxRate       string `json:"max_rate"`
			MaxChangeRate string `json:"max_change_rate"`
		} `json:"commission_rates"`
		UpdateTime time.Time `json:"update_time"`
	} `json:"commission"`
	MinSelfDelegation string `json:"min_self_delegation"`
}
