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
	CoinsTotal        uint64    `json:"coins_total" sql:"default:0"`
	CoinsSpendable    uint64    `json:"coins_spendable" sql:"default:0"`
	CoinsDelegated    uint64    `json:"coins_delegated" sql:"default:0"`
	CoinsUndelegated  uint64    `json:"coins_undelegated" sql:"default:0" sql:"default:0"`
	CoinsRewards      uint64    `json:"coins_rewards"`
	CoinsCommission   uint64    `json:"coins_commission" sql:"default:0"`
	CoinsVesting      uint64    `json:"coins_vesting" sql:"default:0"`
	CoinsVested       uint64    `json:"coins_vested" sql:"default:0"`
	CoinsFailedVested uint64    `json:"coins_failed_vested" sql:"default:0"`
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
