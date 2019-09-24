package utils

import (
	"github.com/tendermint/tendermint/libs/bech32"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ConsAddrFromConsPubkey converts validator consensus public key to proposer address format
func ConsAddrFromConsPubkey(consensusPubKey string) string {
	pk, _ := sdk.GetConsPubKeyBech32(consensusPubKey)
	proposerAddress := pk.Address().String()
	return proposerAddress
}

// AccAddressFromOperatorAddress converts operator address to cosmos address
func AccAddressFromOperatorAddress(operatorAddress string) string {
	_, decoded, _ := bech32.DecodeAndConvert(operatorAddress)
	address, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixAccAddr, decoded)
	return address
}
