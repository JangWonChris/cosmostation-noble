package db

import (
	"os"
	"testing"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"

	"github.com/stretchr/testify/require"
)

var db *Database

func TestMain(m *testing.M) {
	config := config.ParseConfig()
	db = Connect(&config.DB)

	os.Exit(m.Run())
}

func TestInsertOrUpdate(t *testing.T) {
	err := db.Ping()
	require.NoError(t, err)

}

func TestUpdate_Validator(t *testing.T) {
	err := db.Ping()
	require.NoError(t, err)

	val := schema.Validator{
		Address: "cosmos1clpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9q",
		Rank:    5,
	}

	validator, err := db.QueryValidator(val.Address)
	require.NoError(t, err)

	result, err := db.Model(&validator).
		Set("rank = ?", val.Rank).
		Where("id = ?", validator.ID).
		Update()

	require.NoError(t, err)
	require.Equal(t, 1, result.RowsAffected())
}

func TestQuery_LatestBlockHeight(t *testing.T) {
	err := db.Ping()
	require.NoError(t, err)

	height, err := db.QueryLatestBlockHeight()
	require.NoError(t, err)

	require.NotNil(t, height)
}
func TestCreate_Indexes(t *testing.T) {
	err := db.Ping()
	require.NoError(t, err)

	testIndex := "CREATE INDEX account_account_address_idx ON account USING btree(account_address);"

	_, err = db.Model(schema.Block{}).Exec(testIndex)
	require.NoError(t, err)
}

func TestCreate_Tables(t *testing.T) {
	err := db.Ping()
	require.NoError(t, err)

	tables := []interface{}{
		(*schema.Block)(nil),
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
