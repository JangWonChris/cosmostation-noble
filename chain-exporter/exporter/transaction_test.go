package exporter

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	sdkcodec "github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdktypestx "github.com/cosmos/cosmos-sdk/types/tx"
	ceCodec "github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
)

// TestGetTxsChunk decodes transactions in a block and return a format of database transaction.
func TestGetTxsChunk(t *testing.T) {
	// 13030, 272247
	block, err := ex.client.GetBlock(13030)
	if err != nil {
		log.Println(err)
	}
	txResps, err := ex.client.GetTxs(block)
	if err != nil {
		log.Println(err)
	}

	// assume that following expression is for inserting db
	jsonString, err := InsertJSONStringToDB(txResps)
	if err != nil {
		log.Println(err)
		return
	}

	// decoding from db
	err = JSONStringUnmarshal(jsonString)
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func InsertJSONStringToDB(txResps []*sdktypes.TxResponse) ([]string, error) {
	jsonString := make([]string, len(txResps), len(txResps))
	for i, txResp := range txResps {
		chunk, err := ceCodec.AppCodec.MarshalJSON(txResp)
		if err != nil {
			log.Println(err)
		}
		jsonString[i] = string(chunk)
		// show result
		fmt.Println(jsonString[i])
	}

	return jsonString, nil
}

func JSONStringUnmarshal(jsonString []string) error {
	txResps := make([]sdktypes.TxResponse, len(jsonString), len(jsonString))
	for i, js := range jsonString {
		err := ceCodec.AppCodec.UnmarshalJSON([]byte(js), &txResps[i])
		if err != nil {
			log.Println(err)
			return err
		}
		// show result
		fmt.Println("decode:", txResps[i].String())
	}

	txI := txResps[0].GetTx()

	getMsgs := txI.GetMsgs()
	b, err := json.Marshal(getMsgs)
	if err != nil {
		return err
	}
	fmt.Println("GetMsgs.type():", getMsgs[0].Type())
	fmt.Println("GetMsgs():", txResps[0].GetTx().GetMsgs())
	fmt.Println("byte():", b)
	fmt.Println("string():", string(b))

	tx, ok := txI.(*sdktypestx.Tx)
	if !ok {
		return nil
	}
	getMessages := tx.GetBody().GetMessages()
	anyb, err := sdkcodec.MarshalAny(ceCodec.EncodingConfig.Marshaler, getMessages[0])
	if err != nil {
		return err
	}
	fmt.Println("anyb byte:", anyb)
	fmt.Println("anyb string:", string(anyb))
	return nil
}
