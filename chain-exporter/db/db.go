package db

import (
	"context"
	"fmt"
	"strings"

	mblconfig "github.com/cosmostation/mintscan-backend-library/config"
	mddb "github.com/cosmostation/mintscan-database/db"
	"github.com/cosmostation/mintscan-database/schema"
	mdschema "github.com/cosmostation/mintscan-database/schema"

	sdktypes "github.com/cosmos/cosmos-sdk/types"

	pg "github.com/go-pg/pg/v10"
)

var (
	// columnLength is the column length of varchar type in every table.
	// This needs to be considered again to set it to what specific length is needed, but right now set it to 99999.
	columnLength = 99999
)

// Database implements a wrapper of golang ORM with focus on PostgreSQL.
type Database struct {
	*mddb.Database
}

// Connect opens a database connections with the given database connection info from config.
func Connect(dbcfg *mblconfig.DatabaseConfig) *Database {
	db := mddb.Connect(dbcfg.Host, dbcfg.Port, dbcfg.User, dbcfg.Password, dbcfg.DBName, dbcfg.CommonSchema, dbcfg.ChainSchema, dbcfg.Timeout)
	return &Database{db}
}

// CreateTables creates database tables using ORM (Object Relational Mapper).
func (db *Database) CreateTablesAndIndexes() {
	// 생성 오류 시 패닉
	db.CreateTables()
}

func (db *Database) InsertRefineRealTimeData(e *schema.BasicData) error {
	err := db.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		err := db.InsertBlock(tx, e.Block)
		if err != nil {
			return err
		}

		if len(e.Transactions) > 0 {
			for i := range e.Transactions {
				if e.Block.ID != 0 {
					e.Transactions[i].BlockID = e.Block.ID
				} else {
					return fmt.Errorf("failed to insert result txs, can not get block.id")
				}
			}
			err := db.InsertTransaction(tx, e.Transactions, e.TMAs)
			if err != nil {
				return err
			}
		}

		return nil
	})

	// Roll back if any insertion fails.
	if err != nil {
		return err
	}

	return nil
}

// QueryAccountMobile queries account information
func (db *Database) QueryAccountMobile(address string) (*mdschema.AccountMobile, error) {
	var account *mdschema.AccountMobile
	_ = db.Model(&account).
		Where("address = ?", address).
		Select()

	return account, nil
}

// QueryAlarmTokens queries user's alarm tokens
func (db *Database) QueryAlarmTokens(address string) ([]string, error) {
	var accounts []mdschema.AccountMobile
	_ = db.Model(&accounts).
		Column("alarm_token").
		Where("address = ?", address).
		Select()

	var result []string
	if len(accounts) > 0 {
		for _, account := range accounts {
			result = append(result, account.AlarmToken)
		}
	}

	return result, nil
}

// QueryValidatorByAnyAddr returns a validator information by any type of address format
func (db *Database) QueryValidatorByAnyAddr(anyAddr string) (mdschema.Validator, error) {
	var val mdschema.Validator
	var err error

	switch {
	// jeonghwan
	case strings.HasPrefix(anyAddr, sdktypes.GetConfig().GetBech32ConsensusPubPrefix()): // Bech32 prefix for validator public key
		err = db.Model(&val).
			Where("consensus_pubkey = ?", anyAddr).
			Limit(1).
			Select()
	case strings.HasPrefix(anyAddr, sdktypes.GetConfig().GetBech32ValidatorAddrPrefix()): // Bech32 prefix for validator address
		err = db.Model(&val).
			Where("operator_address = ?", anyAddr).
			Limit(1).
			Select()
	case strings.HasPrefix(anyAddr, sdktypes.GetConfig().GetBech32AccountAddrPrefix()): // Bech32 prefix for account address
		err = db.Model(&val).
			Where("address = ?", anyAddr).
			Limit(1).
			Select()
	case len(anyAddr) == 40: // Validator consensus address in hex
		anyAddr := strings.ToUpper(anyAddr)
		err = db.Model(&val).
			Where("proposer = ?", anyAddr).
			Limit(1).
			Select()
	default:
		err = db.Model(&val).
			Where("moniker = ?", anyAddr). // Validator moniker
			Limit(1).
			Select()
	}

	if err != nil {
		if err == pg.ErrNoRows {
			return mdschema.Validator{}, nil
		}
		return mdschema.Validator{}, err
	}

	return val, nil
}
