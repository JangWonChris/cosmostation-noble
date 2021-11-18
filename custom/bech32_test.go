package custom

import (
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	SetAppConfig()

	os.Exit(m.Run())
}

func TestBech32PrefixesToAcctAddr(t *testing.T) {
	_, addr1, err := bech32.DecodeAndConvert("kava140g8fnnl46mlvfhygj3zvjqlku6x0fwuhfj3uf")
	require.NoError(t, err)

	_, addr2, err := bech32.DecodeAndConvert("kavavaloper140g8fnnl46mlvfhygj3zvjqlku6x0fwu6lgey7")
	require.NoError(t, err)

	_, addr3, err := bech32.DecodeAndConvert("cosmosvalconspub1zcjduepq5mhsvc5685267fg2ee5uv30srjjxzetp8msfs3h983vz724496lqtaz884")
	require.NoError(t, err)

	bech1, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32AccountAddrPrefix(), addr1)
	require.NoError(t, err)

	bech2, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32AccountAddrPrefix(), addr2)
	require.NoError(t, err)

	// bech3, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32AccountAddrPrefix(), addr3)
	bech3, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32ConsensusPubPrefix(), addr3)
	require.NoError(t, err)

	require.NotNil(t, bech1)
	require.NotNil(t, bech2)
	require.NotNil(t, bech3)
}
