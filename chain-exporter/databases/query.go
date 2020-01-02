package databases

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
)

// QueryValidatorInfo returns validator information
func (db *Database) QueryValidatorInfo(address string) (schema.ValidatorInfo, error) {
	var validatorInfo schema.ValidatorInfo
	switch {
	case strings.HasPrefix(address, sdk.GetConfig().GetBech32ConsensusPubPrefix()):
		err := db.Model(&validatorInfo).
			Where("address = ?", address).
			Limit(1).
			Select()
		if err != nil {
			return validatorInfo, err
		}
	case strings.HasPrefix(address, sdk.GetConfig().GetBech32ValidatorAddrPrefix()):
		err := db.Model(&validatorInfo).
			Where("operator_address = ?", address).
			Limit(1).
			Select()
		if err != nil {
			return validatorInfo, err
		}
	case strings.HasPrefix(address, sdk.GetConfig().GetBech32AccountAddrPrefix()):
		err := db.Model(&validatorInfo).
			Where("consensus_pubkey = ?", address).
			Limit(1).
			Select()
		if err != nil {
			return validatorInfo, err
		}
	}
	return validatorInfo, nil
}

// QueryIDValidatorSetInfo returns id of a validator from validator_set_infos table
func (db *Database) QueryIDValidatorSetInfo(proposer string) (schema.ValidatorSetInfo, error) {
	var validatorSetInfo schema.ValidatorSetInfo
	err := db.Model(&validatorSetInfo).
		Column("id_validator", "voting_power").
		Where("proposer = ?", proposer).
		Order("id DESC"). // Lastly input data
		Limit(1).
		Select()
	if err != nil {
		return validatorSetInfo, err
	}
	return validatorSetInfo, nil
}

// QueryHighestIDValidatorNum returns highest id of a validator from validator_set_infos table
func (db *Database) QueryHighestIDValidatorNum() (int, error) {
	var validatorSetInfo schema.ValidatorSetInfo
	err := db.Model(&validatorSetInfo).
		Column("id_validator").
		Order("id_validator DESC").
		Limit(1).
		Select()
	if err != nil {
		return 0, err
	}
	return validatorSetInfo.IDValidator, nil
}

// QueryAccount queries account information
func (db *Database) QueryAccount(address string) (types.Account, error) {
	var account types.Account
	_ = db.Model(&account).
		Where("address = ?", address).
		Select()

	return account, nil
}
