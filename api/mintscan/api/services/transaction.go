package services

import (
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	ctypes "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/sync"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/gorilla/mux"
	resty "gopkg.in/resty.v1"

	"github.com/go-pg/pg"

	"github.com/tendermint/tendermint/rpc/client"
)

// GetTxs returns latest transactions
func GetTxs(Codec *codec.Codec, RPCClient *client.HTTP, DB *pg.DB, w http.ResponseWriter, r *http.Request) error {
	// Declare default variables
	limit := int(10)
	from := int(1)

	// Check limit param
	if len(r.URL.Query()["limit"]) > 0 {
		limit, _ = strconv.Atoi(r.URL.Query()["limit"][0])
	}

	// Max limit
	if limit > 20 {
		errors.ErrOverMaxLimit(w, http.StatusRequestEntityTooLarge)
		return nil
	}

	// Check from param
	if len(r.URL.Query()["from"]) > 0 {
		from, _ = strconv.Atoi(r.URL.Query()["from"][0])
	} else {
		// Check current height in db
		var blocks []ctypes.BlockInfo
		_ = DB.Model(&blocks).
			Order("height DESC").
			Limit(1).
			Select()
		if len(blocks) > 0 {
			from = int(blocks[0].Height)
		}
	}

	// Query a number of txs
	transactionInfos := make([]*ctypes.TransactionInfo, 0)
	_ = DB.Model(&transactionInfos).
		Where("height <= ?", from).
		Limit(limit).
		Order("height DESC").
		Select()

	// Check if any transaction exists
	if len(transactionInfos) <= 0 {
		return json.NewEncoder(w).Encode(transactionInfos)
	}

	resultTransactionInfo := make([]*models.ResultTransactionInfo, 0)
	for _, transactionInfo := range transactionInfos {
		tempResultTransactionInfo := &models.ResultTransactionInfo{
			Height: transactionInfo.Height,
			TxHash: transactionInfo.TxHash,
			Time:   transactionInfo.Time,
		}
		resultTransactionInfo = append(resultTransactionInfo, tempResultTransactionInfo)
	}
	return json.NewEncoder(w).Encode(resultTransactionInfo)
}

// GetTx receives transaction hash and returns that transaction
func GetTx(Codec *codec.Codec, RPCClient *client.HTTP, DB *pg.DB, Config *config.Config, w http.ResponseWriter, r *http.Request) error {
	// Get transaction hash
	vars := mux.Vars(r)
	txHexStr := vars["hash"]

	// If txHexStr contains 0x, remove it
	if strings.Contains(txHexStr, "0x") {
		txHexStr = txHexStr[2:]
	}

	// Transaction length check
	if len(txHexStr) != 64 {
		errors.ErrInvalidFormat(w, http.StatusBadRequest)
	}

	// Query LCD
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, _ := resty.R().Get(Config.Node.LCDURL + "/txs/" + txHexStr)

	// Parse struct
	var generalTx models.GeneralTx
	err := json.Unmarshal(resp.Body(), &generalTx)
	if err != nil {
		fmt.Printf("GeneralTx unmarshal error - %v\n", err)
	}

	return json.NewEncoder(w).Encode(generalTx)
}

// BroadcastTx sends transaction
func BroadcastTx(Codec *codec.Codec, RPCClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	// Get transaction hash
	vars := mux.Vars(r)
	txHexStr := vars["hash"]

	// If txHexStr contains 0x, remove it
	if strings.Contains(txHexStr, "0x") {
		txHexStr = txHexStr[2:]
	}

	// Convert from hexadecimal string to bytes
	txByteStr, err := hex.DecodeString(txHexStr)
	if err != nil {
		errors.ErrFailedConversion(w, http.StatusBadRequest)
		return nil
	}

	// Unmarshalling JSON
	var stdTx auth.StdTx
	err = Codec.UnmarshalJSON(txByteStr, &stdTx)
	if err != nil {
		errors.ErrFailedUnmarshalJSON(w, http.StatusBadRequest)
		return nil
	}

	// Marshalling binary length prefix
	bz, err := Codec.MarshalBinaryLengthPrefixed(stdTx)
	if err != nil {
		errors.ErrFailedMarshalBinaryLengthPrefixed(w, http.StatusBadRequest)
		return nil
	}

	// Broadcast transaction
	result, err := RPCClient.BroadcastTxCommit(bz)
	if err != nil {
		return nil
	}

	return json.NewEncoder(w).Encode(result)
}
