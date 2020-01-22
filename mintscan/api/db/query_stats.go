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
