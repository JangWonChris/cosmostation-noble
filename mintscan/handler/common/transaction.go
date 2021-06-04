package common

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"

	mdschema "github.com/cosmostation/mintscan-database/schema"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

// GetTransactions returns transactions with given parameters.
func GetTransactions(rw http.ResponseWriter, r *http.Request) {
	from, limit, err := model.ParseHTTPArgs(r)
	if err != nil {
		zap.S().Debug("failed to parse HTTP args ", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "request is invalid")
		return
	}

	// if limit > 100 {
	// 	zap.L().Debug("request is over max limit ", zap.Int("request limit", limit))
	// 	errors.ErrOverMaxLimit(rw, http.StatusUnauthorized)
	// 	return
	// }

	txs, err := s.DB.QueryTransactions(from, limit)
	if err != nil {
		zap.L().Error("failed to query txs", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	if len(txs) <= 0 {
		zap.L().Debug("found no transactions in database")
		model.Respond(rw, []mdschema.Transaction{})
		return
	}

	result := make([]*model.ResultTx, 0)

	result = model.ParseTransactions(txs)

	model.Respond(rw, result)
	return
}

func GetTransactionsList(rw http.ResponseWriter, r *http.Request) {
	var reqTxs model.TxList
	var respTxs model.TxList

	if err := json.NewDecoder(r.Body).Decode(&reqTxs); err != nil {
		errors.ErrInvalidFormat(rw, http.StatusBadRequest)
		zap.L().Debug("failed to decode tx list", zap.Error(err))
		return
	}

	if len(reqTxs.TxHash) == 0 {
		errors.ErrInvalidFormat(rw, http.StatusBadRequest)
		zap.L().Debug("received empty tx hash list")
		return
	}
	for i := range reqTxs.TxHash {

		if strings.Contains(reqTxs.TxHash[i], "0x") {
			reqTxs.TxHash[i] = reqTxs.TxHash[i][2:]
		}
		if len(reqTxs.TxHash[i]) != 64 {
			zap.L().Debug("tx hash length is invalid", zap.String("txHashStr", reqTxs.TxHash[i]))
			continue
		}
		reqTxs.TxHash[i] = strings.ToUpper(reqTxs.TxHash[i])
		respTxs.TxHash = append(respTxs.TxHash, reqTxs.TxHash[i])
	}

	txs, err := s.DB.QueryTransactionByTxHashes(respTxs.TxHash)
	if err != nil {
		errors.ErrNotFound(rw, http.StatusNotFound)
		zap.L().Error("failed to get transactions by tx hashes", zap.Error(err))
		return
	}

	model.ParseTransactions(txs)
	result := model.ParseTransactions(txs)

	model.Respond(rw, result)
	return
}

// GetTransaction receives transaction hash and returns that transaction
func GetTransactionByID(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Query transction by transaction id if the request param has id; otherwise query with transaction hash.
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		errors.ErrFailedConversion(rw, http.StatusInternalServerError)
		return
	}

	tx, err := s.DB.QueryTransactionByID(id)
	if err != nil {
		zap.L().Error("failed to get transaction by tx id", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusInternalServerError)
		return
	}

	result := model.ParseTransaction(tx)

	model.Respond(rw, result)
	return
}

// GetTransaction receives transaction hash and returns that transaction
func GetTransactionByHash(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hashStr := vars["hash"]

	hashStr = strings.ToUpper(hashStr)

	if strings.Contains(hashStr, "0x") {
		hashStr = hashStr[2:]
	}

	if len(hashStr) != 64 {
		zap.L().Debug("tx hash length is invalid", zap.String("txHashStr", hashStr))
		errors.ErrInvalidFormat(rw, http.StatusBadRequest)
		return
	}

	tx, err := s.DB.QueryTransactionByTxHash(hashStr)
	if err != nil {
		zap.L().Error("failed to get transaction by tx hash", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusInternalServerError)
		return
	}

	result := model.ParseTransaction(tx)

	model.Respond(rw, result)
	return
}

// BroadcastTx receives signed transaction and broadcast it to the active network.
func BroadcastTx(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	signedTx := vars["signed_tx"]

	if strings.Contains(signedTx, "0x") {
		signedTx = signedTx[2:]
	}
	// jeonghwan 오류
	result, err := s.Client.CliCtx.BroadcastTx([]byte(signedTx))
	if err != nil {
		zap.L().Error("failed to broadcast transaction", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	model.Respond(rw, result)
	return
}
