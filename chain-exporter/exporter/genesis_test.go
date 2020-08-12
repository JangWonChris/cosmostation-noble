package exporter

import (
	_ "fmt"
	_ "testing"

	_ "github.com/stretchr/testify/require"

	_ "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/x/auth/exported"
)

/*
func TestGenesisAccounts(t *testing.T) {
	addr1, _ := sdkTypes.AccAddressFromBech32("kava1gt8nvvwfzu7g5dxu4wyha99m37dpw6e6uqz8wu")
	addr2, _ := sdkTypes.AccAddressFromBech32("kava1cdsplflzkwcyx8kz26v07m6q3ucrttfps3qp8v")
	addr3, _ := sdkTypes.AccAddressFromBech32("kava1uxn53yf9w3kuus8w77thuucahjaxtwmfv6g2rd")
	addr4, _ := sdkTypes.AccAddressFromBech32("kava1e26tszfuwjyt4snjmpl8uetavlrkfmtewwrf6z")
	addr5, _ := sdkTypes.AccAddressFromBech32("kava1spyhfn250qwtxusqwc6yrh3mmtg3j83ce57vc9")
	addr6, _ := sdkTypes.AccAddressFromBech32("kava1v5hrqlv8dqgzvy0pwzqzg0gxy899rm4khzp04j")

	testCases := []struct {
		accAddr     sdkTypes.AccAddress
		accAddrType string
		period      string
	}{
		{addr1, "PeriodicVestingAccount", "Vesting All Left"},
		{addr2, "PeriodicVestingAccount", "Vesting Half Left"},
		{addr3, "PeriodicVestingAccount", "Vesting Over"},
		{addr4, "ValidatorVestingAccount", "Cosmostation"},
		{addr5, "ValidatorVestingAccount", "StakeWithUs"},
		{addr6, "ValidatorVestingAccount", "DokiaCapital"},
	}

	var genesisAccts exported.GenesisAccounts

	for _, tc := range testCases {
		var genesisAcct exported.GenesisAccount
		err := genesisAcct.SetAddress(tc.accAddr)
		require.NoError(t, err)

		genesisAccts = append(genesisAccts, genesisAcct)
	}

	accounts, err := ex.getGenesisAccounts(genesisAccts)
	require.NoError(t, err)

	for _, acc := range accounts {
		fmt.Println(acc)
	}
}
*/
