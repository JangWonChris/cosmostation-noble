package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

// GetAccount returns general account information.
func GetAccount(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["accAddr"]

	err := model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}

	resp, err := s.client.HandleResponseHeight("/auth/accounts/" + accAddr)
	if err != nil {
		zap.L().Error("failed to get account information", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, resp)
	return
}

// GetAccountBalance returns account balance.
func GetAccountBalance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["accAddr"]

	err := model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}

	resp, err := s.client.HandleResponseHeight("/bank/balances/" + accAddr)
	if err != nil {
		zap.L().Error("failed to get account balance", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, resp)
	return
}

// GetDelegatorDelegations returns all delegations from a delegator.
// TODO: This API uses 3 REST API requests.
// Don't need to be handled immediately, but if this ever slows down or gives burden to our
// REST server, change to use RPC to see if it gets better.
func GetDelegatorDelegations(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["accAddr"]

	err := model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}

	// Query all delegations from a delegator
	resp, err := s.client.HandleResponseHeight("/staking/delegators/" + accAddr + "/delegations")
	if err != nil {
		zap.L().Error("failed to get delegators delegations", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	delegations := make([]model.Delegations, 0)

	err = json.Unmarshal(resp.Result, &delegations)
	if err != nil {
		zap.L().Error("failed to unmarshal delegations", zap.Error(err))
		errors.ErrFailedUnmarshalJSON(rw, http.StatusInternalServerError)
		return
	}

	resultDelegations := make([]model.ResultDelegations, 0)

	if len(delegations) > 0 {
		for _, delegation := range delegations {
			// Query a delegation reward
			rewardsResp, err := s.client.HandleResponseHeight("/distribution/delegators/" + accAddr + "/rewards/" + delegation.ValidatorAddress)
			if err != nil {
				zap.L().Error("failed to get a delegation reward", zap.Error(err))
				errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
				return
			}

			var rewards []model.Coin
			err = json.Unmarshal(rewardsResp.Result, &rewards)
			if err != nil {
				zap.L().Error("failed to unmarshal rewards", zap.Error(err))
				errors.ErrFailedUnmarshalJSON(rw, http.StatusInternalServerError)
				return
			}

			resultRewards := make([]model.Coin, 0)

			denom, err := s.client.GetBondDenom()
			if err != nil {
				return
			}

			// Exception: reward is null when the fee of delegator's validator is 100%
			if len(rewards) > 0 {
				for _, reward := range rewards {
					tempReward := &model.Coin{
						Denom:  reward.Denom,
						Amount: reward.Amount,
					}
					resultRewards = append(resultRewards, *tempReward)
				}
			} else {
				tempReward := &model.Coin{
					Denom:  denom,
					Amount: "0",
					// Amount: sdk.ZeroInt(), //"0" value is modified because cointype is changed from model.coin to sdk.coin
				}
				resultRewards = append(resultRewards, *tempReward)
			}

			// Query the information from a single validator
			valResp, err := s.client.HandleResponseHeight("/staking/validators/" + delegation.ValidatorAddress)
			if err != nil {
				zap.L().Error("failed to get a delegation reward", zap.Error(err))
				errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
				return
			}

			var validator model.Validator
			err = json.Unmarshal(valResp.Result, &validator)
			if err != nil {
				zap.L().Error("failed to unmarshal validator", zap.Error(err))
				errors.ErrFailedUnmarshalJSON(rw, http.StatusInternalServerError)
				return
			}

			// Calculate the amount of ukava, which should divide validator's token divide delegator_shares
			tokens, _ := strconv.ParseFloat(validator.Tokens, 64)
			delegatorShares, _ := strconv.ParseFloat(validator.DelegatorShares, 64)
			uatom := tokens / delegatorShares
			shares, _ := strconv.ParseFloat(delegation.Shares, 64)
			amount := fmt.Sprintf("%f", shares*uatom)

			temp := &model.ResultDelegations{
				DelegatorAddress: delegation.DelegatorAddress,
				ValidatorAddress: delegation.ValidatorAddress,
				Moniker:          validator.Description.Moniker,
				Shares:           delegation.Shares,
				Balance:          delegation.Balance,
				Amount:           amount,
				Rewards:          resultRewards,
			}
			resultDelegations = append(resultDelegations, *temp)
		}
	}

	model.Respond(rw, resultDelegations)
	return
}

// GetDelegationsRewards returns total amount of rewards from a delegator's delegations.
func GetDelegationsRewards(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["accAddr"]

	err := model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}

	resp, err := s.client.HandleResponseHeight("/distribution/delegators/" + accAddr + "/rewards")
	if err != nil {
		zap.L().Error("failed to get account delegators rewards", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, resp)
	return
}

// GetDelegatorUnbondingDelegations returns unbonding delegations from a delegator
func GetDelegatorUnbondingDelegations(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["accAddr"]

	err := model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "acount address is invalid")
		return
	}

	resp, err := s.client.HandleResponseHeight("/staking/delegators/" + accAddr + "/unbonding_delegations")
	if err != nil {
		zap.L().Error("failed to get account delegators rewards", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	unbondingDelegations := make([]model.UnbondingDelegations, 0)

	err = json.Unmarshal(resp.Result, &unbondingDelegations)
	if err != nil {
		zap.L().Error("failed to unmarshal unbonding delegations", zap.Error(err))
		errors.ErrFailedUnmarshalJSON(rw, http.StatusInternalServerError)
		return
	}

	result := make([]*model.UnbondingDelegations, 0)

	for _, u := range unbondingDelegations {
		val, err := s.db.QueryValidatorByValAddr(u.ValidatorAddress)
		if err != nil {
			zap.L().Debug("failed to query validator information", zap.Error(err))
		}

		temp := &model.UnbondingDelegations{
			DelegatorAddress: u.DelegatorAddress,
			ValidatorAddress: u.ValidatorAddress,
			Moniker:          val.Moniker,
			Entries:          u.Entries,
		}

		result = append(result, temp)
	}

	model.Respond(rw, result)
	return
}

// GetValidatorCommission returns a validator's commission information.
func GetValidatorCommission(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["accAddr"]

	err := model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}

	valAddr, err := model.ConvertValAddrFromAccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate validator address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "validator address is invalid")
		return
	}

	comm, err := s.client.GetValidatorCommission(valAddr)
	if err != nil {
		zap.L().Error("failed to get validator commission", zap.Error(err))
	}

	model.Respond(rw, comm)
	return
}

// GetAccountTxs returns transactions that are sent by an account.
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

	err = model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}

	valAddr, err := model.ConvertValAddrFromAccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to convert validator address from account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "validator address is invalid")
		return
	}

	// Query transactions that are made by the account.
	txs, err := s.db.QueryTransactionsByAddr(accAddr, valAddr, before, after, limit)
	if err != nil {
		zap.L().Error("failed to query txs", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	if len(txs) <= 0 {
		model.Respond(rw, []model.ResultTx{})
		return
	}

	result, err := model.ParseTransactions(txs)
	if err != nil {
		zap.L().Error("failed to parse txs", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	model.Respond(rw, result)
	return
}

// GetAccountTransferTxs returns transfer txs (MsgSend and MsgMultiSend) that are sent by an account.
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
		denom, err = s.client.GetBondDenom()
		if err != nil {
			return
		}
	}

	err = model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}

	txs, err := s.db.QueryTransferTransactionsByAddr(accAddr, denom, before, after, limit)
	if err != nil {
		zap.L().Error("failed to query txs", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	result, err := model.ParseTransactions(txs)
	if err != nil {
		zap.L().Error("failed to parse txs", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	model.Respond(rw, result)
	return
}

// GetTxsBetweenDelegatorAndValidator returns transactions that are made between an account and his delegated validator.
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

	err = model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}

	err = model.VerifyBech32ValAddr(valAddr)
	if err != nil {
		zap.L().Debug("failed to validate validator address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "validator address is invalid")
		return
	}

	txs, err := s.db.QueryTransactionsBetweenAccountAndValidator(accAddr, valAddr, before, after, limit)
	if err != nil {
		zap.L().Error("failed to query txs", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	result, err := model.ParseTransactions(txs)
	if err != nil {
		zap.L().Error("failed to parse txs", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	model.Respond(rw, result)
	return
}
