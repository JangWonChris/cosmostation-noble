package mobile

import (
	"net/http"
	"strconv"

	//internal
	"github.com/cosmostation/cosmostation-cosmos/mintscan/db"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/handler"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"

	//mbl
	mbltypes "github.com/cosmostation/mintscan-backend-library/types"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

func PrePareMsgExp() {
	// transfer tx - send, multisend
	db.PrepareTransferMsgExp(
		handler.MessageTypeMap[mbltypes.BankMsgSend],
		handler.MessageTypeMap[mbltypes.BankMsgMultiSend],
	)
	// validator - delegator 간 tx
	db.PrepareStakingMsgExp(
		handler.MessageTypeMap[mbltypes.DistributionMsgSetWithdrawAddress],
		handler.MessageTypeMap[mbltypes.DistributionMsgWithdrawDelegatorReward],
		handler.MessageTypeMap[mbltypes.DistributionMsgWithdrawValidatorCommission],
		handler.MessageTypeMap[mbltypes.SlashingMsgUnjail],
		handler.MessageTypeMap[mbltypes.StakingMsgCreateValidator],
		handler.MessageTypeMap[mbltypes.StakingMsgEditValidator],
		handler.MessageTypeMap[mbltypes.StakingMsgDelegate],
		handler.MessageTypeMap[mbltypes.StakingMsgBeginRedelegate],
		handler.MessageTypeMap[mbltypes.StakingMsgUndelegate],
	)
}

// ParseHTTPArgsWithBeforeAfterLimit parses the request's URL and returns all arguments pairs.
// It separates page and limit used for pagination where a default limit can be provided.
func ParseHTTPArgsWithTxID(r *http.Request) (txID int64, err error) {
	txidStr := r.FormValue("txid")
	if txidStr == "" {
		return txID, nil
	}

	txID, err = strconv.ParseInt(txidStr, 10, 64)
	if err != nil {
		return txID, err
	}

	if txID < 0 {
		return 0, nil
	}

	return txID, nil
}

// GetAccountTxsHistory returns transactions that are sent by an account
// 주어진 txID 보다 작은 Transaction account history 반환
func GetAccountTxsHistory(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["accAddr"]

	txID, err := ParseHTTPArgsWithTxID(r)
	if err != nil {
		zap.S().Debug("failed to parse HTTP args ", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "request is invalid")
		return
	}

	// if limit > 100 {
	// 	zap.S().Debug("failed to query with this limit ", zap.Int("request limit", limit))
	// 	errors.ErrOverMaxLimit(rw, http.StatusUnauthorized)
	// 	return
	// }

	err = mbltypes.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}

	// Query transactions that are made by the account.
	txs, err := s.DB.QueryTransactionsByAddr(txID, accAddr)
	if err != nil {
		zap.L().Error("failed to query txs", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	if len(txs) <= 0 {
		model.Respond(rw, []model.ResultTx{})
		return
	}

	result := model.ParseTransactions(txs)

	model.Respond(rw, result)
	return
}

// GetAccountTransferTxsHistory returns transfer txs (MsgSend and MsgMultiSend) that are sent by an account
// 전달 받은 txID 보다 작은, send, multisend 메세지만 리턴
func GetAccountTransferTxsHistory(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["accAddr"]

	txID, err := ParseHTTPArgsWithTxID(r)
	if err != nil {
		zap.S().Debug("failed to parse HTTP args ", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "request is invalid")
		return
	}

	// var denom string

	// if len(r.URL.Query()["denom"]) > 0 {
	// 	denom = r.URL.Query()["denom"][0]
	// }

	// if denom == "" {
	// 	denom, err = s.Client.GRPC.GetBondDenom(context.Background())
	// 	if err != nil {
	// 		return
	// 	}
	// }

	err = mbltypes.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}
	txs, err := s.DB.QueryTransferTransactionsByAddr(txID, accAddr)
	if err != nil {
		zap.L().Error("failed to query txs", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	result := model.ParseTransactions(txs)

	model.Respond(rw, result)
	return
}

// GetTxsHistoryBetweenDelegatorAndValidator returns transactions that are made between an account and his delegated validator
// 주어진 tx_id보다 작은, 특정 검증인과 특정 위임자의 tx history 반환
func GetTxsHistoryBetweenDelegatorAndValidator(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["accAddr"]
	valAddr := vars["valAddr"]

	txID, err := ParseHTTPArgsWithTxID(r)
	if err != nil {
		zap.S().Debug("failed to parse HTTP args ", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "request is invalid")
		return
	}
	// before, after, limit, err := model.ParseHTTPArgsWithBeforeAfterLimit(r, model.DefaultBefore, model.DefaultAfter, model.DefaultLimit)
	// if err != nil {
	// 	zap.S().Debug("failed to parse HTTP args ", zap.Error(err))
	// 	errors.ErrInvalidParam(rw, http.StatusBadRequest, "request is invalid")
	// 	return
	// }

	// if limit > 100 {
	// 	zap.S().Debug("failed to query with this limit ", zap.Int("request limit", limit))
	// 	errors.ErrOverMaxLimit(rw, http.StatusUnauthorized)
	// 	return
	// }

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

	txs, err := s.DB.QueryTransactionsBetweenAccountAndValidator(txID, accAddr, valAddr)
	if err != nil {
		zap.L().Error("failed to query txs", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	result := model.ParseTransactions(txs)

	model.Respond(rw, result)
	return
}
