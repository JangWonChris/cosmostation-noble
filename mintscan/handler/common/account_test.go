package common

import (
	"context"
	"fmt"
	"os"
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	app "github.com/cosmos/gaia/v3/app"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/client"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/db"
	"github.com/cosmostation/mintscan-backend-library/codec"
	mintscanconfig "github.com/cosmostation/mintscan-backend-library/config"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	tlog "github.com/tendermint/tendermint/libs/log"

	tdb "github.com/tendermint/tm-db"
)

var iclient *client.Client
var idb *db.Database

func TestMain(m *testing.M) {
	config := mintscanconfig.ParseConfig()
	iclient = client.NewClient(&config.Client)
	idb = db.Connect(&config.DB)

	os.Exit(m.Run())
}

type gaiaInit struct{}

func (g gaiaInit) Get(s string) interface{} {
	return false
}

func TestModuleAccounts(t *testing.T) {
	tmdb := tdb.NewMemDB()
	gapp := app.NewGaiaApp(tlog.NewTMLogger(tlog.NewSyncWriter(os.Stdout)), tmdb, nil, true, map[int64]bool{}, "", uint(1), codec.EncodingConfig, gaiaInit{})
	// gapp := app.NewGaiaApp(tlog.NewTMLogger(tlog.NewSyncWriter(os.Stdout)), tmdb, nil, true, uint(1))
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
		account, err := iclient.CliCtx.GetAccount(mAccAddr)
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
	ctx := context.Background()
	coin, err := iclient.GRPC.GetBalance(ctx, "umuon", address)
	require.NoError(t, err)
	fmt.Println("getbalance :", coin.Denom)
	fmt.Println("getbalance :", coin.Amount)

	coins, err := iclient.GRPC.GetAllBalances(ctx, address, 100)
	require.NoError(t, err)

	for _, coin := range coins {
		fmt.Println("getallbalances :", coin.Denom)
		fmt.Println("getallbalances :", coin.Amount)
	}
}
