package codec

import (
	_ "github.com/cosmos/cosmos-sdk/codec"
	_ "github.com/cosmos/cosmos-sdk/simapp"
)

// Codec is the application-wide Amino codec for serializing interfaces and data
// var Codec *codec.Codec

/*
// initializes upon package loading
func init() {
	Codec = simapp.MakeCodec()
	Codec.Seal()
}
*/
