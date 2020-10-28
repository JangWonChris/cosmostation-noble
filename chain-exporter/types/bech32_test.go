package types

import (
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/stretchr/testify/require"
	// "github.com/tendermint/tendermint/libs/bech32"
)

func TestMain(m *testing.M) {
	// SetAppConfig()

	os.Exit(m.Run())
}

func TestConvertConsAddrFromConsPubkey(t *testing.T) {
	consAddr, err := ConvertConsAddrFromConsPubkey("kavavalconspub1zcjduepqvtvkhh22hgfvp865tj4uwltv0hu7fs3vwmxwrl0n2mdpfuzj0p0qes2k9e")
	require.NoError(t, err)

	require.Equal(t, consAddr, "3D6468FCB5EC366714EF86E5263C0B30C11734FB")
}

func TestConvertAccAddrFromValAddr(t *testing.T) {
	accAddr, err := ConvertAccAddrFromValAddr("kavavaloper140g8fnnl46mlvfhygj3zvjqlku6x0fwu6lgey7")
	require.NoError(t, err)

	require.Equal(t, accAddr, "kava140g8fnnl46mlvfhygj3zvjqlku6x0fwuhfj3uf")
}

func TestConvertValAddrFromAccAddr(t *testing.T) {
	valAddr, err := ConvertValAddrFromAccAddr("kavavaloper140g8fnnl46mlvfhygj3zvjqlku6x0fwu6lgey7")
	require.NoError(t, err)

	require.Equal(t, valAddr, "kavavaloper140g8fnnl46mlvfhygj3zvjqlku6x0fwu6lgey7")
}

func TestVerifyAccAddr(t *testing.T) {
	ok, err := VerifyBech32AccAddr("kava140g8fnnl46mlvfhygj3zvjqlku6x0fwuhfj3uf")
	require.NoError(t, err)

	require.Equal(t, true, ok)
}

func TestVerifyValAddr(t *testing.T) {
	ok, err := VerifyBech32ValAddr("kavavaloper140g8fnnl46mlvfhygj3zvjqlku6x0fwu6lgey7")
	require.NoError(t, err)

	require.Equal(t, true, ok)
}

func TestBech32PrefixesToAcctAddr(t *testing.T) {
	_, addr1, err := bech32.DecodeAndConvert("kava140g8fnnl46mlvfhygj3zvjqlku6x0fwuhfj3uf")
	require.NoError(t, err)

	_, addr2, err := bech32.DecodeAndConvert("kavavaloper140g8fnnl46mlvfhygj3zvjqlku6x0fwu6lgey7")
	require.NoError(t, err)

	_, addr3, err := bech32.DecodeAndConvert("")
	require.NoError(t, err)

	bech1, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32AccountAddrPrefix(), addr1)
	require.NoError(t, err)

	bech2, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32AccountAddrPrefix(), addr2)
	require.NoError(t, err)

	bech3, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32AccountAddrPrefix(), addr3)
	require.NoError(t, err)

	require.NotNil(t, bech1)
	require.NotNil(t, bech2)
	require.NotNil(t, bech3)
}
