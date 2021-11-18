package custom

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

var (
	// config로 빼자
	NonNativeAssets = []string{}

	//ex) 18자리인 경우 : PowerReduction = sdktypes.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	// PowerReduction = sdktypes.DefaultPowerReduction // 1e6;
	PowerReduction = sdktypes.NewIntFromUint64(1000000) // 1e6
)
