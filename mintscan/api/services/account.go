package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/client/common"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/db"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/utils"

	"github.com/tendermint/tendermint/rpc/client"

	"github.com/gorilla/mux"
	resty "gopkg.in/resty.v1"
)

// GetBalance returns balance of an anddress
func GetBalance(codec *codec.Codec, config *config.Config, db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	accAddress := vars["accAddress"]

	if !utils.VerifyAddress(accAddress) {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	resp, _ := resty.R().Get(config.Node.LCDURL + "/bank/balances/" + accAddress)

	var balances []models.Coin
	err := json.Unmarshal(models.ReadRespWithHeight(resp).Result, &balances)
	if err != nil {
		fmt.Printf("failed to unmarshal balances: %t\n", err)
	}

	result := make([]models.Coin, 0)

	for _, balance := range balances {
		tempBalance := &models.Coin{
			Denom:  balance.Denom,
			Amount: balance.Amount,
		}
		result = append(result, *tempBalance)
	}

	utils.Respond(w, result)
	return nil
}

// GetDelegationsRewards returns total amount of rewards
func GetDelegationsRewards(codec *codec.Codec, config *config.Config, db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	accAddress := vars["accAddress"]

	if !utils.VerifyAddress(accAddress) {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	resp, _ := resty.R().Get(config.Node.LCDURL + "/distribution/delegators/" + accAddress + "/rewards")

	var resultRewards models.ResultRewards
	err := json.Unmarshal(models.ReadRespWithHeight(resp).Result, &resultRewards)
	if err != nil {
		fmt.Printf("failed to unmarshal resultRewards: %t\n", err)
	}

	result := make([]models.Rewards, 0)

	for _, reward := range resultRewards.Rewards {
		coins := make([]models.Coin, 0)

		if len(reward.Reward) > 0 {
			for _, reward := range reward.Reward {
				tempCoin := &models.Coin{
					Denom:  reward.Denom,
					Amount: reward.Amount,
				}
				coins = append(coins, *tempCoin)
			}
		} else {
			tempCoin := &models.Coin{
				Denom:  config.Denom,
				Amount: "0",
			}
			coins = append(coins, *tempCoin)
		}

		tempReward := &models.Rewards{
			ValidatorAddress: reward.ValidatorAddress,
			Reward:           coins,
		}

		result = append(result, *tempReward)
	}

	utils.Respond(w, result)
	return nil
}

// GetDelegations returns all delegations from an address
func GetDelegations(codec *codec.Codec, config *config.Config, db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	accAddress := vars["accAddress"]

	if !utils.VerifyAddress(accAddress) {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Query delegations and each delegator's rewards
	resp, _ := resty.R().Get(config.Node.LCDURL + "/staking/delegators/" + accAddress + "/delegations")

	delegations := make([]models.Delegations, 0)
	err := json.Unmarshal(models.ReadRespWithHeight(resp).Result, &delegations)
	if err != nil {
		fmt.Printf("failed to unmarshal delegations: %t\n", err)
	}

	resultDelegations := make([]models.ResultDelegations, 0)

	if len(delegations) > 0 {
		for _, delegation := range delegations {
			rewardsResp, _ := resty.R().Get(config.Node.LCDURL + "/distribution/delegators/" + accAddress + "/rewards/" + delegation.ValidatorAddress)

			var rewards []models.Coin
			err = json.Unmarshal(models.ReadRespWithHeight(rewardsResp).Result, &rewards)
			if err != nil {
				fmt.Printf("failed to unmarshal rewards: %t\n", err)
			}

			resultRewards := make([]models.Coin, 0)

			// Exception: reward is null when the fee of delegator's validator is 100%
			if len(rewards) > 0 {
				for _, reward := range rewards {
					tempReward := &models.Coin{
						Denom:  reward.Denom,
						Amount: reward.Amount,
					}
					resultRewards = append(resultRewards, *tempReward)
				}
			} else {
				tempReward := &models.Coin{
					Denom:  config.Denom,
					Amount: "0",
				}
				resultRewards = append(resultRewards, *tempReward)
			}

			validatorResp, _ := resty.R().Get(config.Node.LCDURL + "/staking/validators/" + delegation.ValidatorAddress)

			var validator models.Validator
			err = json.Unmarshal(models.ReadRespWithHeight(validatorResp).Result, &validator)
			if err != nil {
				fmt.Printf("failed to unmarshal validator: %t\n", err)
			}

			// Calculate the amount of uatom, which should divide validator's token divide delegator_shares
			tokens, _ := strconv.ParseFloat(validator.Tokens, 64)
			delegatorShares, _ := strconv.ParseFloat(validator.DelegatorShares, 64)
			uatom := tokens / delegatorShares
			shares, _ := strconv.ParseFloat(delegation.Shares, 64)
			amount := fmt.Sprintf("%f", shares*uatom)

			tempResultDelegations := &models.ResultDelegations{
				DelegatorAddress: delegation.DelegatorAddress,
				ValidatorAddress: delegation.ValidatorAddress,
				Moniker:          validator.Description.Moniker,
				Shares:           delegation.Shares,
				Balance:          delegation.Balance,
				Amount:           amount,
				Rewards:          resultRewards,
			}
			resultDelegations = append(resultDelegations, *tempResultDelegations)
		}
	}

	utils.Respond(w, resultDelegations)
	return nil
}

// GetCommission returns commission information for validator's address
func GetCommission(codec *codec.Codec, config *config.Config, db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	accAddress := vars["accAddress"]

	if !utils.VerifyAddress(accAddress) {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	operAddr := utils.ValAddressFromAccAddress(accAddress)

	commission := make([]models.Coin, 0)
	if operAddr != "" {
		ctx := context.NewCLIContext().WithCodec(codec).WithClient(rpcClient)
		valAddr, _ := sdk.ValAddressFromBech32(operAddr)
		result, _ := common.QueryValidatorCommission(ctx, distr.QuerierRoute, valAddr)

		var valCom distr.ValidatorAccumulatedCommission
		ctx.Codec.MustUnmarshalJSON(result, &valCom)

		if valCom != nil {
			tempCommission := &models.Coin{
				Denom:  valCom[0].Denom,
				Amount: valCom[0].Amount.String(),
			}
			commission = append(commission, *tempCommission)
		}
	}

	utils.Respond(w, commission)
	return nil
}

// GetUnbondingDelegations returns unbonding delegations from an address
func GetUnbondingDelegations(codec *codec.Codec, config *config.Config, db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	accAddress := vars["accAddress"]

	if !utils.VerifyAddress(accAddress) {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Query unbonding delegations
	unbondingDelegationsResp, _ := resty.R().Get(config.Node.LCDURL + "/staking/delegators/" + accAddress + "/unbonding_delegations")

	unbondingDelegations := make([]models.UnbondingDelegations, 0)
	err := json.Unmarshal(models.ReadRespWithHeight(unbondingDelegationsResp).Result, &unbondingDelegations)
	if err != nil {
		fmt.Printf("failed to unmarshal unbondingDelegations: %t\n", err)
	}

	result := make([]models.UnbondingDelegations, 0)
	for _, unbondingDelegation := range unbondingDelegations {
		validator, _ := db.QueryValidatorByOperAddr(unbondingDelegation.ValidatorAddress)

		tempUnbondingDelegations := &models.UnbondingDelegations{
			DelegatorAddress: unbondingDelegation.DelegatorAddress,
			ValidatorAddress: unbondingDelegation.ValidatorAddress,
			Moniker:          validator.Moniker,
			Entries:          unbondingDelegation.Entries,
		}
		result = append(result, *tempUnbondingDelegations)
	}

	utils.Respond(w, result)
	return nil
}
