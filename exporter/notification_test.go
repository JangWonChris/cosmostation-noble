package exporter

import (
	"github.com/cosmostation/cosmostation-cosmos/app"
	"testing"
	"time"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/stretchr/testify/require"
)

var (
	// Bank module staking hashes
	SampleMsgSendTxHash      = "A80ADDA7929801AF3B1E6957BE9C63C30B5A0B9F903E760C555CAC19D2FC0DFC"
	SampleMsgMultiSendTxHash = ""
)

func TestParseMsgSend(t *testing.T) {
	_, stdTx, err := commonTxParser(SampleMsgSendTxHash)
	require.NoError(t, err)

	for _, msg := range stdTx.GetMsgs() {
		msgSend, ok := msg.(*banktypes.MsgSend)
		require.Equal(t, true, ok)
		require.NotNil(t, msgSend)
	}
}

// TODO: no available tx hash in mainnet
func TestParseMsgMultiSend(t *testing.T) {
	_, _, err := commonTxParser(SampleMsgMultiSendTxHash)
	require.NoError(t, err)
}

func TestProposalAlarm(t *testing.T){
	chainEx := app.NewApp("chain-exporter")
	ex = NewExporter(chainEx)
	go ex.ProposalNotificationToSlack(51)
	time.Sleep(5*time.Second)
}