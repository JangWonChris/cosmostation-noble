package common

import (
	"net/http"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/errors"
	"github.com/cosmostation/cosmostation-cosmos/handler"
	"github.com/cosmostation/cosmostation-cosmos/model"
	mdschema "github.com/cosmostation/mintscan-database/schema"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

// GetBlocks returns blocks with given params.
func GetBlocks(rw http.ResponseWriter, r *http.Request) {
	from, limit, err := model.ParseHTTPArgs(r)
	if err != nil {
		zap.S().Debug("failed to parse HTTP args ", zap.Error(err))
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "request is invalid")
		return
	}

	blocks, err := s.DB.QueryBlocks(from, limit)
	if err != nil {
		zap.S().Debug("failed to get blocks ", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	if len(blocks) <= 0 {
		errors.ErrNotFound(rw, http.StatusNotFound)
		return
	}

	result := make([]*model.ResultBlock, 0)

	for _, block := range blocks {
		validator, err := s.DB.QueryValidatorByAnyAddr(block.Proposer)
		if err != nil {
			zap.S().Error("failed to query validator by proposer", zap.Error(err))
			return
		}

		b := &model.ResultBlock{
			ID:              block.ID,
			ChainID:         handler.ChainNumMap[block.ChainInfoID],
			Height:          block.Height,
			Proposer:        block.Proposer,
			OperatorAddress: validator.OperatorAddress,
			Moniker:         validator.Moniker,
			BlockHash:       block.Hash,
			Identity:        validator.Identity,
			NumSignatures:   block.NumSignatures,
			NumTxs:          block.NumTxs,
			Txs:             nil,
			Timestamp:       block.Timestamp,
		}

		result = append(result, b)
	}

	model.Respond(rw, result)
	return
}

// GetBlocks returns blocks with given params.
func GetBlocksByID(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.S().Debug("failed to parse int args ", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	blocks, err := s.DB.QueryBlockByID(id)
	if err != nil {
		zap.S().Debug("failed to get blocks ", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	if len(blocks) <= 0 {
		zap.S().Error("failed to get block length : 0", zap.Error(err))
		errors.ErrNotFound(rw, http.StatusNotFound)
		return
	}

	result := make([]*model.ResultBlock, 0)

	for _, block := range blocks {
		validator, err := s.DB.QueryValidatorByAnyAddr(block.Proposer)
		if err != nil {
			zap.S().Error("failed to query validator by proposer", zap.Error(err))
			return
		}

		var txResps []*model.ResultTx
		if block.NumTxs > 0 {
			txs, err := s.DB.QueryTransactionsByBlockID(block.ID)
			if err != nil {
				zap.L().Error("failed to query txs", zap.Error(err))
				errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
				return
			}

			txResps = model.ParseTransactions(txs)
		}

		b := &model.ResultBlock{
			ID:              block.ID,
			ChainID:         handler.ChainNumMap[block.ChainInfoID],
			Height:          block.Height,
			Proposer:        block.Proposer,
			OperatorAddress: validator.OperatorAddress,
			Moniker:         validator.Moniker,
			BlockHash:       block.Hash,
			Identity:        validator.Identity,
			NumSignatures:   block.NumSignatures,
			NumTxs:          block.NumTxs,
			Txs:             txResps,
			Timestamp:       block.Timestamp,
		}

		result = append(result, b)
	}

	model.Respond(rw, result)
	return
}

// GetBlocks returns blocks with given params.
func GetBlocksByHash(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hashStr := vars["hash"]

	blocks, err := s.DB.QueryBlockByHash(hashStr)
	if err != nil {
		zap.S().Debug("failed to get blocks ", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	if len(blocks) <= 0 {
		errors.ErrNotFound(rw, http.StatusNotFound)
		return
	}

	result := make([]*model.ResultBlock, 0)

	for _, block := range blocks {
		validator, err := s.DB.QueryValidatorByAnyAddr(block.Proposer)
		if err != nil {
			zap.S().Error("failed to query validator by proposer", zap.Error(err))
			return
		}

		var txResps []*model.ResultTx
		if block.NumTxs > 0 {
			txs, err := s.DB.QueryTransactionsByBlockID(block.ID)
			if err != nil {
				zap.L().Error("failed to query txs", zap.Error(err))
				errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
				return
			}

			txResps = model.ParseTransactions(txs)
		}

		b := &model.ResultBlock{
			ID:              block.ID,
			ChainID:         handler.ChainNumMap[block.ChainInfoID],
			Height:          block.Height,
			Proposer:        block.Proposer,
			OperatorAddress: validator.OperatorAddress,
			Moniker:         validator.Moniker,
			BlockHash:       block.Hash,
			Identity:        validator.Identity,
			NumSignatures:   block.NumSignatures,
			NumTxs:          block.NumTxs,
			Txs:             txResps,
			Timestamp:       block.Timestamp,
		}

		result = append(result, b)
	}

	model.Respond(rw, result)
	return
}

func GetBlockByChainIDHeight(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chainIDStr := vars["chainid"]
	heightStr := vars["height"]

	height, err := strconv.ParseInt(heightStr, 10, 64)
	if err != nil {
		zap.S().Debug("failed to parse int block height", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	blocks, err := s.DB.QueryBlockByChainIDHeight(handler.ChainIDMap[chainIDStr], height)
	if err != nil {
		zap.S().Debug("failed to get blocks ", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	if len(blocks) <= 0 {
		errors.ErrNotFound(rw, http.StatusNotFound)
		return
	}

	result := make([]*model.ResultBlock, 0)

	for _, block := range blocks {
		validator, err := s.DB.QueryValidatorByAnyAddr(block.Proposer)
		if err != nil {
			zap.S().Error("failed to query validator by proposer", zap.Error(err))
			return
		}

		var txResps []*model.ResultTx
		if block.NumTxs > 0 {
			txs, err := s.DB.QueryTransactionsByBlockID(block.ID)
			if err != nil {
				zap.L().Error("failed to query txs", zap.Error(err))
				errors.ErrServerUnavailable(rw, http.StatusServiceUnavailable)
				return
			}

			txResps = model.ParseTransactions(txs)
		}

		b := &model.ResultBlock{
			ID:              block.ID,
			ChainID:         handler.ChainNumMap[block.ChainInfoID],
			Height:          block.Height,
			Proposer:        block.Proposer,
			OperatorAddress: validator.OperatorAddress,
			Moniker:         validator.Moniker,
			BlockHash:       block.Hash,
			Identity:        validator.Identity,
			NumSignatures:   block.NumSignatures,
			NumTxs:          block.NumTxs,
			Txs:             txResps,
			Timestamp:       block.Timestamp,
		}

		result = append(result, b)
	}

	model.Respond(rw, result)
	return
}

// GetBlocksByProposer returns blocks that are proposed by the proposer.
func GetBlocksByProposer(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proposer := vars["proposer"]

	from, limit, err := model.ParseHTTPArgs(r)
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
	val, err := s.DB.QueryValidatorByAnyAddr(proposer)
	if err != nil {
		zap.S().Errorf("failed to query validator information: %s", err)
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	if val.Proposer == "" {
		errors.ErrNotExist(rw, http.StatusNotFound)
		return
	}

	blocks, err := s.DB.QueryBlocksByProposer(val.Proposer, from, limit)
	if err != nil {
		zap.L().Error("failed to query blocks", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	if len(blocks) <= 0 {
		model.Respond(rw, []mdschema.Block{})
		return
	}

	totalNum, err := s.DB.CountProposedBlocksByProposer(val.Proposer)
	if err != nil {
		zap.L().Error("failed to count proposed blocks by proposer", zap.Error(err))
		errors.ErrInternalServer(rw, http.StatusInternalServerError)
		return
	}

	result := make([]*model.ResultBlock, 0)

	for _, b := range blocks {

		b := &model.ResultBlock{
			ID:                     b.ID,
			ChainID:                handler.ChainNumMap[b.ChainInfoID],
			Height:                 b.Height, // 사용 값
			Proposer:               val.Proposer,
			OperatorAddress:        val.OperatorAddress,
			Moniker:                val.Moniker,
			BlockHash:              b.Hash, // 사용 값
			Identity:               val.Identity,
			NumTxs:                 b.NumTxs, // 사용 값
			TotalNumProposerBlocks: totalNum,
			Txs:                    nil,
			Timestamp:              b.Timestamp, // 사용 값
		}

		result = append(result, b)
	}

	model.Respond(rw, result)
	return
}
