package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/models"
	dbtypes "github.com/cosmostation/cosmostation-cosmos/mintscan/api/models/types"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/utils"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// GetBlocks returns latest blocks
func GetBlocks(db *pg.DB, w http.ResponseWriter, r *http.Request) error {
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
		var blockInfo []dbtypes.BlockInfo
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
	var blockInfos []dbtypes.BlockInfo
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
		var validatorInfo dbtypes.ValidatorInfo
		_ = db.Model(&validatorInfo).
			Where("proposer = ?", blockInfo.Proposer).
			Select()

		// query transactions in the block
		var transactionInfos []dbtypes.TransactionInfo
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
func GetProposedBlocks(db *pg.DB, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	address := vars["address"]

	// convert to proposer address format
	validatorInfo, _ := utils.ConvertToProposerSlice(address, db)

	if len(validatorInfo) <= 0 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	address = validatorInfo[0].Proposer

	limit := int(100)
	from := int(0)

	// check limit param
	if len(r.URL.Query()["limit"]) > 0 {
		limit, _ = strconv.Atoi(r.URL.Query()["limit"][0])
	}

	// check max limit
	if limit > 100 {
		errors.ErrOverMaxLimit(w, http.StatusRequestedRangeNotSatisfiable)
		return nil
	}

	// check from param
	if len(r.URL.Query()["from"]) > 0 {
		from, _ = strconv.Atoi(r.URL.Query()["from"][0])
	} else {
		// Query the lastest block height
		var blockInfo []dbtypes.BlockInfo
		_ = db.Model(&blockInfo).
			Column("height").
			Order("id DESC").
			Limit(1).
			Select()
		if len(blockInfo) > 0 {
			from = int(blockInfo[0].Height)
		}
	}

	// query blocks
	blockInfos := make([]*dbtypes.BlockInfo, 0)
	_ = db.Model(&blockInfos).
		Where("height < ? AND proposer = ?", from, address).
		Limit(limit).
		Order("height DESC").
		Select()

	// check if blocks exists
	if len(blockInfos) <= 0 {
		return json.NewEncoder(w).Encode(blockInfos)
	}

	// query total number of proposer blocks
	totalNumProposerBlocks, _ := db.Model(&blockInfos).
		Where("proposer = ?", address).
		Count()

	resultBlocksByOperatorAddress := make([]*models.ResultBlocksByOperatorAddress, 0)
	for _, blockInfo := range blockInfos {
		var validatorInfo dbtypes.ValidatorInfo
		_ = db.Model(&validatorInfo).
			Where("proposer = ?", blockInfo.Proposer).
			Select()

		// query a number of txs
		var transactionInfos []dbtypes.TransactionInfo
		_ = db.Model(&transactionInfos).
			Column("tx_hash").
			Where("height = ?", blockInfo.Height).
			Select()

		// append transactions
		var txData models.TxData
		for _, transactionInfo := range transactionInfos {
			txData.Txs = append(txData.Txs, transactionInfo.TxHash)
		}

		tempResultBlocksByOperatorAddress := &models.ResultBlocksByOperatorAddress{
			Height:                 blockInfo.Height,
			Proposer:               blockInfo.Proposer,
			OperatorAddress:        validatorInfo.OperatorAddress,
			Moniker:                validatorInfo.Moniker,
			BlockHash:              blockInfo.BlockHash,
			Identity:               validatorInfo.Identity,
			NumTxs:                 blockInfo.NumTxs,
			TotalNumProposerBlocks: totalNumProposerBlocks,
			TxData:                 txData,
			Time:                   blockInfo.Time,
		}
		resultBlocksByOperatorAddress = append(resultBlocksByOperatorAddress, tempResultBlocksByOperatorAddress)
	}

	utils.Respond(w, resultBlocksByOperatorAddress)
	return nil
}
