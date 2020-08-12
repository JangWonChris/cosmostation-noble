package exporter

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/client"
	ceCodec "github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/db"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/log"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/cosmos/cosmos-sdk/codec"
)

// Exporter implemnts a wrapper around configuration for this project
type Exporter struct {
	config *config.Config
	cdc    *codec.Codec
	client *client.Client
	db     *db.Database
}

// NewExporter initializes the required config
func NewExporter() *Exporter {
	// Create custom logger with a combination of using uber/zap and lumberjack.v2.
	l, _ := log.NewCustomLogger()
	zap.ReplaceGlobals(l)
	defer l.Sync()

	cfg := config.ParseConfig()

	client, err := client.NewClient(cfg.Node, cfg.KeybaseURL)
	if err != nil {
		zap.L().Error("failed to create new client", zap.Error(err))
		return &Exporter{}
	}

	// Connect to database
	// Ping database to verify connection is succeeded
	db := db.Connect(&cfg.DB)
	err = db.Ping()
	if err != nil {
		zap.L().Error("failed to ping database", zap.Error(err))
	}

	// Setup database tables
	db.CreateTables()

	// Set Bech32 address prefixes and BIP44 coin type for Cosmos
	// types.SetAppConfig()

	return &Exporter{cfg, ceCodec.Codec, client, db}
}

/*
// Start creates database tables and indexes using Postgres ORM library go-pg and
// starts syncing blockchain.
func (ex *Exporter) Start() error {
	// Store data initially
	ex.client.SaveBondedValidators()
	ex.client.SaveUnbondingAndUnBondedValidators()
	ex.client.SaveProposals()

	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		for {
			fmt.Println("start - sync blockchain")
			err := ex.sync()
			if err != nil {
				fmt.Printf("error - sync blockchain: %v\n", err)
			}
			fmt.Println("finish - sync blockchain")
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(7 * time.Second)
			c1 <- "sync governance and validators via LCD"
		}
	}()

	go func() {
		for {
			time.Sleep(20 * time.Minute)
			c2 <- "parsing from keybase server using keybase identity"
		}
	}()

	for {
		select {
		case msg1 := <-c1:
			fmt.Println("start - ", msg1)
			ex.client.SaveBondedValidators()
			ex.client.SaveUnbondingAndUnBondedValidators()
			ex.client.SaveProposals()
			fmt.Println("finish - ", msg1)
		case msg2 := <-c2:
			fmt.Println("start - ", msg2)
			ex.SaveValidatorKeyBase()
			fmt.Println("finish - ", msg2)
		}
	}
}

// sync compares block height between the height saved in your database and
// latest block height on the active chain and calls process to start ingesting blocks.
func (ex *Exporter) sync() error {
	// Query latest block height that is saved in your database
	// Synchronizing blocks from the scratch will return 0 and will ingest accordingly.
	dbHeight, err := ex.db.QueryLatestBlockHeight()
	if dbHeight == -1 {
		fmt.Println(errors.Wrap(err, "failed to query the latest block height from database."))
	}

	// Query latest block height on the active network
	latestBlockHeight, err := ex.client.LatestBlockHeight()
	if latestBlockHeight == -1 {
		fmt.Println(errors.Wrap(err, "failed to query the latest block height on the active network."))
	}

	// skip the first block since it has no pre-commits
	if dbHeight == 0 {
		dbHeight = 1
	}

	// Ingest all blocks up to the best height
	for i := dbHeight + 1; i <= latestBlockHeight; i++ {
		err = ex.process(i)
		if err != nil {
			return err
		}
		fmt.Printf("synced block %d/%d \n", i, latestBlockHeight)
	}

	return nil
}

// sync queries the block at the given height-1 from the node and ingests its metadata (Block,evidence)
// into the database. It also queries the next block to access the commits and stores the missed signatures.
func (ex Exporter) process(height int64) error {
	block, err := ex.client.Block(height)
	if err != nil {
		return fmt.Errorf("failed to query block using rpc client: %t", err)
	}

	nextBlock, err := ex.client.Block(height + 1)
	if err != nil {
		return fmt.Errorf("failed to query block using rpc client: %t", err)
	}

	prevBlock, err := ex.client.Block(block.Block.LastCommit.Height())
	if err != nil {
		return fmt.Errorf("failed to query block using rpc client: %t", err)
	}

	vals, err := ex.client.Validators(block.Block.LastCommit.Height())
	if err != nil {
		return fmt.Errorf("failed to query validators using rpc client: %t", err)
	}

	txs, err := ex.client.Txs(block)
	if err != nil {
		return fmt.Errorf("failed to get transactions for block: %t", err)
	}

	resultBlock, err := ex.getBlock(block)
	if err != nil {
		return fmt.Errorf("failed to get block: %t", err)
	}

	resultEvidence, err := ex.getEvidence(block, nextBlock)
	if err != nil {
		return fmt.Errorf("failed to get evidence: %t", err)
	}

	resultGenesisValSet, err := ex.getGenesisValidatorSet(block, vals)
	if err != nil {
		return fmt.Errorf("failed to get genesis validator set: %t", err)
	}

	resultTxs, err := ex.getTxs(txs)
	if err != nil {
		return fmt.Errorf("failed to get txs: %t", err)
	}

	resultTxIndex, err := ex.getTxIndex(txs)
	if err != nil {
		return fmt.Errorf("failed to get txs: %t", err)
	}

	resultMissingBlocks, resultAccumMissingBlocks, resultMisssingBlocksDetail, err := ex.getPowerEventHistory(prevBlock, block, vals)
	if err != nil {
		return fmt.Errorf("failed to get missing blocks: %t", err)
	}

	resultVote, resultDeposit, resultProposal, resultValidatorSet, err := ex.getTransactions(block)
	if err != nil {
		return fmt.Errorf("failed to get transactions: %t", err)
	}

	// Insert data into database
	err = ex.db.InsertExportedData(resultBlock, resultEvidence, resultGenesisValSet, resultMissingBlocks, resultAccumMissingBlocks,
		resultMisssingBlocksDetail, resultTxs, resultTxIndex, resultVote, resultDeposit, resultProposal, resultValidatorSet)

	if err != nil {
		return fmt.Errorf("failed to insert exported data: %t", err)
	}

	return nil
}

*/

// MIGRATION

// Start starts to synchronize blockchain data
func (ex *Exporter) Start() {
	// Store data initially
	ex.saveValidators()
	ex.saveProposals()

	restServerCh := make(chan string)
	keybaseCh := make(chan string)

	go func() {
		for {
			zap.S().Info("start - sync blockchain")
			err := ex.sync()
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
func (ex *Exporter) sync() error {
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

	// Ingest all blocks up to the latest height
	for i := dbHeight + 1; i <= latestBlockHeight; i++ {
		err = ex.process(i)
		if err != nil {
			return err
		}
		zap.S().Infof("synced block %d/%d", i, latestBlockHeight)
	}

	return nil
}

// process ingests chain data, such as block, transaction, validator, evidence information and
// save them in database.
func (ex *Exporter) process(height int64) error {
	block, err := ex.client.GetBlock(height)
	if err != nil {
		return fmt.Errorf("failed to query block: %s", err)
	}

	// First block has no previous block and no pre-commits.
	// Handle this to save first block information
	// var resultGenesisAccounts []schema.Account
	prevBlock := new(tmctypes.ResultBlock)
	height = block.Block.LastCommit.Height()
	if height == 0 {
		prevBlock = block
		height = 1
	} else {
		prevBlock, err = ex.client.GetBlock(block.Block.LastCommit.Height())
		if err != nil {
			return fmt.Errorf("failed to query previous block: %s", err)
		}
	}

	// if height == 1 {
	// var genesisAccts exported.GenesisAccounts
	// genesisAccts, err = ex.client.GetGenesisAccounts()
	// if err != nil {
	// 	return fmt.Errorf("failed to get genesis accounts: %s", err)
	// }

	// resultGenesisAccounts, err = ex.getGenesisAccounts(genesisAccts)
	// if err != nil {
	// 	return fmt.Errorf("failed to get block: %s", err)
	// }
	// }

	vals, err := ex.client.GetValidators(height, types.DefaultQueryValidatorsPage, types.DefaultQueryValidatorsPerPage)
	if err != nil {
		return fmt.Errorf("failed to query validators: %s", err)
	}

	txs, err := ex.client.GetTxs(block)
	if err != nil {
		return fmt.Errorf("failed to get transactions for block: %s", err)
	}

	resultBlock, err := ex.getBlock(block)
	if err != nil {
		return fmt.Errorf("failed to get block: %s", err)
	}

	resultGenesisValidatorsSet, err := ex.getGenesisValidatorsSet(block, vals)
	if err != nil {
		return fmt.Errorf("failed to get genesis validator set: %s", err)
	}

	resultAccounts, err := ex.getAccounts(block, txs)
	if err != nil {
		return fmt.Errorf("failed to get accounts: %s", err)
	}

	resultMissBlocks, resultAccumulatedMissBlocks, resultMissDetailBlocks, err := ex.getValidatorsUptime(prevBlock, block, vals)
	if err != nil {
		return fmt.Errorf("failed to get missing blocks: %s", err)
	}

	resultEvidence, err := ex.getEvidence(block)
	if err != nil {
		return fmt.Errorf("failed to get evidence: %s", err)
	}

	resultTxs, err := ex.getTxs(block, txs)
	if err != nil {
		return fmt.Errorf("failed to get txs: %s", err)
	}

	resultProposals, resultDeposits, resultVotes, err := ex.getGovernance(block, txs)
	if err != nil {
		return fmt.Errorf("failed to get governance: %s", err)
	}

	resultValidatorsPowerEventHistory, err := ex.getPowerEventHistory(block, txs)
	if err != nil {
		return fmt.Errorf("failed to get transactions: %s", err)
	}

	// TODO: is this right place to be?
	if ex.config.Alarm.Switch {
		ex.handlePushNotification(block, txs)
	}

	err = ex.db.InsertExportedData(schema.ExportData{
		ResultAccounts: resultAccounts,
		ResultBlock:    resultBlock,
		// ResultGenesisAccounts:             resultGenesisAccounts,
		ResultGenesisAccounts:             nil,
		ResultTxs:                         resultTxs,
		ResultEvidence:                    resultEvidence,
		ResultMissBlocks:                  resultMissBlocks,
		ResultMissDetailBlocks:            resultMissDetailBlocks,
		ResultAccumulatedMissBlocks:       resultAccumulatedMissBlocks,
		ResultProposals:                   resultProposals,
		ResultDeposits:                    resultDeposits,
		ReusltVotes:                       resultVotes,
		ResultGenesisValidatorsSet:        resultGenesisValidatorsSet,
		ResultValidatorsPowerEventHistory: resultValidatorsPowerEventHistory,
	})

	if err != nil {
		return err
	}

	return nil
}
