package exporter

import (
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
)

// getBlockInfo provides block information
func (ce *ChainExporter) getBlockInfo(height int64) ([]*schema.BlockInfo, error) {
	blockInfo := make([]*schema.BlockInfo, 0)

	// query current block
	block, err := ce.rpcClient.Block(&height)
	if err != nil {
		return nil, err
	}

	tempBlockInfo := &schema.BlockInfo{
		BlockHash: block.BlockMeta.BlockID.Hash.String(),
		Proposer:  block.Block.ProposerAddress.String(),
		Height:    block.Block.Height,
		TotalTxs:  block.Block.TotalTxs,
		NumTxs:    block.Block.NumTxs,
		Time:      block.BlockMeta.Header.Time,
	}

	blockInfo = append(blockInfo, tempBlockInfo)

	return blockInfo, nil
}
