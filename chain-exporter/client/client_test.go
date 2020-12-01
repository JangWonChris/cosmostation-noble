package client

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	sdktypestx "github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var cli *Client

func TestMain(m *testing.M) {
	log.Println("testmain start")
	cfg := config.ParseConfig()
	var err error
	cli, err = NewClient(cfg.Node, cfg.KeybaseURL)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	os.Exit(m.Run())

	log.Println("testmain end")
}

func TestGetAccount(t *testing.T) {
	address := "cosmos1pvzrncl89w5z9psr8ch90057va9tc23pehpd2t"
	sdkaddr, err := sdktypes.AccAddressFromBech32(address)
	if err != nil {
		log.Println(err)
	}
	ar := authtypes.AccountRetriever{}
	acc, err := ar.GetAccount(cli.cliCtx, sdkaddr)
	if err != nil {
		log.Println(err)
	}

	log.Println(acc.GetAddress())
	log.Println(acc.GetPubKey())
}

func TestGetAccountBalance(t *testing.T) {

	// address := "cosmos1x5wgh6vwye60wv3dtshs9dmqggwfx2ldnqvev0"
	address := "cosmos1emaa7mwgpnpmc7yptm728ytp9quamsvuz92x5u"
	// log.Println(cli.GetAccount("cosmos1x5wgh6vwye60wv3dtshs9dmqggwfx2ldnqvev0"))
	sdkaddr, err := sdktypes.AccAddressFromBech32(address)
	if err != nil {
		log.Println(err)
	}
	b := banktypes.NewQueryBalanceRequest(sdkaddr, "umuon")
	log.Println(b)
	bankClient := banktypes.NewQueryClient(cli.grpcClient)
	var header metadata.MD
	blockHeight := header.Get(grpctypes.GRPCBlockHeightHeader)
	log.Println("blockHeight :", blockHeight)
	// header.Set(k string, vals ...string)
	// header.Append(grpctypes.GRPCBlockHeightHeader, "1")
	// header.Set(grpctypes.GRPCBlockHeightHeader, "1")
	// bankRes, err := bankClient.Balance(
	// 	metadata.AppendToOutgoingContext(context.Background(), grpctypes.GRPCBlockHeightHeader, "1"), // Add metadata to request
	// 	b,
	// 	grpc.Header(&header),
	// )
	bankRes, err := bankClient.Balance(
		context.Background(),
		b,
		grpc.Header(&header), // Also fetch grpc header
	)
	if err != nil {
		log.Println(err)
	}
	if bankRes.GetBalance() != nil {
		log.Println(*bankRes.GetBalance())
	}
	blockHeight = header.Get(grpctypes.GRPCBlockHeightHeader)
	log.Println("blockHeight :", blockHeight)
}
func TestGetNetworkChainID(t *testing.T) {
	n, err := cli.GetNetworkChainID()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(n)
}

func TestGetBlock(t *testing.T) {
	log.Println(cli.GetBlock(11111))
}

func TestGetTx(t *testing.T) {
	sendTx := "A80ADDA7929801AF3B1E6957BE9C63C30B5A0B9F903E760C555CAC19D2FC0DFC"
	withdrawAllRewardsTx := "53A036CC53FD3AD8C4B66C11BBB20DC63A5B606144F6655EC9D9E327AB9BA3D9"
	delegateTx := "DDA04447F569B402D96E7CCCC9ACF0C76D3581EC9B056818CED7913DECA6F10A"
	unknownTx := "30B43BB887FA6F56E5302B6CCB9C439A6C2AF29CFADA1465C0174EE6C62E3D28"
	_, _, _, _ = sendTx, withdrawAllRewardsTx, delegateTx, unknownTx
	txhash := unknownTx
	txResp, err := cli.GetTx(txhash)
	if err != nil {
		log.Fatal(err)
	}

	tx := txResp.GetTx()
	ta, ok := tx.(*sdktypestx.Tx)
	log.Println(ok)
	if ok {
		a, err := cli.cliCtx.JSONMarshaler.MarshalJSON(ta.GetBody())
		if err != nil {
			log.Println(err)
		}
		log.Println("message :", string(a))
		a, err = cli.cliCtx.JSONMarshaler.MarshalJSON(ta.GetAuthInfo().GetFee())
		if err != nil {
			log.Println(err)
		}
		log.Println("fee :", string(a))
		a, err = cli.cliCtx.JSONMarshaler.MarshalJSON(ta.GetAuthInfo())
		if err != nil {
			log.Println(err)
		}
		log.Println("authinfo :", string(a))
		a, err = cli.cliCtx.JSONMarshaler.MarshalJSON(ta.GetAuthInfo().GetSignerInfos()[0])
		if err != nil {
			log.Println(err)
		}
		log.Println("signerinfo[0] :", string(a))
		for _, addr := range ta.GetSigners() {
			log.Println("getsigners addr :", addr)
		}
		a, err = cli.cliCtx.JSONMarshaler.MarshalJSON(ta.GetAuthInfo().GetSignerInfos()[0].GetPublicKey())
		if err != nil {
			log.Println(err)
		}
		log.Println("pubkey[0] :", ta.GetAuthInfo().GetSignerInfos()[0].GetPublicKey().GetValue())
		log.Println("pubkey[0] :", string(a))
		sig := ta.GetSignatures()
		// json.Unmarshal(sig[0], &i)
		log.Println("signatures :", sig[0])
		log.Println("memo:", ta.GetBody().Memo)
	}

	msgs := txResp.GetTx().GetMsgs()
	for i, m := range msgs {
		switch t := m.(type) {
		case *banktypes.MsgSend:
			log.Println("banktypes :", t)
			log.Println(t.FromAddress)
			log.Println(t.ToAddress)
			log.Println(t.Amount)
			log.Println(t.Type())
		default:
			log.Println(i, t)
		}
	}
}

func TestValidatorByStatus(t *testing.T) {
	// this is return empty result
	// v, err := cli.GetValidatorsByStatus(stakingtypes.Bonded)
	// require.NoError(t, err)
	// log.Println(v)

	// stakingtypes에 정의 됨
	// 0: "BOND_STATUS_UNSPECIFIED",
	// 1: "BOND_STATUS_UNBONDED",
	// 2: "BOND_STATUS_UNBONDING",
	// 3: "BOND_STATUS_BONDED",

	bonded := stakingtypes.BondStatus_name[int32(stakingtypes.Bonded)]
	// unbonded := stakingtypes.BondStatus_name[int32(stakingtypes.Unbonded)]
	// unbonding := stakingtypes.BondStatus_name[int32(stakingtypes.Unbonding)]

	queryClient := stakingtypes.NewQueryClient(cli.grpcClient)
	request := stakingtypes.QueryValidatorsRequest{Status: bonded}
	resp, err := queryClient.Validators(context.Background(), &request)
	require.NoError(t, err)
	log.Println("bonded :", len(resp.Validators))
	consAddr, err := resp.Validators[0].GetConsAddr()
	require.NoError(t, err)
	log.Println("consaddr :", consAddr)

	// request = stakingtypes.QueryValidatorsRequest{Status: unbonded}
	// resp, err = queryClient.Validators(context.Background(), &request)
	// require.NoError(t, err)
	// log.Println("unbonded :", len(resp.Validators))

	// request = stakingtypes.QueryValidatorsRequest{Status: unbonding}
	// resp, err = queryClient.Validators(context.Background(), &request)
	// require.NoError(t, err)
	// log.Println("unboding :", len(resp.Validators))
}
