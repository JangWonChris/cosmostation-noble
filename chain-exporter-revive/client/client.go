package client

import (
	"github.com/cosmos/cosmos-sdk/codec"

	ceCodec "github.com/cosmostation/cosmostation-cosmos/chain-exporter-revive/codec"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// Client implements a wrapper around both a Tendermint RPC client and a
// Cosmos SDK REST client that allows for essential data queries.
type Client struct {
	rpcClient  rpcclient.Client // Tendermint RPC node
	clientNode string           // Full node
	cdc        *codec.Codec
}

func NewClient(rpcNode, clientNode string) (Client, error) {
	rpcClient := rpcclient.NewHTTP(rpcNode, "/websocket")

	if err := rpcClient.Start(); err != nil {
		return Client{}, err
	}

	return Client{rpcClient, clientNode, ceCodec.Codec}, nil
}

// Stop defers the node stop execution to the RPC client.
func (c Client) Stop() error {
	return c.rpcClient.Stop()
}
