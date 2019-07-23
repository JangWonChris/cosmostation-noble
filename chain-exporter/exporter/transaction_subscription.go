package exporter

// Handle new event data
// func (ces *ChainExporterService) handleNewTransactions(eventData ctypes.ResultEvent) error {
// 	var eventDataTx tmtypes.EventDataTx
// 	bytes, _ := ces.Codec.MarshalJSON(eventData.Data)   // // MarshalJSON to []byte format
// 	err := ces.Codec.UnmarshalJSON(bytes, &eventDataTx) // []byte로 들어와야 해서 위에 MarshalJSON을 한번 해줘야 된다
// 	if err != nil {
// 		return errors.New("UnmarshalJSON cannot decode eventData bytes")
// 	}

// 	// Tx hash
// 	txHash := hex.EncodeToString(eventDataTx.Tx.Hash())
// 	fmt.Println("txHash: ", txHash)

// 	// msg_index type에 에러가 발생한다. 하지만 msg_index를 사용하지 않을 거기 때문에 일단 패쓰
// 	// err : json: cannot unmarshal string into Go struct field ABCIMessageLog.msg_index of type int
// 	logs, _ := sdk.ParseABCILogs(eventDataTx.TxResult.Result.Log)

// 	// Handle txs that are succesfully included in a block
// 	for _, log := range logs {
// 		if log.Success == true {
// 			var stdTx auth.StdTx
// 			err = ces.Codec.UnmarshalBinaryLengthPrefixed(eventDataTx.Tx, &stdTx) // Tx{} Prefix 포함하고 있기 때문에 UnmarshalBinaryLengthPrefixed 사용
// 			if err != nil {
// 				return errors.New("UnmarshalJSON cannot decode eventDataTx.Tx bytes")
// 			}

// 			// Handle standard transaction's messages
// 			for _, msg := range stdTx.Msgs {
// 				switch msg.Type() {
// 				case "send":
// 					var sendTx bank.MsgSend
// 					_ = ces.Codec.UnmarshalJSON(msg.GetSignBytes(), &sendTx)

// 					// Convert to bech32 cosmos address format
// 					fromAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixAccAddr, sendTx.FromAddress)
// 					toAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixAccAddr, sendTx.ToAddress)

// 					fmt.Println("fromAddress: ", fromAddress)
// 					fmt.Println("toAddress: ", toAddress)
// 				default:
// 					return nil
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }
