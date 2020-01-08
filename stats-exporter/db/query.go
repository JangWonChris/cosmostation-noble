package db

import (
	"fmt"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/schema"
)

// QueryLatestBlocks queries the latest number of block information
func (db *Database) QueryLatestBlocks(num int) ([]*schema.BlockInfo, error) {
	var blocks []*schema.BlockInfo
	err := db.Model(&blocks).
		Column("time").
		Order("height DESC").
		Limit(num).
		Select()
	if err != nil {
		fmt.Printf("failed to query block time: %v \n", err)
		return blocks, nil
	}

	return blocks, nil
}

// QueryValidatorsByRank queries 125 validators order by their rank in an ascending way
func (db *Database) QueryValidatorsByRank(num int) ([]*schema.ValidatorInfo, error) {
	var validators []*schema.ValidatorInfo
	err := db.Model(&validators).
		Order("rank ASC").
		Limit(num).
		Select()
	if err != nil {
		fmt.Printf("failed to query validators: %v \n", err)
		return validators, nil
	}
	return validators, nil
}
