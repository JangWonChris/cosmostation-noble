package lcd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/utils"

	"github.com/go-pg/pg"
	"github.com/rs/zerolog/log"
	resty "gopkg.in/resty.v1"
)

// SaveBondedValidators saves bonded validators information in database
func SaveBondedValidators(db *pg.DB, config *config.Config) {
	resp, _ := resty.R().Get(config.Node.LCDURL + "/staking/validators?status=bonded")

	var bondedValidators []*types.Validator
	err := json.Unmarshal(types.ReadRespWithHeight(resp).Result, &bondedValidators)
	if err != nil {
		log.Info().Str(types.Service, types.LogValidator).Str(types.Method, "SaveBondedValidators").Err(err).Msg("unmarshal bondedValidators error")
	}

	// sort out bondedValidators by highest tokens
	sort.Slice(bondedValidators[:], func(i, j int) bool {
		tempToken1, _ := strconv.Atoi(bondedValidators[i].Tokens)
		tempToken2, _ := strconv.Atoi(bondedValidators[j].Tokens)
		return tempToken1 > tempToken2
	})

	// bondedValidator information for our database table
	validatorInfo := make([]*types.ValidatorInfo, 0)
	for i, bondedValidator := range bondedValidators {
		tempValidatorInfo := &types.ValidatorInfo{
			Rank:                 i + 1,
			OperatorAddress:      bondedValidator.OperatorAddress,
			Address:              utils.AccAddressFromOperatorAddress(bondedValidator.OperatorAddress),
			ConsensusPubkey:      bondedValidator.ConsensusPubkey,
			Proposer:             utils.ConsAddrFromConsPubkey(bondedValidator.ConsensusPubkey),
			Jailed:               bondedValidator.Jailed,
			Status:               bondedValidator.Status,
			Tokens:               bondedValidator.Tokens,
			DelegatorShares:      bondedValidator.DelegatorShares,
			Moniker:              bondedValidator.Description.Moniker,
			Identity:             bondedValidator.Description.Identity,
			Website:              bondedValidator.Description.Website,
			Details:              bondedValidator.Description.Details,
			UnbondingHeight:      bondedValidator.UnbondingHeight,
			UnbondingTime:        bondedValidator.UnbondingTime,
			CommissionRate:       bondedValidator.Commission.CommissionRates.Rate,
			CommissionMaxRate:    bondedValidator.Commission.CommissionRates.MaxRate,
			CommissionChangeRate: bondedValidator.Commission.CommissionRates.MaxChangeRate,
			MinSelfDelegation:    bondedValidator.MinSelfDelegation,
			UpdateTime:           bondedValidator.Commission.UpdateTime,
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
			Set("commission_change_rate = EXCLUDED.commission_change_rate").
			Set("update_time = EXCLUDED.update_time").
			Set("min_self_delegation = EXCLUDED.min_self_delegation").
			Insert()
		if err != nil {
			fmt.Printf("error - sync validators: %v\n", err)
		}
	}
}

// SaveUnbondingAndUnBondedValidators saves unbonding and unbonded validators information in database
func SaveUnbondingAndUnBondedValidators(db *pg.DB, config *config.Config) {
	resp, _ := resty.R().Get(config.Node.LCDURL + "/staking/validators?status=unbonding")

	var unbondingValidators []*types.Validator
	err := json.Unmarshal(types.ReadRespWithHeight(resp).Result, &unbondingValidators)
	if err != nil {
		log.Info().Str(types.Service, types.LogValidator).Str(types.Method, "SaveUnbondingAndUnBondedValidators").Err(err).Msg("unmarshal unbondingValidators error")
	}

	// sort out bondedValidators by highest tokens
	sort.Slice(unbondingValidators[:], func(i, j int) bool {
		tempToken1, _ := strconv.Atoi(unbondingValidators[i].Tokens)
		tempToken2, _ := strconv.Atoi(unbondingValidators[j].Tokens)
		return tempToken1 > tempToken2
	})

	// validators information for our database table
	validatorInfo := make([]*types.ValidatorInfo, 0)
	if len(unbondingValidators) > 0 {
		for _, unbondingValidator := range unbondingValidators {
			tempValidatorInfo := &types.ValidatorInfo{
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
	var rankInfo types.ValidatorInfo
	_ = db.Model(&rankInfo).
		Where("status = ?", 2).
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
			Set("commission_change_rate = EXCLUDED.commission_change_rate").
			Set("update_time = EXCLUDED.update_time").
			Set("min_self_delegation = EXCLUDED.min_self_delegation").
			Insert()

		// save unbonded validators after succesfully saved unbonding validators
		saveUnbondedValidators(db, config)

		if err != nil {
			fmt.Printf("error - save and update validatorinfo: %v\n", err)
		}
	}
}

// saveUnbondedValidators saves unbonded validators information in database
func saveUnbondedValidators(db *pg.DB, config *config.Config) {
	resp, _ := resty.R().Get(config.Node.LCDURL + "/staking/validators?status=unbonded")

	var unbondedValidators []*types.Validator
	err := json.Unmarshal(types.ReadRespWithHeight(resp).Result, &unbondedValidators)
	if err != nil {
		log.Info().Str(types.Service, types.LogValidator).Str(types.Method, "saveUnbondedValidators").Err(err).Msg("unmarshal unbondedValidators error")
	}

	// sort out bondedValidators by highest tokens
	sort.Slice(unbondedValidators[:], func(i, j int) bool {
		tempToken1, _ := strconv.Atoi(unbondedValidators[i].Tokens)
		tempToken2, _ := strconv.Atoi(unbondedValidators[j].Tokens)
		return tempToken1 > tempToken2
	})

	// validators information for our database table
	validatorInfo := make([]*types.ValidatorInfo, 0)
	if len(unbondedValidators) > 0 {
		for _, unbondedValidator := range unbondedValidators {
			tempValidatorInfo := &types.ValidatorInfo{
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
	var rankInfo types.ValidatorInfo
	_ = db.Model(&rankInfo).
		Where("status = ?", 1).
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
			Set("commission_change_rate = EXCLUDED.commission_change_rate").
			Set("update_time = EXCLUDED.update_time").
			Set("min_self_delegation = EXCLUDED.min_self_delegation").
			Insert()
		if err != nil {
			fmt.Printf("error - save and update validatorinfo: %v\n", err)
		}
	}
}
