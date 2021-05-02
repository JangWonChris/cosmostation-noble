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
	PowerReduction  = sdktypes.PowerReduction // 1e6;
)
