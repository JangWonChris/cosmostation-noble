package databases

import (
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	dtypes "github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// ConnectDatabase connects to PostgreSQL
func ConnectDatabase(Config *config.Config) *pg.DB {
	database := pg.Connect(&pg.Options{
		Addr:     Config.DB.Host,
		User:     Config.DB.User,
		Password: Config.DB.Password,
		Database: Config.DB.Table,
	})
	return database
}

// CreateSchema creates database tables using ORM
func CreateSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*dtypes.BlockInfo)(nil), (*dtypes.EvidenceInfo)(nil), (*dtypes.MissInfo)(nil),
		(*dtypes.MissDetailInfo)(nil), (*dtypes.ProposalInfo)(nil), (*dtypes.ValidatorSetInfo)(nil), (*dtypes.ValidatorInfo)(nil),
		(*dtypes.TransactionInfo)(nil), (*dtypes.VoteInfo)(nil), (*dtypes.DepositInfo)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
			Varchar:     5000, // replaces PostgreSQL data type `text` to `varchar(n)`
		})
		if err != nil {
			return err
		}
	}

	// RunInTransaction runs a function in a transaction.
	// If function returns an error transaction is rollbacked, otherwise transaction is committed.
	err := db.RunInTransaction(func(tx *pg.Tx) error {
		// Create indexes to reduce the cost of lookup queries in case of server traffic jams (B-Tree Index)
		_, err := db.Model(dtypes.BlockInfo{}).Exec(`CREATE INDEX block_info_height_idx ON block_infos USING btree(height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(dtypes.ValidatorInfo{}).Exec(`CREATE INDEX validator_info_rank_idx ON validator_infos USING btree(rank);`)
		if err != nil {
			return err
		}
		_, err = db.Model(dtypes.MissDetailInfo{}).Exec(`CREATE INDEX miss_detail_info_height_idx ON miss_detail_infos USING btree(height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(dtypes.MissInfo{}).Exec(`CREATE INDEX miss_info_start_height_idx ON miss_infos USING btree(start_height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(dtypes.TransactionInfo{}).Exec(`CREATE INDEX transaction_info_height_idx ON transaction_infos USING btree(height);`)
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
