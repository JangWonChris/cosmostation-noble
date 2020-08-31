package exporter

import (
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// getBlock exports block information
func (ex *Exporter) getBlock(block *tmctypes.ResultBlock) (schema.Block, error) {
	b := schema.NewBlock(schema.Block{
		ChainID:       block.Block.Header.ChainID,
		Height:        block.Block.Height,
		Proposer:      block.Block.ProposerAddress.String(),
		BlockHash:     block.BlockMeta.BlockID.Hash.String(),
		ParentHash:    block.Block.Header.LastBlockID.Hash.String(),
		NumSignatures: int64(len(block.Block.LastCommit.Precommits)),
		NumTxs:        int64(len(block.Block.Data.Txs)),
		Timestamp:     block.Block.Time,
	})

	return *b, nil
}
