package db

import (
	"fmt"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/schema"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var (
	// columnLength is the column length of varchar type in every table.
	// This needs to be considered again to set it to what specific length is needed, but right now set it to 99999.
	columnLength = 99999
)

// Database implements a wrapper of golang ORM with focus on PostgreSQL.
type Database struct {
	*pg.DB
}

// Connect opens a database connections with the given database connection info from config.
func Connect(config *config.Database) *Database {
	db := pg.Connect(&pg.Options{
		Addr:     config.Host + ":" + config.Port,
		User:     config.User,
		Password: config.Password,
		Database: config.Table,
	})

	// Disable pluralization
	orm.SetTableNameInflector(func(s string) string {
		return s
	})

	return &Database{db}
}

// Ping returns a database connection handle or an error if the connection fails.
func (db *Database) Ping() error {
	_, err := db.Exec("SELECT 1")
	if err != nil {
		return err
	}

	return nil
}

// CreateTables creates database tables using object relational mapping (ORM)
func (db *Database) CreateTables() error {
	for _, model := range []interface{}{(*schema.StatsMarket5M)(nil), (*schema.StatsMarket1H)(nil), (*schema.StatsMarket1D)(nil),
		(*schema.StatsValidators1H)(nil), (*schema.StatsValidators1D)(nil),
		(*schema.StatsNetwork1H)(nil), (*schema.StatsNetwork1D)(nil)} {

		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
			Varchar:     columnLength, // replaces PostgreSQL data type `text` to `varchar(n)`
		})

		if err != nil {
			return err
		}
	}
	return nil
}

// --------------------
// Query
// --------------------

// QueryLatestTwoBlocks returns the latest two blocks for block time calculation.
func (db *Database) QueryLatestTwoBlocks() (blocks []schema.Block, err error) {
	err = db.Model(&blocks).
		Order("height DESC").
		Limit(2).
		Select()

	if err != nil {
		return []schema.Block{}, err
	}

	return blocks, nil
}

// QueryValidatorsByStatus returns a validator by querying with bonding status in an ascending rank order.
func (db *Database) QueryValidatorsByStatus(status int) (validators []schema.Validator, err error) {
	err = db.Model(&validators).
		Where("status = ?", status).
		Order("rank ASC").
		Select()

	if err == pg.ErrNoRows {
		return []schema.Validator{}, fmt.Errorf("found no rows in table: %s", err)
	}

	if err != nil {
		return []schema.Validator{}, err
	}

	return validators, nil
}

// QueryTotalTransactionNum returns a total number of transactions.
func (db *Database) QueryTotalTransactionNum() int {
	var tx schema.Transaction
	_ = db.Model(&tx).
		Order("id DESC").
		Limit(1).
		Select()

	return int(tx.ID)
}

// --------------------
// Insert
// --------------------

// InsertMarket5M inserts market data.
func (db *Database) InsertMarket5M(data *schema.StatsMarket5M) error {
	_, err := db.Model(data).Insert()
	if err != nil {
		return err
	}
	return nil
}

// InsertMarket1H inserts market data.
func (db *Database) InsertMarket1H(data *schema.StatsMarket1H) error {
	_, err := db.Model(data).Insert()
	if err != nil {
		return err
	}
	return nil
}

// InsertMarket1D inserts StatsMarket1D.
func (db *Database) InsertMarket1D(data *schema.StatsMarket1D) error {
	_, err := db.Model(data).Insert()
	if err != nil {
		return err
	}
	return nil
}

// InsertNetworkStats1H inserts StatsNetwork1H.
func (db *Database) InsertNetworkStats1H(data *schema.StatsNetwork1H) error {
	_, err := db.Model(data).Insert()
	if err != nil {
		return err
	}
	return nil
}

// InsertNetworkStats1D inserts StatsNetwork1D.
func (db *Database) InsertNetworkStats1D(data *schema.StatsNetwork1D) error {
	_, err := db.Model(data).Insert()
	if err != nil {
		return err
	}
	return nil
}

// InsertValidatorStats1H inserts StatsValidators1H.
func (db *Database) InsertValidatorStats1H(data []schema.StatsValidators1H) error {
	_, err := db.Model(&data).Insert()
	if err != nil {
		return nil
	}
	return nil
}

// InsertValidatorStats1D inserts StatsValidators1D.
func (db *Database) InsertValidatorStats1D(data []schema.StatsValidators1D) error {
	_, err := db.Model(&data).Insert()
	if err != nil {
		return err
	}
	return nil
}
