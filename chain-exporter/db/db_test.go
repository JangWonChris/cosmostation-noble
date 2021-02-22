package db

import (
	"os"
	"testing"

	//mbl
	"github.com/cosmostation/mintscan-backend-library/config"
	"github.com/cosmostation/mintscan-backend-library/db/schema"

	"github.com/go-pg/pg"

	"github.com/stretchr/testify/require"
)

var db *Database

func TestMain(m *testing.M) {
	// types.SetAppConfig()

	fileBaseName := "chain-exporter"
	cfg := config.ParseConfig(fileBaseName)
	db = Connect(&cfg.DB)

	os.Exit(m.Run())
}

func TestInsertOrUpdate(t *testing.T) {
	err := db.Ping()
	require.NoError(t, err)

}

func TestExistAccount(t *testing.T) {
	err := db.Ping()
	require.NoError(t, err)

	address := "kava140g8fnnl46mlvfhygj3zvjqlku6x0fwuhfj3uf"

	exist, err := db.ExistAccount(address)
	require.NoError(t, err)

	require.Equal(t, true, exist)
}

func TestUpdate_Validator(t *testing.T) {
	err := db.Ping()
	require.NoError(t, err)

	val := &schema.Validator{
		Address: "kava1ulzzxuvghlv04sglkzyxv94rvl7c2llhs098ju",
		Rank:    5,
	}

	validator, err := db.QueryValidatorByAnyAddr(val.Address)
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

func TestQuery_Account(t *testing.T) {
	err := db.Ping()
	require.NoError(t, err)

	acct := &schema.AccountCoin{AccountAddress: "kava1m36xddywe0yneykv34az8smzhtxy3nyc6v9jdj"}

	account, err := db.QueryAccount(acct.AccountAddress)
	require.NoError(t, err)

	require.NotNil(t, account)
}

func TestConnection(t *testing.T) {
	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	require.NoError(t, err)

	require.Equal(t, n, 1, "failed to ping database")
}
