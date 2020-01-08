package exporter

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/types"

	resty "gopkg.in/resty.v1"
)

// SaveValidatorsStats1H saves validator statistics every hour
func (ses *StatsExporterService) SaveValidatorsStats1H() {
	log.Println("Save Validator Stats 1H")

	// query all validators order by their tokens
	var validators []schema.ValidatorInfo
	err := ses.db.Model(&validators).
		Order("rank ASC").
		Select()
	if err != nil {
		fmt.Printf("ValidatorInfo DB error - %v\n", err)
	}

	validatorStats := make([]*types.StatsValidators1H, 0)
	for _, validator := range validators {
		// get self-bonded amount by querying the current delegation between a delegator and a validator
		selfBondedResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/staking/delegators/" + validator.Address + "/delegations/" + validator.OperatorAddress)

		var responseWithHeight types.ResponseWithHeight
		_ = json.Unmarshal(selfBondedResp.Body(), &responseWithHeight)

		var delegatorDelegation types.DelegatorDelegation
		err := json.Unmarshal(responseWithHeight.Result, &delegatorDelegation)
		if err != nil { // "error": "{"codespace":"staking","code":102,"message":"no delegation for this (address, validator) pair"}"
			fmt.Printf("unmarshal delegatorDelegation error - %v, validator address - %v, operator address - %v\n", err, validator.Address, validator.OperatorAddress)
		}

		// get all validator's delegations
		valiDelegationResp, err := resty.R().Get(ses.config.Node.LCDURL + "/staking/validators/" + validator.OperatorAddress + "/delegations")

		var responseWithHeight2 types.ResponseWithHeight
		_ = json.Unmarshal(valiDelegationResp.Body(), &responseWithHeight2)

		var validatorDelegations []types.ValidatorDelegation
		err = json.Unmarshal(responseWithHeight2.Result, &validatorDelegations)
		if err != nil {
			fmt.Printf("unmarshal validatorDelegations error - %v\n", err)
		}

		// initialize variables, otherwise throws an error if there is no delegations
		var selfBondedAmount float64
		var othersAmount float64
		var totalDelegationAmount float64

		// LCD 요청 값이 없을 경우 아래와 같이 에러가 발생하므로 아래 if 문으로 에러 처리. 추후에 Cosmos SDK 에서 에러 처리를 변경 할 걸로 보인다.
		// {  %!v(PANIC=Format method: runtime error: invalid memory address or nil pointer dereference)}
		if delegatorDelegation.DelegatorAddress != "" {
			selfBondedAmount, _ = strconv.ParseFloat(delegatorDelegation.Shares, 64)
		}

		// delegator shares
		if len(validatorDelegations) > 0 {
			for _, validatorDelegation := range validatorDelegations {
				if validatorDelegation.DelegatorAddress != validator.Address {
					shares, _ := strconv.ParseFloat(validatorDelegation.Shares, 64)
					othersAmount += shares
				}
			}
		}

		totalDelegationAmount = selfBondedAmount + othersAmount

		// delegator numbers
		delegatorNum := len(validatorDelegations)

		tempValidatorStats := &types.StatsValidators1H{
			Moniker:          validator.Moniker,
			OperatorAddress:  validator.OperatorAddress,
			Address:          validator.Address,
			Proposer:         validator.Proposer,
			ConsensusPubkey:  validator.ConsensusPubkey,
			TotalDelegations: totalDelegationAmount,
			SelfBonded:       selfBondedAmount,
			Others:           othersAmount,
			DelegatorNum:     delegatorNum,
			Time:             time.Now(),
		}
		validatorStats = append(validatorStats, tempValidatorStats)
	}

	_, err = ses.db.Model(&validatorStats).Insert()
	if err != nil {
		fmt.Printf("save ValidatorStats1H error - %v\n", err)
	}
}

// SaveValidatorsStats24H saves validator statistics 24 hours
func (ses *StatsExporterService) SaveValidatorsStats24H() {
	log.Println("Save Validator Stats 24H")

	// query all validators order by their tokens
	var validators []schema.ValidatorInfo
	err := ses.db.Model(&validators).
		Order("rank ASC").
		Select()
	if err != nil {
		fmt.Printf("ValidatorInfo DB error - %v\n", err)
	}

	validatorStats := make([]*types.StatsValidators24H, 0)
	for _, validator := range validators {
		// get self-bonded amount by querying the current delegation between a delegator and a validator
		selfBondedResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/staking/delegators/" + validator.Address + "/delegations/" + validator.OperatorAddress)

		var responseWithHeight types.ResponseWithHeight
		_ = json.Unmarshal(selfBondedResp.Body(), &responseWithHeight)

		var delegatorDelegation types.DelegatorDelegation
		err := json.Unmarshal(responseWithHeight.Result, &delegatorDelegation)
		if err != nil {
			fmt.Printf("unmarshal delegatorDelegation error - %v, validator address - %v, operator address - %v\n", err, validator.Address, validator.OperatorAddress)
		}

		// get all validator's delegations
		valiDelegationResp, err := resty.R().Get(ses.config.Node.LCDURL + "/staking/validators/" + validator.OperatorAddress + "/delegations")

		var responseWithHeight2 types.ResponseWithHeight
		_ = json.Unmarshal(valiDelegationResp.Body(), &responseWithHeight2)

		var validatorDelegations []types.ValidatorDelegation
		err = json.Unmarshal(responseWithHeight2.Result, &validatorDelegations)
		if err != nil {
			fmt.Printf("unmarshal validatorDelegations error - %v\n", err)
		}

		// initialize variables, otherwise throws an error if there is no delegations
		var selfBondedAmount float64
		var othersAmount float64
		var totalDelegationAmount float64

		// LCD 요청 값이 없을 경우 아래와 같이 에러가 발생하므로 아래 if 문으로 에러 처리. 추후에 Cosmos SDK 에서 에러 처리를 변경 할 걸로 보인다.
		// {  %!v(PANIC=Format method: runtime error: invalid memory address or nil pointer dereference)}
		if delegatorDelegation.DelegatorAddress != "" {
			selfBondedAmount, _ = strconv.ParseFloat(delegatorDelegation.Shares, 64)
		}

		// delegator shares
		if len(validatorDelegations) > 0 {
			for _, validatorDelegation := range validatorDelegations {
				if validatorDelegation.DelegatorAddress != validator.Address {
					shares, _ := strconv.ParseFloat(validatorDelegation.Shares, 64)
					othersAmount += shares
				}
			}
		}

		totalDelegationAmount = selfBondedAmount + othersAmount

		// delegator numbers
		delegatorNum := len(validatorDelegations)

		tempValidatorStats := &types.StatsValidators24H{
			Moniker:          validator.Moniker,
			OperatorAddress:  validator.OperatorAddress,
			Address:          validator.Address,
			Proposer:         validator.Proposer,
			ConsensusPubkey:  validator.ConsensusPubkey,
			TotalDelegations: totalDelegationAmount,
			SelfBonded:       selfBondedAmount,
			Others:           othersAmount,
			DelegatorNum:     delegatorNum,
			Time:             time.Now(),
		}
		validatorStats = append(validatorStats, tempValidatorStats)
	}

	_, err = ses.db.Model(&validatorStats).Insert()
	if err != nil {
		fmt.Printf("save ValidatorStats24H error - %v\n", err)
	}
}
