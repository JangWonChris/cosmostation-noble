package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmostation/mintscan-backend-library/config"
	"github.com/cosmostation/mintscan-backend-library/db/schema"

	"github.com/go-pg/pg"
)

var db *Database

func TestMain(m *testing.M) {
	fileBaseName := "mintscan"
	config := config.ParseConfig(fileBaseName)
	db = Connect(&config.DB)

	os.Exit(m.Run())
}

func TestQueryOptions(t *testing.T) {
	var result []struct {
		Option string
		Count  int
	}

	// select option, count(option) from vote where proposal_id = 29 group by option
	err := db.Model(&schema.Vote{}).
		Column("option").
		ColumnExpr("COUNT('option') AS count").
		Where("proposal_id = ?", 29).
		Group("option").
		Select(&result)

	require.NoError(t, err)

	require.NotNil(t, result)
}

func TestConnection(t *testing.T) {
	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	require.NoError(t, err)

	require.Equal(t, n, 1, "failed to ping database")
}

func TestTimeDiffLatestTwoBlocks(t *testing.T) {
	diff, err := db.QueryTimeDiffLastestTwoBlocks()
	require.NoError(t, err)
	t.Log("time diff : ", diff)
}
