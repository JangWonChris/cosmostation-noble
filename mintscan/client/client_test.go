package client

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/config"
)

var client *Client

func TestMain(m *testing.M) {
	// model.SetAppConfig()

	config := config.ParseConfig()
	client, _ = NewClient(config.Node, config.Market)

	os.Exit(m.Run())
}

func TestGetBlock(t *testing.T) {
	height := int64(67270)

	block, err := client.GetBlock(height)
	require.NoError(t, err)

	require.NotNil(t, block)
}

func TestGetCoinDenom(t *testing.T) {
	bondDenom, err := client.GetBondDenom()
	require.NoError(t, err)

	require.NotNil(t, bondDenom)
}

func TestGetValidatorCommission(t *testing.T) {
	truncatedValCommission, err := client.GetValidatorCommission("kavavaloper1j26c4k2jj9tv95whdhva3e8v2fcm4s3dsgstd2")
	require.NoError(t, err)

	require.NotNil(t, truncatedValCommission)
}

func TestParseTxResponse(t *testing.T) {
	hash := "A8A272A277213D17339B900B1EA2A634CBA33049327E6591648EDA8DA86AF7F2"

	txResponse, err := client.GetTx(hash)
	require.NoError(t, err)
	require.Equal(t, false, txResponse.Empty(), "tx hash has empty txResponse %s", hash)

	stdTx, ok := txResponse.Tx.(auth.StdTx)
	require.Equal(t, false, !ok, "unsupported tx type: %s", txResponse.Tx)

	msgsBz, err := client.cdc.MarshalJSON(stdTx.GetMsgs())
	require.NoError(t, err, "failed to unmarshal transaction messages")

	feeBz, err := client.cdc.MarshalJSON(stdTx.Fee)
	require.NoError(t, err, "failed to unmarshal tx fee")

	logsBz, err := client.cdc.MarshalJSON(txResponse.Logs)
	require.NoError(t, err, "failed to unmarshal tx logs")

	sigs := make([]auth.StdSignature, len(stdTx.GetSignatures()), len(stdTx.GetSignatures()))
	for i, s := range stdTx.GetSignatures() {
		sigs[i].Signature = s.Signature
		sigs[i].PubKey = s.PubKey
	}

	// for i, pk := range stdTx.GetPubKeys() {
	// 	sigs[i].PubKey = pk
	// }

	sigsBz, err := client.cdc.MarshalJSON(sigs)
	require.NoError(t, err, "failed to unmarshal tx signatures")

	fmt.Println("string(msgsBz): ", string(msgsBz))

	require.NotNil(t, string(msgsBz), "Messages")
	require.NotNil(t, string(feeBz), "Fees")
	require.NotNil(t, string(sigsBz), "Signatures")
	require.NotNil(t, string(logsBz), "Logs")
}
