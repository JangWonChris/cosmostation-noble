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

	"github.com/cosmostation/cosmostation-cosmos/mintscan/codec"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"

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
		// WithCodec(codec.Codec).
		WithCodec(authtypes.ModuleCdc).
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

// GetCliContext returns client CLIContext.
func (c *Client) GetCliContext() context.CLIContext {
	return c.cliCtx
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

// GetTx queries for a single transaction by a hash string in hex format.
// An error is returned if the transaction does not exist or cannot be queried.
func (c *Client) GetTx(hash string) (sdkTypes.TxResponse, error) {
	txResponse, err := sdkUtils.QueryTx(c.cliCtx, hash) // use RPC under the hood
	if err != nil {
		return sdkTypes.TxResponse{}, fmt.Errorf("failed to query tx hash: %s", err)
	}

	if txResponse.Empty() {
		return sdkTypes.TxResponse{}, fmt.Errorf("tx hash has empty tx response: %s", err)
	}

	return txResponse, nil
}

// GetTendermintTx queries for a transaction by hash.
// An error is returned if the query fails.
func (c *Client) GetTendermintTx(hash string) (*tmcTypes.ResultTx, error) {
	hashRaw, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	return c.rpcClient.Tx(hashRaw, false)
}

// GetAccount checks account type and returns account interface.
func (c *Client) GetAccount(address string) (exported.Account, error) {
	accAddr, err := sdkTypes.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}

	acc, err := auth.NewAccountRetriever(c.cliCtx).GetAccount(accAddr)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

// GetValidatorCommission queries validator's commission and returns the coins with truncated decimals and the change.
func (c *Client) GetValidatorCommission(address string) (sdkTypes.Coins, error) {
	valAddr, err := sdkTypes.ValAddressFromBech32(address)
	if err != nil {
		return sdkTypes.Coins{}, err
	}

	res, err := common.QueryValidatorCommission(c.cliCtx, distr.QuerierRoute, valAddr)
	if err != nil {
		return sdkTypes.Coins{}, err
	}

	var valCom distr.ValidatorAccumulatedCommission
	c.cliCtx.Codec.MustUnmarshalJSON(res, &valCom)

	truncatedCoins, _ := valCom.TruncateDecimal()

	return truncatedCoins, nil
}

// BroadcastTx broadcasts transaction to the active network
func (c *Client) BroadcastTx(signedTx string) (*tmcTypes.ResultBroadcastTxCommit, error) {
	txByteStr, err := hex.DecodeString(signedTx)
	if err != nil {
		return &tmcTypes.ResultBroadcastTxCommit{}, err
	}

	var stdTx auth.StdTx
	err = c.cdc.UnmarshalJSON(txByteStr, &stdTx)
	if err != nil {
		return &tmcTypes.ResultBroadcastTxCommit{}, err
	}

	bz, err := c.cdc.MarshalBinaryLengthPrefixed(stdTx)
	if err != nil {
		return &tmcTypes.ResultBroadcastTxCommit{}, err
	}

	return c.rpcClient.BroadcastTxCommit(bz)
}

// --------------------
// REST SERVER APIs
// --------------------

// GetTxAPIClient queries for a transaction from the REST client and decodes it into a sdkTypes.Tx [Another way to query a transaction.]
// if the transaction exists. An error is returned if the tx doesn't exist or
// decoding fails.
func (c *Client) GetTxAPIClient(hash string) (txResponse sdkTypes.TxResponse, err error) {
	resp, err := c.apiClient.R().Get("/txs/" + hash)
	if err != nil {
		return sdkTypes.TxResponse{}, fmt.Errorf("failed to request tx hash: %s", err)
	}

	if err := c.cdc.UnmarshalJSON(resp.Body(), &txResponse); err != nil {
		return sdkTypes.TxResponse{}, fmt.Errorf("failed to unmarshal tx hash: %s", err)
	}

	return txResponse, nil
}

// GetTxs returns result of query the REST Server.
func (c *Client) GetTxs(hash string, tx *sdkTypes.TxResponse) (err error) {
	resp, err := c.apiClient.R().Get("/txs/" + hash)
	if err != nil {
		return err
	}

	if err := c.cdc.UnmarshalJSON(resp.Body(), tx); err != nil {
		return err
	}

	return nil
}

// HandleResponseHeight is general request API from REST Server and
// return without any modification
func (c *Client) HandleResponseHeight(reqParam string) (model.ResponseWithHeight, error) {
	resp, err := c.apiClient.R().Get(reqParam)
	if err != nil {
		return model.ResponseWithHeight{}, err
	}

	return model.ReadRespWithHeight(resp), nil
}

// GetCoinGeckoMarketData returns current market data from CoinGecko API based upon params
func (c *Client) GetCoinGeckoMarketData(id string) (model.CoinGeckoMarketData, error) {
	queryStr := "/coins/" + id + "?localization=false&tickers=false&community_data=false&developer_data=false&sparkline=false"

	resp, err := c.coinGeckoClient.R().Get(queryStr)
	if err != nil {
		return model.CoinGeckoMarketData{}, err
	}

	if resp.IsError() {
		return model.CoinGeckoMarketData{}, fmt.Errorf("failed to respond: %s", err)
	}

	var data model.CoinGeckoMarketData
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return model.CoinGeckoMarketData{}, err
	}

	return data, nil
}

// GetCoinGeckoCoinPrice returns simple coin price
func (c *Client) GetCoinGeckoCoinPrice(id string) (json.RawMessage, error) {
	queryStr := "/simple/price?ids=" + id + "&vs_currencies=usd&include_market_cap=false&include_last_updated_at=true"

	resp, err := c.coinGeckoClient.R().Get(queryStr)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to respond: %s", err)
	}

	var rawData json.RawMessage
	err = json.Unmarshal(resp.Body(), &rawData)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}
