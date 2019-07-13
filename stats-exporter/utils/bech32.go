package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tendermint/libs/bech32"
)

// Convert Cosmos Address to Opeartor Address
func ConvertOperatorAddressToAddress(operatorAddress string) string {
	_, decoded, _ := bech32.DecodeAndConvert(operatorAddress)
	address, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixAccAddr, decoded)
	return address
}

// Convert Opeartor Address to Cosmos Address
func ConvertAddressToOperatorAddress(address string) string {
	_, decoded, _ := bech32.DecodeAndConvert(address)
	valiOperatorAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixValAddr, decoded)
	return valiOperatorAddress
}
