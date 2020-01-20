package db

import (
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/schema"
	"github.com/go-pg/pg"
)

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

// QueryValidatorID queries validator index id
// reference chain exporter project on how validator power event data is stored
func (db *Database) QueryValidatorID(address string) (int, error) {
	var validatorSetInfo schema.ValidatorSetInfo
	err := db.Model(&validatorSetInfo).
		Column("id_validator").
		Where("proposer = ?", address).
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

	return validatorSetInfo.IDValidator, nil
}

// QueryValidatorInfoByProposer queries validator information by proposer address format
func (db *Database) QueryValidatorInfoByProposer(proposer string) (schema.ValidatorInfo, error) {
	var validatorInfo schema.ValidatorInfo
	_ = db.Model(&validatorInfo).
		Where("proposer = ?", proposer).
		Select()

	return validatorInfo, nil
}

// QueryValidatorPowerEvents queries validator's power events by limit/offset pagination
func (db *Database) QueryValidatorPowerEvents(validatorID int, limit int, offset int) ([]schema.ValidatorSetInfo, error) {
	validatorSetInfo := make([]schema.ValidatorSetInfo, 0)
	_ = db.Model(&validatorSetInfo).
		Where("id_validator = ?", validatorID).
		Limit(limit).
		Offset(offset).
		Order("id DESC").
		Select()

	return validatorSetInfo, nil
}

// QueryBlocksByProposer queries blocks by proposer
func (db *Database) QueryBlocksByProposer(address string, limit int, offset int) ([]schema.BlockInfo, error) {
	blocks := make([]schema.BlockInfo, 0)
	_ = db.Model(&blocks).
		Where("proposer = ?", address).
		Limit(limit).
		Offset(offset).
		Order("height DESC").
		Select()

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

// QueryTransactions queries transactions
func (db *Database) QueryTransactions(height int64) ([]schema.TransactionInfo, error) {
	var txInfos []schema.TransactionInfo
	_ = db.Model(&txInfos).
		Column("tx_hash").
		Where("height = ?", height).
		Select()

	return txInfos, nil
}
