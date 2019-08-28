package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/bech32"
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

// ValAddressFromAccAddress converts account address to validator address
func ValAddressFromAccAddress(consAddress string) string {
	_, decoded, _ := bech32.DecodeAndConvert(consAddress)
	address, _ := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), decoded)
	return address
}
