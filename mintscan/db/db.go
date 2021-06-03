package db

import (
	"fmt"
	"time"

	mblconfig "github.com/cosmostation/mintscan-backend-library/config"
	mddb "github.com/cosmostation/mintscan-database/db"
	mdschema "github.com/cosmostation/mintscan-database/schema"
	pg "github.com/go-pg/pg/v10"
)

// Database implements a wrapper of golang ORM with focus on PostgreSQL.
type Database struct {
	*mddb.Database
}

// Connect opens a database connections with the given database connection info from config.
func Connect(dbcfg *mblconfig.DatabaseConfig) *Database {
	db := mddb.Connect(dbcfg.Host, dbcfg.Port, dbcfg.User, dbcfg.Password, dbcfg.DBName, dbcfg.Schema, dbcfg.Timeout)
	mdschema.SetCommonSchema("refine")

	return &Database{db}
}

// QueryLatestTwoBlocks() 대체용 함수 - 응답 구조체 파싱 오류
/*
	Error:      	Received unexpected error:
	            	pg: Model(unsupported *time.Time)
*/
func (db *Database) QueryBlockTimeDiff() (string, error) {
	var BlockTimeDiff time.Time
	query := fmt.Sprintf("select timestamp-lead(timestamp) over (order by timestamp desc) as block_time_diff from (select timestamp from %s.block order by id desc limit 2) a limit 1", mdschema.GetCommonSchema())
	_, err := db.Query(&BlockTimeDiff, query)

	if err != nil {
		if err == pg.ErrNoRows {
			return BlockTimeDiff.String(), nil
		}
		return BlockTimeDiff.String(), err
	}

	return BlockTimeDiff.String(), nil
}

// QueryBondedRateIn1D return bonded rate in network from 1 day network stats table.
// func (db *Database) QueryBondedRateIn1D() ([]mdschema.StatsNetwork1D, error) {
// 	var networkStats []mdschema.StatsNetwork1D
// 	err := db.Model(&networkStats).
// 		Order("id DESC").
// 		Limit(2).
// 		Select()

// 	if err != nil {
// 		return networkStats, err
// 	}

// 	return networkStats, nil
// }

// QueryValidatorBondedInfo returns a validator's bonded information.
// sdk 에서 제공하는 IsBonded 함수가 존재한다.
// 이 함수가 필요한 이유는, 최초 본딩 된 날짜를 알기 위함임(제네시스인지, 그 이후 생성 된 검증인 인지)
// func (db *Database) QueryValidatorBondedInfo(address string) (peh schema.PowerEventHistory, err error) {
// 	msgType := "create_validator"

// 	err = db.Model(&peh).
// 		Where("proposer = ? AND msg_type = ?", address, msgType).
// 		Limit(1).
// 		Select()

// 	if err != nil {
// 		return schema.PowerEventHistory{}, err
// 	}

// 	return peh, nil
// }

// QueryValidatorVotingPowerEventHistory returns a validator's voting power events
// func (db *Database) QueryValidatorVotingPowerEventHistory(address string, before, after, limit int) ([]schema.PowerEventHistory, error) {
// 	var peh []schema.PowerEventHistory
// 	var err error

// 	switch {
// 	case before > 0:
// 		err = db.Model(&peh).
// 			Where("operator_address = ? AND height < ?", address, before).
// 			Limit(limit).
// 			Order("id DESC").
// 			Select()
// 	case after > 0:
// 		err = db.Model(&peh).
// 			Where("operator_address = ? AND height > ?", address, after).
// 			Limit(limit).
// 			Order("id ASC").
// 			Select()
// 	default:
// 		err = db.Model(&peh).
// 			Where("operator_address = ?", address).
// 			Limit(limit).
// 			Order("id DESC").
// 			Select()
// 	}

// 	if err != nil {
// 		if err == pg.ErrNoRows {
// 			return []schema.PowerEventHistory{}, nil
// 		}
// 		return []schema.PowerEventHistory{}, err
// 	}

// 	return peh, nil
// }
