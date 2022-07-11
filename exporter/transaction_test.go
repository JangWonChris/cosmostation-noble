package exporter

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	//internal
	"github.com/cosmostation/cosmostation-cosmos/custom"

	//cosmos-sdk
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdktypestx "github.com/cosmos/cosmos-sdk/types/tx"
	authvestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

func TestMap(t *testing.T) {
	m := make(map[string]struct{})

	key1 := ""
	key2 := "abcd"

	m[key1] = struct{}{}
	m[key2] = struct{}{}

	for k, v := range m {
		t.Log("key :", k, " value :", v)
	}
}
func TestAuthz(t *testing.T) {
	// 13030, 272247
	// 122499 (multi msg type)

	txs := []string{
		"FF939BEF6CE78ACF7D6AFFA1F81AACFA7531E27A5739C2E258A7BEEA30D7C8D6",
	}
	// block, err := ex.Client.RPC.GetBlock(2390950)
	// if err != nil {
	// 	t.Log(err)
	// }
	for k := range txs {
		txResp, err := ex.Client.CliCtx.GetTx(txs[k])
		require.NoError(t, err)

		// txResps, err := ex.Client.CliCtx.GetTxs(block)
		// if err != nil {
		// 	t.Log(err)
		// }

		msgs := txResp.GetTx().GetMsgs()

		for i := range msgs {

			switch m := msgs[i].(type) {
			case *authztypes.MsgExec:
				// t.Log(m.Msgs)
				for j := range m.Msgs {
					c := m.Msgs[j].GetCachedValue()
					// t.Log(m.Msgs[j].GetValue())
					switch im := (c).(type) {

					case *authvestingtypes.MsgCreateVestingAccount:
						t.Log(im.Type())

					// authz v0.44.0에서 추가
					case *authztypes.MsgGrant:
						t.Log(im.Type())
					case *authztypes.MsgRevoke:
						t.Log(im.Type())
					case *authztypes.MsgExec:
						t.Log(im.Type())

					//bank (2)
					case *banktypes.MsgSend:
						t.Log(im.Type())
						t.Log(im.Amount)
					case *banktypes.MsgMultiSend:
						t.Log(im.Type())

					//crisis (1)
					case *crisistypes.MsgVerifyInvariant:
						t.Log(im.Type())

					//distribution (4)
					case *distributiontypes.MsgSetWithdrawAddress:
						t.Log(im.Type())
					case *distributiontypes.MsgWithdrawDelegatorReward:
						t.Log(im.Type())
					case *distributiontypes.MsgWithdrawValidatorCommission:
						t.Log(im.Type())
					case *distributiontypes.MsgFundCommunityPool:
						t.Log(im.Type())

					//evidence (1)
					case *evidencetypes.MsgSubmitEvidence:
						t.Log(im.Type())

					// freegrant v0.44.0에서 추가
					case *feegrant.MsgGrantAllowance:
						t.Log(im.Type())
					case *feegrant.MsgRevokeAllowance:
						t.Log(im.Type())

					//gov (3)
					case *govtypes.MsgSubmitProposal:
						t.Log(im.Type())
					case *govtypes.MsgVote:
						t.Log(im.Type())
						t.Log(im.ProposalId)
						t.Log(im.Voter)
						t.Log(im.Option)
					case *govtypes.MsgDeposit:
						t.Log(im.Type())

					//slashing (1)
					case *slashingtypes.MsgUnjail:
						t.Log(im.Type())

					//staking (5)
					case *stakingtypes.MsgCreateValidator:
						t.Log(im.Type())
					case *stakingtypes.MsgEditValidator:
						t.Log(im.Type())
					case *stakingtypes.MsgDelegate:
						t.Log(im.Type())
					case *stakingtypes.MsgBeginRedelegate:
						t.Log(im.Type())
					case *stakingtypes.MsgUndelegate:
						t.Log(im.Type())
					default:
						t.Logf("not found : %v", im)

					}

				}
			default:
				t.Log("not authz :", m)
			}

		}

	}
}
func TestAuthz2(t *testing.T) {
	// 13030, 272247
	// 122499 (multi msg type)

	txs := []string{
		"FF939BEF6CE78ACF7D6AFFA1F81AACFA7531E27A5739C2E258A7BEEA30D7C8D6",
	}
	// block, err := ex.Client.RPC.GetBlock(2390950)
	// if err != nil {
	// 	t.Log(err)
	// }
	for k := range txs {
		txResp, err := ex.Client.CliCtx.GetTx(txs[k])
		require.NoError(t, err)

		// txResps, err := ex.Client.CliCtx.GetTxs(block)
		// if err != nil {
		// 	t.Log(err)
		// }

		msgs := txResp.GetTx().GetMsgs()

		for i := range msgs {

			switch m := msgs[i].(type) {
			case *authztypes.MsgExec:
				// t.Log(m.Msgs)
				for j := range m.Msgs {
					c := m.Msgs[j].GetCachedValue().(sdktypes.Tx)
					for _, msg := range c.GetMsgs() {
						mt, addrs := custom.AccountExporterFromCustomTxMsg(&msg, txResp.TxHash)
						t.Log(mt, addrs)
					}
				}
			default:
				t.Log("not authz :", m)
			}

		}

	}
}

type msgType string
type account string
type MsgParser interface {
	Parse() (msgType, account, error)
}
