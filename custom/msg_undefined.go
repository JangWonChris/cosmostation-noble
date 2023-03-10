package custom

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap"
)

func AccountExporterFromUndefinedTxMsg(msg *sdktypes.Msg, txHash string) (msgType string, accounts []string) {
	switch msg := (*msg).(type) {

	default:
		msgType = proto.MessageName(msg)
		zap.S().Errorf("Undefined msg Type : %T(hash = %s)\n", msg, txHash)
	}

	return
}
