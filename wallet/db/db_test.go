package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmostation/cosmostation-cosmos/wallet/config"
	"github.com/cosmostation/cosmostation-cosmos/wallet/schema"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var db *Database

func TestMain(m *testing.M) {
	config := config.ParseConfig()
	db = Connect(config.DB)

	os.Exit(m.Run())
}

func TestCreate_Indexes(t *testing.T) {
	err := db.Ping()
	require.NoError(t, err)

	testIndex := "CREATE INDEX account_account_address_idx ON account USING btree(account_address);"

	_, err = db.Model(schema.AppAccount{}).Exec(testIndex)
	require.NoError(t, err)
}

func TestCreate_Tables(t *testing.T) {
	err := db.Ping()
	require.NoError(t, err)

	tables := []interface{}{
		(*schema.AppAccount)(nil),
		(*schema.AppVersion)(nil),
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
