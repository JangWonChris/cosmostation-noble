package exporter

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/client"
	ceCodec "github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/db"

	"github.com/cosmos/cosmos-sdk/codec"
)

// Exporter implemnts a wrapper around configuration for this project
type Exporter struct {
	cfg    *config.Config
	cdc    *codec.Codec
	client client.Client
	db     *db.Database
}

// NewExporter initializes the required config
func NewExporter() Exporter {
	cfg := config.ParseConfig()

	client, err := client.NewClient(cfg.Node.RPCNode, cfg.Node.LCDEndpoint)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to connect client."))
	}

	// Connect to database
	db := db.Connect(&cfg.DB)

	// Ping database to verify connection is succeeded
	err = db.Ping()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to ping database."))
	}

	// Setup database tables
	db.CreateTables()

	return Exporter{cfg, ceCodec.Codec, client, db}
}

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
		log.Fatal(errors.Wrap(err, "failed to query the latest block height from database."))
	}

	// Query latest block height on the active network
	latestBlockHeight, err := ex.client.LatestBlockHeight()
	if latestBlockHeight == -1 {
		log.Fatal(errors.Wrap(err, "failed to query the latest block height on the active network."))
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
		resultMisssingBlocksDetail, resultTxs, resultVote, resultDeposit, resultProposal, resultValidatorSet)

	if err != nil {
		return fmt.Errorf("failed to insert exported data: %t", err)
	}

	return nil
}
