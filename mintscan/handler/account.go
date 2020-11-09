package handler

import (
	"fmt"
	"net/http"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cosmosvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	clienttypes "github.com/cosmostation/cosmostation-cosmos/mintscan/client/types"
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

	resp, err := s.client.RequestWithRestServer(clienttypes.PrefixAuth + "/accounts/" + accAddr)
	if err != nil {
		zap.L().Error("failed to get account information", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	var ar authtypes.QueryAccountResponse
	if err = s.client.GetCliContext().JSONMarshaler.UnmarshalJSON(resp, &ar); err != nil {
		zap.L().Error("failed to get unmarshal given response", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	ai, ok := ar.GetAccount().GetCachedValue().(authtypes.AccountI)
	if !ok {
		zap.S().Info("Unsupported account type")
	}
	switch aType := ai.(type) {
	case *authtypes.ModuleAccount:
		zap.S().Info("module account :", aType)
	case *authtypes.BaseAccount:
		zap.S().Info("base account :", aType)
	default:
		zap.S().Info("Unknown account type :", aType)
	}
	// zap.S().Info("account :", ai.GetAddress())
	// zap.S().Info("account :", ai.GetPubKey())
	// zap.S().Info("account :", ai.GetAccountNumber())
	// zap.S().Info("account :", ai.GetSequence())

	model.Respond(rw, ai)
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

	resp, err := s.client.RequestWithRestServer(clienttypes.PrefixBank + "/balances/" + accAddr)
	if err != nil {
		zap.L().Error("failed to get account balance", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	var abr banktypes.QueryAllBalancesResponse
	if err = s.client.GetCliContext().JSONMarshaler.UnmarshalJSON(resp, &abr); err != nil {
		zap.L().Error("failed to get unmarshal given response", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, &abr)
	return
}

// GetDelegatorDelegations returns all delegations from a delegator.
// TODO: This API uses 3 REST API requests.
// Don't need to be handled immediately, but if this ever slows down or gives burden to our
// REST server, change to use RPC to see if it gets better.
func GetDelegatorDelegations(rw http.ResponseWriter, r *http.Request) {
	/*
		이 함수의 기능은 특정 위임자(주어진 주소)가 위임한 모든 검증인의 상세 데이터를 출력하는 것임
		- 위임자 주소 (staking/delegations/delAddr)
		- 검증인 주소 (staking/delegations/delAddr)
		- 모니커
		- 지분(shares) (staking/delegations/delAddr)
		- 잔고 (staking/delegations/delAddr)
		- 지분으로 계산한 잔고 (staking/delegations/delAddr)
		- 리워드
	*/

	vars := mux.Vars(r)
	accAddr := vars["accAddr"]
	//jeonghwan todo 금요일날 하다가 중단
	err := model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}

	// Query all delegations from a delegator
	// https://lcd-office.cosmostation.io/stargate-4/cosmos/staking/v1beta1/delegations/cosmos1x5wgh6vwye60wv3dtshs9dmqggwfx2ldnqvev0
	// /cosmos/staking/v1beta1/delegators/{delegator_addr}/validators
	resp, err := s.client.RequestWithRestServer(clienttypes.PrefixStaking + "/delegations/" + accAddr)
	if err != nil {
		zap.L().Error("failed to get delegators delegations", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	var ddr stakingtypes.QueryDelegatorDelegationsResponse
	if err = s.client.GetCliContext().JSONMarshaler.UnmarshalJSON(resp, &ddr); err != nil {
		zap.L().Error("failed to get unmarshal given response", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	// delegations := make([]model.Delegations, 0)

	// err = json.Unmarshal(resp, &delegations)
	// // err = json.Unmarshal(resp.Result, &delegations)
	// if err != nil {
	// 	zap.L().Error("failed to unmarshal delegations", zap.Error(err))
	// 	errors.ErrFailedUnmarshalJSON(rw, http.StatusInternalServerError)
	// 	return
	// }

	resultDelegations := make([]model.ResultDelegations, 0)

	if len(ddr.DelegationResponses) > 0 {
		for _, delegation := range ddr.DelegationResponses {
			zap.S().Info("deletation.Balance.Denom :", delegation.Balance.Denom)
			zap.S().Info("deletation.Balance.Amount :", delegation.Balance.Amount)
			zap.S().Info("deletation.Delegation.DelegatorAddress :", delegation.Delegation.DelegatorAddress)
			zap.S().Info("deletation.Delegation.ValidatorAddress :", delegation.Delegation.ValidatorAddress)
			zap.S().Info("deletation.Delegation.Shares.String() :", delegation.Delegation.Shares.String())
			// Query a delegation reward
			resp, err := s.client.RequestWithRestServer(clienttypes.PrefixDistribution + "/delegators/" + delegation.Delegation.DelegatorAddress + "/rewards/" + delegation.Delegation.ValidatorAddress)
			if err != nil {
				zap.L().Error("failed to get delegator rewards", zap.Error(err))
				errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
				return
			}

			// var dwar distributiontypes.QueryDelegatorTotalRewardsResponse
			var drr distributiontypes.QueryDelegationRewardsResponse
			if err = s.client.GetCliContext().JSONMarshaler.UnmarshalJSON(resp, &drr); err != nil {
				zap.L().Error("failed to get unmarshal given response", zap.Error(err))
				errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
				return
			}

			// rewardsResp, err := s.client.RequestWithRestServer(clienttypes.PrefixDistribution + "/delegators/" + accAddr + "/rewards/" + delegation.Delegation.ValidatorAddress)
			// if err != nil {
			// 	zap.L().Error("failed to get a delegation reward", zap.Error(err))
			// 	errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			// 	return
			// }

			// var rewards []model.Coin
			// err = json.Unmarshal(rewardsResp, &rewards)
			// // err = json.Unmarshal(rewardsResp.Result, &rewards)
			// if err != nil {
			// 	zap.L().Error("failed to unmarshal rewards", zap.Error(err))
			// 	errors.ErrFailedUnmarshalJSON(rw, http.StatusInternalServerError)
			// 	return
			// }

			resultRewards := make([]model.Coin, 0)

			denom, err := s.client.GetBondDenom()
			if err != nil {
				return
			}

			// Exception: reward is null when the fee of delegator's validator is 100%
			if len(drr.Rewards) > 0 {
				for _, reward := range drr.Rewards {
					tempReward := &model.Coin{
						Denom:  reward.Denom,
						Amount: reward.Amount.String(),
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

			// 위임한 검증인의 모니커 조회
			resp, err = s.client.RequestWithRestServer(clienttypes.PrefixStaking + "/validators/" + delegation.Delegation.ValidatorAddress)
			if err != nil {
				zap.L().Error("failed to get delegations from a validator", zap.Error(err))
				errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
				return
			}

			var vr stakingtypes.QueryValidatorResponse
			err = s.client.GetCliContext().JSONMarshaler.UnmarshalJSON(resp, &vr)
			// err = json.Unmarshal(resp, &delegations)
			if err != nil {
				zap.L().Error("failed to unmarshal delegations", zap.Error(err))
				errors.ErrFailedUnmarshalJSON(rw, http.StatusInternalServerError)
				return
			}
			// Query the information from a single validator
			// valResp, err := s.client.RequestWithRestServer(clienttypes.PrefixStaking + "/validators/" + delegation.Delegation.ValidatorAddress)
			// if err != nil {
			// 	zap.L().Error("failed to get a delegation reward", zap.Error(err))
			// 	errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			// 	return
			// }

			// var validator model.Validator
			// err = json.Unmarshal(valResp, &validator)
			// // err = json.Unmarshal(valResp.Result, &validator)
			// if err != nil {
			// 	zap.L().Error("failed to unmarshal validator", zap.Error(err))
			// 	errors.ErrFailedUnmarshalJSON(rw, http.StatusInternalServerError)
			// 	return
			// }

			// Calculate the amount of ukava, which should divide validator's token divide delegator_shares
			// tokens, _ := strconv.ParseFloat(vr.Validator.Tokens.String(), 64)
			// delegatorShares, _ := strconv.ParseFloat(vr.Validator.DelegatorShares.String(), 64)
			// uatom := tokens / delegatorShares
			// shares, _ := strconv.ParseFloat(delegation.Delegation.Shares.String(), 64)
			// amount := fmt.Sprintf("%f", shares*uatom)

			temp := &model.ResultDelegations{
				DelegatorAddress: delegation.Delegation.DelegatorAddress,
				ValidatorAddress: delegation.Delegation.ValidatorAddress,
				Moniker:          vr.Validator.Description.Moniker,
				Shares:           delegation.Delegation.Shares.String(),
				Balance:          delegation.Balance.Amount.String(),
				Amount:           delegation.Balance.Amount.String(),
				// Amount:           amount,
				Rewards: resultRewards,
			}
			resultDelegations = append(resultDelegations, *temp)
		}
	}

	model.Respond(rw, resultDelegations)
	return
}

// GetDelegationsRewards returns total amount of rewards from a delegator's delegations.
func GetDelegationsRewards(rw http.ResponseWriter, r *http.Request) {
	//jeonghwan
	GetTotalRewardsFromDelegator(rw, r)
	return
}

// GetDelegatorUnbondingDelegations returns unbonding delegations from a delegator
func GetDelegatorUnbondingDelegations(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["delAddr"]

	fmt.Println(accAddr)
	err := model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "acount address is invalid")
		return
	}

	resp, err := s.client.RequestWithRestServer(clienttypes.PrefixStaking + "/delegators/" + accAddr + "/unbonding_delegations")
	if err != nil {
		zap.L().Error("failed to get account delegators rewards", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	var dudr stakingtypes.QueryDelegatorUnbondingDelegationsResponse
	// if err = json.Unmarshal(resp, &dudr); err != nil {
	if err = s.client.GetCliContext().JSONMarshaler.UnmarshalJSON(resp, &dudr); err != nil {
		zap.L().Error("failed to get unmarshal given response", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	result := make([]*model.UnbondingDelegations, 0)

	for _, u := range dudr.UnbondingResponses {
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

	// result, err := model.ParseTransactions(txs)
	// if err != nil {
	// 	zap.L().Error("failed to parse txs", zap.Error(err))
	// 	errors.ErrInternalServer(rw, http.StatusInternalServerError)
	// 	return
	// }

	model.Respond(rw, txs)
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

	// result, err := model.ParseTransactions(txs)
	// if err != nil {
	// 	zap.L().Error("failed to parse txs", zap.Error(err))
	// 	errors.ErrInternalServer(rw, http.StatusInternalServerError)
	// 	return
	// }

	model.Respond(rw, txs)
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

	// result, err := model.ParseTransactions(txs)
	// if err != nil {
	// 	zap.L().Error("failed to parse txs", zap.Error(err))
	// 	errors.ErrInternalServer(rw, http.StatusInternalServerError)
	// 	return
	// }

	model.Respond(rw, txs)
	return
}

// GetTotalBalance returns account's kava total, available, vesting, delegated, unbondings, rewards, deposited, incentive, and commussion.
func GetTotalBalance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["accAddr"]

	err := model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.S().Debugf("failed to validate account address: %s", err)
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}

	latestBlock, err := s.client.GetLatestBlockHeight()
	if err != nil {
		zap.S().Errorf("failed to get the latest block height: %s", err)
		return
	}

	block, err := s.client.GetBlock(latestBlock)
	if err != nil {
		zap.S().Errorf("failed to get block information: %s", err)
		return
	}

	denom, err := s.client.GetBondDenom()
	if err != nil {
		zap.S().Errorf("failed to get staking denom: %s", err)
		return
	}

	// Initialize all variables
	total := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
	available := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
	delegated := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
	undelegated := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
	rewards := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
	vesting := sdktypes.NewCoin(denom, sdktypes.NewInt(0)) // vesting 된 것 중에 delegatable 한 수량
	vested := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
	commission := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
	// failedVested := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
	// incentive := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
	// deposited := sdktypes.NewCoin(denom, sdktypes.NewInt(0))

	account, err := s.client.GetAccount(accAddr)
	if err != nil {
		zap.S().Debugf("failed to get account information: %s", err)
		errors.ErrNotFound(rw, http.StatusNotFound)
		return
	}

	// available

	coins, err := s.client.GetAccountBalance(accAddr)
	if err != nil {
		zap.S().Debugf("failed to get account balance: %s", err)
		errors.ErrNotFound(rw, http.StatusNotFound)
		return
	}

	if coins != nil {
		if coins.Denom == denom {
			available = available.Add(*coins)
		}
	}

	// Delegated
	delegations, err := s.client.GetDelegatorDelegations(accAddr)
	if err != nil {
		zap.S().Errorf("failed to get delegator's delegations: %s", err)
		return
	}

	if len(delegations) > 0 {
		for _, delegation := range delegations {
			delegated = delegated.Add(delegation.Balance)
		}
	}

	// Undelegated
	undelegations, err := s.client.GetDelegatorUndelegations(accAddr)
	if err != nil {
		zap.S().Errorf("failed to get delegator's undelegations: %s", err)
		return
	}

	if len(undelegations) > 0 {
		for _, undelegation := range undelegations {
			for _, e := range undelegation.Entries {
				undelegated = undelegated.Add(sdktypes.NewCoin(denom, e.Balance))
			}
		}
	}

	// Rewards
	totalRewards, err := s.client.GetDelegatorTotalRewards(accAddr)
	if err != nil {
		zap.S().Errorf("failed to get get delegator's total rewards: %s", err)
		return
	}

	if len(totalRewards.Rewards) > 0 {
		for _, tr := range totalRewards.Rewards {
			for _, reward := range tr.Reward {
				if reward.Denom == denom {
					truncatedRewards, _ := reward.TruncateDecimal()
					rewards = rewards.Add(truncatedRewards)
				}
			}
		}
	}

	valAddr, err := model.ConvertValAddrFromAccAddr(accAddr)
	if err != nil {
		zap.S().Errorf("failed to convert validator address from account address: %s", err)
		return
	}

	// Commission
	commissions, err := s.client.GetValidatorCommission(valAddr)
	if err != nil {
		zap.S().Errorf("failed to get validator's commission: %s", err)
		return
	}

	if len(commissions) > 0 {
		for _, c := range commissions {
			commission = commission.Add(c)
		}
	}

	// Vesting, vested, failed vested
	switch account.(type) {
	case *cosmosvesting.PeriodicVestingAccount:
		acct := account.(*cosmosvesting.PeriodicVestingAccount)

		vestingCoins := acct.GetVestingCoins(block.Block.Time)
		vestedCoins := acct.GetVestedCoins(block.Block.Time)
		delegatedVesting := acct.GetDelegatedVesting()

		// When total vesting amount is greater than or equal to delegated vesting amount, then
		// there is still a room to delegate. Otherwise, vesting should be zero.
		if len(vestingCoins) > 0 {
			if vestingCoins.IsAllGTE(delegatedVesting) {
				vestingCoins = vestingCoins.Sub(delegatedVesting)
				for _, vc := range vestingCoins {
					if vc.Denom == denom {
						vesting = vesting.Add(vc)
						available = available.Sub(vc) // available should deduct vesting amount
					}
				}
			}
		}

		if len(vestedCoins) > 0 {
			for _, vc := range vestedCoins {
				if vc.Denom == denom {
					vested = vested.Add(vc)
				}
			}
		}
	}

	// Sum up all
	total = total.Add(available).
		Add(delegated).
		Add(undelegated).
		Add(rewards).
		Add(commission).
		Add(vesting)

	result := &model.ResultTotalBalance{
		Total:       total,
		Available:   available,
		Delegated:   delegated,
		Undelegated: undelegated,
		Rewards:     rewards,
		Commission:  commission,
		Vesting:     vesting,
		Vested:      vested,
	}

	model.Respond(rw, result)
	return
}
