package exporter

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cosmostation/cosmostation-cosmos/api/wallet/config"

	gaiaApp "github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/go-pg/pg"
	"github.com/tendermint/tendermint/libs/bech32"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
)

// Monitor wraps Tendermint RPC client and PostgreSQL database
type ChainExporterService struct {
	cmn.BaseService
	// private fields
	Codec     *codec.Codec
	Config    *config.Config
	DB        *pg.DB
	WsCtx     context.Context
	WsOut     <-chan ctypes.ResultEvent
	RPCClient *client.HTTP
}

// Initializes all the required configs
func NewChainExporterService(config *config.Config) *ChainExporterService {
	ces := &ChainExporterService{
		// initialization
		Codec:     gaiaApp.MakeCodec(), // Register Cosmos SDK codecs
		Config:    config,
		WsCtx:     context.Background(),
		RPCClient: client.NewHTTP(config.Node.GaiadURL, "/websocket"), // Connect to Tendermint RPC client
	}

	// ces.RPCClient.SetLogger(logger)
	// ces.RPCClient.WSEvents.SetLogger(logger)

	// Register a service that can be started, stopped, and reset
	ces.BaseService = *cmn.NewBaseService(logger, "ChainExporterService", ces)

	return ces
}

// Override method for BaseService, which starts a service
func (ces *ChainExporterService) OnStart() error {
	// OnStart both service and rpc client
	ces.BaseService.OnStart()
	ces.RPCClient.OnStart()

	// initialize private fields and start subroutines, etc.
	// ces.WsOut, _ = ces.RPCClient.Subscribe(ces.WsCtx, "new block", "tm.event = 'NewBlock'", 1)
	ces.WsOut, _ = ces.RPCClient.Subscribe(ces.WsCtx, "new tx", "tm.event = 'Tx'", 1)

	// start routine
	for {
		fmt.Println("start - subscribe tx from full node")
		select {
		case eventData, ok := <-ces.WsOut:
			if ok {
				ces.handleNewEventData(eventData) // returns Data, Query, Tags
			}
		}
		fmt.Println("finish - subscribe tx from full node")
		fmt.Println("")
	}
}

// Override method for BaseService, which stops a service
func (ces *ChainExporterService) OnStop() {
	ces.BaseService.OnStop()
	ces.RPCClient.OnStop()
}

// Handle new event data
func (a *ChainExporterService) handleNewEventData(eventData ctypes.ResultEvent) error {
	// MarshalJSON to []byte format
	bytes, _ := a.Codec.MarshalJSON(eventData.Data)

	var eventDataTx tmtypes.EventDataTx
	err := a.Codec.UnmarshalJSON(bytes, &eventDataTx) // []byte로 들어와야 해서 위에 MarshalJSON을 한번 해줘야 된다
	if err != nil {
		return errors.New("UnmarshalJSON cannot decode eventData bytes")
	}

	// Tx hash
	txHash := hex.EncodeToString(eventDataTx.Tx.Hash())

	// fmt.Println("height: ", eventDataTx.Height)
	// fmt.Println("txHash: ", strings.ToUpper(txHash))
	// fmt.Println("eventDataTx.TxResult.Result: ", eventDataTx.TxResult.Result)
	// fmt.Println("eventDataTx.TxResult.Result.Log: ", eventDataTx.TxResult.Result.Log)
	// fmt.Println("")

	// msg_index type에 에러가 발생한다. 하지만 msg_index를 사용하지 않을 거기 때문에 일단 패쓰
	// err : json: cannot unmarshal string into Go struct field ABCIMessageLog.msg_index of type int
	logs, _ := sdk.ParseABCILogs(eventDataTx.TxResult.Result.Log)

	// Handle txs that are succesfully included in a block
	for _, log := range logs {
		if log.Success == true {
			var stdTx auth.StdTx
			err = a.Codec.UnmarshalBinaryLengthPrefixed(eventDataTx.Tx, &stdTx) // Tx{} Prefix 포함하고 있기 때문에 UnmarshalBinaryLengthPrefixed 사용
			if err != nil {
				return errors.New("UnmarshalJSON cannot decode eventDataTx.Tx bytes")
			}

			// Handle standard transaction's messages
			for _, msg := range stdTx.Msgs {
				switch msg.Type() {
				case "send":
					var sendTx bank.MsgSend
					_ = a.Codec.UnmarshalJSON(msg.GetSignBytes(), &sendTx)

					// Convert to bech32 cosmos address format
					fromAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixAccAddr, sendTx.FromAddress)
					toAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixAccAddr, sendTx.ToAddress)

					fmt.Println("=======================================[send]")
					fmt.Println("height: ", eventDataTx.Height)
					fmt.Println("txHash: ", strings.ToUpper(txHash))
					fmt.Println("fromAddress: ", fromAddress)
					fmt.Println("toAddress: ", toAddress)
					fmt.Println(sendTx.Amount)
					fmt.Println("=======================================")

				case "multisend":
					var multiSendTx bank.MsgMultiSend
					_ = a.Codec.UnmarshalJSON(msg.GetSignBytes(), &multiSendTx)

					fmt.Println("=======================================[multisend]")
					fmt.Println("height: ", eventDataTx.Height)
					fmt.Println("txHash: ", strings.ToUpper(txHash))
					fmt.Println(multiSendTx.Inputs)
					fmt.Println(multiSendTx.Outputs)
					fmt.Println("=======================================")

				default:
					return nil
				}
			}
		}
	}

	return nil
}
