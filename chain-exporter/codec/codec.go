package codec

import (
	"github.com/cosmos/cosmos-sdk/codec"
	gaia "github.com/cosmos/gaia/app"
)

// Codec is the application-wide Amino codec and is initialized upon package loading.
var (
	AppCodec   codec.Marshaler
	AminoCodec *codec.LegacyAmino
)

func init() {
	encodingConfig := gaia.MakeEncodingConfig()
	AppCodec = encodingConfig.Marshaler
	AminoCodec = encodingConfig.Amino
}
