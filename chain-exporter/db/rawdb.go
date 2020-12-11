package db

import (
	"fmt"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Database implements a wrapper of golang ORM with focus on PostgreSQL.
type RawDatabase struct {
	*pg.DB
}

const (
	// Define PostgreSQL database indexes to improve the speed of data retrieval operations on a database tables.
	indexTransactionHeight = "CREATE INDEX transaction_height_idx ON transaction USING btree(height);"
	indexTransactionHash   = "CREATE INDEX transaction_tx_hash_idx ON transaction USING btree(tx_hash);"
)

// Connect opens a database connections with the given database connection info from config.
func RawDBConnect(config *config.Database) *RawDatabase {
	db := pg.Connect(&pg.Options{
		Addr:     config.Host + ":" + config.Port,
		User:     config.User,
		Password: config.Password,
		Database: config.Table,
	})

	// Disable pluralization
	orm.SetTableNameInflector(func(s string) string {
		return s
	})

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
		(*schema.Transaction)(nil)} {

		// Disable pluralization
		orm.SetTableNameInflector(func(s string) string {
			return s
		})

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
		_, err := db.Model(schema.Transaction{}).Exec(indexTransactionHeight)
		if err != nil {
			return fmt.Errorf("failed to create tx hash index: %s", err)
		}
		_, err = db.Model(schema.Transaction{}).Exec(indexTransactionHash)
		if err != nil {
			return fmt.Errorf("failed to create tx hash index: %s", err)
		}

		return nil
	})

	return nil
}

// InsertExportedData saves exported blockchain data
// if function returns an error transaction is rollbacked, otherwise transaction is committed.
func (db *RawDatabase) InsertExportedData(e *schema.ExportRawData) error {
	err := db.RunInTransaction(func(tx *pg.Tx) error {
		// if e.ResultBlock.BlockHash != "" {
		// 	err := tx.Insert(&e.ResultBlock)
		// 	if err != nil {
		// 		return fmt.Errorf("failed to insert result block: %s", err)
		// 	}
		// }

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
	var tx schema.Transaction
	err := db.Model(&tx).
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

	return tx.Height, nil
}
