package exporter

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/stretchr/testify/require"
)

var (
	// Bank module staking hashes
	SampleMsgSendTxHash      = "AFB944BA8230912E04363CBCC450F1F1EC9A2405B2D791C928E0597739093224"
	SampleMsgMultiSendTxHash = ""
)

func TestParseMsgSend(t *testing.T) {
	_, stdTx, err := commonTxParser(SampleMsgSendTxHash)
	require.NoError(t, err)

	for _, msg := range stdTx.Msgs {
		msgSend, ok := msg.(bank.MsgSend)
		require.Equal(t, true, ok)
		require.NotNil(t, msgSend)
	}
}

// TODO: no available tx hash in mainnet
func TestParseMsgMultiSend(t *testing.T) {
	_, _, err := commonTxParser(SampleMsgMultiSendTxHash)
	require.NoError(t, err)
}
