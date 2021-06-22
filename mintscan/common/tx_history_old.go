package common

import (
	"net/http"

	//internal
	"github.com/cosmostation/cosmostation-cosmos/app"
	"github.com/cosmostation/cosmostation-cosmos/errors"
	"github.com/cosmostation/cosmostation-cosmos/model"

	//mbl
	mbltypes "github.com/cosmostation/mintscan-backend-library/types"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

// GetAccountTxsHistory returns transactions that are sent by an account
// 주어진 txID 보다 작은 Transaction account history 반환
func GetAccountTxs(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accAddr := vars["accAddr"]

		from, limit, err := model.ParseHTTPArgs(r)
		if err != nil {
			zap.S().Debug("failed to parse HTTP args ", zap.Error(err))
			errors.ErrInvalidParam(rw, http.StatusBadRequest, "request is invalid")
			return
		}

		err = mbltypes.VerifyBech32AccAddr(accAddr)
		if err != nil {
			zap.L().Debug("failed to validate account address", zap.Error(err))
			errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
			return
		}

		// Query transactions that are made by the account.
		txs, err := a.DB.GetAccountTransactionHistory(from, limit, accAddr)
		if err != nil {
			zap.L().Error("failed to query txs", zap.Error(err))
			errors.ErrInternalServer(rw, http.StatusInternalServerError)
			return
		}

		if len(txs) <= 0 {
			model.Respond(rw, []model.ResultTx{})
			return
		}

		result := model.ParseTransactions(a, txs)

		model.Respond(rw, result)
		return
	}
}

// GetAccountTransferTxsHistory returns transfer txs (MsgSend and MsgMultiSend) that are sent by an account
// 전달 받은 txID 보다 작은, send, multisend 메세지만 리턴
func GetAccountTransferTxs(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accAddr := vars["accAddr"]

		from, limit, err := model.ParseHTTPArgs(r)
		if err != nil {
			zap.S().Debug("failed to parse HTTP args ", zap.Error(err))
			errors.ErrInvalidParam(rw, http.StatusBadRequest, "request is invalid")
			return
		}

		err = mbltypes.VerifyBech32AccAddr(accAddr)
		if err != nil {
			zap.L().Debug("failed to validate account address", zap.Error(err))
			errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
			return
		}
		txs, err := a.DB.GetTransferTransactionHistory(from, limit, accAddr)
		if err != nil {
			zap.L().Error("failed to query txs", zap.Error(err))
			errors.ErrInternalServer(rw, http.StatusInternalServerError)
			return
		}

		result := model.ParseTransactions(a, txs)

		model.Respond(rw, result)
		return
	}
}

// GetTxsHistoryBetweenDelegatorAndValidator returns transactions that are made between an account and his delegated validator
// 주어진 tx_id보다 작은, 특정 검증인과 특정 위임자의 tx history 반환
func GetTxsBetweenDelegatorAndValidator(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accAddr := vars["accAddr"]
		valAddr := vars["valAddr"]

		from, limit, err := model.ParseHTTPArgs(r)
		if err != nil {
			zap.S().Debug("failed to parse HTTP args ", zap.Error(err))
			errors.ErrInvalidParam(rw, http.StatusBadRequest, "request is invalid")
			return
		}

		err = mbltypes.VerifyBech32AccAddr(accAddr)
		if err != nil {
			zap.L().Debug("failed to validate account address", zap.Error(err))
			errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
			return
		}

		err = mbltypes.VerifyBech32ValAddr(valAddr)
		if err != nil {
			zap.L().Debug("failed to validate validator address", zap.Error(err))
			errors.ErrInvalidParam(rw, http.StatusBadRequest, "validator address is invalid")
			return
		}

		txs, err := a.DB.GetTransactionHistoryBetweenAccAndVal(from, limit, accAddr, valAddr)
		if err != nil {
			zap.L().Error("failed to query txs", zap.Error(err))
			errors.ErrInternalServer(rw, http.StatusInternalServerError)
			return
		}

		result := model.ParseTransactions(a, txs)

		model.Respond(rw, result)
		return
	}
}
