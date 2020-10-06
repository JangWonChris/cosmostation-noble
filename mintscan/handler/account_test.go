package handler

import (
	"fmt"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/client"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/db"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/bech32"
	tlog "github.com/tendermint/tendermint/libs/log"

	tdb "github.com/tendermint/tm-db"
)

var iclient *client.Client
var idb *db.Database

func TestMain(m *testing.M) {
	config := config.ParseConfig()
	iclient, _ = client.NewClient(config.Node, config.Market)
	idb = db.Connect(config.DB)

	os.Exit(m.Run())
}

func TestModuleAccounts(t *testing.T) {
	tmdb := tdb.NewMemDB()
	// gapp := app.NewGaiaApp(tlog.NewTMLogger(tlog.NewSyncWriter(os.Stdout)), tmdb, nil, true, 0)
	sapp := simapp.NewSimApp(tlog.NewTMLogger(tlog.NewSyncWriter(os.Stdout)), tmdb, nil, true, 0)

	// modAccAddrs := gapp.ModuleAccountAddrs()
	modAccAddrs := sapp.ModuleAccountAddrs()
	authtypes.ModuleCdc.RegisterConcrete(&auth.BaseAccount{}, "cosmos-sdk/ModuleAccount", nil)

	a := crypto.AddressHash([]byte("fee_collector"))
	s, _ := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32AccountAddrPrefix(), a.Bytes())
	fmt.Println("fee_collector:", s)

	for mAccAddr, permission := range modAccAddrs {
		require.Equal(t, true, permission)
		t.Log(mAccAddr)
		t.Log(permission)
		account, err := iclient.GetAccount(mAccAddr)
		t.Log(account)
		require.NoError(t, err)

		acc, ok := account.(exported.ModuleAccountI)
		require.Equal(t, true, ok)

		require.NotNil(t, acc.GetAddress().String(), "Module Account Address")
		require.NotNil(t, acc.GetCoins(), "Module Account Balance")
		require.NotNil(t, acc.GetName(), "Module Account Name")
		require.NotNil(t, acc.GetPermissions(), "Module Account Permissions")
	}
}
