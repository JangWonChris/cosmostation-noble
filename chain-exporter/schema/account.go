package schema

import (
	"time"
)

// // Account has account address information.
// type Account struct {
// 	ID                int64     `json:"id" sql:",pk"`
// 	ChainID           string    `json:"chain_id" sql:",notnull"`
// 	AccountAddress    string    `json:"account_address"`
// 	AccountNumber     uint64    `json:"account_number" sql:"default:0"`
// 	AccountType       string    `json:"account_type"`
// 	CoinsTotal        string    `json:"coins_total" sql:"type:jsonb, default: '[]'::jsonb"` // jsonb type since multiple coins
// 	CoinsSpendable    string    `json:"coins_Spendable" sql:"type:jsonb, default: '[]'::jsonb"`
// 	CoinsDelegated    string    `json:"coins_delegated"`
// 	CoinsUndelegated  string    `json:"coins_undelegated"`
// 	CoinsRewards      string    `json:"coins_rewards"`
// 	CoinsCommission   string    `json:"coins_commission" sql:"type:jsonb, default: '[]'::jsonb"`
// 	CoinsVesting      string    `json:"coins_vesting" sql:"type:jsonb, default: '[]'::jsonb"`
// 	CoinsVested       string    `json:"coins_vested" sql:"type:jsonb, default: '[]'::jsonb"`
// 	CoinsFailedVested string    `json:"coins_failed_vested" sql:"type:jsonb, default: '[]'::jsonb"`
// 	LastTx            string    `json:"tx"`
// 	LastTxTime        string    `json:"last_tx_time"`
// 	CreationTime      string    `json:"creation_time"`
// 	Timestamp         time.Time `json:"timestamp" sql:"default:now()"`
// }

// Account has account address information.
type Account struct {
	ID                int64     `json:"id" sql:",pk"`
	ChainID           string    `json:"chain_id" sql:",notnull"`
	AccountAddress    string    `json:"account_address"`
	AccountNumber     uint64    `json:"account_number" sql:"default:0"`
	AccountType       string    `json:"account_type"`
	CoinsTotal        string    `json:"coins_total"` // jsonb type since multiple coins
	CoinsSpendable    string    `json:"coins_Spendable"`
	CoinsDelegated    string    `json:"coins_delegated"`
	CoinsUndelegated  string    `json:"coins_undelegated"`
	CoinsRewards      string    `json:"coins_rewards"`
	CoinsCommission   string    `json:"coins_commission"`
	CoinsVesting      string    `json:"coins_vesting"`
	CoinsVested       string    `json:"coins_vested"`
	CoinsFailedVested string    `json:"coins_failed_vested"`
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
