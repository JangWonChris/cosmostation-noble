package client

import (
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/codec"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/config"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpcclienthttp "github.com/tendermint/tendermint/rpc/client/http"

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
	grpcClient, err := grpc.Dial(nodeCfg.GRPCEndpoint,
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second*10),
		grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	rpcClient, err := rpcclienthttp.NewWithTimeout(nodeCfg.RPCNode, "/websocket", 10)
	if err != nil {
		return nil, err
	}

	apiClient := resty.New().
		SetHostURL(nodeCfg.LCDEndpoint).
		SetTimeout(time.Duration(10 * time.Second))

	coinGeckoClient := resty.New().
		SetHostURL(marketCfg.CoinGeckoEndpoint).
		SetTimeout(time.Duration(10 * time.Second))

	cliCtx := client.Context{}.
		WithNodeURI(nodeCfg.RPCNode).
		WithJSONMarshaler(codec.EncodingConfig.Marshaler).
		WithLegacyAmino(codec.EncodingConfig.Amino).
		WithTxConfig(codec.EncodingConfig.TxConfig).
		WithInterfaceRegistry(codec.EncodingConfig.InterfaceRegistry).
		WithClient(rpcClient). //tendermint-rc6 에서 추가하지 않으면, 기본 offline mode로 지정이 된다.
		WithAccountRetriever(authtypes.AccountRetriever{})

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
