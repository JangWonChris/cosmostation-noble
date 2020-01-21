package db

import (
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/schema"
	"github.com/go-pg/pg"
)

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
func (db *Database) QueryValidatorPowerEvents(validatorID int, limit int, before int, after int, offset int) ([]schema.ValidatorSetInfo, error) {
	validatorSetInfo := make([]schema.ValidatorSetInfo, 0)

	switch {
	case before > 0:
		_ = db.Model(&validatorSetInfo).
			Where("id_validator = ? AND height < ?", validatorID, before).
			Limit(limit).
			Order("id DESC").
			Select()
	case after > 0:
		_ = db.Model(&validatorSetInfo).
			Where("id_validator = ? AND height > ? ", validatorID, after).
			Limit(limit).
			Order("id ASC").
			Select()

	case offset >= 0:
		_ = db.Model(&validatorSetInfo).
			Where("id_validator = ?", validatorID).
			Limit(limit).
			Offset(offset).
			Order("id DESC").
			Select()
	}

	return validatorSetInfo, nil
}

// QueryUnjailedValidatorsNum queries how many validators are not jailed
func (db *Database) QueryUnjailedValidatorsNum() int {
	var validatorInfo schema.ValidatorInfo
	num, _ := db.Model(&validatorInfo).
		Where("status = ?", 2).
		Count()

	return num
}

// QueryJailedValidatorsNum queries how many validators are not either unbonding or unbonded
func (db *Database) QueryJailedValidatorsNum() int {
	var validatorInfo schema.ValidatorInfo
	num, _ := db.Model(&validatorInfo).
		Where("status = ? OR status = ?", 0, 1).
		Count()

	return num
}
