package exporter

import (
	"encoding/json"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/models"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/schema"
	"go.uber.org/zap"
)

// TODO: REST API 사용보다는 RPC로 해보는 건 어떤지. Delegations API 요청 시 800개가 넘는 Delegations 검증인들이 꾀나 되기 때문에
// 요청도 오래걸리고 계산 로직도 오래걸린다. 현재는 client 요청 시 timeout을 5초에서 10초로 늘려놔서 오래걸리더라도 문제는 없다.
// SaveValidatorsStats1H saves validators statstics every hour.
func (ex *Exporter) SaveValidatorsStats1H() {
	result := make([]schema.StatsValidators1H, 0)

	vals, err := ex.db.QueryValidatorsByStatus(models.BondedValidatorStatus)
	if err != nil {
		zap.S().Errorf("failed to query bonded validators: %s", err)
		return
	}

	if len(vals) <= 0 {
		zap.S().Info("found no validators in database")
		return
	}

	for _, val := range vals {
		var selfBondedAmount float64
		var othersAmount float64

		// Get the current delegation between a delegator and a validator (self-bonded).
		selfBondedResp, err := ex.client.RequestAPIFromLCDWithRespHeight("/staking/delegators/" + val.Address + "/delegations/" + val.OperatorAddress)
		if err != nil {
			zap.S().Errorf("failed to get the current delegation between a delegator and a validator: %s", err)
			return
		}

		var selfDelegation models.SelfDelegation
		err = json.Unmarshal(selfBondedResp.Result, &selfDelegation)
		if err != nil {
			zap.S().Errorf("failed to unmarshal self-bonded delegation: %s", err)
			return
		}

		// LCD 요청 값이 없을 경우 아래와 같이 에러가 발생하므로 아래 if 문으로 에러 처리. 추후에 Cosmos SDK 에서 에러 처리를 변경 할 걸로 보인다.
		// {  %!v(PANIC=Format method: runtime error: invalid memory address or nil pointer dereference)}.
		if selfDelegation.DelegatorAddress != "" {
			selfBondedAmount, _ = strconv.ParseFloat(selfDelegation.Shares, 64)
		}

		// Get the information from a single validator.
		valResp, err := ex.client.RequestAPIFromLCDWithRespHeight("/staking/validators/" + val.OperatorAddress)
		if err != nil {
			zap.S().Errorf("failed to get the validator information: %s", err)
			return
		}

		var valInfo models.Validator
		err = json.Unmarshal(valResp.Result, &valInfo)
		if err != nil {
			zap.S().Errorf("failed to unmarshal the validator information: %s", err)
			return
		}

		othersAmount, _ = strconv.ParseFloat(valInfo.DelegatorShares, 64)
		othersAmount = othersAmount - selfBondedAmount

		// Get all delegations from a validator.
		delgationsResp, err := ex.client.RequestAPIFromLCDWithRespHeight("/staking/validators/" + val.OperatorAddress + "/delegations")
		if err != nil {
			zap.S().Errorf("failed to get all delegations from the validator: %s", err)
			return
		}

		var valDelegations []models.ValidatorDelegation
		err = json.Unmarshal(delgationsResp.Result, &valDelegations)
		if err != nil {
			zap.S().Errorf("failed to unmarshal delegations for the validator: %s", err)
			return
		}

		sv := &schema.StatsValidators1H{
			Moniker:          val.Moniker,
			OperatorAddress:  val.OperatorAddress,
			Address:          val.Address,
			Proposer:         val.Proposer,
			ConsensusPubkey:  val.ConsensusPubkey,
			TotalDelegations: selfBondedAmount + othersAmount,
			SelfBonded:       selfBondedAmount,
			Others:           othersAmount,
			DelegatorNum:     len(valDelegations),
		}

		result = append(result, *sv)
	}

	err = ex.db.InsertValidatorStats1H(result)
	if err != nil {
		zap.S().Errorf("failed to save validators data: %s", err)
		return
	}

	zap.S().Info("successfully saved ValidatorsStats1H")
	return
}

// SaveValidatorsStats1D saves validators statstics every day.
func (ex *Exporter) SaveValidatorsStats1D() {
	result := make([]schema.StatsValidators1D, 0)

	vals, err := ex.db.QueryValidatorsByStatus(models.BondedValidatorStatus)
	if err != nil {
		zap.S().Errorf("failed to query bonded validators: %s", err)
		return
	}

	if len(vals) <= 0 {
		zap.S().Info("found no validators in database")
		return
	}

	for _, val := range vals {
		var selfBondedAmount float64
		var othersAmount float64

		// Get the current delegation between a delegator and a validator (self-bonded).
		selfBondedResp, err := ex.client.RequestAPIFromLCDWithRespHeight("/staking/delegators/" + val.Address + "/delegations/" + val.OperatorAddress)
		if err != nil {
			zap.S().Errorf("failed to get the current delegation between a delegator and a validator: %s", err)
			return
		}

		var selfDelegation models.SelfDelegation
		err = json.Unmarshal(selfBondedResp.Result, &selfDelegation)
		if err != nil {
			zap.S().Errorf("failed to unmarshal self-bonded delegation: %s", err)
			return
		}

		// LCD 요청 값이 없을 경우 아래와 같이 에러가 발생하므로 아래 if 문으로 에러 처리. 추후에 Cosmos SDK 에서 에러 처리를 변경 할 걸로 보인다.
		// {  %!v(PANIC=Format method: runtime error: invalid memory address or nil pointer dereference)}.
		if selfDelegation.DelegatorAddress != "" {
			selfBondedAmount, _ = strconv.ParseFloat(selfDelegation.Shares, 64)
		}

		// Get the information from a single validator.
		valResp, err := ex.client.RequestAPIFromLCDWithRespHeight("/staking/validators/" + val.OperatorAddress)
		if err != nil {
			zap.S().Errorf("failed to get the validator information: %s", err)
			return
		}

		var valInfo models.Validator
		err = json.Unmarshal(valResp.Result, &valInfo)
		if err != nil {
			zap.S().Errorf("failed to unmarshal the validator information: %s", err)
			return
		}

		othersAmount, _ = strconv.ParseFloat(valInfo.DelegatorShares, 64)
		othersAmount = othersAmount - selfBondedAmount

		// Get all delegations from a validator
		delgationsResp, err := ex.client.RequestAPIFromLCDWithRespHeight("/staking/validators/" + val.OperatorAddress + "/delegations")
		if err != nil {
			zap.S().Errorf("failed to get all delegations from the validator: %s", err)
			return
		}

		var valDelegations []models.ValidatorDelegation
		err = json.Unmarshal(delgationsResp.Result, &valDelegations)
		if err != nil {
			zap.S().Errorf("failed to unmarshal delegations for the validator: %s", err)
			return
		}

		sv := &schema.StatsValidators1D{
			Moniker:          val.Moniker,
			OperatorAddress:  val.OperatorAddress,
			Address:          val.Address,
			Proposer:         val.Proposer,
			ConsensusPubkey:  val.ConsensusPubkey,
			TotalDelegations: selfBondedAmount + othersAmount,
			SelfBonded:       selfBondedAmount,
			Others:           othersAmount,
			DelegatorNum:     len(valDelegations),
		}

		result = append(result, *sv)
	}

	err = ex.db.InsertValidatorStats1D(result)
	if err != nil {
		zap.S().Errorf("failed to save validators data: %s", err)
		return
	}

	zap.S().Info("successfully saved ValidatorsStats1D")
	return
}
