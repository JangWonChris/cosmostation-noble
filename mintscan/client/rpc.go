package client

import (
	"context"
	"encoding/hex"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// GetStatus queries for status on the active chain.
func (c *Client) GetStatus() (*tmctypes.ResultStatus, error) {
	return c.rpcClient.Status(context.Background())
}

// GetNetworkChainID returns network chain id.
func (c *Client) GetNetworkChainID() (string, error) {
	status, err := c.GetStatus()
	if err != nil {
		return "", err
	}

	return status.NodeInfo.Network, nil
}

// GetLatestBlockHeight returns the latest block height on the active network.
func (c *Client) GetLatestBlockHeight() (int64, error) {
	status, err := c.GetStatus()
	if err != nil {
		return -1, err
	}

	return status.SyncInfo.LatestBlockHeight, nil
}

// GetBlock queries for a block with height.
func (c *Client) GetBlock(height int64) (*tmctypes.ResultBlock, error) {
	return c.rpcClient.Block(context.Background(), &height)
}

// GetTendermintTx queries for a transaction by hash.
// An error is returned if the query fails.
func (c *Client) GetTendermintTx(hash string) (*tmctypes.ResultTx, error) {
	hashRaw, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	return c.rpcClient.Tx(context.Background(), hashRaw, false)
}

// BroadcastTx broadcasts transaction to the active network
func (c *Client) BroadcastTx(signedTx string) (*tmctypes.ResultBroadcastTxCommit, error) {
	txBytes, err := hex.DecodeString(signedTx)
	if err != nil {
		return &tmctypes.ResultBroadcastTxCommit{}, err
	}

	//todo : jeonghwan

	// var stdTx sdktypes.Tx
	// stdTx, err = c.cliCtx.TxConfig.TxJSONDecoder()(signedTxStr)
	// if err != nil {
	// 	return &tmctypes.ResultBroadcastTxCommit{}, err
	// }

	// txBytes, err := c.cliCtx.TxConfig.TxJSONEncoder()(stdTx)
	// if err != nil {
	// 	return &tmctypes.ResultBroadcastTxCommit{}, err
	// }

	// BroadcastBlock mode will wait tx
	// c.cliCtx.WithBroadcastMode(clientflags.BroadcastBlock).BroadcastTx(bz)
	return c.rpcClient.BroadcastTxCommit(context.Background(), txBytes)
}
