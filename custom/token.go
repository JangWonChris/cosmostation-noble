package custom

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	Cosmos           = "uatom" // 사용 안하는 중
	CoinGeckgoCoinID = "cosmos"
	Currency         = "usd"
)

var (
	// config로 빼자
	NonNativeAssets = []string{}

	//ex) 18자리인 경우 : PowerReduction = sdktypes.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	PowerReduction = sdktypes.PowerReduction // 1e6;
)
