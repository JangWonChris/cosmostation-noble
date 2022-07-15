package custom

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	mbltypes "github.com/cosmostation/mintscan-backend-library/types"
	liquiditytypes "github.com/gravity-devs/liquidity/x/liquidity/types"
)

const (
	LiquidityMsgCreatePool          = "liquidity/create_pool"
	LiquidityMsgDepositWithinBatch  = "liquidity/deposit_within_batch"
	LiquidityMsgWithdrawWithinBatch = "liquidity/withdraw_within_batch"
	LiquidityMsgSwapWithinBatch     = "liquidity/swap_within_batch"
)

func AccountExporterFromCustomTxMsg(msg *sdktypes.Msg, txHash string) (msgType string, accounts []string) {
	switch msg := (*msg).(type) {

	case *liquiditytypes.MsgCreatePool:
		msgType = LiquidityMsgCreatePool
		accounts = mbltypes.AddNotNullAccount(msg.PoolCreatorAddress)
	case *liquiditytypes.MsgDepositWithinBatch:
		msgType = LiquidityMsgDepositWithinBatch
		accounts = mbltypes.AddNotNullAccount(msg.DepositorAddress)
	case *liquiditytypes.MsgWithdrawWithinBatch:
		msgType = LiquidityMsgWithdrawWithinBatch
		accounts = mbltypes.AddNotNullAccount(msg.WithdrawerAddress)
	case *liquiditytypes.MsgSwapWithinBatch:
		msgType = LiquidityMsgSwapWithinBatch
		accounts = mbltypes.AddNotNullAccount(msg.SwapRequesterAddress)

	default:
		// AccountExporterFromCustomTxMsg() 에서 처리
	}

	return
}
