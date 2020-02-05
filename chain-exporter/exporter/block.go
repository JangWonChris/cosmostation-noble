package exporter

import (
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// getBlock provides block information
func (ex *Exporter) getBlock(block *tmctypes.ResultBlock) ([]*schema.BlockCosmoshub3, error) {
	resultBlock := make([]*schema.BlockCosmoshub3, 0)

	tempBlock := &schema.BlockCosmoshub3{
		Height:        block.Block.Height,
		Proposer:      block.Block.ProposerAddress.String(),
		BlockHash:     block.BlockMeta.BlockID.Hash.String(),
		ParentHash:    block.BlockMeta.Header.LastBlockID.Hash.String(),
		NumPrecommits: int64(len(block.Block.LastCommit.Precommits)),
		NumTxs:        block.Block.NumTxs,
		TotalTxs:      block.Block.TotalTxs,
		Timestamp:     block.Block.Time,
	}

	resultBlock = append(resultBlock, tempBlock)

	return resultBlock, nil
}
