package exporter

import (
	dtypes "github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
)

func (ces *ChainExporterService) getBlockInfo(height int64) ([]*dtypes.BlockInfo, error) {
	blockInfo := make([]*dtypes.BlockInfo, 0)

	// Query the current block
	block, err := ces.rpcClient.Block(&height)
	if err != nil {
		return nil, err
	}

	// Parse blockinfo & height needs to be previous height for the first block
	tempBlockInfo := &dtypes.BlockInfo{
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
