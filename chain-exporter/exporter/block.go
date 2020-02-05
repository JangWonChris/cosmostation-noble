package exporter

import (
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
)

// getBlock provides block information
func (ex *Exporter) getBlock(height int64) ([]*schema.BlockInfo, error) {
	blockInfo := make([]*schema.BlockInfo, 0)

	// query current block
	block, err := ex.client.Block(height)
	if err != nil {
		return nil, err
	}

	tempBlock := &schema.BlockInfo{
		BlockHash: block.BlockMeta.BlockID.Hash.String(),
		Proposer:  block.Block.ProposerAddress.String(),
		Height:    block.Block.Height,
		TotalTxs:  block.Block.TotalTxs,
		NumTxs:    block.Block.NumTxs,
		Time:      block.BlockMeta.Header.Time,
	}

	blockInfo = append(blockInfo, tempBlock)

	return blockInfo, nil
}
