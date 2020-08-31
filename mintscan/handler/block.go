package handler

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/schema"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

// GetBlocks returns latest blocks
func GetBlocks(rw http.ResponseWriter, r *http.Request) {
	before, after, limit, err := model.ParseHTTPArgsWithBeforeAfterLimit(r, model.DefaultBefore, model.DefaultAfter, model.DefaultLimit)
	if err != nil {
		zap.S().Debug("failed to parse HTTP args ", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "request is invalid")
		return
	}

	if limit > 100 {
		zap.S().Debug("failed to query with this limit ", zap.Int("request limit", limit))
		errors.ErrOverMaxLimit(rw, http.StatusUnauthorized)
		return
	}

	blocks, _ := s.db.QueryBlocks(before, after, limit)
	if len(blocks) <= 0 {
		model.Respond(rw, []model.ResultBlock{})
		return
	}

	result := make([]*model.ResultBlock, 0)

	for _, block := range blocks {
		validator, err := s.db.QueryValidatorByAny(block.Proposer)
		if err != nil {
			zap.L().Error("failed to query validator by proposer", zap.Error(err))
			return
		}

		txs, err := s.db.QueryTransactionsByBlockHeight(block.Height)
		if err != nil {
			zap.L().Error("failed to query txs", zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
			return
		}

		var txData model.TxData
		for _, tx := range txs {
			txData.Txs = append(txData.Txs, tx.TxHash)
		}

		b := &model.ResultBlock{
			Height:          block.Height,
			Proposer:        block.Proposer,
			OperatorAddress: validator.OperatorAddress,
			Moniker:         validator.Moniker,
			BlockHash:       block.BlockHash,
			Identity:        validator.Identity,
			NumSignatures:   block.NumSignatures,
			NumTxs:          block.NumTxs,
			TxData:          txData,
			Timestamp:       block.Timestamp,
		}

		result = append(result, b)
	}

	model.Respond(rw, result)
	return
}

// GetBlocksByProposer returns blocks that are proposed by a proposer
func GetBlocksByProposer(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proposer := vars["proposer"]

	before, after, limit, err := model.ParseHTTPArgsWithBeforeAfterLimit(r, model.DefaultBefore, model.DefaultAfter, model.DefaultLimit)
	if err != nil {
		zap.S().Debug("failed to parse HTTP args ", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "request is invalid")
		return
	}

	if limit > 100 {
		zap.S().Debug("failed to query with this limit ", zap.Int("request limit", limit))
		errors.ErrOverMaxLimit(rw, http.StatusUnauthorized)
		return
	}

	// Query validator information by any type of bech32 address, even moniker.
	val, _ := s.db.QueryValidatorByAny(proposer)
	if val.Proposer == "" {
		errors.ErrNotExist(rw, http.StatusNotFound)
		return
	}

	blocks, err := s.db.QueryBlocksByProposer(val.Proposer, before, after, limit)
	if err != nil {
		zap.L().Error("failed to query blocks", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	if len(blocks) <= 0 {
		model.Respond(rw, []schema.Block{})
		return
	}

	totalNum := len(blocks)

	result := make([]*model.ResultBlock, 0)

	for _, b := range blocks {
		val, err := s.db.QueryValidatorByAny(b.Proposer)
		if err != nil {
			zap.L().Error("failed to query validator", zap.Error(err))
			errors.ErrInternalServer(rw, http.StatusInternalServerError)
			return
		}

		txs, err := s.db.QueryTransactionsByBlockHeight(b.Height)
		if err != nil {
			zap.L().Error("failed to query transactions in a block", zap.Error(err))
			errors.ErrInternalServer(rw, http.StatusInternalServerError)
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

	model.Respond(rw, result)
	return
}
