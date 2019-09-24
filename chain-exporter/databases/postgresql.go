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
		err := db.CreateTable(model, &orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			return err
		}
	}
	return nil
}
