package client

import (
	"context"

	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// GetStakingQueryClient returns a object of queryClient
func (c *Client) GetStakingQueryClient() types.QueryClient {
	return types.NewQueryClient(c.grpcClient)
}

// GetBondDenom returns bond denomination for the network.
func (c *Client) GetBondDenom() (string, error) {
	queryClient := c.GetStakingQueryClient()
	res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
	if err != nil {
		return "", err
	}

	return res.Params.BondDenom, nil
}

// GetDelegatorDelegations returns a list of delegations made by a certain delegator address
func (c *Client) GetDelegatorDelegations(address string) (*types.QueryDelegatorDelegationsResponse, error) {
	queryClient := c.GetStakingQueryClient()
	request := types.QueryDelegatorDelegationsRequest{DelegatorAddr: address}
	res, err := queryClient.DelegatorDelegations(context.Background(), &request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetDelegatorUnbondingDelegations returns a list of undelegations made by a certain delegator address
func (c *Client) GetDelegatorUnbondingDelegations(address string) (*types.QueryDelegatorUnbondingDelegationsResponse, error) {
	queryClient := c.GetStakingQueryClient()
	request := types.QueryDelegatorUnbondingDelegationsRequest{DelegatorAddr: address}
	res, err := queryClient.DelegatorUnbondingDelegations(context.Background(), &request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetValidator queries validator's information of given validator address
func (c *Client) GetValidator(address string) (*types.QueryValidatorResponse, error) {
	queryClient := c.GetStakingQueryClient()
	request := types.QueryValidatorRequest{ValidatorAddr: address}
	res, err := queryClient.Validator(context.Background(), &request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetRedelegations queries validator's information of given validator address
func (c *Client) GetRedelegations(delAddr, srcValidatorAddress, dstValidatorAddress string) (*types.QueryRedelegationsResponse, error) {
	queryClient := c.GetStakingQueryClient()
	request := types.QueryRedelegationsRequest{DelegatorAddr: delAddr, SrcValidatorAddr: srcValidatorAddress, DstValidatorAddr: dstValidatorAddress}
	res, err := queryClient.Redelegations(context.Background(), &request)
	if err != nil {
		return nil, err
	}

	return res, nil
}
