package types

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/bech32"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ConvertConsAddrFromConsPubkey converts validator consensus public key to proposer address format
func ConvertConsAddrFromConsPubkey(consPubKey string) (string, error) {
	return ConsAddrFromConsPubkey(consPubKey)
}

// ConsAddrFromConsPubkey converts validator consensus public key to proposer address format
func ConsAddrFromConsPubkey(consensusPubKey string) (string, error) {
	pk, err := sdk.GetConsPubKeyBech32(consensusPubKey)
	if err != nil {
		return "", fmt.Errorf("failed to get pubkey from bech32: %s", err)
	}
	return pk.Address().String(), nil
}

// ConvertAccAddrFromValAddr converts validator operator address to account address.
func ConvertAccAddrFromValAddr(valAddr string) (string, error) {
	return AccAddressFromOperatorAddress(valAddr)
}

// AccAddressFromOperatorAddress converts operator address to cosmos address
func AccAddressFromOperatorAddress(operatorAddress string) (string, error) {
	// _, decoded, _ := bech32.DecodeAndConvert(operatorAddress)
	// address, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixAccAddr, decoded)
	// return address
	_, decoded, err := bech32.DecodeAndConvert(operatorAddress)
	if err != nil {
		return "", fmt.Errorf("failed to decode and convert: %s", err)
	}

	accAddr, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32AccountAddrPrefix(), decoded)
	if err != nil {
		return "", fmt.Errorf("failed to convert and encode: %s", err)
	}

	return accAddr, nil
}

// ConvertValAddrFromAccAddr converts account address to validator operator address.
func ConvertValAddrFromAccAddr(accAddr string) (string, error) {
	_, decoded, err := bech32.DecodeAndConvert(accAddr)
	if err != nil {
		return "", fmt.Errorf("failed to decode and convert: %s", err)
	}

	valAddr, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), decoded)
	if err != nil {
		return "", fmt.Errorf("failed to convert and encode: %s", err)
	}

	return valAddr, nil
}
