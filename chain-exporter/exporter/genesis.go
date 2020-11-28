package exporter

import (
	"fmt"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
	"go.uber.org/zap"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

const (
	// startingHeight is used to extract genesis accounts and parse their assets.
	startingHeight = int64(1)
)

func (ex *Exporter) getGenesisAccounts(genesisAccts authtypes.GenesisAccounts) (accounts []schema.Account, err error) {
	chainID, err := ex.client.GetNetworkChainID()
	if err != nil {
		return []schema.Account{}, err
	}

	block, err := ex.client.GetBlock(startingHeight)
	if err != nil {
		return []schema.Account{}, err
	}

	denom, err := ex.client.GetBondDenom()
	if err != nil {
		return []schema.Account{}, err
	}

	for i, account := range genesisAccts {
		switch account.(type) {
		case *authtypes.BaseAccount:
			zap.S().Infof("Account type: %s | Synced account %d/%d", types.BaseAccount, i, len(genesisAccts))

			acc := account.(*authtypes.BaseAccount)

			spendable, rewards, commission, delegated, undelegated, err := ex.client.GetBaseAccountTotalAsset(acc.GetAddress().String())
			if err != nil {
				return []schema.Account{}, err
			}

			total := sdk.NewCoin(denom, sdk.NewInt(0))

			// Sum up all coins that exist in an account.
			total = total.Add(spendable).
				Add(delegated).
				Add(undelegated).
				Add(rewards).
				Add(commission)

			acct := &schema.Account{
				ChainID:          chainID,
				AccountAddress:   acc.Address,
				AccountNumber:    acc.AccountNumber,
				AccountType:      types.BaseAccount,
				CoinsTotal:       total.Amount.String(),
				CoinsSpendable:   spendable.Amount.String(),
				CoinsRewards:     rewards.Amount.String(),
				CoinsCommission:  commission.Amount.String(),
				CoinsDelegated:   delegated.Amount.String(),
				CoinsUndelegated: undelegated.Amount.String(),
				// CoinsTotal:       *total.Amount.BigInt(),
				// CoinsSpendable:   *spendable.Amount.BigInt(),
				// CoinsRewards:     *rewards.Amount.BigInt(),
				// CoinsCommission:  *commission.Amount.BigInt(),
				// CoinsDelegated:   *delegated.Amount.BigInt(),
				// CoinsUndelegated: *undelegated.Amount.BigInt(),
				CreationTime: block.Block.Time.String(),
			}

			accounts = append(accounts, *acct)

		case *authtypes.ModuleAccount:
			zap.S().Infof("Account type: %s | Synced account %d/%d", types.BaseAccount, i, len(genesisAccts))

			acc := account.(authtypes.ModuleAccountI)

			spendable, rewards, commission, delegated, undelegated, err := ex.client.GetBaseAccountTotalAsset(acc.GetAddress().String())
			if err != nil {
				return []schema.Account{}, err
			}

			total := sdk.NewCoin(denom, sdk.NewInt(0))

			// Sum up all coins that exist in an account.
			total = total.Add(spendable).
				Add(delegated).
				Add(undelegated).
				Add(rewards).
				Add(commission)

			acct := &schema.Account{
				ChainID:          chainID,
				AccountAddress:   acc.GetAddress().String(),
				AccountNumber:    acc.GetAccountNumber(),
				AccountType:      types.ModuleAccount,
				CoinsTotal:       total.Amount.String(),
				CoinsSpendable:   spendable.Amount.String(),
				CoinsRewards:     rewards.Amount.String(),
				CoinsCommission:  commission.Amount.String(),
				CoinsDelegated:   delegated.Amount.String(),
				CoinsUndelegated: undelegated.Amount.String(),
				// CoinsTotal:       *total.Amount.BigInt(),
				// CoinsSpendable:   *spendable.Amount.BigInt(),
				// CoinsRewards:     *rewards.Amount.BigInt(),
				// CoinsCommission:  *commission.Amount.BigInt(),
				// CoinsDelegated:   *delegated.Amount.BigInt(),
				// CoinsUndelegated: *undelegated.Amount.BigInt(),
				CreationTime: block.Block.Time.String(),
			}

			accounts = append(accounts, *acct)

		case *authvestingtypes.PeriodicVestingAccount:
			zap.S().Infof("Account type: %s | Synced account %d/%d", types.BaseAccount, i, len(genesisAccts))

			acc := account.(*authvestingtypes.PeriodicVestingAccount)

			spendable, rewards, commission, delegated, undelegated, err := ex.client.GetBaseAccountTotalAsset(acc.GetAddress().String())
			if err != nil {
				return []schema.Account{}, err
			}

			vesting := sdk.NewCoin(denom, sdk.NewInt(0))
			vested := sdk.NewCoin(denom, sdk.NewInt(0))

			vestingCoins := acc.GetVestingCoins(block.Block.Time)
			vestedCoins := acc.GetVestedCoins(block.Block.Time)
			delegatedVesting := acc.GetDelegatedVesting()

			// When total vesting amount is greater than or equal to delegated vesting amount, then
			// there is still a room to delegate. Otherwise, vesting should be zero.
			if len(vestingCoins) > 0 {
				if vestingCoins.IsAllGTE(delegatedVesting) {
					vestingCoins = vestingCoins.Sub(delegatedVesting)
					for _, vc := range vestingCoins {
						if vc.Denom == denom {
							vesting = vesting.Add(vc)
						}
					}
				}
			}

			if len(vestedCoins) > 0 {
				for _, vc := range vestedCoins {
					if vc.Denom == denom {
						vested = vested.Add(vc)
					}
				}
			}

			total := sdk.NewCoin(denom, sdk.NewInt(0))

			// Sum up all coins that exist in an account.
			total = total.Add(spendable).
				Add(delegated).
				Add(undelegated).
				Add(rewards).
				Add(commission).
				Add(vesting)

			acct := &schema.Account{
				ChainID:          chainID,
				AccountAddress:   acc.Address,
				AccountNumber:    acc.AccountNumber,
				AccountType:      types.PeriodicVestingAccount,
				CoinsTotal:       total.Amount.String(),
				CoinsSpendable:   spendable.Amount.String(),
				CoinsRewards:     rewards.Amount.String(),
				CoinsCommission:  commission.Amount.String(),
				CoinsDelegated:   delegated.Amount.String(),
				CoinsUndelegated: undelegated.Amount.String(),
				// CoinsTotal:       *total.Amount.BigInt(),
				// CoinsSpendable:   *spendable.Amount.BigInt(),
				// CoinsRewards:     *rewards.Amount.BigInt(),
				// CoinsCommission:  *commission.Amount.BigInt(),
				// CoinsDelegated:   *delegated.Amount.BigInt(),
				// CoinsUndelegated: *undelegated.Amount.BigInt(),
				// CoinsVesting:     *vesting.Amount.BigInt(),
				// CoinsVested:      *vested.Amount.BigInt(),
				CreationTime: block.Block.Time.String(),
			}

			accounts = append(accounts, *acct)

		default:
			return []schema.Account{}, fmt.Errorf("unrecognized account type: %T", account)
		}
	}

	return accounts, nil
}

// getGenesisValidatorsSet returns validator set in genesis.
func (ex *Exporter) getGenesisValidatorsSet(block *tmctypes.ResultBlock, vals *tmctypes.ResultValidators) ([]schema.PowerEventHistory, error) {
	genesisValsSet := make([]schema.PowerEventHistory, 0)

	if block.Block.Height != 1 {
		return []schema.PowerEventHistory{}, nil
	}

	denom, err := ex.client.GetBondDenom()
	if err != nil {
		return []schema.PowerEventHistory{}, err
	}

	// Get genesis validator set (block height 1).
	for i, val := range vals.Validators {
		gvs := schema.NewPowerEventHistoryForGenesisValidatorSet(schema.PowerEventHistory{
			IDValidator:          i + 1,
			Height:               block.Block.Height,
			Moniker:              "",
			OperatorAddress:      "",
			Proposer:             val.Address.String(),
			VotingPower:          float64(val.VotingPower),
			MsgType:              types.TypeMsgCreateValidator,
			NewVotingPowerAmount: float64(val.VotingPower),
			NewVotingPowerDenom:  denom,
			TxHash:               "",
			Timestamp:            block.Block.Header.Time,
		})

		genesisValsSet = append(genesisValsSet, *gvs)
	}

	return genesisValsSet, nil
}
