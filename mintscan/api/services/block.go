package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/db"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/schema"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/utils"

	"github.com/gorilla/mux"
)

// GetBlocks returns latest blocks
func GetBlocks(db *db.Database, w http.ResponseWriter, r *http.Request) error {
	limit := int(100)
	afterBlock := int(1)

	if len(r.URL.Query()["limit"]) > 0 {
		limit, _ = strconv.Atoi(r.URL.Query()["limit"][0])
	}

	// check max limit
	if limit > 100 {
		errors.ErrOverMaxLimit(w, http.StatusUnauthorized)
		return nil
	}

	if len(r.URL.Query()["afterBlock"]) > 0 {
		afterBlock, _ = strconv.Atoi(r.URL.Query()["afterBlock"][0])
	} else {
		latestBlockHeight, _ := db.QueryLatestBlockHeight()
		if latestBlockHeight >= 0 {
			afterBlock = latestBlockHeight - limit // afterBlock should be latestHeight saved in DB minus limit
		}
	}

	// Query blocks in an ascending order
	blocks, _ := db.QueryBlocks(afterBlock, limit)

	if len(blocks) <= 0 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	result := make([]*models.ResultBlock, 0)

	for _, block := range blocks {
		validator, _ := db.QueryValidatorByProposer(block.Proposer)

		// query transactions in the block
		txInfos, _ := db.QueryTransactions(block.Height)

		var txData models.TxData
		for _, txInfo := range txInfos {
			txData.Txs = append(txData.Txs, txInfo.TxHash)
		}

		tempBlockInfo := &models.ResultBlock{
			Height:          block.Height,
			Proposer:        block.Proposer,
			OperatorAddress: validator.OperatorAddress,
			Moniker:         validator.Moniker,
			BlockHash:       block.BlockHash,
			Identity:        validator.Identity,
			NumTxs:          block.NumTxs,
			TxData:          txData,
			Time:            block.Time,
		}
		result = append(result, tempBlockInfo)
	}

	utils.Respond(w, result)
	return nil
}

// GetProposedBlocksByAddress returns proposed blocks by querying any type of address
func GetProposedBlocks(db *db.Database, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	address := vars["address"]

	// Check if the input validator address exists
	validatorInfo, _ := db.ConvertToProposer(address)
	if validatorInfo.Proposer == "" {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	address = validatorInfo.Proposer

	limit := int(100) // default limit is 100
	before := int(0)
	after := int(0)
	offset := int(0)

	if len(r.URL.Query()["limit"]) > 0 {
		limit, _ = strconv.Atoi(r.URL.Query()["limit"][0])
	}

	if len(r.URL.Query()["before"]) > 0 {
		before, _ = strconv.Atoi(r.URL.Query()["before"][0])
	}

	if len(r.URL.Query()["after"]) > 0 {
		after, _ = strconv.Atoi(r.URL.Query()["after"][0])
	}

	if len(r.URL.Query()["offset"]) > 0 {
		offset, _ = strconv.Atoi(r.URL.Query()["offset"][0])
	}

	if limit > 100 {
		errors.ErrOverMaxLimit(w, http.StatusRequestedRangeNotSatisfiable)
		return nil
	}

	// Query blocks proposed by proposer
	blocks := make([]schema.BlockInfo, 0)

	switch {
	case before > 0:
		blocks, _ = db.QueryBlocksByProposer(address, limit, before, after, offset)
	case after > 0:
		blocks, _ = db.QueryBlocksByProposer(address, limit, before, after, offset)
	case offset >= 0:
		blocks, _ = db.QueryBlocksByProposer(address, limit, before, after, offset)
	}

	if len(blocks) <= 0 {
		return json.NewEncoder(w).Encode(blocks)
	}

	// Query total number of proposed blocks by a proposer
	totalNumProposerBlocks, _ := db.QueryTotalBlocksByProposer(address)

	result := make([]*models.ResultBlocksByOperatorAddress, 0)

	for i, block := range blocks {
		validatorInfo, _ := db.QueryValidatorByProposer(block.Proposer)

		// Query a number of txs
		txInfos, _ := db.QueryTransactions(block.Height)

		var txData models.TxData
		if len(txInfos) > 0 {
			for _, txInfo := range txInfos {
				txData.Txs = append(txData.Txs, txInfo.TxHash)
			}
		}

		tempResultBlocksByOperatorAddress := &models.ResultBlocksByOperatorAddress{
			ID:                     i + 1,
			Height:                 block.Height,
			Proposer:               block.Proposer,
			OperatorAddress:        validatorInfo.OperatorAddress,
			Moniker:                validatorInfo.Moniker,
			BlockHash:              block.BlockHash,
			Identity:               validatorInfo.Identity,
			NumTxs:                 block.NumTxs,
			TotalNumProposerBlocks: totalNumProposerBlocks,
			TxData:                 txData,
			Time:                   block.Time,
		}
		result = append(result, tempResultBlocksByOperatorAddress)
	}

	utils.Respond(w, result)
	return nil
}
