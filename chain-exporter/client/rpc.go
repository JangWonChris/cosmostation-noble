package client

// import "context"

// import (
// 	"context"
// 	"encoding/hex"
// 	"encoding/json"

// 	sdkcodec "github.com/cosmos/cosmos-sdk/codec"
// 	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
// 	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
// )

// // --------------------
// // RPC APIs
// // --------------------

// // GetNetworkChainID returns network chain id.
// func (c *Client) GetNetworkChainID() (string, error) {
// 	status, err := c.rpcClient.Status(context.Background())
// 	if err != nil {
// 		return "", err
// 	}

// 	return status.NodeInfo.Network, nil
// }

// // GetStatus queries for status on the active chain.
// func (c *Client) GetStatus() (*tmctypes.ResultStatus, error) {
// 	return c.rpcClient.Status(context.Background())
// }

// // GetBlock queries for a block with height.
// func (c *Client) GetBlock(height int64) (*tmctypes.ResultBlock, error) {
// 	return c.rpcClient.Block(context.Background(), &height)
// }

// // GetLatestBlockHeight returns the latest block height on the active network.
// func (c *Client) GetLatestBlockHeight() (int64, error) {
// 	status, err := c.rpcClient.Status(context.Background())
// 	if err != nil {
// 		return -1, err
// 	}

// 	return status.SyncInfo.LatestBlockHeight, nil
// }

// // GetValidators returns all the known Tendermint validators for a given block
// // height. An error is returned if the query fails.
// func (c *Client) GetValidators(height int64, page int, perPage int) (*tmctypes.ResultValidators, error) {
// 	return c.rpcClient.Validators(context.Background(), &height, &page, &perPage)
// }

// // GetGenesisAccounts extracts all genesis accounts from genesis file and return them.
// func (c *Client) GetGenesisAccounts() (authtypes.GenesisAccounts, error) {
// 	gen, err := c.rpcClient.Genesis(context.Background())
// 	if err != nil {
// 		return authtypes.GenesisAccounts{}, err
// 	}

// 	// jeonghwan : LegacyAmino 로 풀어도 되는지 확인이 필요
// 	// GetGenesisStateFromAppState() 함수에서 JSONMarshaler 타입만 받음
// 	appState := make(map[string]json.RawMessage)
// 	err = c.cliCtx.LegacyAmino.UnmarshalJSON(gen.Genesis.AppState, &appState)
// 	if err != nil {
// 		return authtypes.GenesisAccounts{}, err
// 	}

// 	genesisState := authtypes.GetGenesisStateFromAppState(c.cliCtx.JSONMarshaler.(sdkcodec.Marshaler), appState)
// 	accs, err := authtypes.UnpackAccounts(genesisState.Accounts)
// 	if err != nil {
// 		return nil, err
// 	}
// 	genesisAccts := authtypes.SanitizeGenesisAccounts(accs)

// 	return genesisAccts, nil
// }

// // GetTendermintTx queries for a transaction by hash.
// func (c *Client) GetTendermintTx(hash string) (*tmctypes.ResultTx, error) {
// 	hashRaw, err := hex.DecodeString(hash)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return c.rpcClient.Tx(context.Background(), hashRaw, false)
// }

// // GetTendermintTxSearch queries for a transaction search by condition.
// // TODO: need more tests. ex:) query := "tx.height=75960",prove := true, page := 1, perPage := 30, orderBy := "asc"
// // If this is not needed for this project, let's just remove.
// func (c *Client) GetTendermintTxSearch(query string, prove bool, page, perPage int, orderBy string) (*tmctypes.ResultTxSearch, error) {
// 	txResp, err := c.rpcClient.TxSearch(context.Background(), query, prove, &page, &perPage, orderBy)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return txResp, nil
// }
