package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/bech32"
)

// Conver from operator address to cosmos address
func OperatorAddressToAddress(operatorAddress string) string {
	_, decoded, _ := bech32.DecodeAndConvert(operatorAddress)
	address, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixAccAddr, decoded)
	return address
}

func ConsensusPubkeyToProposer(consensusPubKey string) string {
	pk, _ := sdk.GetConsPubKeyBech32(consensusPubKey)
	hexAddr := pk.Address().String()
	return hexAddr
}
