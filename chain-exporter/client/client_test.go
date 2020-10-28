package client

import (
	"context"
	"log"
	"os"
	"testing"

	sdkcodec "github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdktypestx "github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/golang/protobuf/proto"
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

	address := "cosmos1x5wgh6vwye60wv3dtshs9dmqggwfx2ldnqvev0"
	// log.Println(cli.GetAccount("cosmos1x5wgh6vwye60wv3dtshs9dmqggwfx2ldnqvev0"))
	sdkaddr, err := sdktypes.AccAddressFromBech32(address)
	if err != nil {
		log.Println(err)
	}
	b := banktypes.NewQueryBalanceRequest(sdkaddr, "umuon")
	log.Println(b)
	bankClient := banktypes.NewQueryClient(cli.grpcClient)
	var header metadata.MD
	bankRes, err := bankClient.Balance(
		context.Background(),
		b,
		grpc.Header(&header), // Also fetch grpc header
	)
	log.Println(*bankRes.GetBalance())
	// blockHeight := header.Get(grpctypes.GRPCBlockHeightHeader)
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
	// z := txResp.Tx.(codectypes.IntoAny)

	log.Println("typeurl :", txResp.Tx.TypeUrl)
	a, err := cli.cliCtx.JSONMarshaler.MarshalJSON(txResp.Tx)
	if err != nil {
		log.Println(err)
	}

	anyb, err := sdkcodec.MarshalAny(codec.EncodingConfig.Marshaler, txResp.Tx)
	if err != nil {
		log.Println(err)
	}
	_ = anyb
	// log.Println("value :", string(txResp.Tx.Value))
	// log.Println("txResp.Tx :", txResp.Tx)
	// log.Println("anyb :", string(anyb))

	// var _ sdk.Msg = &banktypes.MsgSend{}
	var sdkmsg sdktypes.Msg
	var txtype sdktypestx.Tx
	_, _ = sdkmsg, txtype
	var bsend banktypes.MsgSend
	var i interface{}
	// if err := json.Unmarshal(a, &i); err != nil {
	// 	log.Println(err)
	// }
	s, ok := sdkmsg.(proto.Message)
	if !ok {
		log.Println("sdkmsg.(proto.Message) nok")
	}
	log.Println("protomessage s :", s)
	if err := sdkcodec.UnmarshalAny(codec.EncodingConfig.Marshaler, &sdkmsg, txResp.Tx.Value); err != nil {
		log.Println(err)
	}
	// k, ok := i.(sdktypestx.Tx)
	// if !ok {
	// 	log.Println("nok :", ok)
	// }
	switch z := i.(type) {
	case *banktypes.MsgSend:
		log.Println("type msgsend :", z)
	// case sdktypestx.Tx:
	// 	log.Println("type tx :", z)
	// case sdkcodectypes.Any:
	// 	log.Println("type tx :", z)
	default:
		log.Println("default :", z)
	}
	// log.Println(sdkmsg.Type())
	// sdkcodec.UnmarshalAny(codec.EncodingConfig.Marshaler, binf, a)
	log.Println("========================")
	log.Println("bsend :", bsend)
	// log.Println("kk :", k)
	// log.Println("ttt :", ttt)
	// log.Println(binf)
	log.Println("========================")
	// banktypes.MsgSend(txResp.Tx)
	// log.Println(cli.cliCtx.PrintOutput(txResp))
	log.Println("proto message name :", proto.MessageName(txResp.Tx))
	// log.Println("proto message type:", proto.MessageType(proto.MessageName(txResp)))
	// log.Println("proto message reflect:", proto.MessageReflect(txResp).Type())
	// log.Println("Height :", txResp.Height)
	log.Println("string(a) :", string(a))
	// log.Println("txresp.Tx : ", txResp.Tx)
	// "github.com/tendermint/tendermint/types" type.tx
}

func TestGetTx2(t *testing.T) {
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
