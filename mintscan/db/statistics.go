package db

import (
	"github.com/cosmostation/cosmostation-cosmos/mintscan/schema"
)

// QueryValidatorStats1D queries 24 hours of validator's stats
func (db *Database) QueryValidatorStats1D(address string, limit int) ([]schema.StatsValidators1D, error) {
	statsValidators24H := make([]schema.StatsValidators1D, 0)
	_ = db.Model(&statsValidators24H).
		Where("proposer = ?", address).
		Order("id DESC").
		Limit(limit).
		Select()

	return statsValidators24H, nil
}

// QueryPrices1D returns market statistics from 1 hour makret stats table.
func (db *Database) QueryPrices1D(limit int) (stats []schema.StatsMarket1H, err error) {
	err = db.Model(&stats).
		Order("id DESC").
		Limit(limit).
		Select()

	if err != nil {
		return []schema.StatsMarket1H{}, err
	}

	return stats, nil
}

// QueryNetworkStats queries network stats
func (db *Database) QueryNetworkStats(limit int) ([]schema.StatsNetwork1H, error) {
	var networkStats []schema.StatsNetwork1H
	err := db.Model(&networkStats).
		Order("id DESC").
		Limit(limit).
		Select()

	if err != nil {
		return networkStats, err
	}

	return networkStats, nil
}

// QueryBondedRateIn1D queries bonded tokens percentage change for 24 hrs
func (db *Database) QueryBondedRateIn1D() ([]schema.StatsNetwork1D, error) {
	var networkStats []schema.StatsNetwork1D
	err := db.Model(&networkStats).
		Order("id DESC").
		Limit(2).
		Select()

	if err != nil {
		return networkStats, err
	}

	return networkStats, nil
}
