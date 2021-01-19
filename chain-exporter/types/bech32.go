package types

// "github.com/tendermint/tendermint/libs/bech32"

// SetAppConfig creates a new config instance for the SDK configuration.
func SetAppConfig() {
	// 	config := sdk.GetConfig()
	// 	SetBech32AddressPrefixes(config)
	// 	SetBip44CoinType(config)
	// 	config.Seal()
}

// SetBech32AddressPrefixes sets the global prefix to be used when serializing addresses to bech32 strings.
// func SetBech32AddressPrefixes(config *sdk.Config) {
// 	config.SetBech32PrefixForAccount(app.Bech32MainPrefix, app.Bech32MainPrefix+sdk.PrefixPublic)
// 	config.SetBech32PrefixForValidator(app.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator, app.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic)
// 	config.SetBech32PrefixForConsensusNode(app.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus, app.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic)
// }

// SetBip44CoinType sets the global coin type to be used in hierarchical deterministic wallets.
// func SetBip44CoinType(config *sdk.Config) {
// 	config.SetCoinType(app.Bip44CoinType)
// }
