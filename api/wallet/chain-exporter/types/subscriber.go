package types

import (
	"context"
	"encoding/json"
)

type Subscriber interface {
	SetCurrentHeight(height int64) error
	GetNextHeight(ctx context.Context) (int64, error)
	Commit(ctx context.Context, txResult *ElasticsearchTxInfo) error
	Stop() error
	MarsharJSON(o interface{}) (json.RawMessage, error)
	UnmarshalJSON(json []byte, o interface{}) error
	UnmarshalBinaryLengthPrefixed(tx []byte, o interface{}) error
}

type TestSubscriber struct {
	Height     int64
	latestTxid string
}

var _ Subscriber = &TestSubscriber{}

func (sub *TestSubscriber) SetCurrentHeight(height int64) error {
	return nil
}

func (base *TestSubscriber) GetNextHeight(ctx context.Context) (int64, error) {
	return base.Height + 1, nil
}

func (base *TestSubscriber) Commit(ctx context.Context, txResult *ElasticsearchTxInfo) error {
	return nil
}

func (base *TestSubscriber) Stop() error {
	return nil
}

func (sub *TestSubscriber) MarsharJSON(o interface{}) (json.RawMessage, error) {
	return json.RawMessage{}, nil
}

func (sub *TestSubscriber) UnmarshalJSON(json []byte, o interface{}) error {
	return nil
}

func (sub *TestSubscriber) UnmarshalBinaryLengthPrefixed(tx []byte, o interface{}) error {
	return nil
}
