package handler

import (
	"testing"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"
	"go.uber.org/zap"
)

func TestGetBlocksByProposer(t *testing.T) {

	proposer := "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c"

	limit := 5

	// Query validator information by any type of bech32 address, even moniker.
	val, err := tdb.QueryValidatorByAny(proposer)
	if err != nil {
		zap.S().Errorf("failed to query validator information: %s", err)
		return
	}

	if val.Proposer == "" {
		return
	}

	blocks, err := tdb.QueryBlocksByProposer(val.Proposer, 0, 0, limit)
	if err != nil {
		zap.L().Error("failed to query blocks", zap.Error(err))
		return
	}

	if len(blocks) <= 0 {
		return
	}

	totalNum, err := tdb.CountProposedBlocks(val.Proposer)
	if err != nil {
		zap.L().Error("failed to count proposed blocks by proposer", zap.Error(err))
		return
	}

	result := make([]*model.ResultBlock, 0)

	for _, b := range blocks {
		val, err := tdb.QueryValidatorByAny(b.Proposer)
		if err != nil {
			zap.L().Error("failed to query validator", zap.Error(err))
			return
		}

		txs, err := tdb.QueryTransactionsByBlockHeight(b.Height)
		if err != nil {
			zap.L().Error("failed to query transactions in a block", zap.Error(err))
			return
		}

		var txData model.TxData
		if len(txs) > 0 {
			for _, tx := range txs {
				txData.Txs = append(txData.Txs, tx.TxHash)
			}
		}

		b := &model.ResultBlock{
			ID:                     b.ID,
			Height:                 b.Height,
			Proposer:               b.Proposer,
			OperatorAddress:        val.OperatorAddress,
			Moniker:                val.Moniker,
			BlockHash:              b.BlockHash,
			Identity:               val.Identity,
			NumTxs:                 b.NumTxs,
			TotalNumProposerBlocks: totalNum,
			TxData:                 txData,
			Timestamp:              b.Timestamp,
		}

		result = append(result, b)
	}

	t.Log(result)
	return
}

func TestGetBlocksByProposerNew(t *testing.T) {

	operator_address := "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c"

	limit := 5

	// rb := make([]model.ResultBlock, 0)
	var rb []model.ResultBlock

	res, err := tdb.Query(&rb, "select b.id, b.height, b.proposer, v.operator_address, v.moniker, b.block_hash, identity, num_txs, count(*) OVER() AS total_num_proposer_blocks, b.timestamp"+
		" from block as b, "+
		" (select proposer, operator_address, Identity, moniker from validator where operator_address = ? limit 1) as v"+
		" where v.proposer = b.proposer order by height desc limit ?", operator_address, limit)

	if err != nil {
		t.Log(err)
	} else {
		t.Log(res.RowsReturned())
	}

	result := make([]*model.ResultBlock, 0)

	for _, b := range rb {
		var txData model.TxData
		if b.NumTxs > 0 {
			txs, err := tdb.QueryTransactionsByBlockHeight(b.Height)
			if err != nil {
				zap.L().Error("failed to query transactions in a block", zap.Error(err))
				return
			}

			if len(txs) > 0 {
				for _, tx := range txs {
					txData.Txs = append(txData.Txs, tx.TxHash)
				}
			}
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
			TxData:                 txData,
			Timestamp:              b.Timestamp,
		}

		result = append(result, b)
	}

	t.Log(result)
	return
}
