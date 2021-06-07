package common

import (
	"testing"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/handler"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"
	"go.uber.org/zap"
)

func TestGetBlocksByProposerNew(t *testing.T) {
	valAddr := "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c"

	limit := 5

	// rb := make([]model.ResultBlock, 0)
	var rb []model.ResultBlock

	res, err := tdb.Query(&rb, "select b.id, b.height, b.proposer, v.operator_address, v.moniker, b.block_hash, identity, num_txs, count(*) OVER() AS total_num_proposer_blocks, b.timestamp"+
		" from block as b, "+
		" (select proposer, operator_address, Identity, moniker from validator where operator_address = ? limit 1) as v"+
		" where v.proposer = b.proposer order by height desc limit ?", valAddr, limit)

	if err != nil {
		t.Log(err)
	} else {
		t.Log(res.RowsReturned())
	}

	result := make([]*model.ResultBlock, 0)

	for _, b := range rb {
		var txResps []*model.ResultTx
		if b.NumTxs > 0 {
			txs, err := tdb.QueryTransactionsInBlockHeight(handler.ChainIDMap[handler.ChainID], b.Height)
			if err != nil {
				zap.L().Error("failed to query transactions in a block", zap.Error(err))
				return
			}

			txResps = model.ParseTransactions(txs)
		}

		b := &model.ResultBlock{
			ID:                     b.ID,
			Height:                 b.Height,
			Proposer:               b.Proposer,
			OperatorAddress:        b.OperatorAddress,
			Moniker:                b.Moniker,
			BlockHash:              b.BlockHash,
			Identity:               b.Identity,
			NumTxs:                 b.NumTxs,
			TotalNumProposerBlocks: b.TotalNumProposerBlocks,
			Txs:                    txResps,
			Timestamp:              b.Timestamp,
		}

		result = append(result, b)
	}

	t.Log(result)
	return
}
