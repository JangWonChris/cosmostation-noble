package client

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankexported "github.com/cosmos/cosmos-sdk/x/bank/exported"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	tmtypes "github.com/tendermint/tendermint/types"
)

func TestGetGenesisAccount(t *testing.T) {

	var accounts []schema.Account
	// genesisFile := os.Getenv("PWD") + "/genesis.json"
	genesisFile := "/Users/jeonghwan/dev/cosmostation/cosmostation-cosmos/chain-exporter/genesis.json"
	genDoc, err := tmtypes.GenesisDocFromFile(genesisFile)
	if err != nil {
		log.Println(err, "failed to read genesis doc file %s", genesisFile)
	}
	log.Println("genesis_time :", genDoc.GenesisTime)
	log.Println("chainid :", genDoc.ChainID)
	log.Println("initial_height :", genDoc.InitialHeight)

	var appState map[string]json.RawMessage
	if err = json.Unmarshal(genDoc.AppState, &appState); err != nil {
		log.Println(err, "failed to unmarshal genesis state")
	}
	// a := appState[authtypes.ModuleName]
	// log.Println(string(a)) //print message that key is auth {...}
	authGenesisState := authtypes.GetGenesisStateFromAppState(codec.AppCodec, appState)
	// stakingGenesisState := stakingtypes.GetGenesisStateFromAppState(codec.AppCodec, appState)
	// bondDenom := stakingGenesisState.Params.BondDenom
	var distributionGenesisState distributiontypes.GenesisState
	if appState[distributiontypes.ModuleName] != nil {
		codec.AppCodec.MustUnmarshalJSON(appState[distributiontypes.ModuleName], &distributionGenesisState)
		log.Println("abcded :", distributionGenesisState.DelegatorStartingInfos[0].DelegatorAddress)
		log.Println("counts :", len(distributionGenesisState.DelegatorStartingInfos))
	}
	bondDenom := "uatom"

	authAccs := authGenesisState.GetAccounts()
	NumberOfTotalAccounts := len(authAccs)
	accountMapper := make(map[string]*schema.Account, NumberOfTotalAccounts)
	for i, authAcc := range authAccs {
		var ga authtypes.GenesisAccount
		codec.AppCodec.UnpackAny(authAcc, &ga)
		switch ga := ga.(type) {
		case *authtypes.BaseAccount:
		case *authvestingtypes.DelayedVestingAccount:
			/* Endtime 이 지난 vesting account 데이터는 의미 없다. */
			log.Println("DelayedVestingAccount", ga.String())
			log.Println("DelayedVestingAccount", ga.String())
			log.Println("delegated Free :", ga.GetDelegatedFree())
			log.Println("delegated vesting :", ga.GetDelegatedVesting())
			log.Println("vested coins:", ga.GetVestedCoins(time.Unix(1584140400, 0))) // 주어진 시간에 vesting이 풀린 코인
			log.Println("vesting coins :", ga.GetVestingCoins(time.Now()))            // 주어진 시간에 vesting 중인 코인
			log.Println("original vesting :", ga.GetOriginalVesting())
			os.Exit(1)
		case *authvestingtypes.ContinuousVestingAccount:
			log.Println("ContinuousVestingAccount", ga.String())
		case *authvestingtypes.PeriodicVestingAccount:
			log.Println("PeriodicVestingAccount", ga.String())
		}
		// log.Println(authAcc.GetTypeUrl())
		// log.Println(ga.GetAddress().String())
		// log.Println(ga.GetAccountNumber())
		sAcc := schema.Account{
			ChainID:        genDoc.ChainID,
			AccountAddress: ga.GetAddress().String(),
			AccountNumber:  uint64(i),            //account number is set by specified order in genesis file
			AccountType:    authAcc.GetTypeUrl(), //type 변경
			CreationTime:   genDoc.GenesisTime.String(),
		}
		accountMapper[ga.GetAddress().String()] = &sAcc
	}

	balIter := banktypes.GenesisBalancesIterator{}
	balIter.IterateGenesisBalances(cli.cliCtx.JSONMarshaler, appState,
		func(bal bankexported.GenesisBalance) (stop bool) {
			accAddress := bal.GetAddress()
			accCoins := bal.GetCoins()

			// accountMapper[accAddress.String()].CoinsSpendable = *accCoins.AmountOf(bondDenom).String()
			accountMapper[accAddress.String()].CoinsSpendable = accCoins.AmountOf(bondDenom).String()
			return false
		},
	)

	for _, acc := range accountMapper {
		accounts = append(accounts, *acc)
		// log.Println(acc)
	}

}
