package databases

import (
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/config"
	dtypes "github.com/cosmostation/cosmostation-cosmos/stats-exporter/types"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Connect to PostgreSQL
func ConnectDatabase(Config *config.Config) *pg.DB {
	database := pg.Connect(&pg.Options{
		Addr:     Config.DB.Host,
		User:     Config.DB.User,
		Password: Config.DB.Password,
		Database: Config.DB.Table,
	})
	return database
}

// CreateSchema sets up the database using the ORM
func CreateSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*dtypes.CoingeckoMarketStats)(nil), (*dtypes.CoinmarketcapMarketStats)(nil),
		(*dtypes.ValidatorStats)(nil), (*dtypes.NetworkStats)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			return err
		}
	}
	return nil
}
