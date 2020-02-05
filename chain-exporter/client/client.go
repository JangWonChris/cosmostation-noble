package client

import (
	"log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/pkg/errors"

	ceCodec "github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/db"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// Client implements a wrapper around both a Tendermint RPC client and a
// Cosmos SDK REST client that allows for essential data queries.
type Client struct {
	rpcClient  rpcclient.Client // Tendermint RPC node
	clientNode string           // Full node
	cdc        *codec.Codec
	db         *db.Database
}

func NewClient(rpcNode, clientNode string) (Client, error) {
	rpcClient := rpcclient.NewHTTP(rpcNode, "/websocket")

	if err := rpcClient.Start(); err != nil {
		return Client{}, err
	}

	cfg := config.ParseConfig()

	// Connect to database
	db := db.Connect(&cfg.DB)

	// Ping database to verify connection is succeeded
	err := db.Ping()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to ping database."))
	}

	return Client{rpcClient, clientNode, ceCodec.Codec, db}, nil
}
