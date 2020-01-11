package db

import (
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/schema"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type Database struct {
	*pg.DB
}

// Connect connects to PostgreSQL database
func Connect(Config *config.Config) *Database {
	db := pg.Connect(&pg.Options{
		Addr:     Config.DB.Host,
		User:     Config.DB.User,
		Password: Config.DB.Password,
		Database: Config.DB.Table,
	})
	return &Database{db}
}

// CreateSchema creates database tables using ORM
func (db *Database) CreateSchema() error {
	for _, model := range []interface{}{(*schema.StatsCoingeckoMarket1H)(nil), (*schema.StatsCoingeckoMarket24H)(nil),
		(*schema.StatsCoinmarketcapMarket1H)(nil), (*schema.StatsCoinmarketcapMarket24H)(nil),
		(*schema.StatsValidators1H)(nil), (*schema.StatsValidators24H)(nil),
		(*schema.StatsNetwork1H)(nil), (*schema.StatsNetwork24H)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			return err
		}
	}
	return nil
}
