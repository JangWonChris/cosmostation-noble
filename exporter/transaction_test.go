package exporter

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	//internal
	"github.com/cosmostation/cosmostation-cosmos/custom"
	mbltypes "github.com/cosmostation/mintscan-backend-library/types"

	//cosmos-sdk
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdktypestx "github.com/cosmos/cosmos-sdk/types/tx"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/stretchr/testify/require"
)

// TestGetTxsChunk decodes transactions in a block and return a format of database transaction.
func TestGetTxsChunk(t *testing.T) {
	require.NotNil(t, ex.Client)
	// 13030, 272247
	// 122499 (multi msg type)
	block, err := ex.Client.RPC.GetBlock(13030)
	if err != nil {
		log.Println(err)
	}
	txResps, err := ex.Client.CliCtx.GetTxs(block)
	if err != nil {
		log.Println(err)
	}

	tma := ex.disassembleTransaction(txResps)
	log.Println(tma)

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
	block, err := ex.Client.RPC.GetBlock(970957)
	if err != nil {
		t.Log(err)
	}
	txResps, err := ex.Client.CliCtx.GetTxs(block)
	if err != nil {
		t.Log(err)
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
				t.Log(err)
				return
			}
		}
		jsonraws, err := json.Marshal(msgjson)
		t.Log(string(jsonraws))
	}

	return
}

func TestUnmarshalMessageString(t *testing.T) {
	msgStr := "[{\"@type\": \"/cosmos.staking.v1beta1.MsgDelegate\", \"amount\": {\"denom\": \"umuon\", \"amount\": \"18044801\"}, \"delegator_address\": \"cosmos10fyfu7fl78f88a7zhcwu72wk3hjlzdm83yr09k\", \"validator_address\": \"cosmosvaloper10fyfu7fl78f88a7zhcwu72wk3hjlzdm85sh6f9\"}]"

	var jsonRaws []json.RawMessage
	json.Unmarshal([]byte(msgStr), &jsonRaws)

	for _, raw := range jsonRaws {
		t.Log(string(raw))
		var any codectypes.Any
		custom.AppCodec.UnmarshalJSON(raw, &any)
		t.Log(any.TypeUrl)
		// any.GetCachedValue().(type)
		t.Log(any.GetCachedValue())
		b, err := json.Marshal(any)
		require.NoError(t, err)
		t.Log(string(any.Value))

		t.Log(string(b))
	}

}

func TestAuthz(t *testing.T) {
	// 13030, 272247
	// 122499 (multi msg type)

	txs := []string{
		// "A49B11B0DFA2C876533C373BA23687E3033B3AA825BA6D3DCC1682FFBF16163F", //delegate
		"9862D3741C88E7A9A54CB8B2CAC63DC8D7D06BC50A60AAA63B3AC24111F2DE4D", // vote
	}
	// block, err := ex.Client.RPC.GetBlock(2390950)
	// if err != nil {
	// 	t.Log(err)
	// }
	for k := range txs {
		txResp, err := ex.Client.CliCtx.GetTx(txs[k])
		require.NoError(t, err)

		msgs := txResp.GetTx().GetMsgs()

		for i := range msgs {

			switch m := msgs[i].(type) {
			case *authztypes.MsgExec:
				for j := range m.Msgs {
					var msgExecAuthorized sdktypes.Msg
					custom.AppCodec.UnpackAny(m.Msgs[j], &msgExecAuthorized)

					t.Log(">>>>>>>>>>>>> ", msgExecAuthorized)

					msgType, accounts := mbltypes.AccountExporterFromCosmosTxMsg(&msgExecAuthorized)
					t.Log("msg detail :", msgType, accounts)
					t.Log("vote weight: ", sdktypes.OneDec().String())

				}
			default:
				t.Log("not authz :", m)
			}

		}

	}
}
