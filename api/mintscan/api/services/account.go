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
	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models"
	dbtypes "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/types"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/utils"

	"github.com/cosmos/cosmos-sdk/x/distribution/client/common"
	"github.com/tendermint/tendermint/rpc/client"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	resty "gopkg.in/resty.v1"
)

// Balance, Rewards, Commission, Delegations, UnbondingDelegations
func GetAccountInfo(codec *codec.Codec, config *config.Config, db *pg.DB, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	address := vars["address"]

	// Check validity of address
	if !strings.Contains(address, sdk.GetConfig().GetBech32AccountAddrPrefix()) || len(address) != 45 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// ResultAccount Response
	var resultAccountResponse models.ResultAccountResponse

	// Query bank balance
	balance := make([]models.Balance, 0)
	balanceResp, _ := resty.R().Get(config.Node.LCDURL + "/bank/balances/" + address)
	err := json.Unmarshal(balanceResp.Body(), &balance)
	if err != nil {
		fmt.Printf("/bank/balances/ unmarshal error - %v\n", err)
	}
	resultAccountResponse.Balance = balance

	// Query rewards
	rewards := make([]models.Rewards, 0)
	rewardsResp, _ := resty.R().Get(config.Node.LCDURL + "/distribution/delegators/" + address + "/rewards")
	err = json.Unmarshal(rewardsResp.Body(), &rewards)
	if err != nil {
		fmt.Printf("/distribution/delegators/rewards unmarshal error - %v\n", err)
	}
	resultAccountResponse.Rewards = rewards

	// Query commission if an address is validator
	var validatorInfo dbtypes.ValidatorInfo
	err = db.Model(&validatorInfo).
		Column("operator_address").
		Where("address = ?", address).
		Select()

	commission := make([]models.Commission, 0)
	if validatorInfo.OperatorAddress != "" {
		ctx := context.NewCLIContext().WithCodec(codec).WithClient(rpcClient)
		valAddr, _ := sdk.ValAddressFromBech32(validatorInfo.OperatorAddress)
		result, _ := common.QueryValidatorCommission(ctx, codec, distr.QuerierRoute, valAddr)

		var valCom distrTypes.ValidatorAccumulatedCommission
		ctx.Codec.MustUnmarshalJSON(result, &valCom)

		tempCommission := &models.Commission{
			Denom:  valCom[0].Denom,
			Amount: valCom[0].Amount.String(),
		}

		commission = append(commission, *tempCommission)
	}
	resultAccountResponse.Commission = commission

	// Query delegations and each delegator's rewards
	delegations := make([]models.Delegations, 0)
	delegationsResp, _ := resty.R().Get(config.Node.LCDURL + "/staking/delegators/" + address + "/delegations")
	err = json.Unmarshal(delegationsResp.Body(), &delegations)
	if err != nil {
		fmt.Printf("Delegations unmarshal error - %v\n", err)
	}

	resultDelegations := make([]models.Delegations, 0)
	for _, delegation := range delegations {
		var delegatorRewards []models.Rewards
		delegatorRewardsResp, _ := resty.R().Get(config.Node.LCDURL + "/distribution/delegators/" + address + "/rewards/" + delegation.ValidatorAddress)
		err = json.Unmarshal(delegatorRewardsResp.Body(), &delegatorRewards)
		if err != nil {
			fmt.Printf("Distribution Rewards unmarshal error - %v\n", err)
		}

		var validatorInfo dbtypes.ValidatorInfo
		_ = db.Model(&validatorInfo).
			Column("moniker").
			Where("operator_address = ?", delegation.ValidatorAddress).
			Limit(1).
			Select()

		// If the fee of delegator's validator is 100%, then rewards LCD API returns null
		if len(delegatorRewards) > 0 {
			delegation.Rewards.Denom = delegatorRewards[0].Denom
			delegation.Rewards.Amount = delegatorRewards[0].Amount
		} else {
			delegation.Rewards.Denom = "0"
			delegation.Rewards.Amount = "0"
		}

		// Query a validator's information
		var validator models.Validator
		validatorResp, _ := resty.R().Get(config.Node.LCDURL + "/staking/validators/" + delegation.ValidatorAddress)
		err = json.Unmarshal(validatorResp.Body(), &validator)
		if err != nil {
			fmt.Printf("staking/validators/ unmarshal error - %v\n", err)
		}

		// Validator's token divide by delegator_shares equals amount of uatom
		tokens, _ := strconv.ParseFloat(validator.Tokens.String(), 64)
		delegatorShares, _ := strconv.ParseFloat(validator.DelegatorShares.String(), 64)
		uatom := tokens / delegatorShares
		shares, _ := strconv.ParseFloat(delegation.Shares, 64)
		amount := fmt.Sprintf("%f", shares*uatom)

		tempDelegations := &models.Delegations{
			DelegatorAddress: delegation.DelegatorAddress,
			ValidatorAddress: delegation.ValidatorAddress,
			Moniker:          validatorInfo.Moniker,
			Shares:           delegation.Shares,
			Amount:           amount,
			Rewards:          delegation.Rewards,
		}
		resultDelegations = append(resultDelegations, *tempDelegations)
	}

	resultAccountResponse.Delegations = resultDelegations

	// Query unbonding delegations
	unbondingDelegations := make([]models.UnbondingDelegations, 0)
	unbondingDelegationsResp, _ := resty.R().Get(config.Node.LCDURL + "/staking/delegators/" + address + "/unbonding_delegations")
	err = json.Unmarshal(unbondingDelegationsResp.Body(), &unbondingDelegations)
	if err != nil {
		fmt.Printf("UnbondingDelegations unmarshal error - %v\n", err)
	}

	resultUnbondingDelegations := make([]models.UnbondingDelegations, 0)
	for _, unbondingDelegation := range unbondingDelegations {
		var validatorInfo dbtypes.ValidatorInfo
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

	resultAccountResponse.UnbondingDelegations = resultUnbondingDelegations

	utils.Respond(w, resultAccountResponse)
	return nil
}
