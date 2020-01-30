package db

import (
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"

	"github.com/go-pg/pg"
)

// Database implements a wrapper of golang ORM with focus on PostgreSQL
type Database struct {
	*pg.DB
}

// Connect opens a database connections with the given database connection info from config.
func Connect(cfg *config.Config) *Database {
	db := pg.Connect(&pg.Options{
		Addr:     cfg.DB.Host + ":" + cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Database: cfg.DB.Table,
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
