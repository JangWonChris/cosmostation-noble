package client

import (
	"time"

	// cosmos-sdk
	"github.com/cosmos/cosmos-sdk/client"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	//product
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"

	// tendermint

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpc "github.com/tendermint/tendermint/rpc/client/http"

	//etc
	resty "github.com/go-resty/resty/v2"
	"google.golang.org/grpc"
)

// Client implements a wrapper around both Tendermint RPC HTTP client and
// Cosmos SDK REST client that allow for essential data queries.
type Client struct {
	cliCtx        client.Context
	grpcClient    *grpc.ClientConn
	rpcClient     rpcclient.Client
	apiClient     *resty.Client
	keyBaseClient *resty.Client
}

// NewClient creates a new client with the given configuration and
// return Client struct. An error is returned if it fails.
func NewClient(nodeCfg config.Node, keyBaseURL string) (*Client, error) {
	grpcClient, err := grpc.Dial(nodeCfg.GRPCEndpoint,
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second*10),
		grpc.WithInsecure())
	if err != nil {
		return &Client{}, err
	}

	rpcClient, err := rpc.NewWithTimeout(nodeCfg.RPCNode, "/websocket", 10)
	if err != nil {
		return &Client{}, err
	}

	apiClient := resty.New().
		SetHostURL(nodeCfg.LCDEndpoint).
		SetTimeout(time.Duration(10 * time.Second))

	keyBaseClient := resty.New().
		SetHostURL(keyBaseURL).
		SetTimeout(time.Duration(5 * time.Second))

	cliCtx := client.Context{}.
		WithNodeURI(nodeCfg.RPCNode).
		WithJSONMarshaler(codec.EncodingConfig.Marshaler).
		WithLegacyAmino(codec.EncodingConfig.Amino).
		WithTxConfig(codec.EncodingConfig.TxConfig).
		WithInterfaceRegistry(codec.EncodingConfig.InterfaceRegistry).
		WithClient(rpcClient). //tendermint-rc6 에서 추가하지 않으면, 기본 offline mode로 지정이 된다.
		WithAccountRetriever(authtypes.AccountRetriever{})

	return &Client{cliCtx, grpcClient, rpcClient, apiClient, keyBaseClient}, nil
}

// GetGRPCConn returns client object
func (c *Client) GetGRPCConn() *grpc.ClientConn {
	return c.grpcClient
}

// Close close the connection
func (c *Client) Close() {
	c.grpcClient.Close()
}
