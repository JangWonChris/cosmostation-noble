package model

import (
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/libs/bech32"
)

var (
	// SampleBech32PrefixAccAddr defines the Bech32 prefix of an account's address
	SampleBech32PrefixAccAddr = "kava140g8fnnl46mlvfhygj3zvjqlku6x0fwuhfj3uf"

	// SampleBech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	SampleBech32PrefixValAddr = "kavavaloper140g8fnnl46mlvfhygj3zvjqlku6x0fwu6lgey7"

	// SampleBech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	SampleBech32PrefixConsAddr = ""

	// SampleBech32PrefixConsAddrHex defines the Bech32 prefix of a consensus node address
	SampleBech32PrefixConsAddrHex = "3D6468FCB5EC366714EF86E5263C0B30C11734FB"

	// SampleBech32PrefixConsPub defines the Bech32 prefix of a consensus public key
	SampleBech32PrefixConsPub = "kavavalconspub1zcjduepqvtvkhh22hgfvp865tj4uwltv0hu7fs3vwmxwrl0n2mdpfuzj0p0qes2k9e"
)

func TestMain(m *testing.M) {
	SetAppConfig()

	os.Exit(m.Run())
}

func TestConvertConsAddrFromConsPubkey(t *testing.T) {
	consAddr, err := ConvertConsAddrFromConsPubkey(SampleBech32PrefixConsPub)
	require.NoError(t, err)

	require.Equal(t, consAddr, SampleBech32PrefixConsAddrHex)
}

func TestConvertAccAddrFromValAddr(t *testing.T) {
	accAddr, err := ConvertAccAddrFromValAddr(SampleBech32PrefixValAddr)
	require.NoError(t, err)

	require.Equal(t, accAddr, SampleBech32PrefixAccAddr)
}

func TestConvertValAddrFromAccAddr(t *testing.T) {
	valAddr, err := ConvertValAddrFromAccAddr(SampleBech32PrefixValAddr)
	require.NoError(t, err)

	require.Equal(t, valAddr, SampleBech32PrefixValAddr)
}

func TestVerifyAccAddr(t *testing.T) {
	err := VerifyBech32AccAddr(SampleBech32PrefixAccAddr)
	require.NoError(t, err)
}

func TestVerifyValAddr(t *testing.T) {
	err := VerifyBech32ValAddr(SampleBech32PrefixValAddr)
	require.NoError(t, err)
}

func TestBech32PrefixesToAcctAddr(t *testing.T) {
	_, addr1, err := bech32.DecodeAndConvert(SampleBech32PrefixAccAddr)
	require.NoError(t, err)

	_, addr2, err := bech32.DecodeAndConvert(SampleBech32PrefixValAddr)
	require.NoError(t, err)

	_, addr3, err := bech32.DecodeAndConvert(SampleBech32PrefixConsAddr)
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
