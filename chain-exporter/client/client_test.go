package client

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	sdktypestx "github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmostation/cosmostation-cosmos/chain-config/custom"
	mblconfig "github.com/cosmostation/mintscan-backend-library/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var cli *Client

func TestMain(m *testing.M) {
	fileBaseName := "chain-exporter"
	cfg := mblconfig.ParseConfig(fileBaseName)

	cli = NewClient(&cfg.Client)

	os.Exit(m.Run())

}

func TestGetAccount(t *testing.T) {
	// address := "cosmos1pvzrncl89w5z9psr8ch90057va9tc23pehpd2t"
	address := "cosmos1x5wgh6vwye60wv3dtshs9dmqggwfx2ldnqvev0"
	sdkaddr, err := sdktypes.AccAddressFromBech32(address)
	if err != nil {
		log.Println(err)
	}
	ar := authtypes.AccountRetriever{}
	log.Println(cli.CliCtx)
	acc, err := ar.GetAccount(cli.GetCLIContext(), sdkaddr)
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
	bankClient := banktypes.NewQueryClient(cli.GRPC)
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
	n, err := cli.RPC.GetNetworkChainID()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(n)
}

func TestGetBlock(t *testing.T) {
	log.Println(cli.RPC.GetBlock(11111))
}

func TestGetTx(t *testing.T) {
	sendTx := "A80ADDA7929801AF3B1E6957BE9C63C30B5A0B9F903E760C555CAC19D2FC0DFC"
	withdrawAllRewardsTx := "53A036CC53FD3AD8C4B66C11BBB20DC63A5B606144F6655EC9D9E327AB9BA3D9"
	delegateTx := "DDA04447F569B402D96E7CCCC9ACF0C76D3581EC9B056818CED7913DECA6F10A"
	unknownTx := "30B43BB887FA6F56E5302B6CCB9C439A6C2AF29CFADA1465C0174EE6C62E3D28"
	_, _, _, _ = sendTx, withdrawAllRewardsTx, delegateTx, unknownTx
	txhash := unknownTx
	txResp, err := cli.CliCtx.GetTx(txhash)
	if err != nil {
		log.Fatal(err)
	}

	tx := txResp.GetTx()
	ta, ok := tx.(*sdktypestx.Tx)
	log.Println(ok)
	if ok {
		a, err := cli.CliCtx.JSONMarshaler.MarshalJSON(ta.GetBody())
		if err != nil {
			log.Println(err)
		}
		log.Println("message :", string(a))
		a, err = cli.CliCtx.JSONMarshaler.MarshalJSON(ta.GetAuthInfo().GetFee())
		if err != nil {
			log.Println(err)
		}
		log.Println("fee :", string(a))
		a, err = cli.CliCtx.JSONMarshaler.MarshalJSON(ta.GetAuthInfo())
		if err != nil {
			log.Println(err)
		}
		log.Println("authinfo :", string(a))
		a, err = cli.CliCtx.JSONMarshaler.MarshalJSON(ta.GetAuthInfo().GetSignerInfos()[0])
		if err != nil {
			log.Println(err)
		}
		log.Println("signerinfo[0] :", string(a))
		for _, addr := range ta.GetSigners() {
			log.Println("getsigners addr :", addr)
		}
		a, err = cli.CliCtx.JSONMarshaler.MarshalJSON(ta.GetAuthInfo().GetSignerInfos()[0].GetPublicKey())
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

	queryClient := stakingtypes.NewQueryClient(cli.GRPC)
	request := stakingtypes.QueryValidatorsRequest{Status: bonded}
	resp, err := queryClient.Validators(context.Background(), &request)
	require.NoError(t, err)
	t.Log("the number of bonded validators :", len(resp.Validators))

	consAddr, err := resp.Validators[0].GetConsAddr() //expecting cryptotypes.PubKey, got <nil>: invalid type
	t.Log("consaddr :", consAddr)                     // ""

	consPubkey, err := resp.Validators[0].ConsPubKey() //nil
	t.Log("consPubkey:", consPubkey)                   // <nil>

	tmConsPublickey, err := resp.Validators[0].TmConsPublicKey() // nil because validator.ConsPubkey is nil
	t.Log("tmConsPublickey:", tmConsPublickey)                   // {<nil>}

	var pubkey cryptotypes.PubKey
	err = custom.AppCodec.UnpackAny(resp.Validators[0].ConsensusPubkey, &pubkey)
	require.NoError(t, err)

	valconspub_correct, err := sdktypes.Bech32ifyPubKey(sdktypes.Bech32PubKeyTypeConsPub, pubkey)
	require.NoError(t, err)
	t.Log("valconpub", valconspub_correct) //cosmosvalconspub1zcjduepqhv5hmywmedf2j8jpdm2xl9ssyyq0nqf7ak24nex9law4dqtx8drq0xn67q

	ed25519pub, ok := pubkey.(*ed25519.PubKey)
	require.Equal(t, true, ok)
	pb, err := custom.AppCodec.MarshalBinaryBare(ed25519pub)
	require.NoError(t, err)
	valconspub_incorrect1, err := bech32.ConvertAndEncode(sdktypes.Bech32PrefixConsPub, pb)
	require.NoError(t, err)
	t.Log("valconpub1", valconspub_incorrect1) //cosmosvalconspub1pgstk2taj8duk54freqka4r0jcgzzq8esylwm92eunzl7h2ks9nrk3srm3pxx

	valconspub_incorrect2, err := bech32.ConvertAndEncode(sdktypes.Bech32PrefixConsPub, pubkey.Bytes())
	t.Log("valconpub2", valconspub_incorrect2) //cosmosvalconspub1hv5hmywmedf2j8jpdm2xl9ssyyq0nqf7ak24nex9law4dqtx8drqq729uc

	consAddress := sdktypes.ConsAddress(pubkey.Address())
	t.Log("consAddress:", consAddress)
	t.Log("consAddress(string):", consAddress.String()) //types.ConsAddress(pubkey.Bytes()))

	power := resp.Validators[0].ConsensusPower()
	t.Log("power:", power)
	req := stakingtypes.QueryValidatorRequest{
		ValidatorAddr: "cosmosvaloper1x5wgh6vwye60wv3dtshs9dmqggwfx2ldk5cvqu",
	}
	result, err := queryClient.Validator(context.Background(), &req)
	cosmostation := result.Validator

	t.Log(cosmostation.ConsensusPower())
	// sdktypes.TokensToConsensusPower(cosmostation.Tokens.Add(delegated.Tokens))

	// request = stakingtypes.QueryValidatorsRequest{Status: unbonded}
	// resp, err = queryClient.Validators(context.Background(), &request)
	// require.NoError(t, err)
	// log.Println("unbonded :", len(resp.Validators))

	// request = stakingtypes.QueryValidatorsRequest{Status: unbonding}
	// resp, err = queryClient.Validators(context.Background(), &request)
	// require.NoError(t, err)
	// log.Println("unboding :", len(resp.Validators))
}
