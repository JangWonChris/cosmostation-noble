package types

import (
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"

	ltypes "github.com/cosmostation/mintscan-backend-library/types"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// SetAppConfig()

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

func TestPrivValidatorKey(t *testing.T) {
	// bech32hex := "AE61EC6FA0450C6327288B946E233A88683C478A"
	// gaiad tendermint show-address : cosmosvalcons14es7cmaqg5xxxfeg3w2xuge63p5rc3u2vt8ym4
	// gaiad tendermint show-validator : cosmosvalconspub1zcjduepq5mhsvc5685267fg2ee5uv30srjjxzetp8msfs3h983vz724496lqtaz884
	valcon := "cosmosvalcons14es7cmaqg5xxxfeg3w2xuge63p5rc3u2vt8ym4"
	valconpub := "cosmosvalconspub1zcjduepq5mhsvc5685267fg2ee5uv30srjjxzetp8msfs3h983vz724496lqtaz884"
	_, _ = valcon, valconpub

	consAddrStr, err := ltypes.ConvertConsAddrFromConsPubkey(valconpub)
	require.NoError(t, err)

	consAddr, err := sdk.ConsAddressFromHex(consAddrStr)
	require.NoError(t, err)

	madevalcon, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32ConsensusAddrPrefix(), consAddr.Bytes())
	require.NoError(t, err)
	require.Equal(t, valcon, madevalcon)

	pubkey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, valconpub)
	require.NoError(t, err)

	madevalconpub, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, pubkey)
	require.NoError(t, err)
	require.Equal(t, valconpub, madevalconpub)
}
