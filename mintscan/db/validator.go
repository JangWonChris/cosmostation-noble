package db

import (
	"strings"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	mdschema "github.com/cosmostation/mintscan-database/schema"
	pg "github.com/go-pg/pg/v10"
)

// QueryValidatorByAnyAddr returns a validator information by any type of address format
func (db *Database) QueryValidatorByAnyAddr(anyAddr string) (mdschema.Validator, error) {
	var val mdschema.Validator
	var err error

	switch {
	// jeonghwan
	case strings.HasPrefix(anyAddr, sdktypes.GetConfig().GetBech32ConsensusPubPrefix()): // Bech32 prefix for validator public key
		err = db.Model(&val).
			Where("consensus_pubkey = ?", anyAddr).
			Limit(1).
			Select()
	case strings.HasPrefix(anyAddr, sdktypes.GetConfig().GetBech32ValidatorAddrPrefix()): // Bech32 prefix for validator address
		err = db.Model(&val).
			Where("operator_address = ?", anyAddr).
			Limit(1).
			Select()
	case strings.HasPrefix(anyAddr, sdktypes.GetConfig().GetBech32AccountAddrPrefix()): // Bech32 prefix for account address
		err = db.Model(&val).
			Where("address = ?", anyAddr).
			Limit(1).
			Select()
	case len(anyAddr) == 40: // Validator consensus address in hex
		anyAddr := strings.ToUpper(anyAddr)
		err = db.Model(&val).
			Where("proposer = ?", anyAddr).
			Limit(1).
			Select()
	default:
		err = db.Model(&val).
			Where("moniker = ?", anyAddr). // Validator moniker
			Limit(1).
			Select()
	}

	if err != nil {
		if err == pg.ErrNoRows {
			return mdschema.Validator{}, nil
		}
		return mdschema.Validator{}, err
	}

	return val, nil
}

// QueryValidatorStats1D returns validator statistics from 1 day of validator stats table
func (db *Database) QueryValidatorStats1D(proposerHexStr string, limit int) ([]mdschema.StatsValidators1D, error) {
	var stats []mdschema.StatsValidators1D
	err := db.Model(&stats).
		Where("proposer = ?", proposerHexStr).
		Order("id DESC").
		Limit(limit).
		Select()

	if err != nil {
		if err == pg.ErrNoRows {
			return []mdschema.StatsValidators1D{}, nil
		}
		return []mdschema.StatsValidators1D{}, err
	}

	return stats, nil
}

// QueryValidatorBondedInfo returns a validator's bonded information.
// sdk 에서 제공하는 IsBonded 함수가 존재한다.
// 이 함수가 필요한 이유는, 최초 본딩 된 날짜를 알기 위함임(제네시스인지, 그 이후 생성 된 검증인 인지)
func (db *Database) QueryValidatorBondedInfo(address string) (peh mdschema.PowerEventHistory, err error) {
	msgType := "create_validator"

	err = db.Model(&peh).
		Where("proposer = ? AND msg_type = ?", address, msgType).
		Limit(1).
		Select()

	if err != nil {
		return mdschema.PowerEventHistory{}, err
	}

	return peh, nil
}
