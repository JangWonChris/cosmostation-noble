package lcd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	dtypes "github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/utils"

	"github.com/go-pg/pg"
	resty "gopkg.in/resty.v1"
)

// SaveBondedValidators saves bonded validators information in database
func SaveBondedValidators(db *pg.DB, config *config.Config) {
	bondedResp, err := resty.R().Get(config.Node.LCDURL + "/staking/validators?status=bonded")
	if err != nil {
		fmt.Printf("query /staking/validators?status=bonded error - %v\n", err)
	}

	var responseWithHeight dtypes.ResponseWithHeight
	_ = json.Unmarshal(bondedResp.Body(), &responseWithHeight)

	var bondedValidators []*dtypes.Validator
	err = json.Unmarshal(responseWithHeight.Result, &bondedValidators)
	if err != nil {
		fmt.Printf("unmarshal bondedValidators error - %v\n", err)
	}

	// sort out bondedValidators by highest tokens
	sort.Slice(bondedValidators[:], func(i, j int) bool {
		tempToken1, _ := strconv.Atoi(bondedValidators[i].Tokens)
		tempToken2, _ := strconv.Atoi(bondedValidators[j].Tokens)
		return tempToken1 > tempToken2
	})

	// bondedValidator information for our database table
	validatorInfo := make([]*dtypes.ValidatorInfo, 0)
	for i, bondedValidators := range bondedValidators {
		tempValidatorInfo := &dtypes.ValidatorInfo{
			Rank:                 i + 1,
			OperatorAddress:      bondedValidators.OperatorAddress,
			Address:              utils.AccAddressFromOperatorAddress(bondedValidators.OperatorAddress),
			ConsensusPubkey:      bondedValidators.ConsensusPubkey,
			Proposer:             utils.ConsAddrFromConsPubkey(bondedValidators.ConsensusPubkey),
			Jailed:               bondedValidators.Jailed,
			Status:               bondedValidators.Status,
			Tokens:               bondedValidators.Tokens,
			DelegatorShares:      bondedValidators.DelegatorShares,
			Moniker:              bondedValidators.Description.Moniker,
			Identity:             bondedValidators.Description.Identity,
			Website:              bondedValidators.Description.Website,
			Details:              bondedValidators.Description.Details,
			UnbondingHeight:      bondedValidators.UnbondingHeight,
			UnbondingTime:        bondedValidators.UnbondingTime,
			CommissionRate:       bondedValidators.Commission.CommissionRates.Rate,
			CommissionMaxRate:    bondedValidators.Commission.CommissionRates.MaxRate,
			CommissionChangeRate: bondedValidators.Commission.CommissionRates.MaxChangeRate,
			MinSelfDelegation:    bondedValidators.MinSelfDelegation,
			UpdateTime:           bondedValidators.Commission.UpdateTime,
		}
		validatorInfo = append(validatorInfo, tempValidatorInfo)
	}

	// save and update validatorInfo
	if len(validatorInfo) > 0 {
		_, err = db.Model(&validatorInfo).
			OnConflict("(operator_address) DO UPDATE").
			Set("rank = EXCLUDED.rank").
			Set("consensus_pubkey = EXCLUDED.consensus_pubkey").
			Set("proposer = EXCLUDED.proposer").
			Set("jailed = EXCLUDED.jailed").
			Set("status = EXCLUDED.status").
			Set("tokens = EXCLUDED.tokens").
			Set("delegator_shares = EXCLUDED.delegator_shares").
			Set("moniker = EXCLUDED.moniker").
			Set("identity = EXCLUDED.identity").
			Set("website = EXCLUDED.website").
			Set("details = EXCLUDED.details").
			Set("unbonding_height = EXCLUDED.unbonding_height").
			Set("unbonding_time = EXCLUDED.unbonding_time").
			Set("commission_rate = EXCLUDED.commission_rate").
			Set("commission_max_rate = EXCLUDED.commission_max_rate").
			Set("update_time = EXCLUDED.update_time").
			Set("min_self_delegation = EXCLUDED.min_self_delegation").
			Insert()
		if err != nil {
			fmt.Printf("error - sync validators: %v\n", err)
		}
	}
}

// SaveUnbondingValidators saves unbonding validators information in database
func SaveUnbondingValidators(db *pg.DB, config *config.Config) {
	unbondingResp, err := resty.R().Get(config.Node.LCDURL + "/staking/validators?status=unbonding")
	if err != nil {
		fmt.Printf("Query /staking/validators?status=unbonding error - %v\n", err)
	}

	var responseWithHeight dtypes.ResponseWithHeight
	_ = json.Unmarshal(unbondingResp.Body(), &responseWithHeight)

	var unbondingValidators []*dtypes.Validator
	err = json.Unmarshal(responseWithHeight.Result, &unbondingValidators)
	if err != nil {
		fmt.Printf("unmarshal unbondingValidators error - %v\n", err)
	}

	// sort out bondedValidators by highest tokens
	sort.Slice(unbondingValidators[:], func(i, j int) bool {
		tempToken1, _ := strconv.Atoi(unbondingValidators[i].Tokens)
		tempToken2, _ := strconv.Atoi(unbondingValidators[j].Tokens)
		return tempToken1 > tempToken2
	})

	// validators information for our database table
	validatorInfo := make([]*dtypes.ValidatorInfo, 0)
	if len(unbondingValidators) > 0 {
		for _, unbondingValidator := range unbondingValidators {
			tempValidatorInfo := &dtypes.ValidatorInfo{
				OperatorAddress:      unbondingValidator.OperatorAddress,
				Address:              utils.AccAddressFromOperatorAddress(unbondingValidator.OperatorAddress),
				ConsensusPubkey:      unbondingValidator.ConsensusPubkey,
				Proposer:             utils.ConsAddrFromConsPubkey(unbondingValidator.ConsensusPubkey),
				Jailed:               unbondingValidator.Jailed,
				Status:               unbondingValidator.Status,
				Tokens:               unbondingValidator.Tokens,
				DelegatorShares:      unbondingValidator.DelegatorShares,
				Moniker:              unbondingValidator.Description.Moniker,
				Identity:             unbondingValidator.Description.Identity,
				Website:              unbondingValidator.Description.Website,
				Details:              unbondingValidator.Description.Details,
				UnbondingHeight:      unbondingValidator.UnbondingHeight,
				UnbondingTime:        unbondingValidator.UnbondingTime,
				CommissionRate:       unbondingValidator.Commission.CommissionRates.Rate,
				CommissionMaxRate:    unbondingValidator.Commission.CommissionRates.MaxRate,
				CommissionChangeRate: unbondingValidator.Commission.CommissionRates.MaxChangeRate,
				MinSelfDelegation:    unbondingValidator.MinSelfDelegation,
				UpdateTime:           unbondingValidator.Commission.UpdateTime,
			}
			validatorInfo = append(validatorInfo, tempValidatorInfo)
		}
	}

	// ranking
	var rankInfo dtypes.ValidatorInfo
	_ = db.Model(&rankInfo).
		Order("rank DESC").
		Limit(1).
		Select()

	for i, validatorInfo := range validatorInfo {
		validatorInfo.Rank = (rankInfo.Rank + 1 + i)
	}

	// save and update validatorInfo
	if len(validatorInfo) > 0 {
		_, err := db.Model(&validatorInfo).
			OnConflict("(operator_address) DO UPDATE").
			Set("rank = EXCLUDED.rank").
			Set("consensus_pubkey = EXCLUDED.consensus_pubkey").
			Set("proposer = EXCLUDED.proposer").
			Set("jailed = EXCLUDED.jailed").
			Set("status = EXCLUDED.status").
			Set("tokens = EXCLUDED.tokens").
			Set("delegator_shares = EXCLUDED.delegator_shares").
			Set("moniker = EXCLUDED.moniker").
			Set("identity = EXCLUDED.identity").
			Set("website = EXCLUDED.website").
			Set("details = EXCLUDED.details").
			Set("unbonding_height = EXCLUDED.unbonding_height").
			Set("unbonding_time = EXCLUDED.unbonding_time").
			Set("commission_rate = EXCLUDED.commission_rate").
			Set("commission_max_rate = EXCLUDED.commission_max_rate").
			Set("update_time = EXCLUDED.update_time").
			Set("min_self_delegation = EXCLUDED.min_self_delegation").
			Insert()
		if err != nil {
			fmt.Printf("error - save and update validatorinfo: %v\n", err)
		}
	}
}

// SaveUnbondedValidators saves unbonded validators information in database
func SaveUnbondedValidators(db *pg.DB, config *config.Config) {
	unbondedResp, err := resty.R().Get(config.Node.LCDURL + "/staking/validators?status=unbonded")
	if err != nil {
		fmt.Printf("Query /staking/validators?status=unbonded error - %v\n", err)
	}

	var responseWithHeightUnbonded dtypes.ResponseWithHeight
	_ = json.Unmarshal(unbondedResp.Body(), &responseWithHeightUnbonded)

	var unbondedValidators []*dtypes.Validator
	err = json.Unmarshal(responseWithHeightUnbonded.Result, &unbondedValidators)
	if err != nil {
		fmt.Printf("unmarshal unbondedValidators error - %v\n", err)
	}

	// sort out bondedValidators by highest tokens
	sort.Slice(unbondedValidators[:], func(i, j int) bool {
		tempToken1, _ := strconv.Atoi(unbondedValidators[i].Tokens)
		tempToken2, _ := strconv.Atoi(unbondedValidators[j].Tokens)
		return tempToken1 > tempToken2
	})

	// validators information for our database table
	validatorInfo := make([]*dtypes.ValidatorInfo, 0)
	if len(unbondedValidators) > 0 {
		for _, unbondedValidator := range unbondedValidators {
			tempValidatorInfo := &dtypes.ValidatorInfo{
				OperatorAddress:      unbondedValidator.OperatorAddress,
				Address:              utils.AccAddressFromOperatorAddress(unbondedValidator.OperatorAddress),
				ConsensusPubkey:      unbondedValidator.ConsensusPubkey,
				Proposer:             utils.ConsAddrFromConsPubkey(unbondedValidator.ConsensusPubkey),
				Jailed:               unbondedValidator.Jailed,
				Status:               unbondedValidator.Status,
				Tokens:               unbondedValidator.Tokens,
				DelegatorShares:      unbondedValidator.DelegatorShares,
				Moniker:              unbondedValidator.Description.Moniker,
				Identity:             unbondedValidator.Description.Identity,
				Website:              unbondedValidator.Description.Website,
				Details:              unbondedValidator.Description.Details,
				UnbondingHeight:      unbondedValidator.UnbondingHeight,
				UnbondingTime:        unbondedValidator.UnbondingTime,
				CommissionRate:       unbondedValidator.Commission.CommissionRates.Rate,
				CommissionMaxRate:    unbondedValidator.Commission.CommissionRates.MaxRate,
				CommissionChangeRate: unbondedValidator.Commission.CommissionRates.MaxChangeRate,
				MinSelfDelegation:    unbondedValidator.MinSelfDelegation,
				UpdateTime:           unbondedValidator.Commission.UpdateTime,
			}
			validatorInfo = append(validatorInfo, tempValidatorInfo)
		}
	}

	// ranking
	var rankInfo dtypes.ValidatorInfo
	_ = db.Model(&rankInfo).
		Order("rank DESC").
		Limit(1).
		Select()

	for i, validatorInfo := range validatorInfo {
		validatorInfo.Rank = (rankInfo.Rank + 1 + i)
	}

	// save and update validatorInfo
	if len(validatorInfo) > 0 {
		_, err := db.Model(&validatorInfo).
			OnConflict("(operator_address) DO UPDATE").
			Set("rank = EXCLUDED.rank").
			Set("consensus_pubkey = EXCLUDED.consensus_pubkey").
			Set("proposer = EXCLUDED.proposer").
			Set("jailed = EXCLUDED.jailed").
			Set("status = EXCLUDED.status").
			Set("tokens = EXCLUDED.tokens").
			Set("delegator_shares = EXCLUDED.delegator_shares").
			Set("moniker = EXCLUDED.moniker").
			Set("identity = EXCLUDED.identity").
			Set("website = EXCLUDED.website").
			Set("details = EXCLUDED.details").
			Set("unbonding_height = EXCLUDED.unbonding_height").
			Set("unbonding_time = EXCLUDED.unbonding_time").
			Set("commission_rate = EXCLUDED.commission_rate").
			Set("commission_max_rate = EXCLUDED.commission_max_rate").
			Set("update_time = EXCLUDED.update_time").
			Set("min_self_delegation = EXCLUDED.min_self_delegation").
			Insert()
		if err != nil {
			fmt.Printf("error - save and update validatorinfo: %v\n", err)
		}
	}
}
