package utils

import (
	"encoding/hex"

	ctypes "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-pg/pg"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/bech32"
)

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
	var validatorInfo ctypes.ValidatorInfo
	_ = db.Model(&validatorInfo).
		Column("moniker").
		Where("operator_address = ?", cosmosOperAddress).
		Limit(1).
		Select()

	return validatorInfo.Moniker
}
