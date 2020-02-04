package client

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ceCodec "github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
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

// LatestHeight returns the latest block height on the active chain. An error
// is returned if the query fails.
func (c Client) LatestHeight() (int64, error) {
	status, err := c.rpcClient.Status()
	if err != nil {
		return -1, err
	}

	height := status.SyncInfo.LatestBlockHeight
	return height, nil
}

// Block queries for a block by height. An error is returned if the query fails.
func (c Client) Block(height int64) (*tmctypes.ResultBlock, error) {
	return c.rpcClient.Block(&height)
}

// Status queries for status on the active chain
func (c Client) Status() (*tmctypes.ResultStatus, error) {
	return c.rpcClient.Status()
}

// TendermintTx queries for a transaction by hash. An error is returned if the
// query fails.
func (c Client) TendermintTx(hash string) (*tmctypes.ResultTx, error) {
	hashRaw, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	return c.rpcClient.Tx(hashRaw, false)
}

// Validators returns all the known Tendermint validators for a given block
// height. An error is returned if the query fails.
func (c Client) Validators(height int64) (*tmctypes.ResultValidators, error) {
	return c.rpcClient.Validators(&height)
}

// Stop defers the node stop execution to the RPC client.
func (c Client) Stop() error {
	return c.rpcClient.Stop()
}

// SubscribeNewBlocks subscribes to the new block event handler through the RPC
// client with the given subscriber name. An receiving only channel, context
// cancel function and an error is returned. It is up to the caller to cancel
// the context and handle any errors appropriately.
func (c Client) SubscribeNewBlocks(subscriber string) (<-chan tmctypes.ResultEvent, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	eventCh, err := c.rpcClient.Subscribe(ctx, subscriber, "tm.event = 'NewBlock'")
	return eventCh, cancel, err
}

// Tx queries for a transaction from the REST client and decodes it into a sdk.Tx
// if the transaction exists. An error is returned if the tx doesn't exist or
// decoding fails.
func (c Client) Tx(hash string) (sdk.TxResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/txs/%s", c.clientNode, hash))
	if err != nil {
		return sdk.TxResponse{}, err
	}

	defer resp.Body.Close()

	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	var tx sdk.TxResponse

	if err := c.cdc.UnmarshalJSON(bz, &tx); err != nil {
		return sdk.TxResponse{}, err
	}

	return tx, nil
}

// Txs queries for all the transactions in a block. Transactions are returned
// in the sdk.TxResponse format which internally contains an sdk.Tx. An error is
// returned if any query fails.
func (c Client) Txs(block *tmctypes.ResultBlock) ([]sdk.TxResponse, error) {
	txResponses := make([]sdk.TxResponse, len(block.Block.Txs), len(block.Block.Txs))

	for i, tmTx := range block.Block.Txs {
		txResponse, err := c.Tx(fmt.Sprintf("%X", tmTx.Hash()))
		if err != nil {
			return nil, err
		}

		txResponses[i] = txResponse
	}

	return txResponses, nil
}
