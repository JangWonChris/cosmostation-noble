package db

import (
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Database implements a wrapper of golang ORM with focus on PostgreSQL
type Database struct {
	*pg.DB
}

// Connect opens a database connections with the given database connection info from config.
func Connect(cfg *config.DBConfig) *Database {
	db := pg.Connect(&pg.Options{
		Addr:     cfg.Host,
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.Table,
	})
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

// CreateTables creates database tables using object relational mapper (ORM)
func (db *Database) CreateTables() error {
	for _, model := range []interface{}{(*schema.BlockCosmoshub3)(nil), (*schema.Evidence)(nil), (*schema.Miss)(nil),
		(*schema.MissDetail)(nil), (*schema.Proposal)(nil), (*schema.PowerEventHistory)(nil), (*schema.Validator)(nil),
		(*schema.TxCosmoshub3)(nil), (*schema.TxIndex)(nil), (*schema.Vote)(nil), (*schema.Deposit)(nil)} {

		// Disable pluralization
		orm.SetTableNameInflector(func(s string) string {
			return s
		})

		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
			Varchar:     20000, // replaces PostgreSQL data type `text` to `varchar(n)`
		})
		if err != nil {
			return err
		}
	}

	// RunInTransaction runs a function in a transaction.
	// if function returns an error transaction is rollbacked, otherwise transaction is committed.
	err := db.RunInTransaction(func(tx *pg.Tx) error {
		// Create indexes to reduce the cost of lookup queries in case of server traffic jams (B-Tree Index)
		_, err := db.Model(schema.BlockCosmoshub3{}).Exec(`CREATE INDEX block_cosmoshub3_height_idx ON block_cosmoshub3 USING btree(height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.Validator{}).Exec(`CREATE INDEX validator_rank_idx ON validator USING btree(rank);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.PowerEventHistory{}).Exec(`CREATE INDEX power_event_history_height_idx ON power_event_history USING btree(height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.MissDetail{}).Exec(`CREATE INDEX miss_detail_info_height_idx ON miss_detail USING btree(height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.Miss{}).Exec(`CREATE INDEX miss_info_start_height_idx ON miss USING btree(start_height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.TxCosmoshub3{}).Exec(`CREATE INDEX transaction_info_height_idx ON transaction_cosmoshub3 USING btree(height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.TxCosmoshub3{}).Exec(`CREATE INDEX transaction_info_tx_hash_idx ON transaction_cosmoshub3 USING btree(tx_hash);`)
		if err != nil {
			return err
		}

		return nil
	})

	// Roll back if any index creation fails.
	if err != nil {
		return err
	}

	return nil
}
