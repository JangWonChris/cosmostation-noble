package common

import (
	"context"
	"net/http"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	gaiaapp "github.com/cosmos/gaia/v4/app"

	//internal
	"github.com/cosmostation/cosmostation-cosmos/app"
	"github.com/cosmostation/cosmostation-cosmos/errors"
	"github.com/cosmostation/cosmostation-cosmos/model"

	//mbl
	ltypes "github.com/cosmostation/mintscan-backend-library/types"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

var (
	pageLimit      = uint64(100)
	moduleAccounts = make(map[string]string) // map for module account name, address
)

func init() {
	maccPerms := gaiaapp.GetMaccPerms()
	for name := range maccPerms {
		moduleAccounts[name] = authtypes.NewModuleAddress(name).String()
	}
}

func GetModuleAccounts(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		macc := make([]*model.ModuleAccount, 0)

		for name, addr := range moduleAccounts {
			account, err := a.Client.CliCtx.GetAccount(addr)
			if err != nil {
				zap.L().Error("failed to get module account information", zap.Error(err))
				errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			}

			acc, ok := account.(authtypes.ModuleAccountI)
			if !ok {
				zap.L().Error("account type is not module account", zap.Error(err))
			}

			coins, err := a.Client.GRPC.GetAllBalances(context.Background(), addr, pageLimit)
			if err != nil {
				zap.L().Error("failed to get module account balance", zap.Error(err))
				errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			}
			ma := &model.ModuleAccount{
				Address:       addr,
				AccountNumber: acc.GetAccountNumber(),
				Coins:         coins,
				Permissions:   acc.GetPermissions(),
				Name:          name,
			}
			macc = append(macc, ma)
		}

		model.Respond(rw, macc)
	}
}

// GetAuthAccount returns general account information.
func GetAuthAccount(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accAddr := vars["accAddr"]
		err := ltypes.VerifyBech32AccAddr(accAddr)
		if err != nil {
			zap.L().Debug("failed to validate account address", zap.Error(err))
			errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
			return
		}
		account, err := a.Client.CliCtx.GetAccount(accAddr)
		if err != nil {
			zap.L().Error("failed to get account information", zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			return
		}

		var b []byte
		switch account := account.(type) {
		case *authtypes.ModuleAccount:
			b, err = a.Client.GetCLIContext().JSONMarshaler.MarshalJSON(account)
		case *authtypes.BaseAccount:
			b, err = a.Client.GetCLIContext().JSONMarshaler.MarshalJSON(account)
		case *vestingtypes.ContinuousVestingAccount:
			b, err = a.Client.GetCLIContext().JSONMarshaler.MarshalJSON(account)
		case *vestingtypes.DelayedVestingAccount:
			b, err = a.Client.GetCLIContext().JSONMarshaler.MarshalJSON(account)
		case *vestingtypes.PeriodicVestingAccount:
			b, err = a.Client.GetCLIContext().JSONMarshaler.MarshalJSON(account)
		default:
			zap.L().Error("unknown account type :", zap.String("info", account.GetAddress().String()), zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			return
		}

		model.Respond(rw, b)
		return
	}
}

// GetBalance returns account balance.
func GetBalance(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accAddr := vars["accAddr"]
		err := ltypes.VerifyBech32AccAddr(accAddr)
		if err != nil {
			zap.L().Debug("failed to validate account address", zap.Error(err))
			errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
			return
		}

		denom, err := a.Client.GRPC.GetBondDenom(r.Context())
		if err != nil {
			zap.L().Debug("failed to get account balance", zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			return
		}

		res, err := a.Client.GRPC.GetBalance(r.Context(), denom, accAddr)
		if err != nil {
			zap.L().Debug("failed to get account balance", zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			return
		}

		model.Respond(rw, res)
		return
	}
}

// GetAllBalances returns all of bank/balances.
func GetAllBalances(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// type ResultAllBalances struct {
		// 	Balance []Balance `json:"balances"`
		// }

		// type Balance struct {
		// 	Denom       string `json:"denom"`
		// 	Total       string `json:"total"`
		// 	Available   string `json:"available"`
		// 	Delegated   string `json:"delegated"`
		// 	Undelegated string `json:"undelegated"`
		// 	Rewards     string `json:"rewards"`
		// 	Commission  string `json:"commission"`
		// 	Vesting     string `json:"vesting"`
		// 	Vested      string `json:"vested"`
		// }
		vars := mux.Vars(r)
		accAddr := vars["accAddr"]
		err := ltypes.VerifyBech32AccAddr(accAddr)
		if err != nil {
			zap.L().Debug("failed to validate account address", zap.Error(err))
			errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
			return
		}
		res, err := a.Client.GRPC.GetAllBalances(r.Context(), accAddr, pageLimit)
		if err != nil {
			zap.L().Debug("failed to get all account balances", zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			return
		}
		result := make([]model.Balance, len(res))
		for i, coin := range res {
			result[i].Denom = coin.Denom
			result[i].Available = coin.Amount.String()
		}

		model.Respond(rw, result)
		return
	}
}

func GetDelegatorDelegations(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
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
		err := ltypes.VerifyBech32AccAddr(accAddr)
		if err != nil {
			zap.L().Debug("failed to validate account address", zap.Error(err))
			errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
			return
		}
		// Query all delegations from a delegator
		dd, err := a.Client.GRPC.GetDelegatorDelegations(r.Context(), accAddr, pageLimit)
		if err != nil {
			zap.L().Error("failed to get delegators delegations", zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			return
		}

		dtr, err := a.Client.GRPC.GetDelegationTotalRewards(r.Context(), accAddr)
		if err != nil {
			zap.L().Error("failed to get delegator rewards", zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			return
		}

		denom, err := a.Client.GRPC.GetBondDenom(r.Context())
		if err != nil {
			return
		}

		resultDelegations := make([]model.ResultDelegations, 0)
		for i, reward := range dtr.Rewards {

			resultRewards := make(sdktypes.DecCoins, 0)

			tempReward := sdktypes.DecCoin{
				Denom:  denom,
				Amount: reward.Reward.AmountOf(denom),
			}
			resultRewards = append(resultRewards, tempReward)

			// 위임한 검증인의 모니커 조회
			vr, err := a.Client.GRPC.GetValidator(r.Context(), reward.ValidatorAddress)
			if err != nil {
				zap.L().Error("failed to get delegations from a validator", zap.Error(err))
				errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
				return
			}

			temp := &model.ResultDelegations{
				DelegatorAddress: accAddr,
				ValidatorAddress: reward.ValidatorAddress,
				Moniker:          vr.Description.Moniker,
				Shares:           vr.DelegatorShares.String(),
				Amount:           dd.DelegationResponses[i].Balance.Amount.String(),
				Rewards:          resultRewards,
			}
			resultDelegations = append(resultDelegations, *temp)
		}

		model.Respond(rw, resultDelegations)
		return
	}
}

//이력 관리 용
func GetDelegatorUndelegations(a *app.App) http.HandlerFunc {
	return GetDelegatorUnbondingDelegations(a)
}

// GetDelegatorUnbondingDelegations returns unbonding delegations from a delegator
func GetDelegatorUnbondingDelegations(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accAddr := vars["accAddr"]

		err := ltypes.VerifyBech32AccAddr(accAddr)
		if err != nil {
			zap.L().Debug("failed to validate account address", zap.Error(err))
			errors.ErrInvalidParam(rw, http.StatusBadRequest, "acount address is invalid")
			return
		}

		res, err := a.Client.GRPC.GetDelegatorUnbondingDelegations(r.Context(), accAddr, pageLimit)
		if err != nil {
			zap.L().Error("failed to get account delegators rewards", zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			return
		}

		result := make([]*model.UnbondingDelegations, 0)
		for _, u := range res.UnbondingResponses {
			val, err := a.DB.GetValidatorByAnyAddr(u.ValidatorAddress)
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
}

// GetValidatorCommission returns a validator's commission information.
func GetValidatorCommission(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accAddr := vars["accAddr"]

		err := ltypes.VerifyBech32AccAddr(accAddr)
		if err != nil {
			zap.L().Debug("failed to validate account address", zap.Error(err))
			errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
			return
		}

		valAddr, err := ltypes.ConvertValAddrFromAccAddr(accAddr)
		if err != nil {
			zap.L().Debug("failed to validate validator address", zap.Error(err))
			errors.ErrInvalidParam(rw, http.StatusBadRequest, "validator address is invalid")
			return
		}

		comm, err := a.Client.GRPC.GetValidatorCommission(r.Context(), valAddr)
		if err != nil {
			zap.L().Error("failed to get validator commission", zap.Error(err))
		}

		model.Respond(rw, comm)
		return
	}
}

// GetTotalBalance returns account's total, available, vesting, delegated, unbondings, rewards, deposited, incentive, and commission for staking denom.
func GetTotalBalance(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accAddr := vars["accAddr"]

		err := ltypes.VerifyBech32AccAddr(accAddr)
		if err != nil {
			zap.S().Debugf("failed to validate account address: %s", err)
			errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
			return
		}

		denom, err := a.Client.GRPC.GetBondDenom(r.Context())
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

		coins, err := a.Client.GRPC.GetBalance(r.Context(), denom, accAddr)
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
		delegationsResp, err := a.Client.GRPC.GetDelegatorDelegations(r.Context(), accAddr, pageLimit)
		if err != nil {
			zap.S().Errorf("failed to get delegator's delegations: %s", err)
			return
		}
		for _, delegation := range delegationsResp.DelegationResponses {
			delegated = delegated.Add(delegation.Balance)
		}

		// Undelegated
		undelegationsResp, err := a.Client.GRPC.GetDelegatorUnbondingDelegations(r.Context(), accAddr, pageLimit)
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
		totalRewardsResp, err := a.Client.GRPC.GetDelegationTotalRewards(r.Context(), accAddr)
		if err != nil {
			zap.S().Errorf("failed to get get delegator's total rewards: %s", err)
			return
		}
		if totalRewardsResp != nil {
			rewards = rewards.Add(sdktypes.NewCoin(denom, totalRewardsResp.Total.AmountOf(denom).TruncateInt()))
		}

		valAddr, err := ltypes.ConvertValAddrFromAccAddr(accAddr)
		if err != nil {
			zap.S().Errorf("failed to convert validator address from account address: %s", err)
			return
		}
		// Commission
		commissionsResp, err := a.Client.GRPC.GetValidatorCommission(r.Context(), valAddr)
		if err != nil {
			zap.S().Errorf("failed to get validator's commission: %s", err)
			return
		}
		for _, c := range commissionsResp.Commission {
			truncatedCoin, _ := c.TruncateDecimal()
			commission = commission.Add(truncatedCoin)
		}

		account, err := a.Client.CliCtx.GetAccount(accAddr)
		if err != nil {
			zap.S().Debugf("failed to get account information: %s", err)
			errors.ErrNotFound(rw, http.StatusNotFound)
			return
		}

		// Vesting, vested
		switch acct := account.(type) {
		case *vestingtypes.PeriodicVestingAccount:
			vestingCoins := acct.GetVestingCoins(time.Now())
			vestedCoins := acct.GetVestedCoins(time.Now())
			delegatedVesting := acct.GetDelegatedVesting()

			// vesting 수량은 delegate 한 수량은 제외한다. (vesting 중이어도 delegate 한 수량은 delegate에 표시)
			if len(vestingCoins) > 0 {
				if vestingCoins.IsAllGT(delegatedVesting) {
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
}

// GetTotalAllBalances returns account's total, available, vesting, delegated, unbondings, rewards, deposited, incentive, and commission.
func GetTotalAllBalances(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accAddr := vars["accAddr"]

		err := ltypes.VerifyBech32AccAddr(accAddr)
		if err != nil {
			zap.S().Debugf("failed to validate account address: %s", err)
			errors.ErrInvalidParam(rw, http.StatusBadRequest, "account address is invalid")
			return
		}
		denom, err := a.Client.GRPC.GetBondDenom(r.Context())
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
		availableCoins, err := a.Client.GRPC.GetAllBalances(r.Context(), accAddr, pageLimit)
		if err != nil {
			zap.S().Debugf("failed to get account balance: %s", err)
			errors.ErrNotFound(rw, http.StatusNotFound)
			return
		}
		_ = availableCoins

		// Delegated
		delegationsResp, err := a.Client.GRPC.GetDelegatorDelegations(r.Context(), accAddr, pageLimit)
		if err != nil {
			zap.S().Errorf("failed to get delegator's delegations: %s", err)
			return
		}
		for _, delegation := range delegationsResp.DelegationResponses {
			delegated = delegated.Add(delegation.Balance)
		}

		// Undelegated
		undelegationsResp, err := a.Client.GRPC.GetDelegatorUnbondingDelegations(r.Context(), accAddr, pageLimit)
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
		totalRewardsResp, err := a.Client.GRPC.GetDelegationTotalRewards(r.Context(), accAddr)
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

		valAddr, err := ltypes.ConvertValAddrFromAccAddr(accAddr)
		if err != nil {
			zap.S().Errorf("failed to convert validator address from account address: %s", err)
			return
		}

		// Commission
		commissionsResp, err := a.Client.GRPC.GetValidatorCommission(r.Context(), valAddr)
		if err != nil {
			zap.S().Errorf("failed to get validator's commission: %s", err)
			return
		}
		for _, c := range commissionsResp.Commission {
			truncatedCoin, _ := c.TruncateDecimal()
			commission = commission.Add(truncatedCoin)
		}

		account, err := a.Client.CliCtx.GetAccount(accAddr)
		if err != nil {
			zap.S().Debugf("failed to get account information: %s", err)
			errors.ErrNotFound(rw, http.StatusNotFound)
			return
		}

		latestBlock, err := a.Client.RPC.GetLatestBlockHeight()
		if err != nil {
			zap.S().Errorf("failed to get the latest block height: %s", err)
			return
		}
		block, err := a.Client.RPC.GetBlock(latestBlock)
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

			// vesting 양은 delegate 한 양은 제외한다. (vesting 중이어도 delegate 한 수량은 delegate에 표시)
			if len(vestingCoins) > 0 {
				if vestingCoins.IsAllGT(delegatedVesting) {
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
}
