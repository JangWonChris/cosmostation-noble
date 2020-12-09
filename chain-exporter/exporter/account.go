package exporter

import (
	"fmt"

	// gaia

	// cosmos-sdk
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	// internal

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	// tendermint

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"go.uber.org/zap"
)

// getAccounts
func (ex *Exporter) getAccounts(block *tmctypes.ResultBlock, txResps []*sdk.TxResponse) (accounts []schema.Account, err error) {
	if len(txResps) <= 0 {
		return []schema.Account{}, nil
	}

	for _, txResp := range txResps {
		// Other than code equals to 0, it is failed transaction.
		if txResp.Code != 0 {
			return []schema.Account{}, nil
		}

		// stdTx, ok := tx.Tx.(auth.StdTx)
		// if !ok {
		// 	return []schema.Account{}, fmt.Errorf("unsupported tx type: %s", tx.Tx)
		// }

		msgs := txResp.GetTx().GetMsgs()

		for _, msg := range msgs {

			switch m := msg.(type) {
			case *banktypes.MsgSend:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), txResp.TxHash)

				// msgSend := m.(bank.MsgSend)

				fromAcct, err := ex.client.GetAccount(m.FromAddress)
				if err != nil {
					return []schema.Account{}, err
				}

				toAcct, err := ex.client.GetAccount(m.ToAddress)
				if err != nil {
					return []schema.Account{}, err
				}

				exportedAccts := []sdkclient.Account{
					fromAcct, toAcct,
				}

				accounts, err = ex.getAccountAllAssets(exportedAccts, txResp.TxHash, txResp.Timestamp)
				if err != nil {
					return []schema.Account{}, err
				}

			case *banktypes.MsgMultiSend:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), txResp.TxHash)

				// msgMultiSend := m.(bank.MsgMultiSend)

				var exportedAccts []sdkclient.Account

				for _, input := range m.Inputs {
					inputAcct, err := ex.client.GetAccount(input.Address)
					if err != nil {
						return []schema.Account{}, err
					}

					exportedAccts = append(exportedAccts, inputAcct)
				}

				for _, output := range m.Outputs {
					outputAcct, err := ex.client.GetAccount(output.Address)
					if err != nil {
						return []schema.Account{}, err
					}

					exportedAccts = append(exportedAccts, outputAcct)
				}

				accounts, err = ex.getAccountAllAssets(exportedAccts, txResp.TxHash, txResp.Timestamp)
				if err != nil {
					return []schema.Account{}, err
				}

			case *stakingtypes.MsgDelegate:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), txResp.TxHash)

				// msgDelegate := m.(staking.MsgDelegate)

				delegatorAddr, err := ex.client.GetAccount(m.DelegatorAddress)
				if err != nil {
					return []schema.Account{}, err
				}

				valAccAddr, err := types.ConvertAccAddrFromValAddr(m.DelegatorAddress)
				if err != nil {
					return []schema.Account{}, err
				}

				valAddr, err := ex.client.GetAccount(valAccAddr)
				if err != nil {
					return []schema.Account{}, err
				}

				exportedAccts := []sdkclient.Account{
					delegatorAddr, valAddr,
				}

				accounts, err = ex.getAccountAllAssets(exportedAccts, txResp.TxHash, txResp.Timestamp)
				if err != nil {
					return []schema.Account{}, err
				}

			case *stakingtypes.MsgUndelegate:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), txResp.TxHash)

				// msgUndelegate := m.(staking.MsgUndelegate)

				delegatorAddr, err := ex.client.GetAccount(m.DelegatorAddress)
				if err != nil {
					return []schema.Account{}, err
				}

				valAccAddr, err := types.ConvertAccAddrFromValAddr(m.DelegatorAddress)
				if err != nil {
					return []schema.Account{}, err
				}

				valAddr, err := ex.client.GetAccount(valAccAddr)
				if err != nil {
					return []schema.Account{}, err
				}

				exportedAccts := []sdkclient.Account{
					delegatorAddr, valAddr,
				}

				accounts, err = ex.getAccountAllAssets(exportedAccts, txResp.TxHash, txResp.Timestamp)
				if err != nil {
					return []schema.Account{}, err
				}

			case *stakingtypes.MsgBeginRedelegate:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), txResp.TxHash)

				// msgBeginRedelegate := m.(staking.MsgBeginRedelegate)

				delegatorAddr, err := ex.client.GetAccount(m.DelegatorAddress)
				if err != nil {
					return []schema.Account{}, err
				}

				valSrcAccAddr, err := types.ConvertAccAddrFromValAddr(m.ValidatorSrcAddress)
				if err != nil {
					return []schema.Account{}, err
				}

				valDstAccAddr, err := types.ConvertAccAddrFromValAddr(m.ValidatorDstAddress)
				if err != nil {
					return []schema.Account{}, err
				}

				srcAddr, err := ex.client.GetAccount(valSrcAccAddr)
				if err != nil {
					return []schema.Account{}, err
				}

				dstAddr, err := ex.client.GetAccount(valDstAccAddr)
				if err != nil {
					return []schema.Account{}, err
				}

				exportedAccts := []sdkclient.Account{
					delegatorAddr, srcAddr, dstAddr,
				}

				accounts, err = ex.getAccountAllAssets(exportedAccts, txResp.TxHash, txResp.Timestamp)
				if err != nil {
					return []schema.Account{}, err
				}

			default:
				continue
			}
		}
	}

	return accounts, nil
}

func (ex *Exporter) getAccountAllAssets(exportedAccts []sdkclient.Account, txHashStr, txTime string) (accounts []schema.Account, err error) {
	chainID, err := ex.client.GetNetworkChainID()
	if err != nil {
		return []schema.Account{}, err
	}

	denom, err := ex.client.GetBondDenom()
	if err != nil {
		return []schema.Account{}, err
	}

	latestBlockHeight, err := ex.client.GetLatestBlockHeight()
	if err != nil {
		return []schema.Account{}, err
	}

	block, err := ex.client.GetBlock(latestBlockHeight)
	if err != nil {
		return []schema.Account{}, err
	}

	for _, account := range exportedAccts {
		switch acc := account.(type) {
		case *authtypes.BaseAccount:
			zap.S().Infof("Account type: %s | Account: %s", types.BaseAccount, account.GetAddress())

			// acc := account.(*authtypes.BaseAccount)

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
				ChainID:           chainID,
				AccountAddress:    acc.Address,
				AccountNumber:     acc.AccountNumber,
				AccountType:       types.BaseAccount,
				CoinsTotal:        total.Amount.String(),
				CoinsSpendable:    spendable.Amount.String(),
				CoinsRewards:      rewards.Amount.String(),
				CoinsCommission:   commission.Amount.String(),
				CoinsDelegated:    delegated.Amount.String(),
				CoinsUndelegated:  undelegated.Amount.String(),
				CoinsVested:       "0",
				CoinsVesting:      "0",
				CoinsFailedVested: "0",
				LastTx:            txHashStr,
				LastTxTime:        txTime,
				CreationTime:      block.Block.Time.String(),
			}

			accounts = append(accounts, *acct)

		case *authtypes.ModuleAccount:
			zap.S().Infof("Account type: %s | Account: %s", types.ModuleAccount, account.GetAddress().String())

			// acc := account.(authtypes.ModuleAccountI)

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
				ChainID:           chainID,
				AccountAddress:    acc.GetAddress().String(),
				AccountNumber:     acc.GetAccountNumber(),
				AccountType:       types.ModuleAccount,
				CoinsTotal:        total.Amount.String(),
				CoinsSpendable:    spendable.Amount.String(),
				CoinsRewards:      rewards.Amount.String(),
				CoinsCommission:   commission.Amount.String(),
				CoinsDelegated:    delegated.Amount.String(),
				CoinsUndelegated:  undelegated.Amount.String(),
				CoinsVested:       "0",
				CoinsVesting:      "0",
				CoinsFailedVested: "0",
				LastTx:            txHashStr,
				LastTxTime:        txTime,
				CreationTime:      block.Block.Time.String(),
			}

			accounts = append(accounts, *acct)

		case *authvestingtypes.PeriodicVestingAccount:
			zap.S().Infof("Account type: %s | Account: %s", types.PeriodicVestingAccount, account.GetAddress().String())

			// acc := account.(*authvestingtypes.PeriodicVestingAccount)

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
				ChainID:           chainID,
				AccountAddress:    acc.Address,
				AccountNumber:     acc.AccountNumber,
				AccountType:       types.PeriodicVestingAccount,
				CoinsTotal:        total.Amount.String(),
				CoinsSpendable:    spendable.Amount.String(),
				CoinsRewards:      rewards.Amount.String(),
				CoinsCommission:   commission.Amount.String(),
				CoinsDelegated:    delegated.Amount.String(),
				CoinsUndelegated:  undelegated.Amount.String(),
				CoinsVested:       "0",
				CoinsVesting:      "0",
				CoinsFailedVested: "0",
				LastTx:            txHashStr,
				LastTxTime:        txTime,
				CreationTime:      block.Block.Time.String(),
			}

			accounts = append(accounts, *acct)

		case *authvestingtypes.DelayedVestingAccount:
			zap.S().Infof("Account type: %s | Account: %s", types.DelayedVestingAccount, account.GetAddress().String())

			// acc := account.(*authvestingtypes.DelayedVestingAccount)

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
				ChainID:           chainID,
				AccountAddress:    acc.Address,
				AccountNumber:     acc.AccountNumber,
				AccountType:       types.DelayedVestingAccount,
				CoinsTotal:        total.Amount.String(),
				CoinsSpendable:    spendable.Amount.String(),
				CoinsRewards:      rewards.Amount.String(),
				CoinsCommission:   commission.Amount.String(),
				CoinsDelegated:    delegated.Amount.String(),
				CoinsUndelegated:  undelegated.Amount.String(),
				CoinsVested:       "0",
				CoinsVesting:      "0",
				CoinsFailedVested: "0",
				LastTx:            txHashStr,
				LastTxTime:        txTime,
				CreationTime:      block.Block.Time.String(),
			}

			accounts = append(accounts, *acct)

		default:
			return []schema.Account{}, fmt.Errorf("unrecognized account type: %T", account)
		}
	}

	return accounts, nil
}
