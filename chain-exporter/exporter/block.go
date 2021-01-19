package exporter

import (
	"encoding/json"
	"fmt"

	"github.com/cosmostation/mintscan-backend-library/db/schema"
	tmcTypes "github.com/tendermint/tendermint/rpc/core/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// getBlock exports block information.
func (ex *Exporter) getBlock(block *tmcTypes.ResultBlock) (*schema.Block, error) {
	b := schema.Block{
		ChainID:       block.Block.Header.ChainID,
		Height:        block.Block.Height,
		Proposer:      block.Block.ProposerAddress.String(),
		BlockHash:     block.BlockID.Hash.String(),
		ParentHash:    block.Block.Header.LastBlockID.Hash.String(),
		NumSignatures: int64(len(block.Block.LastCommit.Signatures)),
		NumTxs:        int64(len(block.Block.Data.Txs)),
		Timestamp:     block.Block.Time,
	}

	return &b, nil
}

// getBlockJSONChunk decodes transactions in a block and return a format of database transaction.
func (ex *Exporter) getBlockJSONChunk(block *tmctypes.ResultBlock) (schema.RawBlock, error) {
	b := new(schema.RawBlock)

	chunk, err := json.Marshal(block)
	if err != nil {
		return schema.RawBlock{}, fmt.Errorf("failed to marshal block : %s", err)
	}
	b.ChainID = block.Block.ChainID
	b.Height = block.Block.Height
	b.BlockHash = block.BlockID.Hash.String()
	b.NumTxs = int64(len(block.Block.Data.Txs))
	b.Chunk = chunk

	return *b, nil
}
