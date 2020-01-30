package client

import (
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// Validators returns all the known Tendermint validators for a given block
// height. An error is returned if the query fails.
func (c Client) Validators(height int64) (*tmctypes.ResultValidators, error) {
	return c.rpcClient.Validators(&height)
}
