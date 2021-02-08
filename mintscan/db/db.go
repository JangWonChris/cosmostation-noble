package db

import (

	// "github.com/cosmostation/mintscan-backend-library/config"

	// "github.com/cosmostation/mintscan-backend-library/db/schema"

	"github.com/cosmostation/mintscan-backend-library/config"
	ldb "github.com/cosmostation/mintscan-backend-library/db"
	"github.com/cosmostation/mintscan-backend-library/db/schema"
)

// Database implements a wrapper of golang ORM with focus on PostgreSQL.
type Database struct {
	*ldb.Database
}

// Connect opens a database connections with the given database connection info from config.
func Connect(config *config.DatabaseConfig) *Database {
	db := ldb.Connect(config)

	return &Database{db}
}

// QueryBondedRateIn1D return bonded rate in network from 1 day network stats table.
func (db *Database) QueryBondedRateIn1D() ([]schema.StatsNetwork1D, error) {
	var networkStats []schema.StatsNetwork1D
	err := db.Model(&networkStats).
		Order("id DESC").
		Limit(2).
		Select()

	if err != nil {
		return networkStats, err
	}

	return networkStats, nil
}

// QueryValidatorBondedInfo returns a validator's bonded information.
// sdk 에서 제공하는 IsBonded 함수가 존재한다.
// 이 함수가 필요한 이유는, 최초 본딩 된 날짜를 알기 위함임(제네시스인지, 그 이후 생성 된 검증인 인지)
// func (db *Database) QueryValidatorBondedInfo(address string) (peh schema.PowerEventHistory, err error) {
// 	msgType := "create_validator"

// 	err = db.Model(&peh).
// 		Where("proposer = ? AND msg_type = ?", address, msgType).
// 		Limit(1).
// 		Select()

// 	if err != nil {
// 		return schema.PowerEventHistory{}, err
// 	}

// 	return peh, nil
// }

// QueryTransactionsByAddr returns all transactions that are created by an account.
// func (db *Database) QueryTransactionsByAddr(accAddr, valAddr string, before, after, limit int) ([]schema.Transaction, error) {
// 	var txs []schema.Transaction
// 	var err error

// 	switch {
// 	case before > 0:
// 		_, err = db.Query(&txs, "select distinct t.* from transaction as t, transaction_account as ta where ta.account_address = ? and t.tx_hash = ta.tx_hash order by id desc limit ?", accAddr, limit)
// 		// params += " AND (id < ?)"
// 		// err = db.postgres.Model(&txs).
// 		// 	Where(params, accAddr, before).
// 		// 	Limit(limit).
// 		// 	Order("id DESC").
// 		// 	Select()
// 	case after > 0:
// 		_, err = db.Query(&txs, "select distinct t.* from transaction as t, transaction_account as ta where ta.account_address = ? and t.tx_hash = ta.tx_hash order by id desc limit ?", accAddr, limit)
// 		// params += " AND (id > ?)"
// 		// err = db.postgres.Model(&txs).
// 		// 	Where(params, accAddr, after).
// 		// 	Limit(limit).
// 		// 	Order("id ASC").
// 		// 	Select()
// 	default:
// 		_, err = db.Query(&txs, "select distinct t.* from transaction as t, transaction_account as ta where ta.account_address = ? and t.tx_hash = ta.tx_hash order by id desc limit ?", accAddr, limit)
// 		// err = db.postgres.Model(&txs).
// 		// 	Where(params, accAddr).
// 		// 	Limit(limit).
// 		// 	Order("id DESC").
// 		// 	Select()
// 	}

// 	if err != nil {
// 		return []schema.Transaction{}, err
// 	}

// 	return txs, nil
// }

// QueryTransferTransactionsByAddr queries Send / MultiSend transactions that are made by an account
// func (db *Database) QueryTransferTransactionsByAddr(accAddr, denom string, before, after, limit int) ([]schema.Transaction, error) {
// 	var txs []schema.Transaction
// 	var err error

// 	switch {
// 	case before > 0:
// 		_, err = db.Query(&txs, "select t.* from transaction as t left join transaction_account as t1 on t.tx_hash = t1.tx_hash where t1.account_address = ? and (t1.msg_type = 'send' or t1.msg_type = 'multisend') limit ?", accAddr, limit)
// 		// params += " AND (id < ?)"
// 		// err = db.postgres.Model(&txs).
// 		// 	Where(params, before).
// 		// 	Limit(limit).
// 		// 	Order("id DESC").
// 		// 	Select()
// 	case after > 0:
// 		_, err = db.Query(&txs, "select t.* from transaction as t left join transaction_account as t1 on t.tx_hash = t1.tx_hash where t1.account_address = ? and (t1.msg_type = 'send' or t1.msg_type = 'multisend') limit ?", accAddr, limit)
// 		// params += " AND (id > ?)"
// 		// err = db.postgres.Model(&txs).
// 		// 	Where(params, after).
// 		// 	Limit(limit).
// 		// 	Order("id ASC").
// 		// 	Select()
// 	default:
// 		_, err = db.Query(&txs, "select t.* from transaction as t left join transaction_account as t1 on t.tx_hash = t1.tx_hash where t1.account_address = ? and (t1.msg_type = 'send' or t1.msg_type = 'multisend') limit ?", accAddr, limit)
// 		// err = db.postgres.Model(&txs).
// 		// 	Where(params).
// 		// 	Limit(limit).
// 		// 	Order("id DESC").
// 		// 	Select()
// 	}

// 	if err != nil {
// 		return []schema.Transaction{}, err
// 	}

// 	return txs, nil
// }

// QueryTransactionsBetweenAccountAndValidator queries transactions that are made between an account and his delegated validator
// func (db *Database) QueryTransactionsBetweenAccountAndValidator(address, valAddr string, before, after, limit int) ([]schema.Transaction, error) {
// 	var txs []schema.Transaction
// 	var err error

// 	switch {
// 	case before > 0:
// 		_, err = db.Query(&txs, "select t.* from transaction as t, transaction_account as t1, transaction_account as t2 where t1.account_address = ? and t2.account_address = ? and t1.tx_hash = t2.tx_hash and t.tx_hash = t1.tx_hash order by id desc limit ?", address, valAddr, limit)
// 	case after > 0:
// 		_, err = db.Query(&txs, "select t.* from transaction as t, transaction_account as t1, transaction_account as t2 where t1.account_address = ? and t2.account_address = ? and t1.tx_hash = t2.tx_hash and t.tx_hash = t1.tx_hash order by id desc limit ?", address, valAddr, limit)
// 	default:
// 		_, err = db.Query(&txs, "select t.* from transaction as t, transaction_account as t1, transaction_account as t2 where t1.account_address = ? and t2.account_address = ? and t1.tx_hash = t2.tx_hash and t.tx_hash = t1.tx_hash order by id desc limit ?", address, valAddr, limit)
// 	}

// 	if err != nil {
// 		return []schema.Transaction{}, err
// 	}

// 	return txs, nil
// }

// QueryValidatorVotingPowerEventHistory returns a validator's voting power events
// func (db *Database) QueryValidatorVotingPowerEventHistory(address string, before, after, limit int) ([]schema.PowerEventHistory, error) {
// 	var peh []schema.PowerEventHistory
// 	var err error

// 	switch {
// 	case before > 0:
// 		err = db.Model(&peh).
// 			Where("operator_address = ? AND height < ?", address, before).
// 			Limit(limit).
// 			Order("id DESC").
// 			Select()
// 	case after > 0:
// 		err = db.Model(&peh).
// 			Where("operator_address = ? AND height > ?", address, after).
// 			Limit(limit).
// 			Order("id ASC").
// 			Select()
// 	default:
// 		err = db.Model(&peh).
// 			Where("operator_address = ?", address).
// 			Limit(limit).
// 			Order("id DESC").
// 			Select()
// 	}

// 	if err != nil {
// 		if err == pg.ErrNoRows {
// 			return []schema.PowerEventHistory{}, nil
// 		}
// 		return []schema.PowerEventHistory{}, err
// 	}

// 	return peh, nil
// }
