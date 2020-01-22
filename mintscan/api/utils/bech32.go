package utils

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/bech32"
)

// ConsAddrFromConsPubkey converts validator consensus public key to proposer address format
func ConsAddrFromConsPubkey(consensusPubKey string) string {
	pk, _ := sdk.GetConsPubKeyBech32(consensusPubKey)
	proposerAddr := pk.Address().String()

	return proposerAddr
}

// AccAddressFromOperatorAddress converts operator address to account address
func AccAddressFromOperatorAddress(operatorAddress string) string {
	_, decoded, _ := bech32.DecodeAndConvert(operatorAddress)
	accAddr, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixAccAddr, decoded)

	return accAddr
}

// ValAddressFromAccAddress converts account address to validator address
func ValAddressFromAccAddress(address string) string {
	_, decoded, _ := bech32.DecodeAndConvert(address)
	operAddr, _ := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), decoded)

	return operAddr
}

// VerifyAddress verifies address format
func VerifyAddress(address string) bool {
	if len(address) != 45 { // check length
		return false
	}

	if !strings.Contains(address, sdk.Bech32PrefixAccAddr) { // check prefix
		return false
	}
	return true
}

// VerifyValAddress verifies validator operator address format
func VerifyValAddress(address string) bool {
	if !strings.Contains(address, sdk.Bech32PrefixValAddr) { // check prefix
		return false
	}
	return true
}
