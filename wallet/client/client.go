package client

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdkCodec "github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdkUtils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/client/common"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmostation/cosmostation-cosmos/wallet/codec"
	"github.com/cosmostation/cosmostation-cosmos/wallet/config"
	"github.com/cosmostation/cosmostation-cosmos/wallet/model"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tmcTypes "github.com/tendermint/tendermint/rpc/core/types"

	resty "github.com/go-resty/resty/v2"
)

// Client implements a wrapper around both a Tendermint RPC client and a
// Cosmos SDK REST client that allows for essential data queries.
type Client struct {
	cliCtx          context.CLIContext
	cdc             *sdkCodec.Codec
	rpcClient       rpcclient.Client
	apiClient       *resty.Client
	coinGeckoClient *resty.Client
}

// NewClient creates a new client with the given config.
func NewClient(nodeCfg config.NodeConfig, marketCfg config.MarketConfig) (*Client, error) {
	cliCtx := context.NewCLIContext().
		WithCodec(codec.Codec).
		WithNodeURI(nodeCfg.RPCNode).
		WithTrustNode(true)

	rpcClient := rpcclient.NewHTTP(nodeCfg.RPCNode, "/websocket")

	apiClient := resty.New().
		SetHostURL(nodeCfg.LCDEndpoint).
		SetTimeout(time.Duration(5 * time.Second))

	coinGeckoClient := resty.New().
		SetHostURL(marketCfg.CoinGeckoEndpoint).
		SetTimeout(time.Duration(5 * time.Second))

	return &Client{cliCtx, codec.Codec, rpcClient, apiClient, coinGeckoClient}, nil
}

// --------------------
// RPC APIs
// --------------------

// GetNetworkChainID returns network chain id.
func (c *Client) GetNetworkChainID() (string, error) {
	status, err := c.rpcClient.Status()
	if err != nil {
		return "", err
	}

	return status.NodeInfo.Network, nil
}

// GetBondDenom returns bond denomination for the network.
func (c *Client) GetBondDenom() (string, error) {
	route := fmt.Sprintf("custom/%s/%s", stakingTypes.StoreKey, stakingTypes.QueryParameters)
	bz, _, err := c.cliCtx.QueryWithData(route, nil)
	if err != nil {
		return "", err
	}

	var params stakingTypes.Params
	c.cdc.MustUnmarshalJSON(bz, &params)

	return params.BondDenom, nil
}

// GetStatus queries for status on the active chain.
func (c *Client) GetStatus() (*tmcTypes.ResultStatus, error) {
	return c.rpcClient.Status()
}

// GetBlock queries for a block with height.
func (c *Client) GetBlock(height int64) (*tmcTypes.ResultBlock, error) {
	return c.rpcClient.Block(&height)
}

// GetLatestBlockHeight returns the latest block height on the active network.
func (c *Client) GetLatestBlockHeight() (int64, error) {
	status, err := c.rpcClient.Status()
	if err != nil {
		return -1, err
	}

	return status.SyncInfo.LatestBlockHeight, nil
}
