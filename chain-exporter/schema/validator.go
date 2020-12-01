package schema

import "time"

// Validator has validator information.
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
	UnbondingHeight      int64     `json:"unbonding_height"`
	UnbondingTime        time.Time `json:"unbonding_time" sql:"default:null"`
	CommissionRate       string    `json:"rate"`
	CommissionMaxRate    string    `json:"max_rate"`
	CommissionChangeRate string    `json:"max_change_rate"`
	UpdateTime           time.Time `json:"update_time" sql:"default:null"`
	MinSelfDelegation    string    `json:"min_self_delegation"`
	KeybaseURL           string    `json:"keybase_url"`
}

// PowerEventHistory has validator's voting power event history information.
type PowerEventHistory struct {
	ID                   int64     `sql:",pk"`
	IDValidator          int       `json:"id_validator" sql:"default:0"`
	Height               int64     `json:"height"`
	Moniker              string    `json:"moniker"`
	OperatorAddress      string    `json:"operator_address"`
	Proposer             string    `json:"proposer"`
	VotingPower          float64   `json:"voting_power" sql:"default:0"`
	MsgType              string    `json:"msg_type" sql:"default:null"`
	NewVotingPowerAmount float64   `json:"new_voting_power_amount" sql:"new_voting_power_amount"`
	NewVotingPowerDenom  string    `json:"new_voting_power_denom" sql:"new_voting_power_denom"`
	TxHash               string    `json:"tx_hash" sql:"default:null"`
	Timestamp            time.Time `json:"timestamp" sql:"default:now()"`
}

// Miss has a range of every validator's missing blocks information.
type Miss struct {
	ID           int64     `json:"id" sql:",pk"`
	Address      string    `json:"address"`
	StartHeight  int64     `json:"start_height"`
	EndHeight    int64     `json:"end_height"`
	MissingCount int64     `json:"missing_count"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
}

// MissDetail has validator's missing blocks information.
type MissDetail struct {
	ID        int64     `json:"id" sql:",pk"`
	Address   string    `json:"address"`
	Height    int64     `json:"height"`
	Proposer  string    `json:"proposer_address"`
	Timestamp time.Time `json:"start_time" sql:"default:now()"`
}

// Evidence has evidence of slashing information
type Evidence struct {
	ID        int64     `json:"id" sql:",pk"`
	Proposer  string    `json:"proposer"`
	Height    int64     `json:"height"`
	Hash      string    `json:"hash"`
	Timestamp time.Time `json:"timestamp" sql:"default:now()"`
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

// NewPowerEventHistory returns a new PowerEventHistory.
func NewPowerEventHistory(p PowerEventHistory) *PowerEventHistory {
	return &PowerEventHistory{
		IDValidator:          p.IDValidator,
		Height:               p.Height,
		Moniker:              p.Moniker,
		OperatorAddress:      p.OperatorAddress,
		Proposer:             p.Proposer,
		VotingPower:          p.VotingPower,
		MsgType:              p.MsgType,
		NewVotingPowerAmount: p.NewVotingPowerAmount,
		NewVotingPowerDenom:  p.NewVotingPowerDenom,
		TxHash:               p.TxHash,
		Timestamp:            p.Timestamp,
	}
}

// NewPowerEventHistoryForGenesisValidatorSet returns a new PowerEventHistory.
func NewPowerEventHistoryForGenesisValidatorSet(p PowerEventHistory) *PowerEventHistory {
	return &PowerEventHistory{
		IDValidator:          p.IDValidator,
		Height:               p.Height,
		Proposer:             p.Proposer,
		VotingPower:          p.VotingPower,
		MsgType:              p.MsgType,
		NewVotingPowerAmount: p.NewVotingPowerAmount,
		NewVotingPowerDenom:  p.NewVotingPowerDenom,
		Timestamp:            p.Timestamp,
	}
}

// NewMiss returns a new Miss.
func NewMiss(m Miss) *Miss {
	return &Miss{
		Address:      m.Address,
		StartHeight:  m.StartHeight,
		EndHeight:    m.EndHeight,
		MissingCount: m.MissingCount,
		StartTime:    m.StartTime,
		EndTime:      m.EndTime,
	}
}

// NewMissDetail returns a new MissDetail.
func NewMissDetail(m MissDetail) *MissDetail {
	return &MissDetail{
		Address:   m.Address,
		Height:    m.Height,
		Proposer:  m.Proposer,
		Timestamp: m.Timestamp,
	}
}

// NewEvidence returns a new Evidence.
func NewEvidence(e Evidence) *Evidence {
	return &Evidence{
		Proposer:  e.Proposer,
		Height:    e.Height,
		Hash:      e.Hash,
		Timestamp: e.Timestamp,
	}
}
