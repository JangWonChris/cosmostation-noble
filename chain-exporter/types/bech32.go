package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tendermint/libs/bech32"
)

// ConvertConsAddrFromConsPubkey converts validator consensus public key to proposer address in hex string.
func ConvertConsAddrFromConsPubkey(consPubKey string) (string, error) {
	// This method will change to sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, consPubKey) in above v0.38.+.
	pk, err := sdk.GetConsPubKeyBech32(consPubKey)
	if err != nil {
		return "", fmt.Errorf("failed to get pubkey from bech32: %s", err)
	}
	return pk.Address().String(), nil
}

// ConvertAccAddrFromValAddr converts validator operator address to bech32 account address.
func ConvertAccAddrFromValAddr(valAddr string) (string, error) {
	_, decoded, err := bech32.DecodeAndConvert(valAddr)
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

// VerifyBech32AccAddr validates bech32 account address format.
func VerifyBech32AccAddr(accAddr string) error {
	bz, err := sdk.GetFromBech32(accAddr, sdk.GetConfig().GetBech32AccountAddrPrefix())
	if err != nil {
		return err
	}

	return sdk.VerifyAddressFormat(bz)
}

// VerifyBech32ValAddr validates bech32 validator address format.
func VerifyBech32ValAddr(accAddr string) error {
	bz, err := sdk.GetFromBech32(accAddr, sdk.GetConfig().GetBech32ValidatorAddrPrefix())
	if err != nil {
		return err
	}

	return sdk.VerifyAddressFormat(bz)
}
