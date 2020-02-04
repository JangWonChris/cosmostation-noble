package db

import (
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/schema"
	"github.com/go-pg/pg"
)

// QueryBlocks queries blocks with height and limit
func (db *Database) QueryBlocks(height int, limit int) ([]schema.BlockInfo, error) {
	var blocks []schema.BlockInfo
	_ = db.Model(&blocks).
		Where("height > ?", height).
		Limit(limit).
		Order("id ASC").
		Select()

	return blocks, nil
}

// QueryLatestBlockHeight queries the latest block height in database
func (db *Database) QueryLatestBlockHeight() (int, error) {
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

	return int(block.Height), nil
}

// QueryBlocksByProposer queries blocks by proposer
func (db *Database) QueryBlocksByProposer(address string, limit int, before int, after int, offset int) ([]schema.BlockInfo, error) {
	blocks := make([]schema.BlockInfo, 0)

	switch {
	case before > 0:
		_ = db.Model(&blocks).
			Where("proposer = ? AND height < ?", address, before).
			Limit(limit).
			Order("height DESC").
			Select()
	case after >= 0:
		_ = db.Model(&blocks).
			Where("proposer = ? AND height < ?", address, after).
			Limit(limit).
			Order("height ASC").
			Select()
	case offset >= 0:
		_ = db.Model(&blocks).
			Where("proposer = ?", address).
			Limit(limit).
			Offset(offset).
			Order("height DESC").
			Select()
	}

	return blocks, nil
}

// QueryTotalBlocksByProposer queries total number of blocks proposed by a proposer
func (db *Database) QueryTotalBlocksByProposer(address string) (int, error) {
	var blockInfo schema.BlockInfo
	totalNum, _ := db.Model(&blockInfo).
		Where("proposer = ?", address).
		Count()

	return totalNum, nil
}

// QueryLastestTwoBlocks queries lastest two blocks
func (db *Database) QueryLastestTwoBlocks() []schema.BlockInfo {
	var blocks []schema.BlockInfo
	_ = db.Model(&blocks).
		Order("height DESC").
		Limit(2).
		Select()

	return blocks
}

// QueryMissingBlocksInDetail queries how many missing blocks a validator misses in detail
func (db *Database) QueryMissingBlocksInDetail(address string, latestHeight int, count int) ([]schema.MissDetailInfo, error) {
	var missDetailInfo []schema.MissDetailInfo
	_ = db.Model(&missDetailInfo).
		Where("address = ? AND height BETWEEN ? AND ?", address, latestHeight-count, latestHeight).
		Limit(count).
		Order("height DESC").
		Select()

	return missDetailInfo, nil
}

// QueryMissingBlocks queries a range of missing blocks a validator misses
func (db *Database) QueryMissingBlocks(address string, limit int) ([]schema.MissInfo, error) {
	var missInfo []schema.MissInfo
	_ = db.Model(&missInfo).
		Where("address = ?", address).
		Limit(limit).
		Order("start_height DESC").
		Select()

	return missInfo, nil
}

// QueryMissingBlocksCount queries how many missing blocks a validator misses in detail and return total number
func (db *Database) QueryMissingBlocksCount(address string, latestHeight int, count int) (int, error) {
	var missDetailInfo []schema.MissDetailInfo
	missingBlocksCount, _ := db.Model(&missDetailInfo).
		Where("address = ? AND height BETWEEN ? AND ?", address, latestHeight-count, latestHeight).
		Count()

	return missingBlocksCount, nil
}
