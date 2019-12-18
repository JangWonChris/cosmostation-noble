package databases

import (
	"github.com/cosmostation/cosmostation-cosmos/alarm-notification/config"
	dtypes "github.com/cosmostation/cosmostation-cosmos/alarm-notification/types"

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
	for _, model := range []interface{}{(*dtypes.BlockInfo)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			return err
		}
	}
	return nil
}
