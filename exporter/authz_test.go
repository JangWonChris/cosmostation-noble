package exporter

import (
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmostation/cosmostation-cosmos/custom"
	mdschema "github.com/cosmostation/mintscan-database/schema"
	"github.com/stretchr/testify/require"
)

// authz msg 복구용으로 사용
func TestRecoverVoteFromExec(t *testing.T) {
	msgType := "authz/exec"
	limit := 50
	beginID, err := ex.DB.GetBeginTxIDByMsgType(msgType)
	require.NoError(t, err)
	// beginID = int64(5230198)
	t.Log(beginID)

	bid := beginID
	cnt_msgs := 0
	for {
		rawTxs, err := ex.DB.GetTransactionsByMsgType(bid, msgType, limit)
		if err != nil {
			t.Log("error occured beigin id : ", bid)
		}
		require.NoError(t, err)
		if len(rawTxs) == 0 {
			t.Log("len of txs is zero")
			break
		}
		cnt_msgs = 0
		txResps := make([]*sdktypes.TxResponse, 0)
		for i := range rawTxs {
			tx := &sdktypes.TxResponse{}
			err := custom.AppCodec.UnmarshalJSON(rawTxs[i].Chunk, tx)
			require.NoError(t, err)
			txResps = append(txResps, tx)
			cnt_msg := len(tx.GetTx().GetMsgs())
			// if cnt_msg > 1 {
			// 	t.Log(tx.TxHash)
			// }
			cnt_msgs += cnt_msg
		}
		t.Log("bid : ", bid, " len of txs : ", len(txResps), " len of msgs :", cnt_msgs)
		proposals, deposits, votes, err := ex.getGovernance(nil, txResps)
		require.NoError(t, err)
		// powerEvents, err := ex.getPowerEventHistoryNew(txResps)
		// require.NoError(t, err)
		basic := new(mdschema.BasicData)
		basic.Deposits = deposits
		basic.Votes = votes
		// basic.ValidatorsPowerEventHistory = powerEvents
		sum := len(proposals) + len(deposits) + len(votes) //+ len(powerEvents)
		if sum > 0 {
			t.Log("p :", len(proposals), " d :", len(deposits), " v :", len(votes)) //, " power :", len(powerEvents))
			err = ex.DB.InsertExportedData(basic)
			require.NoError(t, err)
		}
		t.Logf("done, next bid : %d", rawTxs[len(rawTxs)-1].ID+1)

		bid = rawTxs[len(rawTxs)-1].ID + 1
		// time.Sleep(200 * time.Millisecond)
	}

}
