package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/schema"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

// GetTransactions returns transactions with given parameters.
func GetTransactions(rw http.ResponseWriter, r *http.Request) {
	before, after, limit, err := model.ParseHTTPArgsWithBeforeAfterLimit(r, model.DefaultBefore, model.DefaultAfter, model.DefaultLimit)
	if err != nil {
		zap.S().Debug("failed to parse HTTP args ", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "request is invalid")
		return
	}

	if limit > 100 {
		zap.L().Debug("request is over max limit ", zap.Int("request limit", limit))
		errors.ErrOverMaxLimit(rw, http.StatusUnauthorized)
		return
	}

	txs, err := s.db.QueryTransactions(before, after, limit)
	if err != nil {
		zap.L().Error("failed to query txs", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	if len(txs) <= 0 {
		zap.L().Debug("found no transactions in database")
		model.Respond(rw, []schema.Transaction{})
		return
	}

	result := make([]*model.ResultTx, 0)

	for _, tx := range txs {
		t := &model.ResultTx{
			ID:        tx.ID,
			Height:    tx.Height,
			TxHash:    tx.TxHash,
			Memo:      tx.Memo,
			Timestamp: tx.Timestamp,
		}

		result = append(result, t)
	}

	model.Respond(rw, result)
	return
}

// GetTransactionsList returns a array of transaction details for each transaction hash including request body.
func GetTransactionsList(rw http.ResponseWriter, r *http.Request) {
	var txList model.TxList

	if err := json.NewDecoder(r.Body).Decode(&txList); err != nil {
		errors.ErrInvalidFormat(rw, http.StatusBadRequest)
		zap.L().Debug("failed to decode tx list", zap.Error(err))
		return
	}

	txResp := make([]sdk.TxResponse, len(txList.TxHash))

	if len(txResp) == 0 {
		errors.ErrInvalidFormat(rw, http.StatusBadRequest)
		zap.L().Debug("received empty tx hash list", zap.Int("len", len(txResp)))
		return
	}

	// Remove if tx hash contains prefix '0x' and check length.
	for i, txHashStr := range txList.TxHash {
		if strings.Contains(txHashStr, "0x") {
			txHashStr = txHashStr[2:]
		}

		if len(txHashStr) != 64 {
			zap.L().Debug("tx hash length is invalid", zap.Int("len", len(txHashStr)), zap.String("txHashStr", txHashStr))
			continue
		}

		err := s.client.GetTxs(txHashStr, &txResp[i])
		if err != nil {
			zap.L().Error("failed to get transaction details", zap.Error(err))
			continue
		}
	}

	model.Respond(rw, txResp)
	return
}

// GetTransaction receives transaction hash and returns that transaction
func GetTransaction(rw http.ResponseWriter, r *http.Request) {
	var txID int64
	var err error

	id := r.FormValue("id")
	txHashStr := r.FormValue("hash")

	// Request param must have either transaction id or transaction hash.
	if id == "" && txHashStr == "" {
		errors.ErrRequiredParam(rw, http.StatusBadRequest, "request must have either transaction id or hash")
		return
	}

	// Query transction by transaction id if the request param has id; otherwise query with transaction hash.
	if id != "" {
		txID, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			errors.ErrFailedConversion(rw, http.StatusInternalServerError)
			return
		}

		tx, err := s.db.QueryTransactionByID(txID)
		if err != nil {
			zap.L().Error("failed to get transaction by tx id", zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusInternalServerError)
			return
		}

		result, _ := model.ParseTransaction(tx)

		model.Respond(rw, result)
		return
	}

	txHashStr = strings.ToUpper(txHashStr)

	if strings.Contains(txHashStr, "0x") {
		txHashStr = txHashStr[2:]
	}

	if len(txHashStr) != 64 {
		zap.L().Debug("tx hash length is invalid", zap.String("txHashStr", txHashStr))
		errors.ErrInvalidFormat(rw, http.StatusBadRequest)
		return
	}

	tx, err := s.db.QueryTransactionByTxHash(txHashStr)
	if err != nil {
		zap.L().Error("failed to get transaction by tx hash", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusInternalServerError)
		return
	}

	result, _ := model.ParseTransaction(tx)

	model.Respond(rw, result)
	return
}

// GetLegacyTransactionFromDB receives transaction hash and returns that transaction
func GetLegacyTransactionFromDB(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txHashStr := vars["hash"]
	var err error

	// Request param must have either transaction id or transaction hash.
	if txHashStr == "" {
		errors.ErrRequiredParam(rw, http.StatusBadRequest, "request must have either transaction id or hash")
		return
	}

	if strings.Contains(txHashStr, "0x") {
		txHashStr = txHashStr[2:]
	}

	if len(txHashStr) != 64 {
		zap.L().Debug("tx hash length is invalid", zap.String("txHashStr", txHashStr))
		errors.ErrInvalidFormat(rw, http.StatusBadRequest)
		return
	}

	tx, err := s.db.QueryTransactionByTxHash(txHashStr)
	if err != nil {
		zap.L().Error("failed to get transaction by tx hash", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusInternalServerError)
		return
	}

	result, _ := model.ParseTransaction(tx)

	model.Respond(rw, result)
	return
}

// GetLegacyTransaction uses RPC API to parse transaction and return.
// [NOT USED]
func GetLegacyTransaction(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txHashStr := vars["hash"]

	if strings.Contains(txHashStr, "0x") {
		txHashStr = txHashStr[2:]
	}

	if len(txHashStr) != 64 {
		zap.L().Debug("tx hash length is invalid", zap.String("txHashStr", txHashStr))
		errors.ErrInvalidFormat(rw, http.StatusBadRequest)
		return
	}

	resp, err := s.client.GetTx(txHashStr)
	if err != nil {
		zap.L().Error("failed to get tx hash info", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	rest.PostProcessResponseBare(rw, s.client.GetCliContext(), resp) // codec marshalling
	return
}

// BroadcastTx receives signed transaction and broadcast it to the active network.
func BroadcastTx(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	signedTx := vars["signed_tx"]

	if strings.Contains(signedTx, "0x") {
		signedTx = signedTx[2:]
	}

	result, err := s.client.BroadcastTx(signedTx)
	if err != nil {
		zap.L().Error("failed to broadcast transaction", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	model.Respond(rw, result)
	return
}
