package handler

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

// GetAccount returns general account information.
func GetAccount(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["accAddr"]

	account, err := s.client.GetAccount(accAddr)
	if err != nil {
		zap.L().Error("failed to get account information", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	var b []byte
	switch account := account.(type) {
	case *authtypes.ModuleAccount, *authtypes.BaseAccount,
		*vestingtypes.ContinuousVestingAccount, *vestingtypes.DelayedVestingAccount, *vestingtypes.PeriodicVestingAccount:
		b, err = s.client.GetCliContext().JSONMarshaler.MarshalJSON(account)
	default:
		zap.L().Error("unknown account type :", zap.String("info", account.String()), zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, b)
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
	queryClient := banktypes.NewQueryClient(s.client.GetCliContext())
	request := banktypes.QueryAllBalancesRequest{Address: accAddr}
	res, err := queryClient.AllBalances(context.Background(), &request)

	//jeonghwan todo :
	// available 외 필요 자산 추가

	model.Respond(rw, res)
	return
}

// By is the type of a "less" function that defines the ordering of its Planet arguments.
type By func(p1, p2 *distributiontypes.DelegationDelegatorReward) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(resp []distributiontypes.DelegationDelegatorReward) {
	ps := &QueryDelegatorTotalRewardsResponseSorter{
		resp: resp,
		by:   by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// no lint
type QueryDelegatorTotalRewardsResponseSorter struct {
	resp []distributiontypes.DelegationDelegatorReward
	by   func(p1, p2 *distributiontypes.DelegationDelegatorReward) bool
}

// Len is part of sort.Interface.
func (s *QueryDelegatorTotalRewardsResponseSorter) Len() int {
	return len(s.resp)
}

// Swap is part of sort.Interface.
func (s *QueryDelegatorTotalRewardsResponseSorter) Swap(i, j int) {
	s.resp[i], s.resp[j] = s.resp[j], s.resp[i]
	// s.resp.reawrds[i], s.planets[j] = s.planets[j], s.planets[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *QueryDelegatorTotalRewardsResponseSorter) Less(i, j int) bool {
	return s.by(&s.resp[i], &s.resp[j])
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
	err := model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}

	// Query all delegations from a delegator
	queryClient := stakingtypes.NewQueryClient(s.client.GetCliContext())
	request := stakingtypes.QueryDelegatorDelegationsRequest{DelegatorAddr: accAddr}
	resps, err := queryClient.DelegatorDelegations(context.Background(), &request)
	if err != nil {
		zap.L().Error("failed to get delegators delegations", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	rewards := func(p1, p2 *distributiontypes.DelegationDelegatorReward) bool {
		return p1.ValidatorAddress < p2.ValidatorAddress
	}
	distributionQueryClient := distributiontypes.NewQueryClient(s.client.GetCliContext())
	totalrequest := distributiontypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: accAddr}
	res2, err2 := distributionQueryClient.DelegationTotalRewards(context.Background(), &totalrequest)
	if err2 != nil {
		zap.L().Error("failed to get delegator rewards", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	// test code
	// 위임자가 위임한 내역 리스트를 모두 뽑는다. (grpc)
	// 위임자가 위임한 내역의 reward 리스트를 요청 (grpc)
	// 둘의 어카운트 순서가 일치하면, 해당 내역으로 짬뽕해서 리턴
	// 아래 결과로만 받을 때는, 소팅이 필요 없어보이고,
	// vali가 일치하는지만 검사해서 리턴하면 될 것 같다.
	for _, r := range resps.DelegationResponses {
		fmt.Println(r)
	}

	fmt.Println("before sort")
	for _, r := range res2.Rewards {
		fmt.Println(r)
	}
	By(rewards).Sort(res2.Rewards)
	fmt.Println("after sort")
	for _, r := range res2.Rewards {
		fmt.Println(r)
	}
	// end testcode

	resultDelegations := make([]model.ResultDelegations, 0)
	for _, resp := range resps.DelegationResponses {
		// Query a delegation reward
		queryClient := distributiontypes.NewQueryClient(s.client.GetCliContext())
		request := distributiontypes.QueryDelegationRewardsRequest{DelegatorAddress: resp.Delegation.DelegatorAddress, ValidatorAddress: resp.Delegation.ValidatorAddress}
		drr, err := queryClient.DelegationRewards(context.Background(), &request)
		if err != nil {
			zap.L().Error("failed to get delegator rewards", zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			return
		}

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
		stakingQueryClient := stakingtypes.NewQueryClient(s.client.GetCliContext())
		stakingQueryRequest := stakingtypes.QueryValidatorRequest{ValidatorAddr: resp.Delegation.ValidatorAddress}
		vr, err := stakingQueryClient.Validator(context.Background(), &stakingQueryRequest)
		if err != nil {
			zap.L().Error("failed to get delegations from a validator", zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			return
		}

		temp := &model.ResultDelegations{
			DelegatorAddress: resp.Delegation.DelegatorAddress,
			ValidatorAddress: resp.Delegation.ValidatorAddress,
			Moniker:          vr.Validator.Description.Moniker,
			Shares:           resp.Delegation.Shares.String(),
			Balance:          resp.Balance.Amount.String(),
			Amount:           resp.Balance.Amount.String(),
			Rewards:          resultRewards,
		}
		resultDelegations = append(resultDelegations, *temp)
	}

	model.Respond(rw, resultDelegations)
	return
}

// GetDelegationsRewards returns total amount of rewards from a delegator's delegations.
func GetDelegationsRewards(rw http.ResponseWriter, r *http.Request) {
	//jeonghwan : 안쓰는 함수
	GetTotalRewardsFromDelegator(rw, r)
	return
}

// GetDelegatorUnbondingDelegations returns unbonding delegations from a delegator
func GetDelegatorUnbondingDelegations(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["delAddr"]

	err := model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "acount address is invalid")
		return
	}

	queryClient := stakingtypes.NewQueryClient(s.client.GetCliContext())
	request := stakingtypes.QueryDelegatorUnbondingDelegationsRequest{DelegatorAddr: accAddr}
	res, err := queryClient.DelegatorUnbondingDelegations(context.Background(), &request)
	if err != nil {
		zap.L().Error("failed to get account delegators rewards", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	result := make([]*model.UnbondingDelegations, 0)
	for _, u := range res.UnbondingResponses {
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

	account, err := s.client.GetAccount(accAddr)
	if err != nil {
		zap.S().Debugf("failed to get account information: %s", err)
		errors.ErrNotFound(rw, http.StatusNotFound)
		return
	}

	// Vesting, vested, failed vested
	switch account.(type) {
	case *vestingtypes.PeriodicVestingAccount:
		acct := account.(*vestingtypes.PeriodicVestingAccount)

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
