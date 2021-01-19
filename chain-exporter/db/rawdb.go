package db

import (
	"fmt"

	//mbl
	lconfig "github.com/cosmostation/mintscan-backend-library/config"
	ldb "github.com/cosmostation/mintscan-backend-library/db"
	lschema "github.com/cosmostation/mintscan-backend-library/db/schema"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
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

// Ping returns a database connection handle or an error if the connection fails.
func (db *RawDatabase) Ping() error {
	_, err := db.Exec("SELECT 1")
	if err != nil {
		return err
	}

	return nil
}

// CreateTables creates database tables using ORM (Object Relational Mapper).
func (db *RawDatabase) CreateTables() error {
	for _, table := range []interface{}{
		(*lschema.RawBlock)(nil),
		(*lschema.RawTransaction)(nil)} {

		err := db.CreateTable(table, &orm.CreateTableOptions{
			IfNotExists: true,
			Varchar:     columnLength, // replaces PostgreSQL data type `text` to `varchar(n)`
		})

		if err != nil {
			return err
		}
	}

	// Create table indexes and roll back if any index creation fails.
	err := db.createIndexes()
	if err != nil {
		return err
	}

	return nil
}

// createIndexes uses RunInTransaction to run a function in a transaction.
// if function returns an error, transaction is rollbacked, otherwise transaction is committed.
// Create B-Tree indexes to reduce the cost of lookup queries
func (db *RawDatabase) createIndexes() error {
	db.RunInTransaction(func(tx *pg.Tx) error {
		_, err := db.Model(lschema.RawTransaction{}).Exec(ldb.GetIndex(lschema.IndexRawTransactionHeight))
		if err != nil {
			return fmt.Errorf("failed to create tx hash index: %s", err)
		}
		_, err = db.Model(lschema.RawTransaction{}).Exec(ldb.GetIndex(lschema.IndexRawTransactionHash))
		if err != nil {
			return fmt.Errorf("failed to create tx hash index: %s", err)
		}

		return nil
	})

	return nil
}

// InsertExportedData saves exported blockchain data
// if function returns an error transaction is rollbacked, otherwise transaction is committed.
func (db *RawDatabase) InsertExportedData(e *lschema.ExportRawData) error {
	err := db.RunInTransaction(func(tx *pg.Tx) error {
		if e.ResultBlockJSONChunk.BlockHash != "" {
			err := tx.Insert(&e.ResultBlockJSONChunk)
			if err != nil {
				return fmt.Errorf("failed to insert result block: %s", err)
			}
		}

		if len(e.ResultTxsJSONChunk) > 0 {
			err := tx.Insert(&e.ResultTxsJSONChunk)
			if err != nil {
				return fmt.Errorf("failed to insert result txs json chunk: %s", err)
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

// QueryLatestBlockHeight queries latest block height in database
func (db *RawDatabase) QueryLatestBlockHeight() (int64, error) {
	var b lschema.RawBlock
	err := db.Model(&b).
		Order("height DESC").
		Limit(1).
		Select()

	// return 0 when there is no row in result set
	if err == pg.ErrNoRows {
		return 0, err
	}

	// return -1 for any type of errors
	if err != nil {
		return -1, err
	}

	return b.Height, nil
}

// QueryLatestBlockHeight queries latest block height in database
func (db *RawDatabase) GetRawBlock(height int64) ([]lschema.RawBlock, error) {
	var b []lschema.RawBlock
	err := db.Model(&b).
		Where("height >= ?", height).
		Order("height ASC").
		Limit(100).
		Select()

	if err != nil {
		return nil, err
	}

	return b, nil
}

func (db *RawDatabase) GetRawTransactions(height int64) ([]lschema.RawTransaction, error) {
	var txs []lschema.RawTransaction
	err := db.Model(&txs).
		Where("height = ?", height).
		Select()

	if err != nil {
		return nil, err
	}

	return txs, nil
}
