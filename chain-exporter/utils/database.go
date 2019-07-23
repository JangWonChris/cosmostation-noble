package utils

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dtypes "github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	"github.com/go-pg/pg"
)

// var (
// 	Bech32AccountAddrColumnName               = "address"          // cosmos
// 	Bech32ValidatorAddrColumnName             = "operator_address" // cosmosvaloper
// 	Bech32Bech32ConsensusPubColumnName        = "consensus_pubkey" // cosmosvalconspub
// 	Bech32Bech32ConsensusPubHexAddrColumnName = "proposer"         // cosmosvalconspub's hex addr format
// )

func QueryValidatorInfo(db *pg.DB, address string) (dtypes.ValidatorInfo, error) {
	var validatorInfo dtypes.ValidatorInfo
	var err error
	switch {
	case strings.HasPrefix(address, sdk.GetConfig().GetBech32ConsensusPubPrefix()):
		err = db.Model(&validatorInfo).
			Where("address = ?", address).
			Limit(1).
			Select()
	case strings.HasPrefix(address, sdk.GetConfig().GetBech32ValidatorAddrPrefix()):
		err = db.Model(&validatorInfo).
			Where("operator_address = ?", address).
			Limit(1).
			Select()
	case strings.HasPrefix(address, sdk.GetConfig().GetBech32AccountAddrPrefix()):
		err = db.Model(&validatorInfo).
			Where("consensus_pubkey = ?", address).
			Limit(1).
			Select()
	}

	if err != nil {
		return validatorInfo, err
	}

	return validatorInfo, nil
}

func QueryIDValidatorSetInfo(db *pg.DB, proposer string) (dtypes.ValidatorSetInfo, error) {
	var validatorSetInfo dtypes.ValidatorSetInfo
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

func QueryHighestIDValidatorNum(db *pg.DB) (int, error) {
	var validatorSetInfo dtypes.ValidatorSetInfo
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
