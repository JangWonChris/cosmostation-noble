package exporter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatorRank(t *testing.T) {
	ex.saveValidators()
}

func TestGetPowerEventHistory(t *testing.T) {

	b, err := ex.Client.RPC.GetBlock(5103)
	require.NoError(t, err)

	stdTx, err := ex.Client.CliCtx.GetTxs(b)
	require.NoError(t, err)
	// b.Block.Data.Txs
	// _, stdTx, err := commonTxParser(SampleMsgCreateValidatorTxHash)
	// require.NoError(t, err)

	peh, err := ex.getPowerEventHistoryNew(stdTx)
	require.NoError(t, err)

	for _, p := range peh {
		t.Log("height:", p.Height)
		t.Log("moniker:", p.Moniker)
		t.Log("operaddr:", p.OperatorAddress)
		t.Log("proposer:", p.Proposer)
		t.Log("txhash:", p.TxHash)
	}
}
