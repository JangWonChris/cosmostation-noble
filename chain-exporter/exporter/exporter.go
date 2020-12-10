package exporter

import (
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/client"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/db"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
)

var (
	// Version is a project's version string.
	Version = "Development"

	// Commit is commit hash of this project.
	Commit = ""
)

// Exporter is
type Exporter struct {
	config *config.Config
	client *client.Client
	db     *db.Database
}

// NewExporter returns new Exporter instance
func NewExporter() *Exporter {
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)
	defer l.Sync()

	// Parse config from configuration file (config.yaml).
	config := config.ParseConfig()

	// Create new client with node configruation.
	// Client is used for requesting any type of network data from RPC full node and REST Server.
	client, err := client.NewClient(config.Node, config.KeybaseURL)
	if err != nil {
		zap.L().Error("failed to create new client", zap.Error(err))
		return &Exporter{}
	}

	// Create connection with PostgreSQL database and
	// Ping database to verify connection is success.
	db := db.Connect(&config.DB)
	err = db.Ping()
	if err != nil {
		zap.L().Error("failed to ping database", zap.Error(err))
		return &Exporter{}
	}

	// Create database tables if not exist already
	db.CreateTables()

	return &Exporter{config, client, db}
}

// Start starts to synchronize blockchain data
func (ex *Exporter) Start(initialHeight int64, chunkOnly bool) {
	zap.S().Info("Starting Chain Exporter...")
	zap.S().Infof("Network Type: %s | Version: %s | Commit: %s", ex.config.Node.NetworkType, Version, Commit)

	//close grpc
	defer ex.client.Close()

	// Store data initially
	ex.saveValidators()
	ex.saveProposals()

	restServerCh := make(chan string)
	keybaseCh := make(chan string)

	go func() {
		for {
			zap.S().Info("start - sync blockchain")
			err := ex.sync(initialHeight, chunkOnly)
			if err != nil {
				zap.S().Infof("error - sync blockchain: %s\n", err)
			}
			zap.S().Info("finish - sync blockchain")

			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(7 * time.Second)
			restServerCh <- "sync governance and validators via LCD"
		}
	}()

	go func() {
		for {
			time.Sleep(20 * time.Minute)
			keybaseCh <- "sync validators keybase identities"
		}
	}()

	for {
		select {
		case msg1 := <-restServerCh:
			zap.S().Infof("start - %s", msg1)
			ex.saveValidators()
			ex.saveProposals()
			zap.S().Infof("finish - %s", msg1)
		case msg2 := <-keybaseCh:
			zap.S().Infof("start - %s", msg2)
			ex.saveValidatorsIdentities()
			zap.S().Infof("finish - %s", msg2)
		}
	}
}

// sync compares block height between the height saved in your database and
// the latest block height on the active chain and calls process to start ingesting data.
func (ex *Exporter) sync(initialHeight int64, chunkOnly bool) error {
	// Query latest block height saved in database
	dbHeight, err := ex.db.QueryLatestBlockHeight()
	if dbHeight == -1 {
		return fmt.Errorf("unexpected error in database: %s", err)
	}

	// Query latest block height on the active network
	latestBlockHeight, err := ex.client.GetLatestBlockHeight()
	if latestBlockHeight == -1 {
		return fmt.Errorf("failed to query the latest block height on the active network: %s", err)
	}

	if dbHeight == 0 && initialHeight != 0 {
		dbHeight = initialHeight - 1
		zap.S().Info("initial Height set : ", initialHeight)
	}

	// Ingest all blocks up to the latest height
	for i := dbHeight + 1; i <= latestBlockHeight; i++ {
		err = ex.process(i, chunkOnly)
		if err != nil {
			return err
		}
		zap.S().Infof("synced block %d/%d", i, latestBlockHeight)
	}

	return nil
}

// process ingests chain data, such as block, transaction, validator, evidence information and
// save them in database.
func (ex *Exporter) process(height int64, chunkOnly bool) error {
	block, err := ex.client.GetBlock(height)
	if err != nil {
		return fmt.Errorf("failed to query block: %s", err)
	}

	exportData := new(schema.ExportData)

	// var prevBlock *tmctypes.ResultBlock
	// var vals *tmctypes.ResultValidators
	// var resultBlock schema.Block
	// var resultAccounts []schema.Account
	// var resultEvidence []schema.Evidence
	// var resultTxs []schema.TransactionLegacy
	// var resultTxsMessages []schema.TransactionMessage
	// var resultProposals []schema.Proposal
	// var resultDeposits []schema.Deposit
	// var resultVotes []schema.Vote
	// var resultValidatorsPowerEventHistory []schema.PowerEventHistory
	// var resultGenesisValidatorsSet []schema.PowerEventHistory
	// var resultMissBlocks, resultAccumulatedMissBlocks []schema.Miss
	// var resultMissDetailBlocks []schema.MissDetail

	// block chunk

	txs, err := ex.client.GetTxs(block)
	if err != nil {
		return fmt.Errorf("failed to get transactions for block: %s", err)
	}

	exportData.ResultTxsJSONChunk, err = ex.getTxsJSONChunk(txs)
	if err != nil {
		return fmt.Errorf("failed to get txs: %s", err)
	}

	if !chunkOnly {

		log.Println("!chunkOnly test")

		if block.Block.LastCommit.Height != 0 {
			prevBlock, err := ex.client.GetBlock(block.Block.LastCommit.Height)
			if err != nil {
				return fmt.Errorf("failed to query previous block: %s", err)
			}

			vals, err := ex.client.GetValidators(block.Block.LastCommit.Height, types.DefaultQueryValidatorsPage, types.DefaultQueryValidatorsPerPage)
			if err != nil {
				return fmt.Errorf("failed to query validators: %s", err)
			}

			exportData.ResultGenesisValidatorsSet, err = ex.getGenesisValidatorsSet(block, vals)
			if err != nil {
				return fmt.Errorf("failed to get genesis validator set: %s", err)
			}
			exportData.ResultMissBlocks, exportData.ResultAccumulatedMissBlocks, exportData.ResultMissDetailBlocks, err = ex.getValidatorsUptime(prevBlock, block, vals)
			if err != nil {
				return fmt.Errorf("failed to get missing blocks: %s", err)
			}
		}

		exportData.ResultBlock, err = ex.getBlock(block)
		if err != nil {
			return fmt.Errorf("failed to get block: %s", err)
		}

		exportData.ResultAccounts, err = ex.getAccounts(block, txs)
		if err != nil {
			return fmt.Errorf("failed to get accounts: %s", err)
		}

		exportData.ResultEvidence, err = ex.getEvidence(block)
		if err != nil {
			return fmt.Errorf("failed to get evidence: %s", err)
		}

		exportData.ResultTxs, err = ex.getTxs(block, txs)
		if err != nil {
			return fmt.Errorf("failed to get txs: %s", err)
		}

		exportData.ResultTxsMessages, err = ex.extractAccount(txs)
		if err != nil {
			return fmt.Errorf("failed to get account by each tx message: %s", err)
		}

		exportData.ResultProposals, exportData.ResultDeposits, exportData.ResultVotes, err = ex.getGovernance(block, txs)
		if err != nil {
			return fmt.Errorf("failed to get governance: %s", err)
		}

		exportData.ResultValidatorsPowerEventHistory, err = ex.getPowerEventHistory(block, txs)
		if err != nil {
			return fmt.Errorf("failed to get transactions: %s", err)
		}

		// TODO: is this right place to be?
		if ex.config.Alarm.Switch {
			ex.handlePushNotification(block, txs)
		}
	}

	err = ex.db.InsertExportedData(exportData)

	// err = ex.db.InsertExportedData(schema.ExportData{
	// 	ResultAccounts:                    resultAccounts,
	// 	ResultBlock:                       resultBlock,
	// 	ResultTxs:                         resultTxs,
	// 	ResultTxsJSONChunk:                resultTxsJSONChunk,
	// 	ResultTxsMessages:                 resultTxsMessages,
	// 	ResultEvidence:                    resultEvidence,
	// 	ResultMissBlocks:                  resultMissBlocks,
	// 	ResultMissDetailBlocks:            resultMissDetailBlocks,
	// 	ResultAccumulatedMissBlocks:       resultAccumulatedMissBlocks,
	// 	ResultProposals:                   resultProposals,
	// 	ResultDeposits:                    resultDeposits,
	// 	ResultVotes:                       resultVotes,
	// 	ResultGenesisValidatorsSet:        resultGenesisValidatorsSet,
	// 	ResultValidatorsPowerEventHistory: resultValidatorsPowerEventHistory,
	// })

	if err != nil {
		return err
	}

	return nil
}
