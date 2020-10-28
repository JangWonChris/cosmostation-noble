package codec

import (
	"github.com/cosmos/cosmos-sdk/codec"
	gaia "github.com/cosmos/gaia/app"
	"github.com/cosmos/gaia/app/params"
)

// Codec is the application-wide Amino codec and is initialized upon package loading.
var (
	AppCodec       codec.Marshaler
	AminoCodec     *codec.LegacyAmino
	EncodingConfig params.EncodingConfig
)

func init() {
	EncodingConfig = gaia.MakeEncodingConfig()
	AppCodec = EncodingConfig.Marshaler
	AminoCodec = EncodingConfig.Amino
}
