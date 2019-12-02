package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/client/common"
	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/types"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/utils"

	"github.com/tendermint/tendermint/rpc/client"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	resty "gopkg.in/resty.v1"
)

// GetBalance returns balance of an anddress
func GetBalance(codec *codec.Codec, config *config.Config, db *pg.DB, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	accAddress := vars["accAddress"]

	if !strings.Contains(accAddress, sdk.GetConfig().GetBech32AccountAddrPrefix()) || len(accAddress) != 45 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	resp, _ := resty.R().Get(config.Node.LCDURL + "/bank/balances/" + accAddress)

	var balances []models.Coin
	err := json.Unmarshal(types.ReadRespWithHeight(resp).Result, &balances)
	if err != nil {
		log.Info().Str(models.Service, models.LogAccount).Str(models.Method, "GetBalance").Err(err).Msg("unmarshal balances error")
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
func GetDelegationsRewards(codec *codec.Codec, config *config.Config, db *pg.DB, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	accAddress := vars["accAddress"]

	if !strings.Contains(accAddress, sdk.GetConfig().GetBech32AccountAddrPrefix()) || len(accAddress) != 45 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	resp, _ := resty.R().Get(config.Node.LCDURL + "/distribution/delegators/" + accAddress + "/rewards")

	var resultRewards models.ResultRewards
	err := json.Unmarshal(types.ReadRespWithHeight(resp).Result, &resultRewards)
	if err != nil {
		log.Info().Str(models.Service, models.LogAccount).Str(models.Method, "GetDelegationsRewards").Err(err).Msg("unmarshal resultRewards error")
	}

	resultDelegatorRewards := make([]models.Rewards, 0)
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

		resultDelegatorRewards = append(resultDelegatorRewards, *tempReward)
	}

	utils.Respond(w, resultDelegatorRewards)
	return nil
}

// GetDelegations returns all delegations from an address
func GetDelegations(codec *codec.Codec, config *config.Config, db *pg.DB, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	accAddress := vars["accAddress"]

	if !strings.Contains(accAddress, sdk.GetConfig().GetBech32AccountAddrPrefix()) || len(accAddress) != 45 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// query delegations and each delegator's rewards
	resp, _ := resty.R().Get(config.Node.LCDURL + "/staking/delegators/" + accAddress + "/delegations")

	delegations := make([]models.Delegations, 0)
	err := json.Unmarshal(types.ReadRespWithHeight(resp).Result, &delegations)
	if err != nil {
		log.Info().Str(models.Service, models.LogAccount).Str(models.Method, "GetDelegations").Err(err).Msg("unmarshal delegations error")
	}

	resultDelegations := make([]models.ResultDelegations, 0)
	if len(delegations) > 0 {
		for _, delegation := range delegations {
			var validatorInfo types.ValidatorInfo
			_ = db.Model(&validatorInfo).
				Column("moniker"). // query validator's moniker
				Where("operator_address = ?", delegation.ValidatorAddress).
				Limit(1).
				Select()

			// query rewards
			rewardsResp, _ := resty.R().Get(config.Node.LCDURL + "/distribution/delegators/" + accAddress + "/rewards/" + delegation.ValidatorAddress)

			var rewards []models.Coin
			err = json.Unmarshal(types.ReadRespWithHeight(rewardsResp).Result, &rewards)
			if err != nil {
				log.Info().Str(models.Service, models.LogAccount).Str(models.Method, "GetDelegations").Err(err).Msg("unmarshal rewards error")
			}

			// if the fee of delegator's validator is 100%, then reward is null
			resultRewards := make([]models.Coin, 0)
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

			// query information of the validator
			validatorResp, _ := resty.R().Get(config.Node.LCDURL + "/staking/validators/" + delegation.ValidatorAddress)

			var validator types.Validator
			err = json.Unmarshal(types.ReadRespWithHeight(validatorResp).Result, &validator)
			if err != nil {
				log.Info().Str(models.Service, models.LogAccount).Str(models.Method, "GetDelegations").Err(err).Msg("unmarshal validator error")
			}

			// validator's token divide by delegator_shares equals amount of uatom
			tokens, _ := strconv.ParseFloat(validator.Tokens, 64)
			delegatorShares, _ := strconv.ParseFloat(validator.DelegatorShares, 64)
			uatom := tokens / delegatorShares
			shares, _ := strconv.ParseFloat(delegation.Shares, 64)
			amount := fmt.Sprintf("%f", shares*uatom)

			tempResultDelegations := &models.ResultDelegations{
				DelegatorAddress: delegation.DelegatorAddress,
				ValidatorAddress: delegation.ValidatorAddress,
				Moniker:          validatorInfo.Moniker,
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
func GetCommission(codec *codec.Codec, config *config.Config, db *pg.DB, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	accAddress := vars["accAddress"]

	if !strings.Contains(accAddress, sdk.GetConfig().GetBech32AccountAddrPrefix()) || len(accAddress) != 45 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	valAddress := utils.ValAddressFromAccAddress(accAddress)

	commission := make([]models.Coin, 0)
	if valAddress != "" {
		ctx := context.NewCLIContext().WithCodec(codec).WithClient(rpcClient)
		valAddr, _ := sdk.ValAddressFromBech32(valAddress)
		result, _ := common.QueryValidatorCommission(ctx, distr.QuerierRoute, valAddr)

		var valCom distrTypes.ValidatorAccumulatedCommission
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
func GetUnbondingDelegations(codec *codec.Codec, config *config.Config, db *pg.DB, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	accAddress := vars["accAddress"]

	if !strings.Contains(accAddress, sdk.GetConfig().GetBech32AccountAddrPrefix()) || len(accAddress) != 45 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// query unbonding delegations
	unbondingDelegationsResp, _ := resty.R().Get(config.Node.LCDURL + "/staking/delegators/" + accAddress + "/unbonding_delegations")

	unbondingDelegations := make([]models.UnbondingDelegations, 0)
	err := json.Unmarshal(types.ReadRespWithHeight(unbondingDelegationsResp).Result, &unbondingDelegations)
	if err != nil {
		log.Info().Str(models.Service, models.LogAccount).Str(models.Method, "GetUnbondingDelegations").Err(err).Msg("unmarshal unbondingDelegations error")
	}

	resultUnbondingDelegations := make([]models.UnbondingDelegations, 0)
	for _, unbondingDelegation := range unbondingDelegations {
		var validatorInfo types.ValidatorInfo
		_ = db.Model(&validatorInfo).
			Column("moniker").
			Where("operator_address = ?", unbondingDelegation.ValidatorAddress).
			Limit(1).
			Select()

		tempUnbondingDelegations := &models.UnbondingDelegations{
			DelegatorAddress: unbondingDelegation.DelegatorAddress,
			ValidatorAddress: unbondingDelegation.ValidatorAddress,
			Moniker:          validatorInfo.Moniker,
			Entries:          unbondingDelegation.Entries,
		}
		resultUnbondingDelegations = append(resultUnbondingDelegations, *tempUnbondingDelegations)
	}

	utils.Respond(w, resultUnbondingDelegations)
	return nil
}
