package db

import (
	"os"
	"testing"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/models"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/schema"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"

	"github.com/stretchr/testify/require"
)

var db *Database

func TestMain(m *testing.M) {
	models.SetAppConfig()

	config := config.ParseConfig()
	db = Connect(&config.DB)

	os.Exit(m.Run())
}

func TestCreate_Tables(t *testing.T) {
	err := db.Ping()
	require.NoError(t, err)

	tables := []interface{}{
		(*schema.StatsMarket5M)(nil),
		(*schema.StatsMarket1H)(nil),
		(*schema.StatsMarket1D)(nil),
		(*schema.StatsNetwork1H)(nil),
		(*schema.StatsNetwork1D)(nil),
		(*schema.StatsValidators1H)(nil),
		(*schema.StatsValidators1D)(nil),
	}

	for _, table := range tables {
		orm.SetTableNameInflector(func(s string) string {
			return s
		})

		err := db.CreateTable(table, &orm.CreateTableOptions{
			IfNotExists: true,
			Varchar:     columnLength,
		})

		require.NoError(t, err)
	}
}

func TestConnection(t *testing.T) {
	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	require.NoError(t, err)

	require.Equal(t, n, 1, "failed to ping database")
}
