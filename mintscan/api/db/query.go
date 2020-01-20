package db

import (
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/schema"
	"github.com/go-pg/pg"
)

// QueryLatestBlockHeight queries the latest block height in database
func (db *Database) QueryLatestBlockHeight() (int64, error) {
	var block schema.BlockInfo
	err := db.Model(&block).
		Order("height DESC").
		Limit(1).
		Select()

	// return 0 when there is no row in result set
	if err == pg.ErrNoRows {
		return 0, err
	}

	// return -1 for any type of errors
	if err != nil {
		return -1, err
	}

	return block.Height, nil
}
