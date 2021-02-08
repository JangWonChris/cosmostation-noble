package common

import (
	"context"
	"fmt"
	"os"
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	app "github.com/cosmos/gaia/v4/app"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/client"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/db"
	mintscanconfig "github.com/cosmostation/mintscan-backend-library/config"
	"github.com/stretchr/testify/require"
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
	maccPerms := app.GetMaccPerms()

	for name := range maccPerms {
		maccAddr := authtypes.NewModuleAddress(name).String()
		// require.Equal(t, true, permission)
		account, err := iclient.CliCtx.GetAccount(maccAddr)
		require.NoError(t, err)
		// t.Log(account)

		acc, ok := account.(authtypes.ModuleAccountI)
		require.Equal(t, true, ok)

		pageLimit := uint64(100)
		coins, err := iclient.GRPC.GetAllBalances(context.Background(), maccAddr, pageLimit)
		require.NoError(t, err)

		t.Log("=======")
		t.Log("Module Account Name : ", acc.GetName())
		t.Log("Module Account Permissions : ", acc.GetPermissions())
		t.Log("Module Account Address : ", acc.GetAddress().String())
		t.Log("Module Account Balance : ", coins)
	}
}

func TestKindOfBalance(t *testing.T) {
	address := "cosmos1x5wgh6vwye60wv3dtshs9dmqggwfx2ldnqvev0"

	// acc, err := iclient.CliCtx.GetAccount(address)

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
