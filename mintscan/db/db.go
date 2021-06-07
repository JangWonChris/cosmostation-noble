package db

import (
	"fmt"
	"time"

	mblconfig "github.com/cosmostation/mintscan-backend-library/config"
	mddb "github.com/cosmostation/mintscan-database/db"
	mdschema "github.com/cosmostation/mintscan-database/schema"
	pg "github.com/go-pg/pg/v10"
)

// Database implements a wrapper of golang ORM with focus on PostgreSQL.
type Database struct {
	*mddb.Database
}

// Connect opens a database connections with the given database connection info from config.
func Connect(dbcfg *mblconfig.DatabaseConfig) *Database {
	db := mddb.Connect(dbcfg.Host, dbcfg.Port, dbcfg.User, dbcfg.Password, dbcfg.DBName, dbcfg.CommonSchema, dbcfg.ChainSchema, dbcfg.Timeout)
	fmt.Println("common schema :", dbcfg.CommonSchema)
	fmt.Println("chain schema :", dbcfg.ChainSchema)

	return &Database{db}
}

// QueryTransactionsInBlockHeight returns transactions that are included in a single block.
//QueryTransactionsByBlockHeight 에서 이름 변경
func (db *Database) QueryTransactionsInBlockHeight(chain_info_id int, height int64) ([]mdschema.Transaction, error) {
	var txs []mdschema.Transaction
	err := db.Model(&txs).
		Column("chunk").
		Where("height = ? and chain_info_id = ?", height, chain_info_id).
		Select()

	if err != nil {
		if err == pg.ErrNoRows {
			return []mdschema.Transaction{}, nil
		}
		return []mdschema.Transaction{}, err
	}

	return txs, nil
}

// QueryLatestTwoBlocks() 대체용 함수 - 응답 구조체 파싱 오류
/*
	Error:      	Received unexpected error:
	            	pg: Model(unsupported *time.Time)
*/
func (db *Database) QueryBlockTimeDiff() (string, error) {
	var BlockTimeDiff time.Time
	query := fmt.Sprintf("select timestamp-lead(timestamp) over (order by timestamp desc) as block_time_diff from (select timestamp from %s.block order by id desc limit 2) a limit 1", mdschema.GetCommonSchema())
	_, err := db.Query(&BlockTimeDiff, query)

	if err != nil {
		if err == pg.ErrNoRows {
			return BlockTimeDiff.String(), nil
		}
		return BlockTimeDiff.String(), err
	}

	return BlockTimeDiff.String(), nil
}

// QueryBondedRateIn1D return bonded rate in network from 1 day network stats table.
// func (db *Database) QueryBondedRateIn1D() ([]mdschema.StatsNetwork1D, error) {
// 	var networkStats []mdschema.StatsNetwork1D
// 	err := db.Model(&networkStats).
// 		Order("id DESC").
// 		Limit(2).
// 		Select()

// 	if err != nil {
// 		return networkStats, err
// 	}

// 	return networkStats, nil
// }
