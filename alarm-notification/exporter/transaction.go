package exporter

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/tendermint/tendermint/libs/bech32"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func (ces *ChainExporterService) startSubscription() {
	for {
		fmt.Println("start - subscribe tx from full node")
		select {
		case eventData, ok := <-ces.wsOut:
			if ok {
				ces.handleNewEventData(eventData) // returns Data, Query, Tags
			}
		}
		fmt.Println("finish - subscribe tx from full node")
		fmt.Println("")
	}
}

// Handle new event data
func (ces *ChainExporterService) handleNewEventData(eventData ctypes.ResultEvent) error {

	// MarshalJSON to []byte format
	bytes, _ := ces.codec.MarshalJSON(eventData.Data)

	var eventDataTx tmtypes.EventDataTx
	err := ces.codec.UnmarshalJSON(bytes, &eventDataTx) // []byte로 들어와야 해서 위에 MarshalJSON을 한번 해줘야 된다
	if err != nil {
		return errors.New("UnmarshalJSON cannot decode eventData bytes")
	}

	// Tx hash
	txHash := hex.EncodeToString(eventDataTx.Tx.Hash())

	// msg_index type에 에러가 발생한다. 하지만 msg_index를 사용하지 않을 거기 때문에 일단 패쓰
	// err : json: cannot unmarshal string into Go struct field ABCIMessageLog.msg_index of type int
	logs, _ := sdk.ParseABCILogs(eventDataTx.TxResult.Result.Log)

	// Handle txs that are succesfully included in a block
	for _, log := range logs {
		if log.Success == true {
			var stdTx auth.StdTx
			err = ces.codec.UnmarshalBinaryLengthPrefixed(eventDataTx.Tx, &stdTx) // Tx{} Prefix 포함하고 있기 때문에 UnmarshalBinaryLengthPrefixed 사용
			if err != nil {
				fmt.Println("UnmarshalJSON eventDataTx.Tx error: ", err)
			}

			for _, msg := range stdTx.Msgs {
				switch msg.Type() {
				// fmt.Println("msg: ", msg)
				case "send":
					var sendTx bank.MsgSend
					err = ces.codec.UnmarshalJSON(msg.GetSignBytes(), &sendTx)
					if err != nil {
						return errors.New("UnmarshalJSON cannot decode MsgSend bytes")
					}

					// Convert to bech32 cosmos address format
					fromAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixAccAddr, sendTx.FromAddress)
					toAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixAccAddr, sendTx.ToAddress)

					fmt.Println("=======================================[send]")
					fmt.Println("height: ", eventDataTx.Height)
					fmt.Println("txHash: ", strings.ToUpper(txHash))
					fmt.Println("fromAddress: ", fromAddress)
					fmt.Println("toAddress: ", toAddress)
					fmt.Println("amount: ", sendTx.Amount)
					fmt.Println("=======================================")

				case "multisend":
					var multiSendTx bank.MsgMultiSend
					err = ces.codec.UnmarshalJSON(msg.GetSignBytes(), &multiSendTx)
					if err != nil {
						fmt.Println("Unmarshal MsgMultiSend JSON Error: ", err)
					}

					fmt.Println("=======================================[multisend]")
					fmt.Println(multiSendTx.Inputs)
					fmt.Println(multiSendTx.Outputs)
					fmt.Println("=======================================")

				default:
					fmt.Println("")
				}
			}

			// fmt.Println("Fee: ", stdTx.Fee)
			// fmt.Println("Signatures: ", stdTx.Signatures)
			// fmt.Println("Memo: ", stdTx.Memo)
			// fmt.Println("")
			// fmt.Println(bz)
		}
	}

	return nil
}
