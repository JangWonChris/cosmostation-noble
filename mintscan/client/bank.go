package client

// import (
// 	"context"

// 	sdktypes "github.com/cosmos/cosmos-sdk/types"
// 	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
// )

// // GetBankQueryClient returns a object of queryClient
// func (c *Client) GetBankQueryClient() banktypes.QueryClient {
// 	return banktypes.NewQueryClient(c.GRPC)
// }

// // GetBalance returns balance of a given account for staking denom
// func (c *Client) GetBalance(address string) (*sdktypes.Coin, error) {
// 	denom, err := c.GetBondDenom()
// 	if err != nil {
// 		return nil, err
// 	}
// 	bankClient := c.GetBankQueryClient()
// 	request := banktypes.QueryBalanceRequest{Address: address, Denom: denom}
// 	bankRes, err := bankClient.Balance(context.Background(), &request)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return bankRes.GetBalance(), nil
// }

// // GetAllBalances returns balances of a given account
// func (c *Client) GetAllBalances(address string) (sdktypes.Coins, error) {
// 	bankClient := c.GetBankQueryClient()
// 	request := banktypes.QueryAllBalancesRequest{Address: address}
// 	bankRes, err := bankClient.AllBalances(context.Background(), &request)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return bankRes.GetBalances(), nil
// }
