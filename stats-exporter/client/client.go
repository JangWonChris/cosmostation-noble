package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdkCodec "github.com/cosmos/cosmos-sdk/codec"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/codec"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/models"

	// rpc "github.com/tendermint/tendermint/rpc/client/http"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	resty "github.com/go-resty/resty/v2"
)

// Client implements a wrapper around both Tendermint RPC HTTP client and
// Cosmos SDK REST client that allow for essential data queries.
type Client struct {
	cliCtx context.CLIContext
	cdc    *sdkCodec.Codec
	// rpcClient       *rpc.HTTP
	rpcClient       rpcclient.Client
	apiClient       *resty.Client
	coinGeckoClient *resty.Client
}

// NewClient creates a new client with the given configuration and
// return Client struct. An error is returned if it fails.
func NewClient(nodeCfg config.Node, marketCfg config.Market) (*Client, error) {
	cliCtx := context.NewCLIContext().
		WithCodec(codec.Codec).
		WithNodeURI(nodeCfg.RPCNode).
		WithTrustNode(true)

	rpcClient := rpcclient.NewHTTP(nodeCfg.RPCNode, "/websocket")
	// rpcClient, err := rpc.NewWithTimeout(nodeCfg.RPCNode, "/websocket", 5)
	// if err != nil {
	// 	return &Client{}, err
	// }

	apiClient := resty.New().
		SetHostURL(nodeCfg.LCDEndpoint).
		SetTimeout(time.Duration(5 * time.Second))

	coinGeckoClient := resty.New().
		SetHostURL(marketCfg.CoinGeckoEndpoint).
		SetTimeout(time.Duration(5 * time.Second))

	return &Client{cliCtx, codec.Codec, rpcClient, apiClient, coinGeckoClient}, nil
}

//-----------------------------------------------------------------------------
// RPC APIs

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
	route := fmt.Sprintf("custom/%s/%s", stakingtypes.StoreKey, stakingtypes.QueryParameters)
	bz, _, err := c.cliCtx.QueryWithData(route, nil)
	if err != nil {
		return "", err
	}

	var params stakingtypes.Params
	c.cdc.MustUnmarshalJSON(bz, &params)

	return params.BondDenom, nil
}

//-----------------------------------------------------------------------------
// REST SERVER APIs

// RequestAPIFromLCDWithRespHeight is general request API from REST Server and
// return without any modification.
func (c *Client) RequestAPIFromLCDWithRespHeight(reqParam string) (models.ResponseWithHeight, error) {
	resp, err := c.apiClient.R().Get(reqParam)
	if err != nil {
		return models.ResponseWithHeight{}, err
	}

	return models.ReadRespWithHeight(resp), nil
}

// GetInflation returns current minting inflation value
func (c *Client) GetInflation() (models.Inflation, error) {
	resp, err := c.apiClient.R().Get("/minting/inflation")
	if err != nil {
		return models.Inflation{}, err
	}

	var inflation models.Inflation
	err = json.Unmarshal(resp.Body(), &inflation)
	if err != nil {
		return models.Inflation{}, err
	}

	return inflation, nil
}

// GetValidators returns validators detail information in Tendemrint validators in active chain
// An error is returns if the query fails.
func (c *Client) GetValidators() ([]*models.Validator, error) {
	resp, err := c.apiClient.R().Get("/staking/validators")
	if err != nil {
		return nil, err
	}

	var vals []*models.Validator
	err = json.Unmarshal(resp.Body(), &vals)
	if err != nil {
		return nil, err
	}

	return vals, nil
}

//-----------------------------------------------------------------------------
// CoinGecko APIs

// CoinMarketData returns coin market data using CoinGecko API based upon coin id.
func (c *Client) CoinMarketData(id string) (models.CoinGeckoMarket, error) {
	queryStr := "/coins/" + id + "?localization=false&tickers=false&community_data=false&developer_data=false&sparkline=false"
	resp, err := c.coinGeckoClient.R().Get(queryStr)
	if err != nil {
		return models.CoinGeckoMarket{}, err
	}

	if resp.IsError() {
		return models.CoinGeckoMarket{}, fmt.Errorf("failed to respond coingecko server: %s", err)
	}

	var data models.CoinGeckoMarket
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return models.CoinGeckoMarket{}, err
	}

	return data, nil
}
