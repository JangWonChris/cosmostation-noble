package db

import (
	"encoding/json"
	"os"
	"testing"

	//mbl
	"github.com/cosmostation/cosmostation-cosmos/custom"
	mblconfig "github.com/cosmostation/mintscan-backend-library/config"
	mdschema "github.com/cosmostation/mintscan-database/schema"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ibcchanneltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"

	pg "github.com/go-pg/pg/v10"

	"github.com/stretchr/testify/require"
)

var db *Database

func TestMain(m *testing.M) {
	// types.SetAppConfig()

	fileBaseName := "mintscan"
	cfg := mblconfig.ParseConfig(fileBaseName)
	db = Connect(&cfg.DB)

	os.Exit(m.Run())
}

func TestInsertOrUpdate(t *testing.T) {
	err := db.Ping()
	require.NoError(t, err)

}

func TestUpdate_Validator(t *testing.T) {
	err := db.Ping()
	require.NoError(t, err)

	val := &mdschema.Validator{
		Address: "kava1ulzzxuvghlv04sglkzyxv94rvl7c2llhs098ju",
		Rank:    5,
	}

	validator, err := db.GetValidatorByAnyAddr(val.Address)
	require.NoError(t, err)

	result, err := db.Model(&validator).
		Set("rank = ?", val.Rank).
		Where("id = ?", validator.ID).
		Update()

	require.NoError(t, err)
	require.Equal(t, 1, result.RowsAffected())
}

func TestConnection(t *testing.T) {
	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	require.NoError(t, err)

	require.Equal(t, n, 1, "failed to ping database")
}

func TestGetTx(t *testing.T) {
	id := int64(4318956)
	txs, err := db.GetTransactionByID(id)
	require.NoError(t, err)

	t.Log(txs.Hash)

	unmarshaler := custom.EncodingConfig.Marshaler.UnmarshalJSON
	var txResp sdktypes.TxResponse
	err = unmarshaler(txs.Chunk, &txResp)
	require.NoError(t, err)

	tx := txResp.GetTx()
	msgs := tx.GetMsgs()

	type IBCRecvPacketData struct {
		Denom    string       `json:"denom,omitempty"`
		Amount   sdktypes.Int `json:"amount"`
		Sender   string       `json:"sender,omitempty"`
		Receiver string       `json:"receiver,omitempty"`
	}

	var pd IBCRecvPacketData

	for _, msg := range msgs {
		switch m := msg.(type) {
		case *ibcchanneltypes.MsgRecvPacket:
			t.Log(m.Packet)
			t.Log(string(m.Packet.GetData()))
			json.Unmarshal(m.Packet.GetData(), &pd)
			t.Log("sender :", pd.Sender)
			t.Log("receiver :", pd.Receiver)
			t.Log("denom :", pd.Denom)
			t.Log("amount :", pd.Amount)
		}
	}

}
