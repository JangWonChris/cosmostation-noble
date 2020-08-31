package exporter

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"

	// cosmosvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	supplyExported "github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	// kavavesting "github.com/kava-labs/kava/x/validator-vesting"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"go.uber.org/zap"
)

// getAccounts
func (ex *Exporter) getAccounts(block *tmctypes.ResultBlock, txs []*sdk.TxResponse) (accounts []schema.Account, err error) {
	if len(txs) <= 0 {
		return []schema.Account{}, nil
	}

	for _, tx := range txs {
		// Other than code equals to 0, it is failed transaction.
		if tx.Code != 0 {
			return []schema.Account{}, nil
		}

		stdTx, ok := tx.Tx.(auth.StdTx)
		if !ok {
			return []schema.Account{}, fmt.Errorf("unsupported tx type: %s", tx.Tx)
		}

		switch stdTx.Msgs[0].(type) {
		case bank.MsgSend:
			zap.S().Infof("MsgType: %s | Hash: %s", stdTx.Msgs[0].Type(), tx.TxHash)

			msgSend := stdTx.Msgs[0].(bank.MsgSend)

			fromAcct, err := ex.client.GetAccount(msgSend.FromAddress.String())
			if err != nil {
				return []schema.Account{}, err
			}

			toAcct, err := ex.client.GetAccount(msgSend.ToAddress.String())
			if err != nil {
				return []schema.Account{}, err
			}

			exportedAccts := []exported.Account{
				fromAcct, toAcct,
			}

			accounts, err = ex.getAccountAllAssets(exportedAccts, tx.TxHash, tx.Timestamp)
			if err != nil {
				return []schema.Account{}, err
			}

		case bank.MsgMultiSend:
			zap.S().Infof("MsgType: %s | Hash: %s", stdTx.Msgs[0].Type(), tx.TxHash)

			msgMultiSend := stdTx.Msgs[0].(bank.MsgMultiSend)

			var exportedAccts []exported.Account

			for _, input := range msgMultiSend.Inputs {
				inputAcct, err := ex.client.GetAccount(input.Address.String())
				if err != nil {
					return []schema.Account{}, err
				}

				exportedAccts = append(exportedAccts, inputAcct)
			}

			for _, output := range msgMultiSend.Outputs {
				outputAcct, err := ex.client.GetAccount(output.Address.String())
				if err != nil {
					return []schema.Account{}, err
				}

				exportedAccts = append(exportedAccts, outputAcct)
			}

			accounts, err = ex.getAccountAllAssets(exportedAccts, tx.TxHash, tx.Timestamp)
			if err != nil {
				return []schema.Account{}, err
			}

		case staking.MsgDelegate:
			zap.S().Infof("MsgType: %s | Hash: %s", stdTx.Msgs[0].Type(), tx.TxHash)

			msgDelegate := stdTx.Msgs[0].(staking.MsgDelegate)

			delegatorAddr, err := ex.client.GetAccount(msgDelegate.DelegatorAddress.String())
			if err != nil {
				return []schema.Account{}, err
			}

			valAccAddr, err := types.ConvertAccAddrFromValAddr(msgDelegate.DelegatorAddress.String())
			if err != nil {
				return []schema.Account{}, err
			}

			valAddr, err := ex.client.GetAccount(valAccAddr)
			if err != nil {
				return []schema.Account{}, err
			}

			exportedAccts := []exported.Account{
				delegatorAddr, valAddr,
			}

			accounts, err = ex.getAccountAllAssets(exportedAccts, tx.TxHash, tx.Timestamp)
			if err != nil {
				return []schema.Account{}, err
			}

		case staking.MsgUndelegate:
			zap.S().Infof("MsgType: %s | Hash: %s", stdTx.Msgs[0].Type(), tx.TxHash)

			msgUndelegate := stdTx.Msgs[0].(staking.MsgUndelegate)

			delegatorAddr, err := ex.client.GetAccount(msgUndelegate.DelegatorAddress.String())
			if err != nil {
				return []schema.Account{}, err
			}

			valAccAddr, err := types.ConvertAccAddrFromValAddr(msgUndelegate.DelegatorAddress.String())
			if err != nil {
				return []schema.Account{}, err
			}

			valAddr, err := ex.client.GetAccount(valAccAddr)
			if err != nil {
				return []schema.Account{}, err
			}

			exportedAccts := []exported.Account{
				delegatorAddr, valAddr,
			}

			accounts, err = ex.getAccountAllAssets(exportedAccts, tx.TxHash, tx.Timestamp)
			if err != nil {
				return []schema.Account{}, err
			}

		case staking.MsgBeginRedelegate:
			zap.S().Infof("MsgType: %s | Hash: %s", stdTx.Msgs[0].Type(), tx.TxHash)

			msgBeginRedelegate := stdTx.Msgs[0].(staking.MsgBeginRedelegate)

			delegatorAddr, err := ex.client.GetAccount(msgBeginRedelegate.DelegatorAddress.String())
			if err != nil {
				return []schema.Account{}, err
			}

			valSrcAccAddr, err := types.ConvertAccAddrFromValAddr(msgBeginRedelegate.ValidatorSrcAddress.String())
			if err != nil {
				return []schema.Account{}, err
			}

			valDstAccAddr, err := types.ConvertAccAddrFromValAddr(msgBeginRedelegate.ValidatorDstAddress.String())
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

			exportedAccts := []exported.Account{
				delegatorAddr, srcAddr, dstAddr,
			}

			accounts, err = ex.getAccountAllAssets(exportedAccts, tx.TxHash, tx.Timestamp)
			if err != nil {
				return []schema.Account{}, err
			}

		default:
			continue
		}
	}

	return accounts, nil
}

func (ex *Exporter) getAccountAllAssets(exportedAccts []exported.Account, txHashStr, txTime string) (accounts []schema.Account, err error) {
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
				LastTx:           txHashStr,
				LastTxTime:       txTime,
				CreationTime:     block.Block.Time.String(),
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
				LastTx:           txHashStr,
				LastTxTime:       txTime,
				CreationTime:     block.Block.Time.String(),
			}

			accounts = append(accounts, *acct)

			/*
				case *cosmosvesting.PeriodicVestingAccount:
					acc := account.(*cosmosvesting.PeriodicVestingAccount)

					spendable, rewards, commission, delegated, undelegated, err := ex.client.GetBaseAccountTotalAsset(acc.GetAddress().String())
					if err != nil {
						return []schema.Account{}, err
					}

					vesting := acc.GetVestingCoins(block.Block.Time)
					vested := acc.GetVestedCoins(block.Block.Time)

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
						LastTx:           txHashStr,
						LastTxTime:       txTime,
						CreationTime:     block.Block.Time.String(),
					}
					accounts = append(accounts, *acct)
			*/

			/*
				case *kavavesting.ValidatorVestingAccount:
					acc := account.(*kavavesting.ValidatorVestingAccount)

					spendable, rewards, commission, delegated, undelegated, err := ex.client.GetBaseAccountTotalAsset(acc.GetAddress().String())
					if err != nil {
						return []schema.Account{}, err
					}

					vesting := acc.GetVestingCoins(block.Block.Time)
					vested := acc.GetVestedCoins(block.Block.Time)
					failedVested := acc.GetFailedVestedCoins()

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
						ChainID:           chainID,
						AccountAddress:    acc.Address.String(),
						AccountNumber:     acc.AccountNumber,
						AccountType:       types.ValidatorVestingAccount,
						CoinsTotal:        total.String(),
						CoinsSpendable:    spendable.String(),
						CoinsRewards:      rewards.String(),
						CoinsCommission:   commission.String(),
						CoinsDelegated:    delegated.String(),
						CoinsUndelegated:  undelegated.String(),
						CoinsVesting:      vesting.String(),
						CoinsVested:       vested.String(),
						CoinsFailedVested: failedVested.String(),
						LastTx:            txHashStr,
						LastTxTime:        txTime,
						CreationTime:      block.Block.Time.String(),
					}

					accounts = append(accounts, *acct)
			*/

		default:
			return []schema.Account{}, fmt.Errorf("unrecognized account type: %T", account)
		}
	}

	return accounts, nil
}
