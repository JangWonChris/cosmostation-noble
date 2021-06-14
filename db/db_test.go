package db

import (
	"os"
	"testing"

	//mbl
	mblconfig "github.com/cosmostation/mintscan-backend-library/config"
	mdschema "github.com/cosmostation/mintscan-database/schema"

	pg "github.com/go-pg/pg/v10"

	"github.com/stretchr/testify/require"
)

var db *Database

func TestMain(m *testing.M) {
	// types.SetAppConfig()

	fileBaseName := "chain-exporter"
	cfg := mblconfig.ParseConfig(fileBaseName)
	db = Connect(&cfg.DB)

	os.Exit(m.Run())
}

func TestInsertOrUpdate(t *testing.T) {
	err := db.Ping()
	require.NoError(t, err)

}

func TestUpdate_Validator(t *testing.T) {
	err := db.Ping()
	require.NoError(t, err)

	val := &mdschema.Validator{
		Address: "kava1ulzzxuvghlv04sglkzyxv94rvl7c2llhs098ju",
		Rank:    5,
	}

	validator, err := db.GetValidatorByAnyAddr(val.Address)
	require.NoError(t, err)

	result, err := db.Model(&validator).
		Set("rank = ?", val.Rank).
		Where("id = ?", validator.ID).
		Update()

	require.NoError(t, err)
	require.Equal(t, 1, result.RowsAffected())
}

func TestConnection(t *testing.T) {
	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	require.NoError(t, err)

	require.Equal(t, n, 1, "failed to ping database")
}
