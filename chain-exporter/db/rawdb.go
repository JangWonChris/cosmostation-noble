package db

import (

	//mbl
	mblconfig "github.com/cosmostation/mintscan-backend-library/config"
	mdrawdb "github.com/cosmostation/mintscan-database/rawdb"
	mdschema "github.com/cosmostation/mintscan-database/schema"
)

// Database implements a wrapper of golang ORM with focus on PostgreSQL.
type RawDatabase struct {
	*mdrawdb.Database
}

// Connect opens a database connections with the given database connection info from config.
func RawDBConnect(dbcfg *mblconfig.DatabaseConfig) *RawDatabase {
	db := mdrawdb.Connect(dbcfg.Host, dbcfg.Port, dbcfg.User, dbcfg.Password, dbcfg.DBName, dbcfg.CommonSchema, dbcfg.ChainSchema, dbcfg.Timeout)
	return &RawDatabase{db}
}

// CreateTables creates database tables using ORM (Object Relational Mapper).
func (db *RawDatabase) CreateTablesAndIndexes() {
	// 생성 오류 시 패닉
	db.CreateTables()
}

// GetBlockByID returns 200 blocks block.id > id(param)
func (db *RawDatabase) GetBlockByID(id int64) ([]mdschema.RawBlock, error) {
	var b []mdschema.RawBlock
	err := db.Model(&b).
		Where("id >= ? AND id <= 1281505", id).
		Order("id ASC").
		Limit(200).
		Select()

	if err != nil {
		return nil, err
	}

	return b, nil
}

// GetTransactions
func (db *RawDatabase) GetTransactionsByID(id int64) ([]mdschema.RawTransaction, error) {
	var txs []mdschema.RawTransaction
	err := db.Model(&txs).
		Where("id >= ? AND id <= 1054686", id).
		Order("id ASC").
		Limit(200).
		Select()

	if err != nil {
		return nil, err
	}

	return txs, nil
}
