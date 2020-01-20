package db

import (
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"

	"github.com/go-pg/pg"
)

type Database struct {
	*pg.DB
}

// Connect connects to PostgreSQL
func Connect(Config *config.Config) *Database {
	db := pg.Connect(&pg.Options{
		Addr:     Config.DB.Host,
		User:     Config.DB.User,
		Password: Config.DB.Password,
		Database: Config.DB.Table,
	})
	return &Database{db}
}
