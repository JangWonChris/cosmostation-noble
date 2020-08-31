package exporter

import (
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	// kavavesting "github.com/kava-labs/kava/x/validator-vesting"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

/*
func (ex *Exporter) getGenesisAccounts(genesisAccts exported.GenesisAccounts) (accounts []schema.Account, err error) {
	chainID, err := ex.client.GetNetworkChainID()
	if err != nil {
		return []schema.Account{}, err
	}

	denom, err := ex.client.GetBondDenom()
	if err != nil {
		return []schema.Account{}, err
	}

	height := int64(1)
	genesisBlock, err := ex.client.GetBlock(height)
	if err != nil {
		return []schema.Account{}, err
	}

	for _, account := range genesisAccts {
		switch account.(type) {
		case *auth.BaseAccount:
			acc := account.(*auth.BaseAccount)

			spendable, rewards, commission, delegated, undelegated, err := ex.client.GetBaseAccountTotalAsset(acc.GetAddress().String())
			if err != nil {
				return []schema.Account{}, err
			}

			// Sum up all coins that exist in an account.
			total := sdkTypes.NewCoins(sdkTypes.NewCoin(denom, spendable.AmountOf(denom).
				Add(rewards.AmountOf(denom)).
				Add(delegated.Amount).
				Add(undelegated.Amount).
				Add(commission.AmountOf(denom))))

			acct := &schema.Account{
				ChainID:          chainID,
				AccountAddress:   acc.Address.String(),
				AccountNumber:    acc.AccountNumber,
				AccountType:      types.BaseAccount,
				CoinsTotal:       total.String(),
				CoinsSpendable:   spendable.String(),
				CoinsRewards:     rewards.String(),
				CoinsCommission:  commission.String(),
				CoinsDelegated:   delegated.String(),
				CoinsUndelegated: undelegated.String(),
				CreationTime:     genesisBlock.Block.Time.String(),
			}

			accounts = append(accounts, *acct)

		case *supply.ModuleAccount:
			acc := account.(supplyExported.ModuleAccountI)

			spendable, rewards, commission, delegated, undelegated, err := ex.client.GetBaseAccountTotalAsset(acc.GetAddress().String())
			if err != nil {
				return []schema.Account{}, err
			}

			// Sum up all coins that exist in an account.
			total := sdkTypes.NewCoins(sdkTypes.NewCoin(denom, spendable.AmountOf(denom).
				Add(rewards.AmountOf(denom)).
				Add(delegated.Amount).
				Add(undelegated.Amount).
				Add(commission.AmountOf(denom))))

			acct := &schema.Account{
				ChainID:          chainID,
				AccountAddress:   acc.GetAddress().String(),
				AccountNumber:    acc.GetAccountNumber(),
				AccountType:      types.ModuleAccount,
				CoinsTotal:       total.String(),
				CoinsSpendable:   spendable.String(),
				CoinsRewards:     rewards.String(),
				CoinsCommission:  commission.String(),
				CoinsDelegated:   delegated.String(),
				CoinsUndelegated: undelegated.String(),
				CreationTime:     genesisBlock.Block.Time.String(),
			}

			accounts = append(accounts, *acct)

		case *cosmosvesting.PeriodicVestingAccount:
			acc := account.(*cosmosvesting.PeriodicVestingAccount)

			spendable, rewards, commission, delegated, undelegated, err := ex.client.GetBaseAccountTotalAsset(acc.GetAddress().String())
			if err != nil {
				return []schema.Account{}, err
			}

			vesting := acc.GetVestingCoins(genesisBlock.Block.Time)
			vested := acc.GetVestedCoins(genesisBlock.Block.Time)

			// Avoid vesting amount to be negative value
			// when delegated vesting amount is greater than current vesting amounts.
			if acc.GetDelegatedVesting().IsAllGT(vesting) {
				vesting = sdkTypes.NewCoins(sdkTypes.NewCoin(denom, sdkTypes.NewInt(0)))
			} else {
				vesting = vesting.Sub(acc.GetDelegatedVesting())
			}

			// Sum up all coins that exist in an account.
			total := sdkTypes.NewCoins(sdkTypes.NewCoin(denom, spendable.AmountOf(denom).
				Add(rewards.AmountOf(denom)).
				Add(delegated.Amount).
				Add(undelegated.Amount).
				Add(commission.AmountOf(denom)).
				Add(vesting.AmountOf(denom))))

			acct := &schema.Account{
				ChainID:          chainID,
				AccountAddress:   acc.Address.String(),
				AccountNumber:    acc.AccountNumber,
				AccountType:      types.PeriodicVestingAccount,
				CoinsTotal:       total.String(),
				CoinsSpendable:   spendable.String(), // only numbers: spendable.AmountOf(denom).String()
				CoinsRewards:     rewards.String(),
				CoinsCommission:  commission.String(),
				CoinsDelegated:   delegated.String(),
				CoinsUndelegated: undelegated.String(),
				CoinsVesting:     vesting.String(),
				CoinsVested:      vested.String(),
				CreationTime:     genesisBlock.Block.Time.String(),
			}

			accounts = append(accounts, *acct)
		// case *kavavesting.ValidatorVestingAccount:
		// 	acc := account.(*kavavesting.ValidatorVestingAccount)

		// 	spendable, rewards, commission, delegated, undelegated, err := ex.client.GetBaseAccountTotalAsset(acc.GetAddress().String())
		// 	if err != nil {
		// 		return []schema.Account{}, err
		// 	}

		// 	vesting := acc.GetVestingCoins(genesisBlock.Block.Time)
		// 	vested := acc.GetVestedCoins(genesisBlock.Block.Time)
		// 	failedVested := acc.GetFailedVestedCoins()

		// 	// Avoid vesting amount to be negative value
		// 	// when delegated vesting amount is greater than current vesting amounts.
		// 	if acc.GetDelegatedVesting().IsAllGT(vesting) {
		// 		vesting = sdkTypes.NewCoins(sdkTypes.NewCoin(denom, sdkTypes.NewInt(0)))
		// 	} else {
		// 		vesting = vesting.Sub(acc.GetDelegatedVesting())
		// 	}

		// 	// Sum up all coins that exist in an account.
		// 	total := sdkTypes.NewCoins(sdkTypes.NewCoin(denom, spendable.AmountOf(denom).
		// 		Add(rewards.AmountOf(denom)).
		// 		Add(delegated.Amount).
		// 		Add(undelegated.Amount).
		// 		Add(commission.AmountOf(denom)).
		// 		Add(vesting.AmountOf(denom))))

		// 	acct := &schema.Account{
		// 		ChainID:           chainID,
		// 		AccountAddress:    acc.Address.String(),
		// 		AccountNumber:     acc.AccountNumber,
		// 		AccountType:       types.ValidatorVestingAccount,
		// 		CoinsTotal:        total.String(),
		// 		CoinsSpendable:    spendable.String(),
		// 		CoinsRewards:      rewards.String(),
		// 		CoinsCommission:   commission.String(),
		// 		CoinsDelegated:    delegated.String(),
		// 		CoinsUndelegated:  undelegated.String(),
		// 		CoinsVesting:      vesting.String(),
		// 		CoinsVested:       vested.String(),
		// 		CoinsFailedVested: failedVested.String(),
		// 		CreationTime:      genesisBlock.Block.Time.String(),
		// 	}

		// 	accounts = append(accounts, *acct)
		default:
			return []schema.Account{}, fmt.Errorf("unrecognized account type: %T", account)
		}
	}

	return accounts, nil
}
*/

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
