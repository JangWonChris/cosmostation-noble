package utils

import (
	"github.com/tendermint/tendermint/libs/bech32"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Convert validator consensus public key to proposer address format
func ConsensusPubkeyToProposer(consensusPubKey string) string {
	pk, _ := sdk.GetConsPubKeyBech32(consensusPubKey)
	proposerAddress := pk.Address().String()
	return proposerAddress
}

// Convert operator address to cosmos address
func OperatorAddressToCosmosAddress(operatorAddress string) string {
	_, decoded, _ := bech32.DecodeAndConvert(operatorAddress)
	cosmosAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixAccAddr, decoded)
	return cosmosAddress
}
