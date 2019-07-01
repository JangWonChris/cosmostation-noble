package utils

import (
	"fmt"
	"strings"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/stats"
	ctypes "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-pg/pg"
	"github.com/tendermint/tendermint/libs/bech32"
)

// ConvertToProposer() converts any type of input address to proposer address
func ConvertToProposer(address string, db *pg.DB) (ctypes.ValidatorInfo, error) {
	var validatorInfo ctypes.ValidatorInfo
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
func ConvertToProposerSlice(address string, db *pg.DB) ([]ctypes.ValidatorInfo, error) {
	var validatorInfo []ctypes.ValidatorInfo
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

// Conver from operator address to cosmos address
func OperatorAddressToCosmosAddress(operatorAddress string) string {
	_, decoded, _ := bech32.DecodeAndConvert(operatorAddress)
	cosmosAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixAccAddr, decoded)

	return cosmosAddress
}
