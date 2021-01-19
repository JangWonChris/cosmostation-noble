package exporter

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	//internal
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	//mbl
	"github.com/cosmostation/mintscan-backend-library/db/schema"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankexported "github.com/cosmos/cosmos-sdk/x/bank/exported"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	gaia "github.com/cosmos/gaia/v3/app"
	tmconfig "github.com/tendermint/tendermint/config"
	tmtypes "github.com/tendermint/tendermint/types"
)

func TestGetGenesisStateFromGenesisFile(t *testing.T) {
	var accounts []schema.AccountCoin
	// genesisFile := os.Getenv("PWD") + "/genesis.json"
	baseConfig := tmconfig.DefaultBaseConfig()
	genesisFile := filepath.Join(gaia.DefaultNodeHome, baseConfig.Genesis)
	// genesisFile := "/Users/jeonghwan/dev/cosmostation/cosmostation-cosmos/chain-exporter/genesis.json"
	log.Println("genesis file path :", genesisFile)
	genesisFile = "../cosmoshub-test-stargate-e.json"
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
	stakingGenesisState := stakingtypes.GetGenesisStateFromAppState(codec.AppCodec, appState)
	bondDenom := stakingGenesisState.Params.BondDenom
	lastValidatorPowers := stakingGenesisState.GetLastValidatorPowers()
	for _, val := range lastValidatorPowers {
		log.Println(val.Power)
	}
	os.Exit(0)
	var distributionGenesisState distributiontypes.GenesisState
	if appState[distributiontypes.ModuleName] != nil {
		codec.AppCodec.MustUnmarshalJSON(appState[distributiontypes.ModuleName], &distributionGenesisState)
		log.Println("abcded :", distributionGenesisState.DelegatorStartingInfos[0].DelegatorAddress)
		log.Println("counts :", len(distributionGenesisState.DelegatorStartingInfos))
	}
	// bondDenom := "uatom"

	authAccs := authGenesisState.GetAccounts()
	NumberOfTotalAccounts := len(authAccs)
	accountMapper := make(map[string]*schema.AccountCoin, NumberOfTotalAccounts)
	for _, authAcc := range authAccs {
		var ga authtypes.GenesisAccount
		codec.AppCodec.UnpackAny(authAcc, &ga)
		switch ga := ga.(type) {
		case *authtypes.BaseAccount:
		case *authvestingtypes.DelayedVestingAccount:
			/* Endtime 이 지난 vesting account 데이터는 의미 없다. */
			// ibc tokens은 delegate이 불가능하다 (stargate-5), bondDenom만 담는 것으로 하자.
			log.Printf("type %T\n", ga)
			log.Println("DelayedVestingAccount", ga.String())
			log.Println("delegated Free :", ga.GetDelegatedFree().AmountOf(bondDenom))
			log.Println("delegated vesting :", ga.GetDelegatedVesting().AmountOf(bondDenom))
			log.Println("vested coins:", ga.GetVestedCoins(time.Now()).AmountOf(bondDenom))    // 주어진 시간에 vesting이 풀린 코인
			log.Println("vesting coins :", ga.GetVestingCoins(time.Now()).AmountOf(bondDenom)) // 주어진 시간에 vesting 중인 코인
			log.Println("original vesting :", ga.GetOriginalVesting().AmountOf(bondDenom))
		case *authvestingtypes.ContinuousVestingAccount:
			log.Println("ContinuousVestingAccount", ga.String())
		case *authvestingtypes.PeriodicVestingAccount:
			log.Println("PeriodicVestingAccount", ga.String())
		}
		// log.Println(authAcc.GetTypeUrl())
		// log.Println(ga.GetAddress().String())
		// log.Println(ga.GetAccountNumber())
		sAcc := schema.AccountCoin{
			// ChainID:        genDoc.ChainID,
			AccountAddress: ga.GetAddress().String(),
			// AccountNumber:  ga.GetAccountNumber(), //account number is set by specified order in genesis file
			// AccountType:    authAcc.GetTypeUrl(),  //type 변경
			Total:        "0",
			Available:    "0",
			Delegated:    "0",
			Rewards:      "0",
			Commission:   "0",
			Undelegated:  "0",
			FailedVested: "0",
			Vested:       "0",
			Vesting:      "0",
			// CreationTime: genDoc.GenesisTime.String(),
		}
		accountMapper[ga.GetAddress().String()] = &sAcc
	}

	balIter := banktypes.GenesisBalancesIterator{}
	balIter.IterateGenesisBalances(codec.AppCodec, appState,
		func(bal bankexported.GenesisBalance) (stop bool) {
			accAddress := bal.GetAddress()
			accCoins := bal.GetCoins()

			// accountMapper[accAddress.String()].CoinsSpendable = *accCoins.AmountOf(bondDenom).String()
			accountMapper[accAddress.String()].Available = accCoins.AmountOf(bondDenom).String()
			return false
		},
	)

	for _, acc := range accountMapper {
		accounts = append(accounts, *acc)
		// log.Println(acc)
	}

}

func TestExporterNil(t *testing.T) {
	s := new(schema.ExportData)
	log.Println(s)
	log.Println(s.ResultAccounts)
	log.Println(len(s.ResultTxs))
	log.Println(s.ResultTxs)
}
