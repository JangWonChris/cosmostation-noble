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

	// check limit param
	if len(r.URL.Query()["limit"]) > 0 {
		limit, _ = strconv.Atoi(r.URL.Query()["limit"][0])
	}

	// check max limit
	if limit > 100 {
		errors.ErrOverMaxLimit(w, http.StatusUnauthorized)
		return nil
	}

	// check afterBlock param
	if len(r.URL.Query()["afterBlock"]) > 0 {
		afterBlock, _ = strconv.Atoi(r.URL.Query()["afterBlock"][0])
	} else {
		// Query the lastest block height
		var blockInfo []schema.BlockInfo
		_ = db.Model(&blockInfo).
			Column("height").
			Order("id DESC").
			Limit(1).
			Select()
		if len(blockInfo) > 0 {
			afterBlock = int(blockInfo[0].Height) - limit // // afterBlock should be latestHeight saved in DB minus limit
		}
	}

	// query a number of blocks in an ascending order
	var blockInfos []schema.BlockInfo
	_ = db.Model(&blockInfos).
		Where("height > ?", afterBlock).
		Limit(limit).
		Order("id ASC").
		Select()

	// check if blocks exists
	if len(blockInfos) <= 0 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	resultBlock := make([]*models.ResultBlock, 0)
	for _, blockInfo := range blockInfos {
		// query validator information using proposer address
		var validatorInfo schema.ValidatorInfo
		_ = db.Model(&validatorInfo).
			Where("proposer = ?", blockInfo.Proposer).
			Select()

		// query transactions in the block
		var transactionInfos []schema.TransactionInfo
		_ = db.Model(&transactionInfos).
			Where("height = ?", blockInfo.Height).
			Select()

		// append transactions
		var txData models.TxData
		for _, transactionInfo := range transactionInfos {
			txData.Txs = append(txData.Txs, transactionInfo.TxHash)
		}

		tempBlockInfo := &models.ResultBlock{
			Height:          blockInfo.Height,
			Proposer:        blockInfo.Proposer,
			OperatorAddress: validatorInfo.OperatorAddress,
			Moniker:         validatorInfo.Moniker,
			BlockHash:       blockInfo.BlockHash,
			Identity:        validatorInfo.Identity,
			NumTxs:          blockInfo.NumTxs,
			TxData:          txData,
			Time:            blockInfo.Time,
		}
		resultBlock = append(resultBlock, tempBlockInfo)
	}

	utils.Respond(w, resultBlock)
	return nil
}

// GetProposedBlocksByAddress returns proposed blocks by querying any type of address
func GetProposedBlocks(db *db.Database, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	address := vars["address"]

	// check if the input validator address exists
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

	resultBlocksByOperatorAddress := make([]*models.ResultBlocksByOperatorAddress, 0)
	for i, block := range blocks {
		validatorInfo, _ := db.QueryValidatorInfoByProposer(block.Proposer)

		// query a number of txs
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
		resultBlocksByOperatorAddress = append(resultBlocksByOperatorAddress, tempResultBlocksByOperatorAddress)
	}

	utils.Respond(w, resultBlocksByOperatorAddress)
	return nil
}
