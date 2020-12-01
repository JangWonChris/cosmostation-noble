package client

import (
	"fmt"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/client/common"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// GetDistributionQueryClient returns a object of queryClient
func (c *Client) GetDistributionQueryClient() distributiontypes.QueryClient {
	return distributiontypes.NewQueryClient(c.grpcClient)
}

// GetValidatorCommission queries validator's commission and returns the coins with truncated decimals and the change.
func (c *Client) GetValidatorCommission(address string) (sdktypes.Coins, error) {
	valAddr, err := sdktypes.ValAddressFromBech32(address)
	if err != nil {
		return sdktypes.Coins{}, err
	}

	res, err := common.QueryValidatorCommission(c.cliCtx, valAddr)
	if err != nil {
		return sdktypes.Coins{}, err
	}

	var valCom distributiontypes.ValidatorAccumulatedCommission
	c.cliCtx.LegacyAmino.MustUnmarshalJSON(res, &valCom)

	truncatedCoins, _ := valCom.Commission.TruncateDecimal()

	return truncatedCoins, nil
}

// GetDelegatorTotalRewards returns the total rewards balance from all delegations by a delegator
func (c *Client) GetDelegatorTotalRewards(address string) (distributiontypes.QueryDelegatorTotalRewardsResponse, error) {
	delAddr, err := sdktypes.AccAddressFromBech32(address)
	if err != nil {
		return distributiontypes.QueryDelegatorTotalRewardsResponse{}, err
	}

	bz, err := c.cliCtx.LegacyAmino.MarshalJSON(distributiontypes.NewQueryDelegatorParams(delAddr))
	if err != nil {
		return distributiontypes.QueryDelegatorTotalRewardsResponse{}, err
	}

	route := fmt.Sprintf("custom/%s/%s", distributiontypes.QuerierRoute, distributiontypes.QueryDelegatorTotalRewards)
	res, _, err := c.cliCtx.QueryWithData(route, bz)
	if err != nil {
		return distributiontypes.QueryDelegatorTotalRewardsResponse{}, err
	}

	var totalRewards distributiontypes.QueryDelegatorTotalRewardsResponse
	if err := c.cliCtx.LegacyAmino.UnmarshalJSON(res, &totalRewards); err != nil {
		return distributiontypes.QueryDelegatorTotalRewardsResponse{}, err
	}

	return totalRewards, nil
}
