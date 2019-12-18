package databases

import (
	"github.com/cosmostation/cosmostation-cosmos/wallet/api/config"
	"github.com/cosmostation/cosmostation-cosmos/wallet/api/models"

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
	for _, model := range []interface{}{(*models.Account)(nil), (*models.AppVersion)(nil)} {
		// disable pluralization
		orm.SetTableNameInflector(func(s string) string {
			return s
		})

		// create tables
		err := DB.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
			Varchar:     999, // replaces PostgreSQL data type `text` with `varchar(n)`
		})
		if err != nil {
			panic(err)
		}
	}
	return nil
}
