package db

import (
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/models"
)

// QueryValidatorsStats24H queries 24 hours of validator's stats
func (db *Database) QueryValidatorStats24H(address string, limit int) ([]models.StatsValidators24H, error) {
	statsValidators24H := make([]models.StatsValidators24H, 0)
	_ = db.Model(&statsValidators24H).
		Where("proposer = ?", address).
		Order("id DESC").
		Limit(limit).
		Select()

	return statsValidators24H, nil
}

// QueryOneDayPrices queries every price of an hour for 24 hours
func (db *Database) QueryOneDayPrices(limit int) ([]models.StatsCoingeckoMarket1H, error) {
	var prices []models.StatsCoingeckoMarket1H
	_ = db.Model(&prices).
		Order("id DESC").
		Limit(limit).
		Select()

	return prices, nil
}

// QueryNetworkStats queries network stats
func (db *Database) QueryNetworkStats(limit int) ([]models.StatsNetwork1H, error) {
	var networkStats []models.StatsNetwork1H
	err := db.Model(&networkStats).
		Order("id DESC").
		Limit(limit).
		Select()

	if err != nil {
		return networkStats, err
	}

	return networkStats, nil
}

// QueryBondedRateIn24H queries bonded tokens percentage change for 24 hrs
func (db *Database) QueryBondedRateIn24H() ([]models.StatsNetwork24H, error) {
	var networkStats []models.StatsNetwork24H
	err := db.Model(&networkStats).
		Order("id DESC").
		Limit(2).
		Select()

	if err != nil {
		return networkStats, err
	}

	return networkStats, nil
}
