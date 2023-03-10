package custom

import (
	"fmt"
	"log"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	chainapp "github.com/strangelove-ventures/noble/app"
)

func init() {
	SetAppConfig()
	if !IsSetAppConfig() {
		panic(fmt.Errorf("bech32 is not set corretly"))
	}
	log.Println("Current bech32 : ", sdktypes.GetConfig())
}

func IsSetAppConfig() bool {
	/* 체인 별 Bech32PrefixAccAddr을 비교하도록 코드를 변경 해야한다 */
	if sdktypes.GetConfig().GetBech32AccountAddrPrefix() != chainapp.AccountAddressPrefix {
		log.Println("bech32 is not identical, will set config ")
		return false
	}
	return true
}

// SetAppConfig creates a new config instance for the SDK configuration.
// 안써도 지우지 않는 이유는 다른 네트워크에서 이 함수를 사용하기 때문임
func SetAppConfig() {
	if !IsSetAppConfig() {
		config := sdktypes.GetConfig()
		SetBech32AddressPrefixes(config)
		SetBip44CoinType(config)
		config.Seal()
	}
}

// SetBech32AddressPrefixes sets the global prefix to be used when serializing addresses to bech32 strings.
func SetBech32AddressPrefixes(config *sdktypes.Config) {
	accountPubKeyPrefix := chainapp.AccountAddressPrefix + "pub"
	validatorAddressPrefix := chainapp.AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := chainapp.AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := chainapp.AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := chainapp.AccountAddressPrefix + "valconspub"

	config.SetBech32PrefixForAccount(chainapp.AccountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	config.SetAddressVerifier(func(bytes []byte) error {
		if len(bytes) == 0 {
			return sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "addresses cannot be empty")
		}
		if len(bytes) > address.MaxAddrLen {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "address max length is %d, got %d", address.MaxAddrLen, len(bytes))
		}
		// TODO: Do we want to allow addresses of lengths other than 20 and 32 bytes?
		if len(bytes) != 20 && len(bytes) != 32 {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "address length must be 20 or 32 bytes, got %d", len(bytes))
		}
		return nil
	})

}

// SetBip44CoinType sets the global coin type to be used in hierarchical deterministic wallets.
func SetBip44CoinType(config *sdktypes.Config) {
	// config.SetCoinType(app.Bip44CoinType)
}
