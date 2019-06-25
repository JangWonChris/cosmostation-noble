package subscribe

import (
	"context"
	"encoding/json"

	gaiaApp "github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/libs/log"
)

type Subscriber struct {
	Logger    log.Logger
	Codec     *codec.Codec
	TxDecoder sdk.TxDecoder
	Height    int64
}

func NewSubscriber(logger log.Logger) *Subscriber {
	// Register Cosmos SDK codecs
	sdkCodec := gaiaApp.MakeCodec()

	subscriber := &Subscriber{
		Logger:    logger,
		Codec:     sdkCodec,
		TxDecoder: auth.DefaultTxDecoder(sdkCodec),
		Height:    0,
	}
	return subscriber
}

func (sub *Subscriber) SetCurrentHeight(height int64) error {
	sub.Height = height
	return nil
}

func (sub *Subscriber) GetNextHeight(ctx context.Context) (int64, error) {
	// height index에서 가져오는게 아니라
	// txIndex에서 가져오고, 지역변수를 계속 +1 하는것으로 바꾸는게 낫겠다.
	return sub.Height + 1, nil
}

func (sub *Subscriber) MarsharJSON(o interface{}) (json.RawMessage, error) {
	return sub.Codec.MarshalJSON(o)
}

func (sub *Subscriber) UnmarshalJSON(json []byte, o interface{}) error {
	return sub.Codec.UnmarshalJSON(json, o)
}

func (sub *Subscriber) UnmarshalBinaryLengthPrefixed(tx []byte, o interface{}) error {
	return sub.Codec.UnmarshalBinaryLengthPrefixed(tx, o)
}

func (sub *Subscriber) Stop() error {
	return nil
}
