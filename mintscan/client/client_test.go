package client

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	sdktypestx "github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	mintscanconfig "github.com/cosmostation/cosmostation-cosmos/mintscan/config"
)

var cli *Client

func TestMain(m *testing.M) {
	config := mintscanconfig.ParseConfig()
	cli, _ = NewClient(config.Node, config.Market)

	os.Exit(m.Run())
}

func TestGetChainID(t *testing.T) {
	chainID, err := cli.GetNetworkChainID()
	require.NoError(t, err)

	require.NotNil(t, chainID)
}

func TestGetBlock(t *testing.T) {
	height := int64(67270)

	block, err := cli.GetBlock(height)
	require.NoError(t, err)

	require.NotNil(t, block)
}

func TestGetCoinDenom(t *testing.T) {
	bondDenom, err := cli.GetBondDenom()
	require.NoError(t, err)

	require.NotNil(t, bondDenom)
}

func TestGetAccountSpendableCoins(t *testing.T) {
	address := "cosmos1emaa7mwgpnpmc7yptm728ytp9quamsvuz92x5u"
	// account, err := cli.GetAccount("cosmos1clpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9q")
	// require.NoError(t, err)
	sdkaddr, err := sdktypes.AccAddressFromBech32(address)
	require.NoError(t, err)

	b := banktypes.NewQueryBalanceRequest(sdkaddr, "umuon")
	log.Println(b)
	bankClient := banktypes.NewQueryClient(cli.grpcClient)
	var header metadata.MD
	blockHeight := header.Get(grpctypes.GRPCBlockHeightHeader)
	bankRes, err := bankClient.Balance(
		context.Background(),
		b,
		grpc.Header(&header), // Also fetch grpc header
	)
	if err != nil {
		log.Println(err)
	}
	require.NotNil(t, *bankRes.GetBalance())

	blockHeight = header.Get(grpctypes.GRPCBlockHeightHeader)
	log.Println("blockHeight :", blockHeight)

}

func TestGetAccountDelegatedCoins(t *testing.T) {
	delegations, err := cli.GetDelegatorDelegations("cosmos1clpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9q")
	require.NoError(t, err)

	for _, delegation := range delegations {
		require.NotNil(t, delegation)
	}
}

func TestGetAccountUndelegatedCoins(t *testing.T) {
	undelegations, err := cli.GetDelegatorUndelegations("cosmos1clpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9q")
	require.NoError(t, err)

	for _, undelegation := range undelegations {
		require.NotNil(t, undelegation)
	}
}

func TestGetAccountTotalRewards(t *testing.T) {
	rewards, err := cli.GetDelegatorTotalRewards("cosmos1clpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9q")
	require.NoError(t, err)

	require.NotNil(t, rewards)
}

func TestGetValidatorCommission(t *testing.T) {
	truncatedValCommission, err := cli.GetValidatorCommission("cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn")
	require.NoError(t, err)

	require.NotNil(t, truncatedValCommission)
}

func TestParseTxResponse(t *testing.T) {
	hash := "A8A272A277213D17339B900B1EA2A634CBA33049327E6591648EDA8DA86AF7F2"

	txResponse, err := cli.GetTx(hash)
	require.NoError(t, err)
	require.Equal(t, false, txResponse.Empty(), "tx hash has empty txResponse %s", hash)

	// stdTx, ok := txResponse.Tx.(auth.StdTx)
	txI := txResponse.GetTx()
	tx, ok := txI.(*sdktypestx.Tx)
	require.Equal(t, false, !ok, "unsupported type")

	msgsBz, err := cli.cliCtx.JSONMarshaler.MarshalJSON(tx.GetBody())
	require.NoError(t, err, "failed to unmarshal transaction messages")

	feeBz, err := cli.cliCtx.JSONMarshaler.MarshalJSON(tx.GetAuthInfo().GetFee())
	require.NoError(t, err, "failed to unmarshal tx fee")

	logsBz, err := json.Marshal(txResponse.Logs)
	require.NoError(t, err, "failed to unmarshal tx logs")

	sigs := make([][]byte, len(tx.GetSignatures()), len(tx.GetSignatures()))
	for i, s := range tx.GetSignatures() {
		sigs[i] = s
		// sigs[i].PubKey = s.PubKey
	}

	sigsBz, err := json.Marshal(sigs)
	require.NoError(t, err, "failed to unmarshal tx signatures")

	require.NotNil(t, string(msgsBz), "Messages")
	require.NotNil(t, string(feeBz), "Fees")
	require.NotNil(t, string(sigsBz), "Signatures")
	require.NotNil(t, string(logsBz), "Logs")
}
