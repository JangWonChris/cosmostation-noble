package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tendermint/libs/bech32"
)

// SetAppConfig creates a new config instance for the SDK configuration.
func SetAppConfig() {
	config := sdk.GetConfig()
	SetBech32AddressPrefixes(config)
	SetBip44CoinType(config)
	config.Seal()
}

// SetBech32AddressPrefixes sets the global prefix to be used when serializing addresses to bech32 strings.
func SetBech32AddressPrefixes(config *sdk.Config) {
	// config.SetBech32PrefixForAccount(app.Bech32MainPrefix, app.Bech32MainPrefix+sdk.PrefixPublic)
	// config.SetBech32PrefixForValidator(app.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator, app.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic)
	// config.SetBech32PrefixForConsensusNode(app.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus, app.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic)
}

// SetBip44CoinType sets the global coin type to be used in hierarchical deterministic wallets.
func SetBip44CoinType(config *sdk.Config) {
	// config.SetCoinType(app.Bip44CoinType)
}

// ConvertConsAddrFromConsPubkey converts validator consensus public key to proposer address format
func ConvertConsAddrFromConsPubkey(consPubKey string) (string, error) {
	// pubKey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, consPubKey)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to get pubkey from bech32: %s", err)
	// }
	// return pubKey.Address().String(), nil
	return ConsAddrFromConsPubkey(consPubKey)
}

// ConsAddrFromConsPubkey converts validator consensus public key to proposer address format
func ConsAddrFromConsPubkey(consensusPubKey string) (string, error) {
	// deprecated for cosmos-sdk 0.37.4
	pk, err := sdk.GetConsPubKeyBech32(consensusPubKey)
	if err != nil {
		return "", fmt.Errorf("failed to get pubkey from bech32: %s", err)
	}
	return pk.Address().String(), nil
}

// ConvertAccAddrFromValAddr converts validator operator address to account address.
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
