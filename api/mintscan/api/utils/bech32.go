package utils

import (
	"fmt"
	"strings"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/stats"
	dbtypes "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-pg/pg"
	"github.com/tendermint/tendermint/libs/bech32"
)

// Conver from operator address to cosmos address
func OperatorAddressToCosmosAddress(operatorAddress string) (string, error) {
	_, decoded, err := bech32.DecodeAndConvert(operatorAddress)
	if err != nil {
		return "", err
	}

	cosmosAddress, err := bech32.ConvertAndEncode(sdk.Bech32PrefixAccAddr, decoded)
	if err != nil {
		return "", err
	}
	return cosmosAddress, nil
}

func ConsensusPubkeyToProposer(consensusPubKey string) (string, error) {
	pk, err := sdk.GetConsPubKeyBech32(consensusPubKey)
	if err != nil {
		return "", err
	}
	hexAddr := pk.Address().String()
	return hexAddr, nil
}

// ConvertCosmosAddressToMoniker() converts from cosmos address to moniker
func ConvertCosmosAddressToMoniker(cosmosAddr string, db *pg.DB) (string, error) {
	// First convert from Cosmos Address to ValiOperatorAddress
	_, decoded, err := bech32.DecodeAndConvert(cosmosAddr)
	if err != nil {
		return "", err
	}
	cosmosOperAddress, err := bech32.ConvertAndEncode(sdk.Bech32PrefixValAddr, decoded)
	if err != nil {
		return "", err
	}

	// Check if the address matches any moniker in our DB
	var validatorInfo dbtypes.ValidatorInfo
	err = db.Model(&validatorInfo).
		Column("moniker").
		Where("operator_address = ?", cosmosOperAddress).
		Limit(1).
		Select()
	if err != nil {
		return "", err
	}

	return validatorInfo.Moniker, nil
}

// ConvertToProposer() converts any type of input address to proposer address
func ConvertToProposer(address string, db *pg.DB) (dbtypes.ValidatorInfo, error) {
	var validatorInfo dbtypes.ValidatorInfo
	switch {
	case strings.HasPrefix(address, sdk.Bech32PrefixAccAddr):
		_, decoded, _ := bech32.DecodeAndConvert(address)
		cosmosOperAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixValAddr, decoded)
		_ = db.Model(&validatorInfo).
			Where("operator_address = ?", cosmosOperAddress).
			Limit(1).
			Select()
	case strings.HasPrefix(address, sdk.Bech32PrefixValAddr):
		_ = db.Model(&validatorInfo).
			Where("operator_address = ?", address).
			Limit(1).
			Select()
	case strings.HasPrefix(address, sdk.Bech32PrefixValPub):
		_ = db.Model(&validatorInfo).
			Where("consensus_pubkey = ?", address).
			Limit(1).
			Select()
	case len(address) == 40:
		fmt.Println("")
		fmt.Println("address: ", address)
		address := strings.ToUpper(address)
		_ = db.Model(&validatorInfo).
			Where("proposer = ?", address).
			Limit(1).
			Select()
	default:
		_ = db.Model(&validatorInfo).
			Where("moniker = ?", address).
			Limit(1).
			Select()
	}
	return validatorInfo, nil
}

// ConvertToProposerSlice() converts any type of input address to proposer address
func ConvertToProposerSlice(address string, db *pg.DB) ([]dbtypes.ValidatorInfo, error) {
	var validatorInfo []dbtypes.ValidatorInfo
	switch {
	case strings.HasPrefix(address, sdk.Bech32PrefixAccAddr):
		_, decoded, _ := bech32.DecodeAndConvert(address)
		cosmosOperAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixValAddr, decoded)
		_ = db.Model(&validatorInfo).
			Where("operator_address = ?", cosmosOperAddress).
			Limit(1).
			Select()
	case strings.HasPrefix(address, sdk.Bech32PrefixValAddr):
		_ = db.Model(&validatorInfo).
			Where("operator_address = ?", address).
			Limit(1).
			Select()
	case strings.HasPrefix(address, sdk.Bech32PrefixValPub):
		_ = db.Model(&validatorInfo).
			Where("consensus_pubkey = ?", address).
			Limit(1).
			Select()
	case len(address) == 40:
		upperCaseAddr := strings.ToUpper(address)
		_ = db.Model(&validatorInfo).
			Where("proposer = ?", upperCaseAddr).
			Limit(1).
			Select()
	default:
		_ = db.Model(&validatorInfo).
			Where("moniker = ?", address).
			Limit(1).
			Select()
	}
	return validatorInfo, nil
}

// ConvertToProposerForValidatorStats() converts any type of input address to proposer address
func ConvertToProposerForValidatorStats(address string, db *pg.DB) (stats.ValidatorStats, error) {
	var validatorStats stats.ValidatorStats
	switch {
	case strings.HasPrefix(address, sdk.Bech32PrefixAccAddr):
		_ = db.Model(&validatorStats).
			Where("moniker = ?", address).
			Limit(1).
			Select()
	case strings.HasPrefix(address, sdk.Bech32PrefixValAddr):
		_ = db.Model(&validatorStats).
			Where("operator_address = ?", address).
			Limit(1).
			Select()
	case strings.HasPrefix(address, sdk.Bech32PrefixValPub):
		_ = db.Model(&validatorStats).
			Where("consensus_pubkey = ?", address).
			Limit(1).
			Select()
	case len(address) == 40:
		upperCaseAddr := strings.ToUpper(address)
		_ = db.Model(&validatorStats).
			Where("proposer = ?", upperCaseAddr).
			Limit(1).
			Select()
	default:
		_ = db.Model(&validatorStats).
			Where("moniker = ?", address).
			Limit(1).
			Select()
	}
	return validatorStats, nil
}
