package databases

import (
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"

	"github.com/go-pg/pg"
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
