package schema

import (
	"time"
)

// Account has account address information.
type Account struct {
	ID                int64     `json:"id" sql:",pk"`
	ChainID           string    `json:"chain_id" sql:",notnull"`
	AccountAddress    string    `json:"account_address"`
	AccountNumber     uint64    `json:"account_number" sql:"default:0"`
	AccountType       string    `json:"account_type"`
	CoinsTotal        string    `json:"coins_total" sql:"type:numeric(255), default:0"`
	CoinsSpendable    string    `json:"coins_spendable" sql:"type:numeric(255), default:0"`
	CoinsDelegated    string    `json:"coins_delegated" sql:"type:numeric(255), default:0"`
	CoinsUndelegated  string    `json:"coins_undelegated" sql:"type:numeric(255), default:0"`
	CoinsRewards      string    `json:"coins_rewards" sql:"type:numeric(255), default:0"`
	CoinsCommission   string    `json:"coins_commission" sql:"type:numeric(255), default:0"`
	CoinsVesting      string    `json:"coins_vesting" sql:"type:numeric(255), default:0"`
	CoinsVested       string    `json:"coins_vested" sql:"type:numeric(255), default:0"`
	CoinsFailedVested string    `json:"coins_failed_vested" sql:"type:numeric(255), default:0"`
	LastTx            string    `json:"tx"`
	LastTxTime        string    `json:"last_tx_time"`
	CreationTime      string    `json:"creation_time"`
	Timestamp         time.Time `json:"timestamp" sql:"default:now()"`
}

// AccountMobile defines an account for our mobile wallet app users.
type AccountMobile struct {
	IdfAccount  uint16    `json:"idf_account" sql:",pk"`
	ChainID     uint16    `json:"chain_id,omitempty" sql:",notnull"`
	DeviceType  string    `json:"device_type,omitempty" sql:",notnull"`
	Address     string    `json:"address" sql:",unique, notnull"`
	AlarmToken  string    `json:"alarm_token" sql:",notnull"`
	AlarmStatus bool      `json:"alarm_status" sql:",notnull"`
	Timestamp   time.Time `json:"timestamp,omitempty" sql:"default:now()"`
}

// NewAccount returns a new Account.
func NewAccount(acc Account) *Account {
	return &Account{
		ChainID:          acc.ChainID,
		AccountAddress:   acc.AccountAddress,
		AccountNumber:    acc.AccountNumber,
		AccountType:      acc.AccountType,
		CoinsTotal:       acc.CoinsTotal,
		CoinsSpendable:   acc.CoinsSpendable,
		CoinsVesting:     acc.CoinsVesting,
		CoinsDelegated:   acc.CoinsDelegated,
		CoinsUndelegated: acc.CoinsUndelegated,
		CoinsRewards:     acc.CoinsRewards,
		CoinsCommission:  acc.CoinsCommission,
		LastTx:           acc.LastTx,
		LastTxTime:       acc.LastTxTime,
		CreationTime:     acc.CreationTime,
	}
}
