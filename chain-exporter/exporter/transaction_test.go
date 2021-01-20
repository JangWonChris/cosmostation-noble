package exporter

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	//internal
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/custom"

	//cosmos-sdk
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdktypestx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/stretchr/testify/require"
)

// TestGetTxsChunk decodes transactions in a block and return a format of database transaction.
func TestGetTxsChunk(t *testing.T) {
	require.NotNil(t, ex.client)
	// 13030, 272247
	// 122499 (multi msg type)
	block, err := ex.client.RPC.GetBlock(13030)
	if err != nil {
		log.Println(err)
	}
	txResps, err := ex.client.CliCtx.GetTxs(block)
	if err != nil {
		log.Println(err)
	}

	tms, err := ex.transactionAccount(block.Block.ChainID, txResps)
	log.Println(tms)
	return

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
		chunk, err := custom.AppCodec.MarshalJSON(txResp)
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
		err := custom.AppCodec.UnmarshalJSON([]byte(js), &txResps[i])
		if err != nil {
			log.Println(err)
			return err
		}
		// show result
		fmt.Println("decode:", txResps[i].String())
	}

	return nil
}

func TestGetMessage(t *testing.T) {
	// 13030, 272247
	// 122499 (multi msg type)
	block, err := ex.client.RPC.GetBlock(122499)
	if err != nil {
		log.Println(err)
	}
	txResps, err := ex.client.CliCtx.GetTxs(block)
	if err != nil {
		log.Println(err)
	}

	for _, txResp := range txResps {
		txI := txResp.GetTx()
		tx, ok := txI.(*sdktypestx.Tx)
		if !ok {
			return
		}
		getMessages := tx.GetBody().GetMessages()
		msgjson := make([]json.RawMessage, len(getMessages), len(getMessages))
		var err error
		for i, msg := range getMessages {
			msgjson[i], err = custom.AppCodec.MarshalJSON(msg)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		jsonraws, err := json.Marshal(msgjson)
		fmt.Println(string(jsonraws))
	}

	return
}

func TestUnmarshalMessageString(t *testing.T) {
	msgStr := "[{\"@type\": \"/cosmos.staking.v1beta1.MsgDelegate\", \"amount\": {\"denom\": \"umuon\", \"amount\": \"18044801\"}, \"delegator_address\": \"cosmos10fyfu7fl78f88a7zhcwu72wk3hjlzdm83yr09k\", \"validator_address\": \"cosmosvaloper10fyfu7fl78f88a7zhcwu72wk3hjlzdm85sh6f9\"}]"

	var jsonRaws []json.RawMessage
	json.Unmarshal([]byte(msgStr), &jsonRaws)

	for _, raw := range jsonRaws {
		fmt.Println(string(raw))
		var any codectypes.Any
		custom.AppCodec.UnmarshalJSON(raw, &any)
		fmt.Println(any.TypeUrl)
		// any.GetCachedValue().(type)
		fmt.Println(any.GetCachedValue())
		b, err := json.Marshal(any)
		require.NoError(t, err)
		fmt.Println(string(any.Value))

		fmt.Println(string(b))
	}

}
