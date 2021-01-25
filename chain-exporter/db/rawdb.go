package db

import (

	//mbl
	lconfig "github.com/cosmostation/mintscan-backend-library/config"
	ldb "github.com/cosmostation/mintscan-backend-library/db"
)

// Database implements a wrapper of golang ORM with focus on PostgreSQL.
type RawDatabase struct {
	*ldb.RawDatabase
}

// Connect opens a database connections with the given database connection info from config.
func RawDBConnect(config *lconfig.DatabaseConfig) *RawDatabase {
	db := ldb.RawDBConnect(config)

	return &RawDatabase{db}
}

// CreateTables creates database tables using ORM (Object Relational Mapper).
func (db *RawDatabase) CreateTablesAndIndexes() {
	// 생성 오류 시 패닉
	db.CreateTables()
}
