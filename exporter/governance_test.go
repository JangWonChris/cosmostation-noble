package exporter

import (
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmostation/cosmostation-cosmos/custom"
	mdschema "github.com/cosmostation/mintscan-database/schema"

	"github.com/stretchr/testify/require"
)

var (
	// Governance module transaction hashes
	SampleMsgSubmitProposalTxHash = "181C8382A5E32F42F3B9E26D93445C6BEB13C71DF02A13B778AB9A53E1F03AD9"
	SampleMsgDepositTxHash        = "0106E424610C7E18AC9985E24F18A6615F43623B3C9026BAC556459552445737"
	SampleMsgVoteTxHash           = "7CD599384D2AB0B8D3EF7B32CD3795C22E02C86FB9BC1091B4BAEDC8ED2518C7"
)

func TestParseMsgSubmitProposal(t *testing.T) {
	txResp, stdTx, err := commonTxParser(SampleMsgSubmitProposalTxHash)
	require.NoError(t, err)

	for _, msg := range stdTx.GetMsgs() {
		msgSubmitProposal, ok := msg.(*govtypes.MsgSubmitProposal)
		require.Equal(t, true, ok)
		require.NotNil(t, msgSubmitProposal)

		// Proposal ID
		for _, log := range txResp.Logs {
			for _, event := range log.Events {
				if event.Type == "submit_proposal" {
					for _, attribute := range event.Attributes {
						if attribute.Key == "proposal_id" {
							require.NotNil(t, attribute.Value)
						}
					}
				}
			}
		}
	}
}

func TestParseMsgDeposit(t *testing.T) {
	_, stdTx, err := commonTxParser(SampleMsgDepositTxHash)
	require.NoError(t, err)

	for _, msg := range stdTx.GetMsgs() {
		msgDeposit, ok := msg.(*govtypes.MsgDeposit)
		require.Equal(t, true, ok)
		require.NotNil(t, msgDeposit)
	}
}

func TestParseMsgVote(t *testing.T) {
	_, stdTx, err := commonTxParser(SampleMsgVoteTxHash)
	require.NoError(t, err)

	for _, msg := range stdTx.GetMsgs() {
		msgVote, ok := msg.(*govtypes.MsgVote)
		require.Equal(t, true, ok)
		require.NotNil(t, msgVote)
	}
}

func TestSaveAllProposals(t *testing.T) {
	ex.saveAllProposals()
}

func TestGetProposalVyStatus(t *testing.T) {
	ex.saveAllProposals()
	dp, err := ex.Client.GetProposalsByStatus(govtypes.StatusDepositPeriod)
	require.NoError(t, err)
	t.Log(dp)
}

func TestExec(t *testing.T) {
	msgType := "authz/exec"
	limit := 50
	beginID, err := ex.DB.GetBeginTxIDByMsgType(msgType)
	require.NoError(t, err)
	// beginID = int64(4701382)
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
		txResps := make([]*sdkTypes.TxResponse, 0)
		for i := range rawTxs {
			tx := &sdktypes.TxResponse{}
			err := custom.AppCodec.UnmarshalJSON(rawTxs[i].Chunk, tx)
			require.NoError(t, err)
			txResps = append(txResps, tx)
			cnt_msg := len(tx.GetTx().GetMsgs())
			cnt_msgs += cnt_msg
		}
		t.Log("bid : ", bid, " len of txs : ", len(txResps), " len of msgs :", cnt_msgs)
		p, d, v, err := ex.getGovernance(nil, txResps)
		require.NoError(t, err)
		if len(v) > 0 {
			if len(p) > 0 || len(d) > 0 {
				t.Log("P :", len(p), " d :", len(d), " v :", len(v))
			} else {
				t.Log("v :", len(v))
			}
			// 	for i := range v {
			// 		t.Log(v[i].Voter, v[i].TxHash, v[i].Option, v[i].ProposalID)
			// 	}

			basic := new(mdschema.BasicData)
			basic.Votes = v
			// t.Log(basic.Votes)
			err = ex.DB.InsertExportedData(basic)
			require.NoError(t, err)
			t.Log("vote updated")
		}

		bid = rawTxs[len(rawTxs)-1].ID + 1
		// time.Sleep(200 * time.Millisecond)
	}

}
