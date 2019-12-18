package utils

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dbtypes "github.com/cosmostation/cosmostation-cosmos/mintscan/api/models/types"

	"github.com/go-pg/pg"

	"github.com/tendermint/tendermint/libs/bech32"
)

// ConvertCosmosAddressToMoniker converts from cosmos address to moniker
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

// ConvertToProposer converts any type of input address to proposer address
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

// ConvertToProposerSlice converts any type of input address to proposer address
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
