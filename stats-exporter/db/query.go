package db

import (
	"fmt"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/schema"
)

// QueryLatestBlocks queries the latest number of block information
func (db *Database) QueryLatestBlocks(num int) ([]*schema.BlockInfo, error) {
	var blockInfo []*schema.BlockInfo
	err := db.Model(&blockInfo).
		Column("time").
		Order("height DESC").
		Limit(num).
		Select()
	if err != nil {
		fmt.Printf("failed to query block time: %v \n", err)
		return blockInfo, nil
	}

	return blockInfo, nil
}
