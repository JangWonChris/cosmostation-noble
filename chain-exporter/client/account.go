package client

import (
	"context"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// GetBaseAccountTotalAsset returns total available, rewards, commission, delegations, and undelegations from a delegator.
func (c *Client) GetBaseAccountTotalAsset(address string) (sdktypes.Coin, sdktypes.Coin, sdktypes.Coin, sdktypes.Coin, sdktypes.Coin, error) {
	denom, err := c.GetBondDenom()
	if err != nil {
		return sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, err
	}

	spendable := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
	delegated := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
	undelegated := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
	rewards := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
	commission := sdktypes.NewCoin(denom, sdktypes.NewInt(0))

	sdkaddr, err := sdktypes.AccAddressFromBech32(address)
	if err != nil {
		return sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, err
	}

	b := banktypes.NewQueryBalanceRequest(sdkaddr, denom)
	bankClient := banktypes.NewQueryClient(c.grpcClient)
	var header metadata.MD
	bankRes, err := bankClient.Balance(
		context.Background(),
		b,
		grpc.Header(&header), // Also fetch grpc header
	)
	if bankRes.GetBalance() != nil {
		spendable = spendable.Add(*bankRes.GetBalance())
	}

	// Get total delegated coins.
	delegations, err := c.GetDelegatorDelegations(address)
	if err != nil {
		return sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, err
	}

	if len(delegations) > 0 {
		for _, delegation := range delegations {
			delegated = delegated.Add(delegation.Balance)
		}
	}

	// Get total undelegated coins.
	undelegations, err := c.GetDelegatorUndelegations(address)
	if err != nil {
		return sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, err
	}

	if len(undelegations) > 0 {
		for _, undelegation := range undelegations {
			for _, e := range undelegation.Entries {
				undelegated = undelegated.Add(sdktypes.NewCoin(denom, e.Balance))
			}
		}
	}

	// Get total rewarded coins.
	totalRewards, err := c.GetDelegatorTotalRewards(address)
	if err != nil {
		return sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, err
	}

	if len(totalRewards.Rewards) > 0 {
		for _, tr := range totalRewards.Rewards {
			for _, reward := range tr.Reward {
				if reward.Denom == denom {
					truncatedRewards, _ := reward.TruncateDecimal()
					rewards = rewards.Add(truncatedRewards)
				}
			}
		}
	}

	valAddr, err := types.ConvertValAddrFromAccAddr(address)
	if err != nil {
		return sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, err
	}

	// Get commission
	commissions, err := c.GetValidatorCommission(valAddr)
	if err != nil {
		return sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, err
	}

	if len(commissions) > 0 {
		for _, c := range commissions {
			commission = commission.Add(c)
		}
	}

	return spendable, delegated, undelegated, rewards, commission, nil
}
