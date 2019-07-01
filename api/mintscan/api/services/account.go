package services

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models"
	u "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/utils"
	dbtypes "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/go-pg/pg"
	resty "gopkg.in/resty.v1"
)

func GetAccountInfo(DB *pg.DB, Config *config.Config, w http.ResponseWriter, r *http.Request) error {
	// Receive address
	vars := mux.Vars(r)
	address := vars["address"]

	// Check the validity of cosmos address
	if !strings.Contains(address, sdk.Bech32PrefixAccAddr) || len(address) != 45 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Final Account Response
	var accountResponse models.AccountResponse

	// Query LCD: Bank Balance
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	balanceResp, _ := resty.R().Get(Config.Node.LCDURL + "/bank/balances/" + address)

	var balance []models.Balance
	err := json.Unmarshal(balanceResp.Body(), &balance)
	if err != nil {
		fmt.Printf("Balance unmarshal error - %v\n", err)
	}

	// Return array result when empty
	if balance == nil {
		accountResponse.Balance = []models.Balance{}
	} else {
		accountResponse.Balance = balance
	}

	// Query LCD: Rewards
	rewardsResp, _ := resty.R().Get(Config.Node.LCDURL + "/distribution/delegators/" + address + "/rewards")

	var rewards []models.Rewards
	err = json.Unmarshal(rewardsResp.Body(), &rewards)
	if err != nil {
		fmt.Printf("Rewards unmarshal error - %v\n", err)
	}

	// Returns empty
	if rewards == nil {
		accountResponse.Rewards = []models.Rewards{}
	} else {
		accountResponse.Rewards = rewards
	}

	// Query LCD: Delegator Rewards
	delegationsResp, _ := resty.R().Get(Config.Node.LCDURL + "/staking/delegators/" + address + "/delegations")

	var delegations []models.Delegations
	err = json.Unmarshal(delegationsResp.Body(), &delegations)
	if err != nil {
		fmt.Printf("Delegations unmarshal error - %v\n", err)
	}

	var resultDelegations []models.Delegations
	for _, delegation := range delegations {
		delegatorRewardsResp, _ := resty.R().Get(Config.Node.LCDURL + "/distribution/delegators/" + address + "/rewards/" + delegation.ValidatorAddress)

		var delegatorRewards []models.Rewards
		err = json.Unmarshal(delegatorRewardsResp.Body(), &delegatorRewards)
		if err != nil {
			fmt.Printf("Distribution Rewards unmarshal error - %v\n", err)
		}

		var validatorInfo dbtypes.ValidatorInfo
		_ = DB.Model(&validatorInfo).
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
		validatorResp, _ := resty.R().Get(Config.Node.LCDURL + "/staking/validators/" + delegation.ValidatorAddress)

		var validator models.Validator
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

	// Returns empty
	if delegations == nil {
		accountResponse.Delegations = []models.Delegations{}
	} else {
		accountResponse.Delegations = resultDelegations
	}

	// Query LCD: Unbonding Delegations
	unbondingDelegationsResp, _ := resty.R().Get(Config.Node.LCDURL + "/staking/delegators/" + address + "/unbonding_delegations")

	var unbondingDelegations []models.UnbondingDelegations
	err = json.Unmarshal(unbondingDelegationsResp.Body(), &unbondingDelegations)
	if err != nil {
		fmt.Printf("UnbondingDelegations unmarshal error - %v\n", err)
	}

	var resultUnbondingDelegations []models.UnbondingDelegations
	for _, unbondingDelegation := range unbondingDelegations {
		var validatorInfo dbtypes.ValidatorInfo
		_ = DB.Model(&validatorInfo).
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

	// Returns empty
	if unbondingDelegations == nil {
		accountResponse.UnbondingDelegations = []models.UnbondingDelegations{}
	} else {
		accountResponse.UnbondingDelegations = resultUnbondingDelegations
	}

	u.Respond(w, accountResponse)
	return nil
}
