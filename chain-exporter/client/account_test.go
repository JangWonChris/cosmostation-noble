package client

import (
	"encoding/json"
	"log"
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankexported "github.com/cosmos/cosmos-sdk/x/bank/exported"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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

	var genesisState map[string]json.RawMessage
	if err = json.Unmarshal(genDoc.AppState, &genesisState); err != nil {
		log.Println(err, "failed to unmarshal genesis state")
	}
	// a := genesisState[authtypes.ModuleName]
	// log.Println(string(a)) //print message that key is auth {...}
	authGenesisState := authtypes.GetGenesisStateFromAppState(codec.AppCodec, genesisState)

	authAccs := authGenesisState.GetAccounts()
	NumberOfTotalAccounts := len(authAccs)
	accountMapper := make(map[string]*schema.Account, NumberOfTotalAccounts)
	for _, authAcc := range authAccs {
		var ga authtypes.GenesisAccount
		codec.AppCodec.UnpackAny(authAcc, &ga)
		// log.Println(authAcc.GetTypeUrl())
		// log.Println(ga.GetAddress().String())
		// log.Println(ga.GetAccountNumber())
		sAcc := schema.Account{
			ChainID:        genDoc.ChainID,
			AccountAddress: ga.GetAddress().String(),
			AccountNumber:  ga.GetAccountNumber(),
			AccountType:    authAcc.GetTypeUrl(), //type 변경
			CreationTime:   genDoc.GenesisTime.String(),
		}
		accountMapper[ga.GetAddress().String()] = &sAcc
	}

	balIter := banktypes.GenesisBalancesIterator{}
	balIter.IterateGenesisBalances(cli.cliCtx.JSONMarshaler, genesisState,
		func(bal bankexported.GenesisBalance) (stop bool) {
			accAddress := bal.GetAddress()
			accCoins := bal.GetCoins()

			accountMapper[accAddress.String()].CoinsSpendable = accCoins.AmountOf("umuon").Uint64()

			return false
		},
	)

	for _, acc := range accountMapper {
		accounts = append(accounts, *acc)
		log.Println(acc)
	}

}
