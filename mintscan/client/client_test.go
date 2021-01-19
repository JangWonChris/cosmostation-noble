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
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	mintscanconfig "github.com/cosmostation/mintscan-backend-library/config"
)

var cli *Client

func TestMain(m *testing.M) {
	config := mintscanconfig.ParseConfig()
	cli = NewClient(&config.Client)

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
	bankClient := cli.GRPC.GetBankQueryClient()
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
	resp, err := cli.GetDelegatorDelegations("cosmos1clpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9q")
	require.NoError(t, err)

	for _, delegation := range resp.DelegationResponses {
		require.NotNil(t, delegation)
	}
}

func TestGetAccountUndelegatedCoins(t *testing.T) {
	res, err := cli.GetDelegatorUnbondingDelegations("cosmos1clpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9q")
	require.NoError(t, err)

	for _, undelegation := range res.UnbondingResponses {
		require.NotNil(t, undelegation)
	}
}

func TestGetAccountTotalRewards(t *testing.T) {
	rewards, err := cli.GetDelegationTotalRewards("cosmos1clpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9q")
	require.NoError(t, err)

	require.NotNil(t, rewards)
}

func TestGetValidatorCommission(t *testing.T) {
	res, err := cli.GetValidatorCommission("cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn")
	require.NoError(t, err)

	require.NotNil(t, res.Commission)
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

	msgsBz, err := cli.CliCtx.JSONMarshaler.MarshalJSON(tx.GetBody())
	require.NoError(t, err, "failed to unmarshal transaction messages")

	feeBz, err := cli.CliCtx.JSONMarshaler.MarshalJSON(tx.GetAuthInfo().GetFee())
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

func TestRPCGetAccount(t *testing.T) {

	address := "cosmos1x5wgh6vwye60wv3dtshs9dmqggwfx2ldnqvev0"
	sdkaddr, err := sdktypes.AccAddressFromBech32(address)
	require.NoError(t, err)

	accGetter := authtypes.AccountRetriever{}
	acc, height, err := accGetter.GetAccountWithHeight(cli.GetCLIContext(), sdkaddr)
	require.NoError(t, err)

	log.Println(acc)
	log.Println(height)

	b, err := json.Marshal(acc)
	require.NoError(t, err)

	log.Println(string(b))
}
