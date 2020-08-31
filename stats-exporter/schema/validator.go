package schema

import "time"

// Validator defines the structure for validator information.
type Validator struct {
	ID                   int64     `sql:",pk"`
	Rank                 int       `json:"rank"`
	Address              string    `json:"address"`
	OperatorAddress      string    `json:"operator_address" sql:",unique"`
	ConsensusPubkey      string    `json:"consensus_pubkey"`
	Proposer             string    `json:"proposer"`
	Jailed               bool      `json:"jailed" sql:"default:false,notnull"`
	Status               int       `json:"status" sql:"default:0"`
	Tokens               string    `json:"tokens"`
	DelegatorShares      string    `json:"delegator_shares"`
	Moniker              string    `json:"moniker"`
	Identity             string    `json:"identity"`
	Website              string    `json:"website"`
	Details              string    `json:"details"`
	UnbondingHeight      string    `json:"unbonding_height"`
	UnbondingTime        time.Time `json:"unbonding_time" sql:"default:null"`
	CommissionRate       string    `json:"rate"`
	CommissionMaxRate    string    `json:"max_rate"`
	CommissionChangeRate string    `json:"max_change_rate"`
	UpdateTime           time.Time `json:"update_time" sql:"default:null"`
	MinSelfDelegation    string    `json:"min_self_delegation"`
	KeybaseURL           string    `json:"keybase_url"`
}

// NewValidator returns a new Validator.
func NewValidator(val Validator) *Validator {
	return &Validator{
		Rank:                 val.Rank,
		Address:              val.Address,
		OperatorAddress:      val.OperatorAddress,
		ConsensusPubkey:      val.ConsensusPubkey,
		Proposer:             val.Proposer,
		Jailed:               val.Jailed,
		Status:               val.Status,
		Tokens:               val.Tokens,
		DelegatorShares:      val.DelegatorShares,
		Moniker:              val.Moniker,
		Identity:             val.Identity,
		Website:              val.Website,
		Details:              val.Details,
		UnbondingHeight:      val.UnbondingHeight,
		UnbondingTime:        val.UnbondingTime,
		CommissionRate:       val.CommissionRate,
		CommissionMaxRate:    val.CommissionMaxRate,
		CommissionChangeRate: val.CommissionChangeRate,
		UpdateTime:           val.UpdateTime,
		MinSelfDelegation:    val.MinSelfDelegation,
		KeybaseURL:           val.KeybaseURL,
	}
}
