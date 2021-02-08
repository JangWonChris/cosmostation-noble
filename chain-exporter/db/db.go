package db

import (

	//mbl
	lconfig "github.com/cosmostation/mintscan-backend-library/config"
	ldb "github.com/cosmostation/mintscan-backend-library/db"
	"github.com/cosmostation/mintscan-backend-library/db/schema"
	"github.com/cosmostation/mintscan-backend-library/types"

	"github.com/go-pg/pg"
)

var (
	// columnLength is the column length of varchar type in every table.
	// This needs to be considered again to set it to what specific length is needed, but right now set it to 99999.
	columnLength = 99999
)

// Database implements a wrapper of golang ORM with focus on PostgreSQL.
type Database struct {
	*ldb.Database
}

// Connect opens a database connections with the given database connection info from config.
func Connect(config *lconfig.DatabaseConfig) *Database {
	db := ldb.Connect(config)

	return &Database{db}
}

// Ping returns a database connection handle or an error if the connection fails.
func (db *Database) Ping() error {
	_, err := db.Exec("SELECT 1")
	if err != nil {
		return err
	}

	return nil
}

// CreateTables creates database tables using ORM (Object Relational Mapper).
func (db *Database) CreateTablesAndIndexes() {
	// 생성 오류 시 패닉
	db.CreateTables()
}

// --------------------
// Query
// --------------------

func (db *Database) QueryTxForPowerEventHistory(beginHeight, endHeight int64) ([]schema.RawTransaction, error) {
	var txs []schema.RawTransaction
	_, err := db.Query(&txs, "select t.* from stargate_final.raw_transaction t where exists ( select 1 from transaction_account as ta where height >= ? and height < ? and msg_type in (?, ?, ?, ?) and t.tx_hash = ta.tx_hash) order by t.height asc ", beginHeight, endHeight, types.StakingMsgCreateValidator, types.StakingMsgDelegate, types.StakingMsgBeginRedelegate, types.StakingMsgUndelegate)
	if err != nil {
		if err == pg.ErrNoRows {
			return txs, nil
		}
		return txs, err
	}

	return txs, nil
}

// QueryAccountMobile queries account information
func (db *Database) QueryAccountMobile(address string) (*schema.AccountMobile, error) {
	var account *schema.AccountMobile
	_ = db.Model(&account).
		Where("address = ?", address).
		Select()

	return account, nil
}

// InsertGenesisAccount insert the given genesis accounts using Copy command, it will faster than insert
// func (db *Database) InsertGenesisAccount(acc []schema.AccountCoin) error {
// 	err := db.RunInTransaction(func(tx *pg.Tx) error {
// 		if len(acc) > 0 {
// 			err := tx.Insert(&acc)
// 			if err != nil {
// 				return fmt.Errorf("failed to insert result genesis accounts: %s", err)
// 			}
// 		}
// 		return nil
// 	})

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
