package custom

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

const ()

func AccountExporterFromCustomTxMsg(msg *sdktypes.Msg, txHash string) (msgType string, accounts []string) {
	switch (*msg).(type) {
	default:
		// AccountExporterFromCustomTxMsg() 에서 처리
	}

	return
}
