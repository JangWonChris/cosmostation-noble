package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	// "github.com/tendermint/tendermint/libs/bech32"
)

func TestConvertConsAddrFromConsPubkey(t *testing.T) {
	consAddr, err := ConvertConsAddrFromConsPubkey("cosmosvalconspub1zcjduepq0dc9apn3pz2x2qyujcnl2heqq4aceput2uaucuvhrjts75q0rv5smjjn7v")
	require.NoError(t, err)

	require.Equal(t, consAddr, "099E2B09583331AFDE35E5FA96673D2CA7DEA316")
}

func TestConvertAccAddrFromValAddr(t *testing.T) {
	accAddr, err := ConvertAccAddrFromValAddr("cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn")
	require.NoError(t, err)

	require.Equal(t, accAddr, "cosmos1clpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9q")
}

func TestConvertValAddrFromAccAddr(t *testing.T) {
	valAddr, err := ConvertValAddrFromAccAddr("cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn")
	require.NoError(t, err)

	require.Equal(t, valAddr, "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn")
}

func TestVerifyAccAddr(t *testing.T) {
	err := VerifyBech32AccAddr("cosmos1clpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9q")
	require.NoError(t, err)
}

func TestVerifyValAddr(t *testing.T) {
	err := VerifyBech32ValAddr("cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn")
	require.NoError(t, err)
}

func TestBech32PrefixesToAcctAddr(t *testing.T) {
	_, addr1, err := bech32.DecodeAndConvert("cosmos1clpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9q")
	require.NoError(t, err)

	_, addr2, err := bech32.DecodeAndConvert("cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn")
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
