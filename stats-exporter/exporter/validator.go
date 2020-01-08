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
	validatorStats := make([]*schema.StatsValidators1H, 0)

	// query 125 validators that are exist in the current network
	validators, _ := ses.db.QueryValidatorsByRank(125)

	for _, validator := range validators {
		var selfBondedAmount float64
		var othersAmount float64

		// reqeusts self-bonded amount by querying the current delegation between a delegator and a validator
		selfBondedResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/staking/delegators/" + validator.Address + "/delegations/" + validator.OperatorAddress)

		var delegatorDelegation types.DelegatorDelegation
		err := json.Unmarshal(types.ReadRespWithHeight(selfBondedResp).Result, &delegatorDelegation)
		if err != nil {
			fmt.Printf("failed to unmarshal delegatorDelegation: %v, \n", err)
			fmt.Printf("valAddr - %v, OperAddr: %v \n", validator.Address, validator.OperatorAddress)
		}

		// LCD 요청 값이 없을 경우 아래와 같이 에러가 발생하므로 아래 if 문으로 에러 처리. 추후에 Cosmos SDK 에서 에러 처리를 변경 할 걸로 보인다.
		// {  %!v(PANIC=Format method: runtime error: invalid memory address or nil pointer dereference)}
		if delegatorDelegation.DelegatorAddress != "" {
			selfBondedAmount, _ = strconv.ParseFloat(delegatorDelegation.Shares, 64)
		}

		// reqeusts validator information
		validatorResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/staking/validators/" + validator.OperatorAddress)

		var validatorInfo types.Validator
		err = json.Unmarshal(types.ReadRespWithHeight(validatorResp).Result, &validatorInfo)
		if err != nil {
			fmt.Printf("failed to unmarshal Validator: %v \n", err)
		}

		othersAmount, _ = strconv.ParseFloat(validatorInfo.DelegatorShares, 64)
		othersAmount = othersAmount - selfBondedAmount

		// reqeusts all validator's delegations to calculate delegatorNum
		valiDelegationResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/staking/validators/" + validator.OperatorAddress + "/delegations")

		var validatorDelegations []types.ValidatorDelegation
		err = json.Unmarshal(types.ReadRespWithHeight(valiDelegationResp).Result, &validatorDelegations)
		if err != nil {
			fmt.Printf("failed to unmarshal ValidatorDelegation: %v \n", err)
			fmt.Printf("OperAddr: %v \n", validator.OperatorAddress)
		}

		tempValidatorStats := &schema.StatsValidators1H{
			Moniker:          validator.Moniker,
			OperatorAddress:  validator.OperatorAddress,
			Address:          validator.Address,
			Proposer:         validator.Proposer,
			ConsensusPubkey:  validator.ConsensusPubkey,
			TotalDelegations: selfBondedAmount + othersAmount,
			SelfBonded:       selfBondedAmount,
			Others:           othersAmount,
			DelegatorNum:     len(validatorDelegations),
			Time:             time.Now(),
		}
		validatorStats = append(validatorStats, tempValidatorStats)
	}

	result, _ := ses.db.InsertValidatorStats1H(validatorStats)
	if result {
		log.Println("succesfully saved ValidatorStats 1H")
	}
}

// SaveValidatorsStats24H saves validator statistics 24 hours
func (ses *StatsExporterService) SaveValidatorsStats24H() {
	validatorStats := make([]*schema.StatsValidators24H, 0)

	// query 125 validators that are exist in the current network
	validators, _ := ses.db.QueryValidatorsByRank(125)

	for _, validator := range validators {
		var selfBondedAmount float64
		var othersAmount float64

		// reqeusts self-bonded amount by querying the current delegation between a delegator and a validator
		selfBondedResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/staking/delegators/" + validator.Address + "/delegations/" + validator.OperatorAddress)

		var delegatorDelegation types.DelegatorDelegation
		err := json.Unmarshal(types.ReadRespWithHeight(selfBondedResp).Result, &delegatorDelegation)
		if err != nil {
			fmt.Printf("failed to unmarshal delegatorDelegation: %v, \n", err)
			fmt.Printf("valAddr - %v, OperAddr: %v \n", validator.Address, validator.OperatorAddress)
		}

		// LCD 요청 값이 없을 경우 아래와 같이 에러가 발생하므로 아래 if 문으로 에러 처리. 추후에 Cosmos SDK 에서 에러 처리를 변경 할 걸로 보인다.
		// {  %!v(PANIC=Format method: runtime error: invalid memory address or nil pointer dereference)}
		if delegatorDelegation.DelegatorAddress != "" {
			selfBondedAmount, _ = strconv.ParseFloat(delegatorDelegation.Shares, 64)
		}

		// reqeusts validator information
		validatorResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/staking/validators/" + validator.OperatorAddress)

		var validatorInfo types.Validator
		err = json.Unmarshal(types.ReadRespWithHeight(validatorResp).Result, &validatorInfo)
		if err != nil {
			fmt.Printf("failed to unmarshal Validator: %v \n", err)
		}

		othersAmount, _ = strconv.ParseFloat(validatorInfo.DelegatorShares, 64)
		othersAmount = othersAmount - selfBondedAmount

		// reqeusts all validator's delegations to calculate delegatorNum
		valiDelegationResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/staking/validators/" + validator.OperatorAddress + "/delegations")

		var validatorDelegations []types.ValidatorDelegation
		err = json.Unmarshal(types.ReadRespWithHeight(valiDelegationResp).Result, &validatorDelegations)
		if err != nil {
			fmt.Printf("failed to unmarshal ValidatorDelegation: %v \n", err)
			fmt.Printf("OperAddr: %v \n", validator.OperatorAddress)
		}

		tempValidatorStats := &schema.StatsValidators24H{
			Moniker:          validator.Moniker,
			OperatorAddress:  validator.OperatorAddress,
			Address:          validator.Address,
			Proposer:         validator.Proposer,
			ConsensusPubkey:  validator.ConsensusPubkey,
			TotalDelegations: selfBondedAmount + othersAmount,
			SelfBonded:       selfBondedAmount,
			Others:           othersAmount,
			DelegatorNum:     len(validatorDelegations),
			Time:             time.Now(),
		}
		validatorStats = append(validatorStats, tempValidatorStats)
	}

	result, _ := ses.db.InsertValidatorStats24H(validatorStats)
	if result {
		log.Println("succesfully saved ValidatorStats 24H")
	}
}
