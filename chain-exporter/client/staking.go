package client

import (
	"context"
	"fmt"
	"log"
	"sort"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
)

// GetStakingQueryClient returns a object of queryClient
func (c *Client) GetStakingQueryClient() stakingtypes.QueryClient {
	return stakingtypes.NewQueryClient(c.grpcClient)
}

// GetBondDenom returns bond denomination for the network.
func (c *Client) GetBondDenom() (string, error) {
	queryClient := c.GetStakingQueryClient()
	res, err := queryClient.Params(context.Background(), &stakingtypes.QueryParamsRequest{})
	if err != nil {
		return "", err
	}

	return res.Params.BondDenom, nil
}

// GetValidatorsByStatus returns validatorset by given status
func (c *Client) GetValidatorsByStatus(status stakingtypes.BondStatus) (validators []schema.Validator, err error) {

	var statusName string

	switch status {
	case stakingtypes.Bonded, stakingtypes.Unbonding, stakingtypes.Unbonded:
		statusName = stakingtypes.BondStatus_name[int32(status)]
	default:
		statusName = stakingtypes.BondStatus_name[int32(stakingtypes.Unspecified)]
	}

	queryClient := stakingtypes.NewQueryClient(c.grpcClient)
	request := stakingtypes.QueryValidatorsRequest{Status: statusName}
	resp, err := queryClient.Validators(context.Background(), &request)

	if len(resp.Validators) <= 0 {
		return []schema.Validator{}, nil
	}

	sort.Slice(resp.Validators[:], func(i, j int) bool {
		return resp.Validators[0].Tokens.GT(resp.Validators[1].Tokens)
		// tempTk1, _ := strconv.Atoi(bondedVals[i].Tokens)
		// tempTk2, _ := strconv.Atoi(bondedVals[j].Tokens)
		// return tempTk1 > tempTk2
	})

	for i, val := range resp.Validators {
		accAddr, err := types.ConvertAccAddrFromValAddr(val.OperatorAddress)
		if err != nil {
			return []schema.Validator{}, fmt.Errorf("failed to convert address from validator Address : %s", err)
		}
		log.Println("conspubkey get cached value : ", val.ConsensusPubkey.GetCachedValue())
		conspubkey, err := val.TmConsPubKey()
		if err != nil {
			return []schema.Validator{}, fmt.Errorf("failed to get consesnsus pubkey : %s", err)
		}
		consAddr, err := types.ConvertConsAddrFromConsPubkey(conspubkey.Address().String())
		if err != nil {
			return []schema.Validator{}, fmt.Errorf("failed to convert cons address from conspubkey : %s", err)
		}

		v := schema.NewValidator(schema.Validator{
			Rank:                 i + 1,
			OperatorAddress:      val.OperatorAddress,
			Address:              accAddr,
			ConsensusPubkey:      conspubkey.Address().String(),
			Proposer:             consAddr,
			Jailed:               val.Jailed,
			Status:               int(val.Status),
			Tokens:               val.Tokens.String(),
			DelegatorShares:      val.DelegatorShares.String(),
			Moniker:              val.Description.Moniker,
			Identity:             val.Description.Identity,
			Website:              val.Description.Website,
			Details:              val.Description.Details,
			UnbondingHeight:      val.UnbondingHeight,
			UnbondingTime:        val.UnbondingTime,
			CommissionRate:       val.Commission.CommissionRates.Rate.String(),
			CommissionMaxRate:    val.Commission.CommissionRates.MaxRate.String(),
			CommissionChangeRate: val.Commission.CommissionRates.MaxChangeRate.String(),
			MinSelfDelegation:    val.MinSelfDelegation.String(),
			UpdateTime:           val.Commission.UpdateTime,
		})

		validators = append(validators, *v)
	}
	// resp, err := c.apiClient.R().Get("/staking/validators?status=bonded")
	// if err != nil {
	// 	return []schema.Validator{}, fmt.Errorf("failed to request bonded vals: %s", err)
	// }

	// var bondedVals []types.Validator
	// err = json.Unmarshal(types.ReadRespWithHeight(resp).Result, &bondedVals)
	// if err != nil {
	// 	return []schema.Validator{}, fmt.Errorf("failed to unmarshal bonded vals: %s", err)
	// }

	// Sort bondedVals by highest token amount
	// sort.Slice(bondedVals[:], func(i, j int) bool {
	// 	tempTk1, _ := strconv.Atoi(bondedVals[i].Tokens)
	// 	tempTk2, _ := strconv.Atoi(bondedVals[j].Tokens)
	// 	return tempTk1 > tempTk2
	// })

	// if len(bondedVals) <= 0 {
	// 	return []schema.Validator{}, nil
	// }

	// for i, val := range bondedVals {
	// 	accAddr, _ := types.ConvertAccAddrFromValAddr(val.OperatorAddress)
	// 	consAddr, _ := types.ConvertConsAddrFromConsPubkey(val.ConsensusPubkey)

	// 	v := schema.NewValidator(schema.Validator{
	// 		Rank:                 i + 1,
	// 		OperatorAddress:      val.OperatorAddress,
	// 		Address:              accAddr,
	// 		ConsensusPubkey:      val.ConsensusPubkey,
	// 		Proposer:             consAddr,
	// 		Jailed:               val.Jailed,
	// 		Status:               val.Status,
	// 		Tokens:               val.Tokens,
	// 		DelegatorShares:      val.DelegatorShares,
	// 		Moniker:              val.Description.Moniker,
	// 		Identity:             val.Description.Identity,
	// 		Website:              val.Description.Website,
	// 		Details:              val.Description.Details,
	// 		UnbondingHeight:      val.UnbondingHeight,
	// 		UnbondingTime:        val.UnbondingTime,
	// 		CommissionRate:       val.Commission.CommissionRates.Rate,
	// 		CommissionMaxRate:    val.Commission.CommissionRates.MaxRate,
	// 		CommissionChangeRate: val.Commission.CommissionRates.MaxChangeRate,
	// 		MinSelfDelegation:    val.MinSelfDelegation,
	// 		UpdateTime:           val.Commission.UpdateTime,
	// 	})

	// 	validators = append(validators, *v)
	// }

	return validators, nil
}

// GetBondedValidators returns all bonded validators
// func (c *Client) GetBondedValidators() (validators []schema.Validator, err error) {
// 	resp, err := c.apiClient.R().Get("/staking/validators?status=bonded")
// 	if err != nil {
// 		return []schema.Validator{}, fmt.Errorf("failed to request bonded vals: %s", err)
// 	}

// 	var bondedVals []types.Validator
// 	err = json.Unmarshal(types.ReadRespWithHeight(resp).Result, &bondedVals)
// 	if err != nil {
// 		return []schema.Validator{}, fmt.Errorf("failed to unmarshal bonded vals: %s", err)
// 	}

// 	// Sort bondedVals by highest token amount
// 	sort.Slice(bondedVals[:], func(i, j int) bool {
// 		tempTk1, _ := strconv.Atoi(bondedVals[i].Tokens)
// 		tempTk2, _ := strconv.Atoi(bondedVals[j].Tokens)
// 		return tempTk1 > tempTk2
// 	})

// 	if len(bondedVals) <= 0 {
// 		return []schema.Validator{}, nil
// 	}

// 	for i, val := range bondedVals {
// 		accAddr, _ := types.ConvertAccAddrFromValAddr(val.OperatorAddress)
// 		consAddr, _ := types.ConvertConsAddrFromConsPubkey(val.ConsensusPubkey)

// 		v := schema.NewValidator(schema.Validator{
// 			Rank:                 i + 1,
// 			OperatorAddress:      val.OperatorAddress,
// 			Address:              accAddr,
// 			ConsensusPubkey:      val.ConsensusPubkey,
// 			Proposer:             consAddr,
// 			Jailed:               val.Jailed,
// 			Status:               val.Status,
// 			Tokens:               val.Tokens,
// 			DelegatorShares:      val.DelegatorShares,
// 			Moniker:              val.Description.Moniker,
// 			Identity:             val.Description.Identity,
// 			Website:              val.Description.Website,
// 			Details:              val.Description.Details,
// 			UnbondingHeight:      val.UnbondingHeight,
// 			UnbondingTime:        val.UnbondingTime,
// 			CommissionRate:       val.Commission.CommissionRates.Rate,
// 			CommissionMaxRate:    val.Commission.CommissionRates.MaxRate,
// 			CommissionChangeRate: val.Commission.CommissionRates.MaxChangeRate,
// 			MinSelfDelegation:    val.MinSelfDelegation,
// 			UpdateTime:           val.Commission.UpdateTime,
// 		})

// 		validators = append(validators, *v)
// 	}

// 	return validators, nil
// }

// GetUnbondingValidators returns unbonding validators
// func (c *Client) GetUnbondingValidators() (validators []schema.Validator, err error) {
// 	resp, err := c.apiClient.R().Get("/staking/validators?status=unbonding")
// 	if err != nil {
// 		return []schema.Validator{}, fmt.Errorf("failed to request unbonding vals: %s", err)
// 	}

// 	var unbondingVals []*types.Validator
// 	err = json.Unmarshal(types.ReadRespWithHeight(resp).Result, &unbondingVals)
// 	if err != nil {
// 		return []schema.Validator{}, fmt.Errorf("failed to unmarshal unbonding vals: %s", err)
// 	}

// 	// Sort bondedValidators by highest token amount
// 	sort.Slice(unbondingVals[:], func(i, j int) bool {
// 		tempTk1, _ := strconv.Atoi(unbondingVals[i].Tokens)
// 		tempTk2, _ := strconv.Atoi(unbondingVals[j].Tokens)
// 		return tempTk1 > tempTk2
// 	})

// 	if len(unbondingVals) <= 0 {
// 		return []schema.Validator{}, nil
// 	}

// 	for _, val := range unbondingVals {
// 		accAddr, _ := types.ConvertAccAddrFromValAddr(val.OperatorAddress)
// 		consAddr, _ := types.ConvertConsAddrFromConsPubkey(val.ConsensusPubkey)

// 		v := schema.NewValidator(schema.Validator{
// 			OperatorAddress:      val.OperatorAddress,
// 			Address:              accAddr,
// 			ConsensusPubkey:      val.ConsensusPubkey,
// 			Proposer:             consAddr,
// 			Jailed:               val.Jailed,
// 			Status:               val.Status,
// 			Tokens:               val.Tokens,
// 			DelegatorShares:      val.DelegatorShares,
// 			Moniker:              val.Description.Moniker,
// 			Identity:             val.Description.Identity,
// 			Website:              val.Description.Website,
// 			Details:              val.Description.Details,
// 			UnbondingHeight:      val.UnbondingHeight,
// 			UnbondingTime:        val.UnbondingTime,
// 			CommissionRate:       val.Commission.CommissionRates.Rate,
// 			CommissionMaxRate:    val.Commission.CommissionRates.MaxRate,
// 			CommissionChangeRate: val.Commission.CommissionRates.MaxChangeRate,
// 			MinSelfDelegation:    val.MinSelfDelegation,
// 			UpdateTime:           val.Commission.UpdateTime,
// 		})

// 		validators = append(validators, *v)
// 	}

// 	return validators, nil
// }

// GetUnbondedValidators returns unbonded validators
// func (c *Client) GetUnbondedValidators() (validators []schema.Validator, err error) {
// 	resp, err := c.apiClient.R().Get("/staking/validators?status=unbonded")
// 	if err != nil {
// 		return []schema.Validator{}, fmt.Errorf("failed to request unbonded vals: %s", err)
// 	}

// 	var unbondedVals []types.Validator
// 	err = json.Unmarshal(types.ReadRespWithHeight(resp).Result, &unbondedVals)
// 	if err != nil {
// 		return []schema.Validator{}, fmt.Errorf("failed to unmarshal unbonded vals: %s", err)
// 	}

// 	// Sort bondedValidators by highest token amount
// 	sort.Slice(unbondedVals[:], func(i, j int) bool {
// 		tempTk1, _ := strconv.Atoi(unbondedVals[i].Tokens)
// 		tempTk2, _ := strconv.Atoi(unbondedVals[j].Tokens)
// 		return tempTk1 > tempTk2
// 	})

// 	if len(unbondedVals) <= 0 {
// 		return []schema.Validator{}, nil
// 	}

// 	for _, val := range unbondedVals {
// 		accAddr, _ := types.ConvertAccAddrFromValAddr(val.OperatorAddress)
// 		consAddr, _ := types.ConvertConsAddrFromConsPubkey(val.ConsensusPubkey)

// 		v := schema.NewValidator(schema.Validator{
// 			OperatorAddress:      val.OperatorAddress,
// 			Address:              accAddr,
// 			ConsensusPubkey:      val.ConsensusPubkey,
// 			Proposer:             consAddr,
// 			Jailed:               val.Jailed,
// 			Status:               val.Status,
// 			Tokens:               val.Tokens,
// 			DelegatorShares:      val.DelegatorShares,
// 			Moniker:              val.Description.Moniker,
// 			Identity:             val.Description.Identity,
// 			Website:              val.Description.Website,
// 			Details:              val.Description.Details,
// 			UnbondingHeight:      val.UnbondingHeight,
// 			UnbondingTime:        val.UnbondingTime,
// 			CommissionRate:       val.Commission.CommissionRates.Rate,
// 			CommissionMaxRate:    val.Commission.CommissionRates.MaxRate,
// 			CommissionChangeRate: val.Commission.CommissionRates.MaxChangeRate,
// 			MinSelfDelegation:    val.MinSelfDelegation,
// 			UpdateTime:           val.Commission.UpdateTime,
// 		})

// 		validators = append(validators, *v)
// 	}

// 	return validators, nil
// }

// GetDelegatorDelegations returns a list of delegations made by a certain delegator address
func (c *Client) GetDelegatorDelegations(address string) (*stakingtypes.QueryDelegatorDelegationsResponse, error) {
	queryClient := c.GetStakingQueryClient()
	request := stakingtypes.QueryDelegatorDelegationsRequest{DelegatorAddr: address}
	res, err := queryClient.DelegatorDelegations(context.Background(), &request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetDelegatorUnbondingDelegations returns a list of undelegations made by a certain delegator address
func (c *Client) GetDelegatorUnbondingDelegations(address string) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error) {
	queryClient := c.GetStakingQueryClient()
	request := stakingtypes.QueryDelegatorUnbondingDelegationsRequest{DelegatorAddr: address}
	res, err := queryClient.DelegatorUnbondingDelegations(context.Background(), &request)
	if err != nil {
		return nil, err
	}

	return res, nil
}
