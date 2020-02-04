package exporter

import (
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
)

// getBlock provides block information
func (ex *Exporter) getBlock(height int64) ([]*schema.Block, error) {
	Block := make([]*schema.Block, 0)

	// query current block
	block, err := ex.client.Block(height)
	if err != nil {
		return nil, err
	}

	tempBlock := &schema.Block{
		BlockHash: block.BlockMeta.BlockID.Hash.String(),
		Proposer:  block.Block.ProposerAddress.String(),
		Height:    block.Block.Height,
		TotalTxs:  block.Block.TotalTxs,
		NumTxs:    block.Block.NumTxs,
		Time:      block.BlockMeta.Header.Time,
	}

	Block = append(Block, tempBlock)

	return Block, nil
}
