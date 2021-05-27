package custom

import (
	"fmt"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"go.uber.org/zap"
)

func AccountExporterFromCustomTxMsg(msg *sdktypes.Msg, txHash string) (msgType string, accounts []string) {
	switch msg := (*msg).(type) {

	// custom msgs (아래 주석은 sample)
	// case *authvestingtypes.MsgCreateVestingAccount:
	// 	msgType = AuthMsgCreateVestingAccount
	// 	accounts = append(accounts, msg.FromAddress, msg.ToAddress)

	default:
		// 전체 case에서 이 msg를 찾지 못하였기 때문에 에러 로깅한다.
		msgType = fmt.Sprintf("%T", msg)
		zap.S().Errorf("Undefined msg Type : %T(hash = %s)\n", msg, txHash)
	}

	return
}
