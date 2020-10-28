package exporter

import (
	"testing"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/stretchr/testify/require"
)

var (
	// Staking module transaction hashes
	SampleMsgCreateValidatorTxHash = "52C69CC3C4D51A5D96D7B451D5F5DEB554609892D9FEB1793AA6BD00EE5D3F9D"
	SampleMsgDelegateTxHash        = "B201EC273AB0D881C52B06217B5A59EFC88FACCB43D5481AD47C4212F3D2F61E"
	SampleMsgUndelegateTxHash      = "86D941EAD011931BDD26021ACDCA9285CB899AB05B6E8B00511BEA51C851945C"
	SampleMsgBeginRedelegateTxHash = "C0281D29FE13C85680D7F77643B01E8AE738E851A2DACC71B86E66B3DCF2765E"
)

func TestParseMsgCreateValidator(t *testing.T) {
	_, stdTx, err := commonTxParser(SampleMsgCreateValidatorTxHash)
	require.NoError(t, err)

	for _, msg := range stdTx.GetMsgs() {
		msgCreateValidator, ok := msg.(*stakingtypes.MsgCreateValidator)
		require.Equal(t, true, ok)
		require.NotNil(t, msgCreateValidator)
	}
}

func TestParseMsgDelegate(t *testing.T) {
	_, stdTx, err := commonTxParser(SampleMsgDelegateTxHash)
	require.NoError(t, err)

	for _, msg := range stdTx.GetMsgs() {
		msgDelegate, ok := msg.(*stakingtypes.MsgDelegate)
		require.Equal(t, true, ok)
		require.NotNil(t, msgDelegate)
	}
}

func TestParseMsgUndelegate(t *testing.T) {
	_, stdTx, err := commonTxParser(SampleMsgUndelegateTxHash)
	require.NoError(t, err)

	for _, msg := range stdTx.GetMsgs() {
		msgUndelegate, ok := msg.(*stakingtypes.MsgUndelegate)
		require.Equal(t, true, ok)
		require.NotNil(t, msgUndelegate)
	}
}

func TestParseMsgBeginRedelegate(t *testing.T) {
	_, stdTx, err := commonTxParser(SampleMsgBeginRedelegateTxHash)
	require.NoError(t, err)

	for _, msg := range stdTx.GetMsgs() {
		msgBeginRedelegate, ok := msg.(*stakingtypes.MsgBeginRedelegate)
		require.Equal(t, true, ok)
		require.NotNil(t, msgBeginRedelegate)
	}
}

func TestValidatorRank(t *testing.T) {
	ex := NewExporter()
	ex.saveValidators()
}
