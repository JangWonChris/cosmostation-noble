package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models"
	u "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/utils"
	dbtypes "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/types"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/utils"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// GetBlocks returns latest blocks
func GetBlocks(DB *pg.DB, w http.ResponseWriter, r *http.Request) error {
	// Declare default variables
	limit := int(100)
	afterBlock := int(1)

	// Check limit param
	if len(r.URL.Query()["limit"]) > 0 {
		limit, _ = strconv.Atoi(r.URL.Query()["limit"][0])
	}

	// Check max limit
	if limit > 100 {
		errors.ErrOverMaxLimit(w, http.StatusUnauthorized)
		return nil
	}

	// Check afterBlock param
	if len(r.URL.Query()["afterBlock"]) > 0 {
		afterBlock, _ = strconv.Atoi(r.URL.Query()["afterBlock"][0])
	} else {
		// Query the lastest block height
		var blockInfo []dbtypes.BlockInfo
		_ = DB.Model(&blockInfo).
			Column("height").
			Order("id DESC").
			Limit(1).
			Select()
		if len(blockInfo) > 0 {
			afterBlock = int(blockInfo[0].Height) - limit // // afterBlock should be latestHeight saved in DB minus limit
		}
	}

	// Query a number of blocks in an ascending order
	var blockInfos []dbtypes.BlockInfo
	_ = DB.Model(&blockInfos).
		Where("height > ?", afterBlock).
		Limit(limit).
		Order("id ASC").
		Select()

	// Check if blocks exists
	if len(blockInfos) <= 0 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	resultBlock := make([]*models.ResultBlock, 0)
	for _, blockInfo := range blockInfos {
		// Query validator information using proposer address
		var validatorInfo dbtypes.ValidatorInfo
		_ = DB.Model(&validatorInfo).
			Where("proposer = ?", blockInfo.Proposer).
			Select()

		// Query transactions in the block
		var transactionInfos []dbtypes.TransactionInfo
		_ = DB.Model(&transactionInfos).
			Where("height = ?", blockInfo.Height).
			Select()

		// Append transactions
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

	u.Respond(w, resultBlock)
	return nil
}

// GetProposedBlocksByAddress returns proposed blocks by querying any type of address
func GetProposedBlocks(DB *pg.DB, w http.ResponseWriter, r *http.Request) error {
	// Receive address
	vars := mux.Vars(r)
	address := vars["address"]

	// Change to proposer address format
	validatorInfo, _ := utils.ConvertToProposerSlice(address, DB)

	// Check the address by length of validatorInfo
	if len(validatorInfo) <= 0 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Proposer Address
	address = validatorInfo[0].Proposer

	// Declare default variables
	limit := int(100)
	from := int(0)

	// Check limit param
	if len(r.URL.Query()["limit"]) > 0 {
		limit, _ = strconv.Atoi(r.URL.Query()["limit"][0])
	}

	// Check max limit
	if limit > 100 {
		errors.ErrOverMaxLimit(w, http.StatusRequestedRangeNotSatisfiable)
		return nil
	}

	// Check from param
	if len(r.URL.Query()["from"]) > 0 {
		from, _ = strconv.Atoi(r.URL.Query()["from"][0])
	} else {
		// Query the lastest block height
		var blockInfo []dbtypes.BlockInfo
		_ = DB.Model(&blockInfo).
			Column("height").
			Order("id DESC").
			Limit(1).
			Select()
		if len(blockInfo) > 0 {
			from = int(blockInfo[0].Height)
		}
	}

	// Query blocks
	blockInfos := make([]*dbtypes.BlockInfo, 0)
	_ = DB.Model(&blockInfos).
		Where("height < ? AND proposer = ?", from, address).
		Limit(limit).
		Order("height DESC").
		Select()

	// Check if blocks exists
	if len(blockInfos) <= 0 {
		return json.NewEncoder(w).Encode(blockInfos)
	}

	// Query total number of proposer blocks
	totalNumProposerBlocks, _ := DB.Model(&blockInfos).
		Where("proposer = ?", address).
		Count()

	resultBlocksByOperatorAddr := make([]*models.ResultBlocksByOperatorAddr, 0)
	for _, blockInfo := range blockInfos {
		// Query validator information
		var validatorInfo dbtypes.ValidatorInfo
		_ = DB.Model(&validatorInfo).
			Where("proposer = ?", blockInfo.Proposer).
			Select()

		// Query a number of txs
		var transactionInfos []dbtypes.TransactionInfo
		_ = DB.Model(&transactionInfos).
			Column("tx_hash").
			Where("height = ?", blockInfo.Height).
			Select()

		// Append transactions
		var txData models.TxData
		for _, transactionInfo := range transactionInfos {
			txData.Txs = append(txData.Txs, transactionInfo.TxHash)
		}

		tempResultBlocksByOperatorAddr := &models.ResultBlocksByOperatorAddr{
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
		resultBlocksByOperatorAddr = append(resultBlocksByOperatorAddr, tempResultBlocksByOperatorAddr)
	}

	u.Respond(w, resultBlocksByOperatorAddr)
	return nil
}
