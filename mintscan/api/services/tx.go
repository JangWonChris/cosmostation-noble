package services

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/db"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/schema"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/utils"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/gorilla/mux"
	resty "gopkg.in/resty.v1"

	"github.com/tendermint/tendermint/rpc/client"
)

// GetTxs returns latest transactions
func GetTxs(codec *codec.Codec, db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	limit := int(10)
	from := int(1)

	// check limit param
	if len(r.URL.Query()["limit"]) > 0 {
		limit, _ = strconv.Atoi(r.URL.Query()["limit"][0])
	}

	// check max limit
	if limit > 20 {
		errors.ErrOverMaxLimit(w, http.StatusRequestEntityTooLarge)
		return nil
	}

	// check from param
	if len(r.URL.Query()["from"]) > 0 {
		from, _ = strconv.Atoi(r.URL.Query()["from"][0])
	} else {
		// Check current height in db
		var blocks []schema.BlockInfo
		_ = db.Model(&blocks).
			Order("height DESC").
			Limit(1).
			Select()
		if len(blocks) > 0 {
			from = int(blocks[0].Height)
		}
	}

	// query a number of txs
	transactionInfos := make([]*schema.TransactionInfo, 0)
	_ = db.Model(&transactionInfos).
		Where("height <= ?", from).
		Limit(limit).
		Order("height DESC").
		Select()

	// check if any transaction exists
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

	utils.Respond(w, resultTransactionInfo)
	return nil
}

// GetTx receives transaction hash and returns that transaction
func GetTx(codec *codec.Codec, config *config.Config, db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	txHexStr := vars["hash"]

	// check if txHexStr contains 0x, remove it
	if strings.Contains(txHexStr, "0x") {
		txHexStr = txHexStr[2:]
	}

	// check tx length
	if len(txHexStr) != 64 {
		errors.ErrInvalidFormat(w, http.StatusBadRequest)
	}

	resp, _ := resty.R().Get(config.Node.LCDEndpoint + "/txs/" + txHexStr)

	var generalTx models.GeneralTx
	err := json.Unmarshal(resp.Body(), &generalTx)
	if err != nil {
		fmt.Printf("GeneralTx unmarshal error - %v\n", err)
	}

	utils.Respond(w, generalTx)
	return nil
}

// BroadcastTx sends transaction
func BroadcastTx(codec *codec.Codec, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	txHexStr := vars["hash"]

	// check if txHexStr contains 0x, remove it
	if strings.Contains(txHexStr, "0x") {
		txHexStr = txHexStr[2:]
	}

	// convert to bytes
	txByteStr, err := hex.DecodeString(txHexStr)
	if err != nil {
		errors.ErrFailedConversion(w, http.StatusBadRequest)
		return nil
	}

	var stdTx auth.StdTx
	err = codec.UnmarshalJSON(txByteStr, &stdTx)
	if err != nil {
		errors.ErrFailedUnmarshalJSON(w, http.StatusBadRequest)
		return nil
	}

	bz, err := codec.MarshalBinaryLengthPrefixed(stdTx)
	if err != nil {
		errors.ErrFailedMarshalBinaryLengthPrefixed(w, http.StatusBadRequest)
		return nil
	}

	result, err := rpcClient.BroadcastTxCommit(bz)
	if err != nil {
		return nil
	}

	utils.Respond(w, result)
	return nil
}
