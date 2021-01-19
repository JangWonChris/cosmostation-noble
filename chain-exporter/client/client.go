package client

import (

	// cosmos-sdk

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	//product
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	lclient "github.com/cosmostation/mintscan-backend-library/client"
	"github.com/cosmostation/mintscan-backend-library/config"
	// tendermint
	//etc
)

// Client implements a wrapper around both Tendermint RPC HTTP client and
// Cosmos SDK REST client that allow for essential data queries.
type Client struct {
	*lclient.Client
}

// NewClient creates a new client with the given configuration and
// return Client struct. An error is returned if it fails.
func NewClient(cfg *config.ClientConfig) (*Client, error) {
	client := lclient.NewClient(cfg)

	client.CliCtx.Context = client.CliCtx.Context.
		WithJSONMarshaler(codec.EncodingConfig.Marshaler).
		WithLegacyAmino(codec.EncodingConfig.Amino).
		WithTxConfig(codec.EncodingConfig.TxConfig).
		WithInterfaceRegistry(codec.EncodingConfig.InterfaceRegistry).
		WithAccountRetriever(authtypes.AccountRetriever{})

	return &Client{client}, nil
}
