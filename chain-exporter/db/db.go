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
	for _, model := range []interface{}{(*schema.BlockInfo)(nil), (*schema.EvidenceInfo)(nil), (*schema.MissInfo)(nil),
		(*schema.MissDetailInfo)(nil), (*schema.ProposalInfo)(nil), (*schema.ValidatorSetInfo)(nil), (*schema.ValidatorInfo)(nil),
		(*schema.TransactionInfo)(nil), (*schema.VoteInfo)(nil), (*schema.DepositInfo)(nil)} {
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
		_, err := db.Model(schema.BlockInfo{}).Exec(`CREATE INDEX block_info_height_idx ON block_infos USING btree(height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.ValidatorInfo{}).Exec(`CREATE INDEX validator_info_rank_idx ON validator_infos USING btree(rank);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.ValidatorInfo{}).Exec(`CREATE INDEX validator_set_info_height_idx ON validator_set_infos USING btree(height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.MissDetailInfo{}).Exec(`CREATE INDEX miss_detail_info_height_idx ON miss_detail_infos USING btree(height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.MissInfo{}).Exec(`CREATE INDEX miss_info_start_height_idx ON miss_infos USING btree(start_height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.TransactionInfo{}).Exec(`CREATE INDEX transaction_info_height_idx ON transaction_infos USING btree(height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.TransactionInfo{}).Exec(`CREATE INDEX transaction_info_tx_hash_idx ON transaction_infos USING btree(tx_hash);`)
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
