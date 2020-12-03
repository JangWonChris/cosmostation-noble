package client

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

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

// GetAuthQueryClient returns a object of queryClient
func (c *Client) GetAuthQueryClient() authtypes.QueryClient {
	return authtypes.NewQueryClient(c.grpcClient)
}

// GetAccount checks account type and returns account interface.
func (c *Client) GetAccount(address string) (client.Account, error) {
	accAddr, err := sdktypes.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}

	ar := authtypes.AccountRetriever{}
	acc, _, err := ar.GetAccountWithHeight(c.cliCtx, accAddr)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
