package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tendermint/libs/bech32"
)

// Convert Cosmos Address to Opeartor Address
func ConvertOperatorAddressToCosmosAddress(address string) string {
	_, decoded, _ := bech32.DecodeAndConvert(address)
	cosmosAddr, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixAccAddr, decoded)
	return cosmosAddr
}

// Convert Opeartor Address to Cosmos Address
func ConvertCosmosAddressToOperatorAddress(address string) string {
	_, decoded, _ := bech32.DecodeAndConvert(address)
	valiOperatorAddr, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixValAddr, decoded)
	return valiOperatorAddr
}
