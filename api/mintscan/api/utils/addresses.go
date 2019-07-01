package utils

import (
	"fmt"
	"strings"
	"encoding/hex"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/stats"
	dbtypes "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-pg/pg"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/bech32"
)

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

// Conver from operator address to cosmos address
func OperatorAddressToCosmosAddress(operatorAddress string) string {
	_, decoded, _ := bech32.DecodeAndConvert(operatorAddress)
	cosmosAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixAccAddr, decoded)

	return cosmosAddress
}

// ConsensusPubkeyToProposer() receives consensusPublKey and returns proposer address
func ConsensusPubkeyToProposer(consensusPubKey string) string {
	_, data, err := bech32.DecodeAndConvert(consensusPubKey)
	if err != nil {
		panic(err)
	}
	encodedData := hex.EncodeToString(data[:])
	subStringLast64 := encodedData[len(encodedData)-64:]

	decodedData, _ := hex.DecodeString(subStringLast64)
	convertedData := [32]byte{}
	copy(convertedData[:], decodedData)

	rePub := ed25519.PubKeyEd25519(convertedData)
	proposer := rePub.Address().String()

	return proposer
}

// ConvertCosmosAddressToMoniker() converts from cosmos address to moniker
func ConvertCosmosAddressToMoniker(cosmosAddr string, db *pg.DB) string {
	// First convert from Cosmos Address to ValiOperatorAddress
	_, decoded, _ := bech32.DecodeAndConvert(cosmosAddr)
	cosmosOperAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixValAddr, decoded)

	// Check if the address matches any moniker in our DB
	var validatorInfo dbtypes.ValidatorInfo
	_ = db.Model(&validatorInfo).
		Column("moniker").
		Where("operator_address = ?", cosmosOperAddress).
		Limit(1).
		Select()

	return validatorInfo.Moniker
}
