package client

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/client/common"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/codec"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpcclienthttp "github.com/tendermint/tendermint/rpc/client/http"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	resty "github.com/go-resty/resty/v2"
)

// Client implements a wrapper around both Tendermint RPC HTTP client and
// Cosmos SDK REST client that allow for essential data queries.
type Client struct {
	cliCtx          client.Context
	grpcClient      *grpc.ClientConn
	rpcClient       rpcclient.Client
	apiClient       *resty.Client
	coinGeckoClient *resty.Client
}

// NewClient creates a new client with the given configuration and
// return Client struct. An error is returned if it fails.
func NewClient(nodeCfg config.NodeConfig, marketCfg config.MarketConfig) (*Client, error) {
	cliCtx := client.Context{}.
		WithNodeURI(nodeCfg.RPCNode).
		WithJSONMarshaler(codec.EncodingConfig.Marshaler).
		WithLegacyAmino(codec.EncodingConfig.Amino).
		WithTxConfig(codec.EncodingConfig.TxConfig).
		WithInterfaceRegistry(codec.EncodingConfig.InterfaceRegistry).
		WithAccountRetriever(authtypes.AccountRetriever{})

	grpcClient, err := grpc.Dial(nodeCfg.GRPCEndpoint,
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second*10),
		grpc.WithInsecure())
	if err != nil {
		return &Client{}, err
	}

	rpcClient, err := rpcclienthttp.NewWithTimeout(nodeCfg.RPCNode, "/websocket", 10)
	if err != nil {
		return &Client{}, err
	}

	apiClient := resty.New().
		SetHostURL(nodeCfg.LCDEndpoint).
		SetTimeout(time.Duration(10 * time.Second))

	coinGeckoClient := resty.New().
		SetHostURL(marketCfg.CoinGeckoEndpoint).
		SetTimeout(time.Duration(10 * time.Second))

	return &Client{cliCtx, grpcClient, rpcClient, apiClient, coinGeckoClient}, nil
}

// Close close the connection
func (c *Client) Close() {
	c.grpcClient.Close()
}

// GetCliContext returns client CLIContext.
func (c *Client) GetCliContext() client.Context {
	return c.cliCtx
}

// --------------------
// RPC APIs
// --------------------

// GetNetworkChainID returns network chain id.
func (c *Client) GetNetworkChainID() (string, error) {
	status, err := c.rpcClient.Status(context.Background())
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
	c.cliCtx.LegacyAmino.MustUnmarshalJSON(bz, &params)

	return params.BondDenom, nil
}

// GetStatus queries for status on the active chain.
func (c *Client) GetStatus() (*tmctypes.ResultStatus, error) {
	return c.rpcClient.Status(context.Background())
}

// GetBlock queries for a block with height.
func (c *Client) GetBlock(height int64) (*tmctypes.ResultBlock, error) {
	return c.rpcClient.Block(context.Background(), &height)
}

// GetLatestBlockHeight returns the latest block height on the active network.
func (c *Client) GetLatestBlockHeight() (int64, error) {
	status, err := c.rpcClient.Status(context.Background())
	if err != nil {
		return -1, err
	}

	return status.SyncInfo.LatestBlockHeight, nil
}

// GetTx queries for a single transaction by a hash string in hex format.
// An error is returned if the transaction does not exist or cannot be queried.
func (c *Client) GetTx(hash string) (*sdktypes.TxResponse, error) {
	txResponse, err := authclient.QueryTx(c.cliCtx, hash) // use RPC under the hood
	if err != nil {
		return &sdktypes.TxResponse{}, fmt.Errorf("failed to query tx hash: %s", err)
	}

	if txResponse.Empty() {
		return &sdktypes.TxResponse{}, fmt.Errorf("tx hash has empty tx response: %s", err)
	}

	return txResponse, nil
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

// GetAccount checks account type and returns account interface.
func (c *Client) GetAccount(address string) (authtypes.AccountI, error) {
	accAddr, err := sdktypes.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}

	ar := authtypes.AccountRetriever{}
	acc, err := ar.GetAccount(c.cliCtx, accAddr)
	if err != nil {
		return nil, err
	}
	// acc, err := auth.NewAccountRetriever(c.cliCtx).GetAccount(accAddr)
	// if err != nil {
	// 	return nil, err
	// }

	return acc, nil
}

// GetDelegatorDelegations returns a list of delegations made by a certain delegator address
func (c *Client) GetDelegatorDelegations(address string) (stakingtypes.DelegationResponses, error) {
	delAddr, err := sdktypes.AccAddressFromBech32(address)
	if err != nil {
		return stakingtypes.DelegationResponses{}, err
	}

	bz, err := c.cliCtx.LegacyAmino.MarshalJSON(stakingtypes.NewQueryDelegatorParams(delAddr))
	if err != nil {
		return stakingtypes.DelegationResponses{}, err
	}

	route := fmt.Sprintf("custom/%s/%s", stakingtypes.QuerierRoute, stakingtypes.QueryDelegatorDelegations)

	res, _, err := c.cliCtx.QueryWithData(route, bz)
	if err != nil {
		return stakingtypes.DelegationResponses{}, err
	}

	var delegations stakingtypes.DelegationResponses
	if err := c.cliCtx.LegacyAmino.UnmarshalJSON(res, &delegations); err != nil {
		return stakingtypes.DelegationResponses{}, err
	}

	return delegations, nil
}

// GetDelegatorUndelegations returns a list of undelegations made by a certain delegator address
func (c *Client) GetDelegatorUndelegations(address string) (stakingtypes.UnbondingDelegations, error) {
	delAddr, err := sdktypes.AccAddressFromBech32(address)
	if err != nil {
		return stakingtypes.UnbondingDelegations{}, err
	}

	bz, err := c.cliCtx.LegacyAmino.MarshalJSON(stakingtypes.NewQueryDelegatorParams(delAddr))
	if err != nil {
		return stakingtypes.UnbondingDelegations{}, err
	}

	route := fmt.Sprintf("custom/%s/%s", stakingtypes.QuerierRoute, stakingtypes.QueryDelegatorUnbondingDelegations)

	res, _, err := c.cliCtx.QueryWithData(route, bz)
	if err != nil {
		return stakingtypes.UnbondingDelegations{}, err
	}

	var undelegations stakingtypes.UnbondingDelegations
	if err := c.cliCtx.LegacyAmino.UnmarshalJSON(res, &undelegations); err != nil {
		return stakingtypes.UnbondingDelegations{}, err
	}

	return undelegations, nil
}

// GetDelegatorTotalRewards returns the total rewards balance from all delegations by a delegator
func (c *Client) GetDelegatorTotalRewards(address string) (distributiontypes.QueryDelegatorTotalRewardsResponse, error) {
	delAddr, err := sdktypes.AccAddressFromBech32(address)
	if err != nil {
		return distributiontypes.QueryDelegatorTotalRewardsResponse{}, err
	}

	bz, err := c.cliCtx.LegacyAmino.MarshalJSON(distributiontypes.NewQueryDelegatorParams(delAddr))
	if err != nil {
		return distributiontypes.QueryDelegatorTotalRewardsResponse{}, err
	}

	route := fmt.Sprintf("custom/%s/%s", distributiontypes.QuerierRoute, distributiontypes.QueryDelegatorTotalRewards)

	res, _, err := c.cliCtx.QueryWithData(route, bz)
	if err != nil {
		return distributiontypes.QueryDelegatorTotalRewardsResponse{}, err
	}

	var totalRewards distributiontypes.QueryDelegatorTotalRewardsResponse
	if err := c.cliCtx.LegacyAmino.UnmarshalJSON(res, &totalRewards); err != nil {
		return distributiontypes.QueryDelegatorTotalRewardsResponse{}, err
	}

	return totalRewards, nil
}

// GetValidatorCommission queries validator's commission and returns the coins with truncated decimals and the change.
func (c *Client) GetValidatorCommission(address string) (sdktypes.Coins, error) {
	valAddr, err := sdktypes.ValAddressFromBech32(address)
	if err != nil {
		return sdktypes.Coins{}, err
	}

	res, err := common.QueryValidatorCommission(c.cliCtx, valAddr)
	if err != nil {
		return sdktypes.Coins{}, err
	}

	var valCom distributiontypes.ValidatorAccumulatedCommission
	c.cliCtx.LegacyAmino.MustUnmarshalJSON(res, &valCom)

	truncatedCoins, _ := valCom.Commission.TruncateDecimal()

	return truncatedCoins, nil
}

// BroadcastTx broadcasts transaction to the active network
func (c *Client) BroadcastTx(signedTx string) (*tmctypes.ResultBroadcastTxCommit, error) {
	txBytes, err := hex.DecodeString(signedTx)
	if err != nil {
		return &tmctypes.ResultBroadcastTxCommit{}, err
	}

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

// --------------------
// REST SERVER APIs
// --------------------

// GetTxAPIClient queries for a transaction from the REST client and decodes it into a sdktypes.Tx [Another way to query a transaction.]
// if the transaction exists. An error is returned if the tx doesn't exist or
// decoding fails.
func (c *Client) GetTxAPIClient(hash string) (txResponse sdktypes.TxResponse, err error) {
	resp, err := c.apiClient.R().Get("/txs/" + hash)
	if err != nil {
		return sdktypes.TxResponse{}, fmt.Errorf("failed to request tx hash: %s", err)
	}

	if err := c.cliCtx.LegacyAmino.UnmarshalJSON(resp.Body(), &txResponse); err != nil {
		return sdktypes.TxResponse{}, fmt.Errorf("failed to unmarshal tx hash: %s", err)
	}

	return txResponse, nil
}

// GetTxs returns result of query the REST Server.
func (c *Client) GetTxs(hash string, tx *sdktypes.TxResponse) (err error) {
	resp, err := c.apiClient.R().Get("/txs/" + hash)
	if err != nil {
		return err
	}

	if err := c.cliCtx.LegacyAmino.UnmarshalJSON(resp.Body(), tx); err != nil {
		return err
	}

	return nil
}

// RequestWithRestServer is general request API from REST Server and
// return without any modification
func (c *Client) RequestWithRestServer(reqParam string) ([]byte, error) {
	resp, err := c.apiClient.R().Get(reqParam)
	if err != nil {
		return nil, err
		// return model.ResponseWithHeight{}, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get respond : %s", resp.Status())
	}

	return resp.Body(), nil
	// return model.ReadRespWithHeight(resp), nil
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

// GetCoinMarketChartData returns current market chart data from CoinGecko API based upon params.
func (c *Client) GetCoinMarketChartData(id string, from string, to string) (data model.CoinGeckoMarketDataChart, err error) {
	resp, err := c.coinGeckoClient.R().Get("/coins/" + id + "/market_chart/range?id=" + id + "&vs_currency=usd&from=" + from + "&to=" + to)
	if err != nil {
		return model.CoinGeckoMarketDataChart{}, err
	}

	if resp.IsError() {
		return model.CoinGeckoMarketDataChart{}, fmt.Errorf("failed to request: %s", err)
	}

	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return model.CoinGeckoMarketDataChart{}, err
	}

	return data, nil
}

// GetAccountBalance returns balances of a given account
func (c *Client) GetAccountBalance(address string) (*sdktypes.Coin, error) {
	sdkaddr, err := sdktypes.AccAddressFromBech32(address)
	if err != nil {
		return &sdktypes.Coin{}, err
	}
	denom, err := c.GetBondDenom()
	if err != nil {
		return &sdktypes.Coin{}, err
	}
	// jeonghwan umuon을 구해올 방법, 혹은 다른 denom들
	b := banktypes.NewQueryBalanceRequest(sdkaddr, denom)
	bankClient := banktypes.NewQueryClient(c.grpcClient)
	var header metadata.MD
	bankRes, err := bankClient.Balance(
		context.Background(),
		b,
		grpc.Header(&header), // Also fetch grpc header
	)

	return bankRes.GetBalance(), nil
}
