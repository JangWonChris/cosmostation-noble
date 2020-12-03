package exporter

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
	"go.uber.org/zap"

	//gaia
	gaia "github.com/cosmos/gaia/app"

	//cosmos-sdk
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankexported "github.com/cosmos/cosmos-sdk/x/bank/exported"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	//tendermint
	tmconfig "github.com/tendermint/tendermint/config"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	// startingHeight is used to extract genesis accounts and parse their assets.
	startingHeight = int64(1)
)

// GetGenesisStateFromGenesisFile get the genesis account information from genesis state ({NODE_HOME}/config/Genesis.json)
func (ex *Exporter) GetGenesisStateFromGenesisFile(genesisPath string) (accounts []schema.Account, err error) {

	// genesisFile := os.Getenv("PWD") + "/genesis.json"
	baseConfig := tmconfig.DefaultBaseConfig()
	genesisFile := filepath.Join(gaia.DefaultNodeHome, baseConfig.Genesis)

	if genesisPath == "" {
		genesisPath = genesisFile
	}
	// genesisFile := "/Users/jeonghwan/dev/cosmostation/cosmostation-cosmos/chain-exporter/genesis.json"
	genDoc, err := tmtypes.GenesisDocFromFile(genesisPath)
	if err != nil {
		log.Println(err, "failed to read genesis doc file %s", genesisPath)
		return
	}

	var genesisState map[string]json.RawMessage
	if err = json.Unmarshal(genDoc.AppState, &genesisState); err != nil {
		log.Println(err, "failed to unmarshal genesis state")
		return
	}
	// a := genesisState[authtypes.ModuleName]
	// log.Println(string(a)) //print message that key is auth {...}
	authGenesisState := authtypes.GetGenesisStateFromAppState(codec.AppCodec, genesisState)
	stakingGenesisState := stakingtypes.GetGenesisStateFromAppState(codec.AppCodec, genesisState)
	bondDenom := stakingGenesisState.GetParams().BondDenom

	authAccs := authGenesisState.GetAccounts()
	NumberOfTotalAccounts := len(authAccs)
	accountMapper := make(map[string]*schema.Account, NumberOfTotalAccounts)
	for i, authAcc := range authAccs {
		var ga authtypes.GenesisAccount
		codec.AppCodec.UnpackAny(authAcc, &ga)
		switch ga := ga.(type) {
		case *authtypes.BaseAccount:
		case *authvestingtypes.DelayedVestingAccount:
			log.Println("DelayedVestingAccount", ga.String())
			log.Println("delegated Free :", ga.GetDelegatedFree())
			log.Println("delegated vesting :", ga.GetDelegatedVesting())
			log.Println("vested coins:", ga.GetVestedCoins(time.Now()))
			log.Println("vesting coins :", ga.GetVestingCoins(time.Now()))
			log.Println("original vesting :", ga.GetOriginalVesting())
		case *authvestingtypes.ContinuousVestingAccount:
			log.Println("ContinuousVestingAccount", ga.String())
		case *authvestingtypes.PeriodicVestingAccount:
			log.Println("PeriodicVestingAccount", ga.String())
		}
		sAcc := schema.Account{
			ChainID:           genDoc.ChainID,
			AccountAddress:    ga.GetAddress().String(),
			AccountNumber:     uint64(i),            //account number is set by specified order in genesis file
			AccountType:       authAcc.GetTypeUrl(), //type 변경
			CoinsTotal:        "0",
			CoinsSpendable:    "0",
			CoinsDelegated:    "0",
			CoinsRewards:      "0",
			CoinsCommission:   "0",
			CoinsUndelegated:  "0",
			CoinsFailedVested: "0",
			CoinsVested:       "0",
			CoinsVesting:      "0",
			CreationTime:      genDoc.GenesisTime.String(),
		}
		accountMapper[ga.GetAddress().String()] = &sAcc
	}

	balIter := banktypes.GenesisBalancesIterator{}
	balIter.IterateGenesisBalances(codec.AppCodec, genesisState,
		func(bal bankexported.GenesisBalance) (stop bool) {
			accAddress := bal.GetAddress()
			accCoins := bal.GetCoins()

			// accountMapper[accAddress.String()].CoinsSpendable = *accCoins.AmountOf(bondDenom).BigInt()
			accountMapper[accAddress.String()].CoinsSpendable = accCoins.AmountOf(bondDenom).String()
			return false
		},
	)

	for _, acc := range accountMapper {
		accounts = append(accounts, *acc)
		log.Println(acc)
		log.Println(acc.CoinsSpendable)
	}

	ex.db.InsertGenesisAccount(accounts)

	return
}

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
