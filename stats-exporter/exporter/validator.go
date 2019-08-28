package exporter

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/types"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/utils"

	resty "gopkg.in/resty.v1"
)

func (ses *StatsExporterService) SaveValidatorsStats1H() {
	log.Println("Save Validator Stats 1H")

	// Query all validators order by their tokens
	var validators []types.ValidatorInfo
	err := ses.db.Model(&validators).
		Order("rank ASC").
		Select()
	if err != nil {
		fmt.Printf("ValidatorInfo DB error - %v\n", err)
	}

	validatorStats := make([]*types.StatsValidators1H, 0)
	for _, validator := range validators {
		// validator's address
		address := utils.ConvertOperatorAddressToAddress(validator.OperatorAddress)

		// get self-bonded amount by querying the current delegation between a delegator and a validator
		var delegatorDelegation types.DelegatorDelegation
		selfBondedResp, err := resty.R().Get(ses.config.Node.LCDURL + "/staking/delegators/" + address + "/delegations/" + validator.OperatorAddress)
		if err != nil {
			fmt.Printf("Query /staking/delegators/{address}/delegations/{operatorAddr} error - %v\n", err)
		}

		err = json.Unmarshal(selfBondedResp.Body(), &delegatorDelegation)
		if err != nil {
			fmt.Printf("Unmarshal delegatorDelegation error - %v\n", err)
		}

		// get all validator's delegations
		var validatorDelegations []types.DelegatorDelegation
		valiDelegationResp, err := resty.R().Get(ses.config.Node.LCDURL + "/staking/validators/" + validator.OperatorAddress + "/delegations")
		err = json.Unmarshal(valiDelegationResp.Body(), &validatorDelegations)
		if err != nil {
			fmt.Printf("Unmarshal validatorDelegations error - %v\n", err)
		}

		// Initialize variables, otherwise throws an error if there is no delegations
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
				if validatorDelegation.DelegatorAddress != address {
					shares, _ := strconv.ParseFloat(validatorDelegation.Shares, 64)
					othersAmount += shares
				}
			}
		}

		// delegator numbers
		delegatorNum := len(validatorDelegations)

		tempValidatorStats := &types.StatsValidators1H{
			Moniker:          validator.Moniker,
			OperatorAddress:  validator.OperatorAddress,
			Address:          address,
			Proposer:         validator.Proposer,
			TotalDelegations: totalDelegationAmount,
			SelfBonded:       selfBondedAmount,
			Others:           othersAmount,
			DelegatorNum:     delegatorNum,
			Time:             time.Now(),
		}
		validatorStats = append(validatorStats, tempValidatorStats)
	}

	// Save
	_, err = ses.db.Model(&validatorStats).Insert()
	if err != nil {
		fmt.Printf("save ValidatorStats error - %v\n", err)
	}
}

func (ses *StatsExporterService) SaveValidatorsStats24H() {
	log.Println("Save Validator Stats 1H")

	// Query all validators order by their tokens
	var validators []types.ValidatorInfo
	err := ses.db.Model(&validators).
		Order("rank ASC").
		Select()
	if err != nil {
		fmt.Printf("ValidatorInfo DB error - %v\n", err)
	}

	validatorStats := make([]*types.StatsValidators24H, 0)
	for _, validator := range validators {
		// validator's address
		address := utils.ConvertOperatorAddressToAddress(validator.OperatorAddress)

		// get self-bonded amount by querying the current delegation between a delegator and a validator
		var delegatorDelegation types.DelegatorDelegation
		selfBondedResp, err := resty.R().Get(ses.config.Node.LCDURL + "/staking/delegators/" + address + "/delegations/" + validator.OperatorAddress)
		if err != nil {
			fmt.Printf("Query /staking/delegators/{address}/delegations/{operatorAddr} error - %v\n", err)
		}

		err = json.Unmarshal(selfBondedResp.Body(), &delegatorDelegation)
		if err != nil {
			fmt.Printf("Unmarshal delegatorDelegation error - %v\n", err)
		}

		// get all validator's delegations
		var validatorDelegations []types.DelegatorDelegation
		valiDelegationResp, err := resty.R().Get(ses.config.Node.LCDURL + "/staking/validators/" + validator.OperatorAddress + "/delegations")
		err = json.Unmarshal(valiDelegationResp.Body(), &validatorDelegations)
		if err != nil {
			fmt.Printf("Unmarshal validatorDelegations error - %v\n", err)
		}

		// Initialize variables, otherwise throws an error if there is no delegations
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
				if validatorDelegation.DelegatorAddress != address {
					shares, _ := strconv.ParseFloat(validatorDelegation.Shares, 64)
					othersAmount += shares
				}
			}
		}

		// delegator numbers
		delegatorNum := len(validatorDelegations)

		tempValidatorStats := &types.StatsValidators24H{
			Moniker:          validator.Moniker,
			OperatorAddress:  validator.OperatorAddress,
			Address:          address,
			Proposer:         validator.Proposer,
			TotalDelegations: totalDelegationAmount,
			SelfBonded:       selfBondedAmount,
			Others:           othersAmount,
			DelegatorNum:     delegatorNum,
			Time:             time.Now(),
		}
		validatorStats = append(validatorStats, tempValidatorStats)
	}

	// Save
	_, err = ses.db.Model(&validatorStats).Insert()
	if err != nil {
		fmt.Printf("save ValidatorStats error - %v\n", err)
	}
}
