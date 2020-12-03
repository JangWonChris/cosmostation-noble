package handler

import (
	"fmt"
	"os"
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/gaia/app"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/client"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/codec"
	mintscanconfig "github.com/cosmostation/cosmostation-cosmos/mintscan/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/db"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	tlog "github.com/tendermint/tendermint/libs/log"

	tdb "github.com/tendermint/tm-db"
)

var iclient *client.Client
var idb *db.Database

func TestMain(m *testing.M) {
	config := mintscanconfig.ParseConfig()
	iclient, _ = client.NewClient(config.Node, config.Market)
	idb = db.Connect(config.DB)

	os.Exit(m.Run())
}

func TestModuleAccounts(t *testing.T) {
	tmdb := tdb.NewMemDB()
	gapp := app.NewGaiaApp(tlog.NewTMLogger(tlog.NewSyncWriter(os.Stdout)), tmdb, nil, true, map[int64]bool{}, "", uint(1), codec.EncodingConfig, nil)
	// sapp := simapp.NewSimApp(tlog.NewTMLogger(tlog.NewSyncWriter(os.Stdout)), tmdb, nil, true, 0)

	modAccAddrs := gapp.ModuleAccountAddrs()
	authtypes.ModuleCdc.RegisterConcrete(&authtypes.BaseAccount{}, "cosmos-sdk/ModuleAccount", nil)

	a := crypto.AddressHash([]byte("fee_collector"))
	s, _ := bech32.ConvertAndEncode(sdktypes.GetConfig().GetBech32AccountAddrPrefix(), a.Bytes())
	fmt.Println("fee_collector:", s)

	for mAccAddr, permission := range modAccAddrs {
		require.Equal(t, true, permission)
		t.Log(mAccAddr)
		t.Log(permission)
		account, err := iclient.GetAccount(mAccAddr)
		t.Log(account)
		require.NoError(t, err)

		acc, ok := account.(authtypes.ModuleAccountI)
		require.Equal(t, true, ok)

		require.NotNil(t, acc.GetAddress().String(), "Module Account Address")
		// require.NotNil(t, acc.GetCoins(), "Module Account Balance")
		require.NotNil(t, acc.GetName(), "Module Account Name")
		require.NotNil(t, acc.GetPermissions(), "Module Account Permissions")
	}
}

func TestKindOfBalance(t *testing.T) {
	address := "cosmos1x5wgh6vwye60wv3dtshs9dmqggwfx2ldnqvev0"
	coin, err := iclient.GetBalance(address)
	require.NoError(t, err)
	fmt.Println("getbalance :", coin.Denom)
	fmt.Println("getbalance :", coin.Amount)

	coins, err := iclient.GetAllBalances(address)
	require.NoError(t, err)

	for _, coin := range coins {
		fmt.Println("getallbalances :", coin.Denom)
		fmt.Println("getallbalances :", coin.Amount)
	}
}
