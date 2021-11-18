package custom

import (
	"github.com/cosmos/cosmos-sdk/codec"
	chainapp "github.com/cosmos/gaia/v6/app"
	chainparams "github.com/cosmos/gaia/v6/app/params"
)

// Codec is the application-wide Amino codec and is initialized upon package loading.
var (
	AppCodec       codec.Codec
	AminoCodec     *codec.LegacyAmino
	EncodingConfig chainparams.EncodingConfig
)

func init() {
	EncodingConfig = chainapp.MakeEncodingConfig()
	AppCodec = EncodingConfig.Marshaler
	AminoCodec = EncodingConfig.Amino
}
