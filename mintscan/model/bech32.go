package model

import (
	"fmt"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

// ConvertConsAddrFromConsPubkey converts validator consensus public key to proposer address format
func ConvertConsAddrFromConsPubkey(consPubKey string) (string, error) {
	pubKey, err := sdktypes.GetPubKeyFromBech32(sdktypes.Bech32PubKeyTypeConsPub, consPubKey)
	if err != nil {
		return "", fmt.Errorf("failed to get pubkey from bech32: %s", err)
	}

	return pubKey.Address().String(), nil
}

// ConvertAccAddrFromValAddr converts validator operator address to account address.
func ConvertAccAddrFromValAddr(valAddr string) (string, error) {
	_, decoded, err := bech32.DecodeAndConvert(valAddr)
	if err != nil {
		return "", fmt.Errorf("failed to decode and convert: %s", err)
	}

	accAddr, err := bech32.ConvertAndEncode(sdktypes.GetConfig().GetBech32AccountAddrPrefix(), decoded)
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

	valAddr, err := bech32.ConvertAndEncode(sdktypes.GetConfig().GetBech32ValidatorAddrPrefix(), decoded)
	if err != nil {
		return "", fmt.Errorf("failed to convert and encode: %s", err)
	}

	return valAddr, nil
}

// VerifyBech32AccAddr validates bech32 account address format.
func VerifyBech32AccAddr(accAddr string) error {
	bz, err := sdktypes.GetFromBech32(accAddr, sdktypes.GetConfig().GetBech32AccountAddrPrefix())
	if err != nil {
		return err
	}

	return sdktypes.VerifyAddressFormat(bz)
}

// VerifyBech32ValAddr validates bech32 validator address format.
func VerifyBech32ValAddr(accAddr string) error {
	bz, err := sdktypes.GetFromBech32(accAddr, sdktypes.GetConfig().GetBech32ValidatorAddrPrefix())
	if err != nil {
		return err
	}

	return sdktypes.VerifyAddressFormat(bz)
}
