package custom

import (
	"fmt"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	mbltypes "github.com/cosmostation/mintscan-backend-library/types"
	liquiditytypes "github.com/gravity-devs/liquidity/x/liquidity/types"
	"go.uber.org/zap"
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
		// 전체 case에서 이 msg를 찾지 못하였기 때문에 에러 로깅한다.
		msgType = fmt.Sprintf("%T", msg)
		zap.S().Errorf("Undefined msg Type : %T(hash = %s)\n", msg, txHash)
	}

	return
}
