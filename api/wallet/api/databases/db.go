package databases

import (
	"github.com/cosmostation/cosmostation-cosmos/api/wallet/api/models"
	"github.com/cosmostation/cosmostation-cosmos/api/wallet/config"

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

// Create tables if it doesn't already exist
func CreateSchema(DB *pg.DB) error {
	for _, model := range []interface{}{(*models.Account)(nil), (*models.Version)(nil)} {
		err := DB.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			panic(err)
		}
	}
	return nil
}
