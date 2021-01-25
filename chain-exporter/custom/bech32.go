package custom

import (
	"log"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

// "github.com/tendermint/tendermint/libs/bech32"

func IsSetBech32() bool {

	/* 체인 별 Bech32PrefixAccAddr을 비교하도록 코드를 변경 해야한다 */
	if sdktypes.GetConfig().GetBech32AccountAddrPrefix() != sdktypes.Bech32PrefixAccAddr {
		log.Println("bech32 is not identical, will set config ")
		return false
	}
	return true
}

// SetAppConfig creates a new config instance for the SDK configuration.
//안써도 지우지 않는 이유는 다른 네트워크에서 이 함수를 사용하기 때문임
func SetAppConfig() {
	if !IsSetBech32() {
		// config := sdk.GetConfig()
		// SetBech32AddressPrefixes(config)
		// SetBip44CoinType(config)
		// config.Seal()
	}
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
