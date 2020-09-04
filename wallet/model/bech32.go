package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// VerifyBech32AccAddr validates bech32 account address format.
func VerifyBech32AccAddr(accAddr string) error {
	bz, err := sdk.GetFromBech32(accAddr, sdk.GetConfig().GetBech32AccountAddrPrefix())
	if err != nil {
		return err
	}

	return sdk.VerifyAddressFormat(bz)
}
