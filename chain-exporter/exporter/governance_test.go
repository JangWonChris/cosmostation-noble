package exporter

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/x/gov"

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

	for _, msg := range stdTx.Msgs {
		msgSubmitProposal, ok := msg.(gov.MsgSubmitProposal)
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

	for _, msg := range stdTx.Msgs {
		msgDeposit, ok := msg.(gov.MsgDeposit)
		require.Equal(t, true, ok)
		require.NotNil(t, msgDeposit)
	}
}

func TestParseMsgVote(t *testing.T) {
	_, stdTx, err := commonTxParser(SampleMsgVoteTxHash)
	require.NoError(t, err)

	for _, msg := range stdTx.Msgs {
		msgVote, ok := msg.(gov.MsgVote)
		require.Equal(t, true, ok)
		require.NotNil(t, msgVote)
	}
}
