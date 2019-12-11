package databases

import (
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/types"

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
	for _, model := range []interface{}{(*types.StatsCoingeckoMarket1H)(nil), (*types.StatsCoingeckoMarket24H)(nil),
		(*types.StatsCoinmarketcapMarket1H)(nil), (*types.StatsCoinmarketcapMarket24H)(nil),
		(*types.StatsValidators1H)(nil), (*types.StatsValidators24H)(nil),
		(*types.StatsNetwork1H)(nil), (*types.StatsNetwork24H)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			return err
		}
	}
	return nil
}
