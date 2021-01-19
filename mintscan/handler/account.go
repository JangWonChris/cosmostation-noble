package handler

import (
	"fmt"
	"net/http"
	"sort"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

var (
	pageLimit = uint64(100)
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
	account, err := s.client.CliCtx.GetAccount(accAddr)
	if err != nil {
		zap.L().Error("failed to get account information", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	var b []byte
	switch account := account.(type) {
	case *authtypes.ModuleAccount:
		b, err = s.client.GetCLIContext().JSONMarshaler.MarshalJSON(account)
	case *authtypes.BaseAccount:
		b, err = s.client.GetCLIContext().JSONMarshaler.MarshalJSON(account)
	case *vestingtypes.ContinuousVestingAccount:
		b, err = s.client.GetCLIContext().JSONMarshaler.MarshalJSON(account)
	case *vestingtypes.DelayedVestingAccount:
		b, err = s.client.GetCLIContext().JSONMarshaler.MarshalJSON(account)
	case *vestingtypes.PeriodicVestingAccount:
		b, err = s.client.GetCLIContext().JSONMarshaler.MarshalJSON(account)
	default:
		zap.L().Error("unknown account type :", zap.String("info", account.GetAddress().String()), zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, b)
	return
}

// GetBalance returns account balance.
func GetBalance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["accAddr"]
	err := model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}
	denom, err := s.client.GRPC.GetBondDenom(r.Context())
	if err != nil {
		zap.L().Debug("failed to get account balance", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}
	res, err := s.client.GRPC.GetBalance(r.Context(), denom, accAddr)
	if err != nil {
		zap.L().Debug("failed to get account balance", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	model.Respond(rw, res)
	return
}

// GetAllBalances returns account balance.
func GetAllBalances(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["accAddr"]
	err := model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}
	res, err := s.client.GRPC.GetAllBalances(r.Context(), accAddr, pageLimit)
	if err != nil {
		zap.L().Debug("failed to get all account balances", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

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
	resps, err := s.client.GRPC.GetDelegatorDelegations(r.Context(), accAddr, pageLimit)
	if err != nil {
		zap.L().Error("failed to get delegators delegations", zap.Error(err))
		errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
		return
	}

	rewards := func(p1, p2 *distributiontypes.DelegationDelegatorReward) bool {
		return p1.ValidatorAddress < p2.ValidatorAddress
	}
	res2, err2 := s.client.GRPC.GetDelegationTotalRewards(r.Context(), accAddr)
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
		drr, err := s.client.GRPC.GetDelegationRewards(r.Context(), resp.Delegation.DelegatorAddress, resp.Delegation.ValidatorAddress)
		if err != nil {
			zap.L().Error("failed to get delegator rewards", zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			return
		}

		resultRewards := make([]model.Coin, 0)

		denom, err := s.client.GRPC.GetBondDenom(r.Context())
		if err != nil {
			return
		}

		// Exception: reward is null when the fee of delegator's validator is 100%
		if len(drr) > 0 {
			for _, reward := range drr {
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
		vr, err := s.client.GRPC.GetValidator(r.Context(), resp.Delegation.ValidatorAddress)
		if err != nil {
			zap.L().Error("failed to get delegations from a validator", zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			return
		}

		temp := &model.ResultDelegations{
			DelegatorAddress: resp.Delegation.DelegatorAddress,
			ValidatorAddress: resp.Delegation.ValidatorAddress,
			Moniker:          vr.Description.Moniker,
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

	res, err := s.client.GRPC.GetDelegatorUnbondingDelegations(r.Context(), accAddr, pageLimit)
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
			UnbondingDelegation: u,
			// DelegatorAddress: u.DelegatorAddress,
			// ValidatorAddress: u.ValidatorAddress,
			// Entries:          u.Entries,
			Moniker: val.Moniker,
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

	comm, err := s.client.GRPC.GetValidatorCommission(r.Context(), valAddr)
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

	err := model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}
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

	err := model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to validate account address", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}
	before, after, limit, err := model.ParseHTTPArgsWithBeforeAfterLimit(r, model.DefaultBefore, model.DefaultAfter, model.DefaultLimit)
	if err != nil {
		zap.S().Debug("failed to parse HTTP args ", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "request is invalid")
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
		denom, err = s.client.GRPC.GetBondDenom(r.Context())
		if err != nil {
			return
		}
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

	err := model.VerifyBech32AccAddr(accAddr)
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

	txs, err := s.db.QueryTransactionsBetweenAccountAndValidator(accAddr, valAddr, before, after, limit)
	if err != nil {
		zap.L().Error("failed to query txs", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	model.Respond(rw, txs)
	return
}

// GetTotalBalance returns account's total, available, vesting, delegated, unbondings, rewards, deposited, incentive, and commission for staking denom.
func GetTotalBalance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["accAddr"]

	err := model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.S().Debugf("failed to validate account address: %s", err)
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}

	denom, err := s.client.GRPC.GetBondDenom(r.Context())
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
	coins, err := s.client.GRPC.GetBalance(r.Context(), denom, accAddr)
	if err != nil {
		zap.S().Debugf("failed to get account balance: %s", err)
		errors.ErrNotFound(rw, http.StatusNotFound)
		return
	}
	// jeonghwan
	// coins nil check가 필요한지, getbalance를 리턴 받을 때, available로 받으면 안되는지 확인 필요
	if coins != nil {
		available = available.Add(*coins)
	}

	// Delegated
	delegationsResp, err := s.client.GRPC.GetDelegatorDelegations(r.Context(), accAddr, pageLimit)
	if err != nil {
		zap.S().Errorf("failed to get delegator's delegations: %s", err)
		return
	}
	for _, delegation := range delegationsResp.DelegationResponses {
		delegated = delegated.Add(delegation.Balance)
	}

	// Undelegated
	undelegationsResp, err := s.client.GRPC.GetDelegatorUnbondingDelegations(r.Context(), accAddr, pageLimit)
	if err != nil {
		zap.S().Errorf("failed to get delegator's undelegations: %s", err)
		return
	}
	for _, undelegation := range undelegationsResp.UnbondingResponses {
		for _, e := range undelegation.Entries {
			undelegated = undelegated.Add(sdktypes.NewCoin(denom, e.Balance))
		}
	}

	// Rewards
	totalRewardsResp, err := s.client.GRPC.GetDelegationTotalRewards(r.Context(), accAddr)
	if err != nil {
		zap.S().Errorf("failed to get get delegator's total rewards: %s", err)
		return
	}
	if totalRewardsResp != nil {
		rewards = rewards.Add(sdktypes.NewCoin(denom, totalRewardsResp.Total.AmountOf(denom).TruncateInt()))
	}
	// totalDecs, _ := totalRewardsResp.Total.TruncateDecimal()
	// for _, reward := range totalDecs {
	// 	if reward.Denom == denom {
	// 		rewards = rewards.Add(reward)
	// 		// 특정 denom에 대한 합계를 구하고 나머지는 사용하지 않음
	// 		break
	// 	}
	// }

	valAddr, err := model.ConvertValAddrFromAccAddr(accAddr)
	if err != nil {
		zap.S().Errorf("failed to convert validator address from account address: %s", err)
		return
	}
	// Commission
	commissionsResp, err := s.client.GRPC.GetValidatorCommission(r.Context(), valAddr)
	if err != nil {
		zap.S().Errorf("failed to get validator's commission: %s", err)
		return
	}
	for _, c := range commissionsResp.Commission {
		truncatedCoin, _ := c.TruncateDecimal()
		commission = commission.Add(truncatedCoin)
	}

	account, err := s.client.CliCtx.GetAccount(accAddr)
	if err != nil {
		zap.S().Debugf("failed to get account information: %s", err)
		errors.ErrNotFound(rw, http.StatusNotFound)
		return
	}

	latestBlock, err := s.client.RPC.GetLatestBlockHeight()
	if err != nil {
		zap.S().Errorf("failed to get the latest block height: %s", err)
		return
	}
	block, err := s.client.RPC.GetBlock(latestBlock)
	if err != nil {
		zap.S().Errorf("failed to get block information: %s", err)
		return
	}
	// Vesting, vested
	switch acct := account.(type) {
	case *vestingtypes.PeriodicVestingAccount:
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

// GetTotalAllBalances returns account's total, available, vesting, delegated, unbondings, rewards, deposited, incentive, and commission.
func GetTotalAllBalances(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accAddr := vars["accAddr"]

	err := model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.S().Debugf("failed to validate account address: %s", err)
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
		return
	}

	denom, err := s.client.GRPC.GetBondDenom(r.Context())
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
	availableCoins, err := s.client.GRPC.GetAllBalances(r.Context(), accAddr, pageLimit)
	if err != nil {
		zap.S().Debugf("failed to get account balance: %s", err)
		errors.ErrNotFound(rw, http.StatusNotFound)
		return
	}
	_ = availableCoins

	// Delegated
	delegationsResp, err := s.client.GRPC.GetDelegatorDelegations(r.Context(), accAddr, pageLimit)
	if err != nil {
		zap.S().Errorf("failed to get delegator's delegations: %s", err)
		return
	}
	for _, delegation := range delegationsResp.DelegationResponses {
		delegated = delegated.Add(delegation.Balance)
	}

	// Undelegated
	undelegationsResp, err := s.client.GRPC.GetDelegatorUnbondingDelegations(r.Context(), accAddr, pageLimit)
	if err != nil {
		zap.S().Errorf("failed to get delegator's undelegations: %s", err)
		return
	}
	for _, undelegation := range undelegationsResp.UnbondingResponses {
		for _, e := range undelegation.Entries {
			undelegated = undelegated.Add(sdktypes.NewCoin(denom, e.Balance))
		}
	}

	// Rewards
	totalRewardsResp, err := s.client.GRPC.GetDelegationTotalRewards(r.Context(), accAddr)
	if err != nil {
		zap.S().Errorf("failed to get get delegator's total rewards: %s", err)
		return
	}
	// rewards, _ := totalRewardsResp.Total.TruncateDecimal()
	for _, tr := range totalRewardsResp.Rewards {
		for _, reward := range tr.Reward {
			if reward.Denom == denom {
				truncatedRewards, _ := reward.TruncateDecimal()
				rewards = rewards.Add(truncatedRewards)
			}
		}
	}

	valAddr, err := model.ConvertValAddrFromAccAddr(accAddr)
	if err != nil {
		zap.S().Errorf("failed to convert validator address from account address: %s", err)
		return
	}

	// Commission
	commissionsResp, err := s.client.GRPC.GetValidatorCommission(r.Context(), valAddr)
	if err != nil {
		zap.S().Errorf("failed to get validator's commission: %s", err)
		return
	}
	for _, c := range commissionsResp.Commission {
		truncatedCoin, _ := c.TruncateDecimal()
		commission = commission.Add(truncatedCoin)
	}

	account, err := s.client.CliCtx.GetAccount(accAddr)
	if err != nil {
		zap.S().Debugf("failed to get account information: %s", err)
		errors.ErrNotFound(rw, http.StatusNotFound)
		return
	}

	latestBlock, err := s.client.RPC.GetLatestBlockHeight()
	if err != nil {
		zap.S().Errorf("failed to get the latest block height: %s", err)
		return
	}
	block, err := s.client.RPC.GetBlock(latestBlock)
	if err != nil {
		zap.S().Errorf("failed to get block information: %s", err)
		return
	}
	// Vesting, vested
	switch acct := account.(type) {
	case *vestingtypes.PeriodicVestingAccount:
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
