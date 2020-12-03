package client

import (
	"context"
	"fmt"
	"sort"

	// cosmos-sdk
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	// internal
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
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

		var conspubkey cryptotypes.PubKey
		codec.AppCodec.UnpackAny(val.ConsensusPubkey, &conspubkey)

		valconspub, err := sdktypes.Bech32ifyPubKey(sdktypes.Bech32PubKeyTypeConsPub, conspubkey)
		if err != nil {
			return []schema.Validator{}, fmt.Errorf("failed to get consesnsus pubkey : %s", err)
		}

		// log.Println("conspubkey get cached value : ", val.ConsensusPubkey.GetCachedValue())
		// conspubkey, err := val.TmConsPubKey()
		// if err != nil {
		// 	return []schema.Validator{}, fmt.Errorf("failed to get consesnsus pubkey : %s", err)
		// }

		v := schema.NewValidator(schema.Validator{
			Rank:                 i + 1,
			OperatorAddress:      val.OperatorAddress,
			Address:              accAddr,
			ConsensusPubkey:      valconspub,
			Proposer:             conspubkey.Address().String(),
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

	return validators, nil
}

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
