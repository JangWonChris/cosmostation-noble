package client

import (
	"context"

	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// GetDistributionQueryClient returns a object of queryClient
func (c *Client) GetDistributionQueryClient() distributiontypes.QueryClient {
	return distributiontypes.NewQueryClient(c.grpcClient)
}

// GetValidatorCommission queries validator's commission and returns the coins with truncated decimals and the change.
func (c *Client) GetValidatorCommission(address string) (distributiontypes.ValidatorAccumulatedCommission /* QueryValidatorCommissionResponse */, error) {
	queryClient := c.GetDistributionQueryClient()
	request := distributiontypes.QueryValidatorCommissionRequest{ValidatorAddress: address}
	res, err := queryClient.ValidatorCommission(context.Background(), &request)
	if err != nil {
		return distributiontypes.ValidatorAccumulatedCommission{}, err
	}

	return res.Commission, nil
}

// GetDelegationRewards returns the rewards from between given delegator and validator
func (c *Client) GetDelegationRewards(delegatorAddr, validatorAddr string) (*distributiontypes.QueryDelegationRewardsResponse, error) {
	queryClient := c.GetDistributionQueryClient()
	request := distributiontypes.QueryDelegationRewardsRequest{DelegatorAddress: delegatorAddr, ValidatorAddress: validatorAddr}
	resp, err := queryClient.DelegationRewards(context.Background(), &request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetDelegationTotalRewards returns the total rewards balance from all delegations by a delegator
func (c *Client) GetDelegationTotalRewards(address string) (*distributiontypes.QueryDelegationTotalRewardsResponse, error) {
	queryClient := c.GetDistributionQueryClient()
	request := distributiontypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: address}
	resp, err := queryClient.DelegationTotalRewards(context.Background(), &request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
