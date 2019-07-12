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

// SaveBondedValidators queries the validators information from LCD and stores them in the database
func SaveBondedValidators(db *pg.DB, config *config.Config) {
	bondedResp, err := resty.R().Get(config.Node.LCDURL + "/staking/validators?status=bonded")
	if err != nil {
		fmt.Printf("LCD resty - %v\n", err)
	}

	// Parse Validator struct
	var bondedValidators []*dtypes.Validator
	err = json.Unmarshal(bondedResp.Body(), &bondedValidators)
	if err != nil {
		fmt.Printf("Unmarshal - %v\n", err)
	}

	// Sort bondedValidators by highest tokens
	sort.Slice(bondedValidators[:], func(i, j int) bool {
		tempToken1, _ := strconv.Atoi(bondedValidators[i].Tokens)
		tempToken2, _ := strconv.Atoi(bondedValidators[j].Tokens)
		return tempToken1 > tempToken2
	})

	// BondedValidator information
	validatorInfo := make([]*dtypes.ValidatorInfo, 0)
	for i, bondedValidators := range bondedValidators {
		tempValidatorInfo := &dtypes.ValidatorInfo{
			Rank:                 i + 1,
			OperatorAddress:      bondedValidators.OperatorAddress,
			Address:        utils.OperatorAddressToAddress(bondedValidators.OperatorAddress),
			ConsensusPubkey:      bondedValidators.ConsensusPubkey,
			Proposer:             utils.ConsensusPubkeyToProposer(bondedValidators.ConsensusPubkey),
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
			CommissionRate:       bondedValidators.Commission.Rate,
			CommissionMaxRate:    bondedValidators.Commission.MaxRate,
			CommissionChangeRate: bondedValidators.Commission.MaxChangeRate,
			MinSelfDelegation:    bondedValidators.MinSelfDelegation,
			UpdateTime:           bondedValidators.Commission.UpdateTime,
		}
		validatorInfo = append(validatorInfo, tempValidatorInfo)
	}

	// Save & Update validatorInfo
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

// SaveUnbondedAndUnbodingValidators queries the validators information from LCD and stores them in the database
func SaveUnbondedAndUnbodingValidators(db *pg.DB, config *config.Config) {
	unbondedResp, _ := resty.R().Get(config.Node.LCDURL + "/staking/validators?status=unbonded")
	unbondingResp, _ := resty.R().Get(config.Node.LCDURL + "/staking/validators?status=unbonding")

	// Parse Unbonded Validator struct
	var unbondedValidators []*dtypes.Validator
	_ = json.Unmarshal(unbondedResp.Body(), &unbondedValidators)

	// Parse Unbonding Validator struct
	var unbondingValidators []*dtypes.Validator
	_ = json.Unmarshal(unbondingResp.Body(), &unbondingValidators)

	// Validators information
	validatorInfo := make([]*dtypes.ValidatorInfo, 0)
	for _, unbondedValidator := range unbondedValidators {
		tempValidatorInfo := &dtypes.ValidatorInfo{
			OperatorAddress:      unbondedValidator.OperatorAddress,
			Address:        utils.OperatorAddressToAddress(unbondedValidator.OperatorAddress),
			ConsensusPubkey:      unbondedValidator.ConsensusPubkey,
			Proposer:             utils.ConsensusPubkeyToProposer(unbondedValidator.ConsensusPubkey),
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
			CommissionRate:       unbondedValidator.Commission.Rate,
			CommissionMaxRate:    unbondedValidator.Commission.MaxRate,
			CommissionChangeRate: unbondedValidator.Commission.MaxChangeRate,
			MinSelfDelegation:    unbondedValidator.MinSelfDelegation,
			UpdateTime:           unbondedValidator.Commission.UpdateTime,
		}
		validatorInfo = append(validatorInfo, tempValidatorInfo)
	}

	for _, unbondingValidator := range unbondingValidators {
		tempValidatorInfo := &dtypes.ValidatorInfo{
			OperatorAddress:      unbondingValidator.OperatorAddress,
			Address:        utils.OperatorAddressToAddress(unbondingValidator.OperatorAddress),
			ConsensusPubkey:      unbondingValidator.ConsensusPubkey,
			Proposer:             utils.ConsensusPubkeyToProposer(unbondingValidator.ConsensusPubkey),
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
			CommissionRate:       unbondingValidator.Commission.Rate,
			CommissionMaxRate:    unbondingValidator.Commission.MaxRate,
			CommissionChangeRate: unbondingValidator.Commission.MaxChangeRate,
			MinSelfDelegation:    unbondingValidator.MinSelfDelegation,
			UpdateTime:           unbondingValidator.Commission.UpdateTime,
		}
		validatorInfo = append(validatorInfo, tempValidatorInfo)
	}

	// Sort bondedValidators by highest tokens
	sort.Slice(validatorInfo[:], func(i, j int) bool {
		tempToken1, _ := strconv.Atoi(validatorInfo[i].Tokens)
		tempToken2, _ := strconv.Atoi(validatorInfo[j].Tokens)
		return tempToken1 > tempToken2
	})

	for i, validatorInfo := range validatorInfo {
		validatorInfo.Rank = (101 + i)
	}

	// Save & Update validatorInfo
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
