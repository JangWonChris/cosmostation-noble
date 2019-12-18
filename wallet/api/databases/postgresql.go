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
func CreateSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*models.Account)(nil), (*models.AppVersion)(nil)} {
		// disable pluralization
		orm.SetTableNameInflector(func(s string) string {
			return s
		})

		// create tables
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
			Varchar:     999, // replaces PostgreSQL data type `text` with `varchar(n)`
		})

		if err != nil {
			panic(err)
		}
	}

	// runs a function in a transaction.
	// if function returns an error transaction is rollbacked, otherwise transaction is committed.
	// create indexes to reduce the cost of lookup queries in case of server traffic jams (B-Tree Index)
	err := db.RunInTransaction(func(tx *pg.Tx) error {
		_, err := db.Model(models.Account{}).Exec(`CREATE INDEX account_address_idx ON account USING btree(address);`)
		if err != nil {
			return err
		}
		_, err = db.Model(models.Account{}).Exec(`CREATE INDEX account_alarm_token_idx ON account USING btree(alarm_token);`)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
