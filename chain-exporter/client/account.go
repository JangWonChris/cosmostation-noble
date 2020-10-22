package client

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
)

// GetBaseAccountTotalAsset returns total available, rewards, commission, delegations, and undelegations from a delegator.
func (c *Client) GetBaseAccountTotalAsset(address string) (sdktypes.Coin, sdktypes.Coin, sdktypes.Coin, sdktypes.Coin, sdktypes.Coin, error) {
	account, err := c.GetAccount(address)
	if err != nil {
		return sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, err
	}

	denom, err := c.GetBondDenom()
	if err != nil {
		return sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, err
	}

	spendable := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
	delegated := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
	undelegated := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
	rewards := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
	commission := sdktypes.NewCoin(denom, sdktypes.NewInt(0))

	// bank.

	// banktypes.GetGenesisStateFromAppState(c.cliCtx.JSONMarshaler, appState map[string]json.RawMessage)
	// Get total spendable coins.
	if len(account.GetCoins()) > 0 {
		for _, coin := range account.GetCoins() {
			if coin.Denom == denom {
				spendable = spendable.Add(coin)
			}
		}
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
