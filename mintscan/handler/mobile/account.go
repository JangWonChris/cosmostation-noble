package mobile

import (
	"context"
	"net/http"

	//internal
	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"

	//mbl
	"github.com/cosmostation/mintscan-backend-library/types"
	// "github.com/cosmostation/mintscan-backend-library/db"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

// GetAccountTxs returns transactions that are sent by an account
func GetAccountTxs(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["accAddr"]

	before, after, limit, err := model.ParseHTTPArgsWithBeforeAfterLimit(r, model.DefaultBefore, model.DefaultAfter, model.DefaultLimit)
	if err != nil {
		zap.S().Debug("failed to parse HTTP args ", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "request is invalid")
		return
	}

	if limit > 100 {
		zap.S().Debug("failed to query with this limit ", zap.Int("request limit", limit))
		errors.ErrOverMaxLimit(rw, http.StatusUnauthorized)
		return
	}

	err = types.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}

	valAddr, err := types.ConvertValAddrFromAccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to convert validator address from account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "validator address is invalid")
		return
	}

	// Query transactions that are made by the account.
	txs, err := s.DB.QueryTransactionsByAddr(accAddr, valAddr, before, after, limit)
	if err != nil {
		zap.L().Error("failed to query txs", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	if len(txs) <= 0 {
		model.Respond(rw, []model.ResultTx{})
		return
	}

	// result, err := model.ParseTransactions(txs)
	// if err != nil {
	// 	zap.L().Error("failed to parse txs", zap.Error(err))
	// 	errors.ErrInternalServer(rw, http.StatusInternalServerError)
	// 	return
	// }

	model.Respond(rw, txs)
	return
}

// GetAccountTransferTxs returns transfer txs (MsgSend and MsgMultiSend) that are sent by an account
func GetAccountTransferTxs(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["accAddr"]

	before, after, limit, err := model.ParseHTTPArgsWithBeforeAfterLimit(r, model.DefaultBefore, model.DefaultAfter, model.DefaultLimit)
	if err != nil {
		zap.S().Debug("failed to parse HTTP args ", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "request is invalid")
		return
	}

	if limit > 100 {
		zap.S().Debug("failed to query with this limit ", zap.Int("request limit", limit))
		errors.ErrOverMaxLimit(rw, http.StatusUnauthorized)
		return
	}

	var denom string

	if len(r.URL.Query()["denom"]) > 0 {
		denom = r.URL.Query()["denom"][0]
	}

	if denom == "" {
		denom, err = s.Client.GRPC.GetBondDenom(context.Background())
		if err != nil {
			return
		}
	}

	err = types.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}

	txs, err := s.DB.QueryTransferTransactionsByAddr(accAddr, denom, before, after, limit)
	if err != nil {
		zap.L().Error("failed to query txs", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	// result, err := model.ParseTransactions(txs)
	// if err != nil {
	// 	zap.L().Error("failed to parse txs", zap.Error(err))
	// 	errors.ErrInternalServer(rw, http.StatusInternalServerError)
	// 	return
	// }

	model.Respond(rw, txs)
	return
}

// GetTxsBetweenDelegatorAndValidator returns transactions that are made between an account and his delegated validator
func GetTxsBetweenDelegatorAndValidator(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["accAddr"]
	valAddr := vars["valAddr"]

	before, after, limit, err := model.ParseHTTPArgsWithBeforeAfterLimit(r, model.DefaultBefore, model.DefaultAfter, model.DefaultLimit)
	if err != nil {
		zap.S().Debug("failed to parse HTTP args ", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "request is invalid")
		return
	}

	if limit > 100 {
		zap.S().Debug("failed to query with this limit ", zap.Int("request limit", limit))
		errors.ErrOverMaxLimit(rw, http.StatusUnauthorized)
		return
	}

	err = types.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}

	err = types.VerifyBech32ValAddr(valAddr)
	if err != nil {
		zap.L().Debug("failed to validate validator address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "validator address is invalid")
		return
	}

	txs, err := s.DB.QueryTransactionsBetweenAccountAndValidator(accAddr, valAddr, before, after, limit)
	if err != nil {
		zap.L().Error("failed to query txs", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	// result, err := model.ParseTransactions(txs)
	// if err != nil {
	// 	zap.L().Error("failed to parse txs", zap.Error(err))
	// 	errors.ErrInternalServer(rw, http.StatusInternalServerError)
	// 	return
	// }

	model.Respond(rw, txs)
	return
}
