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
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/lcd"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"

	"github.com/cosmos/cosmos-sdk/codec"

	resty "gopkg.in/resty.v1"
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

	client, err := client.NewClient(cfg.Node.RPCNode, "/websocket")
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

	// Set timeout for request
	resty.SetTimeout(5 * time.Second)

	return Exporter{cfg, ceCodec.Codec, client, db}
}

// Start creates database tables and indexes using Postgres ORM library go-pg and
// starts syncing blockchain.
func (ex Exporter) Start() error {
	// Store data initially
	lcd.SaveBondedValidators(ex.db, ex.cfg)
	lcd.SaveUnbondingAndUnBondedValidators(ex.db, ex.cfg)
	lcd.SaveProposals(ex.db, ex.cfg)

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
			lcd.SaveBondedValidators(ex.db, ex.cfg)
			lcd.SaveUnbondingAndUnBondedValidators(ex.db, ex.cfg)
			lcd.SaveProposals(ex.db, ex.cfg)
			fmt.Println("finish - ", msg1)
		case msg2 := <-c2:
			fmt.Println("start - ", msg2)
			ex.SaveValidatorKeyBase()
			fmt.Println("finish - ", msg2)
		}
	}
}

// sync synchronizes the block data from connected full node
func (ex Exporter) sync() error {
	var blocks []schema.BlockInfo
	err := ex.db.Model(&blocks).
		Order("height DESC").
		Limit(1).
		Select()
	if err != nil {
		return err
	}

	currentHeight := int64(1)
	if len(blocks) > 0 {
		currentHeight = blocks[0].Height
	}

	// query current height
	status, err := ex.client.Status()
	if err != nil {
		return err
	}
	maxHeight := status.SyncInfo.LatestBlockHeight

	if currentHeight == 1 {
		currentHeight = 0
	}

	// ingest all blocks up to the best height
	for i := currentHeight + 1; i <= maxHeight; i++ {
		err = ex.process(i)
		if err != nil {
			return err
		}
		fmt.Printf("synced block %d/%d \n", i, maxHeight)
	}
	return nil
}

// sync queries the block at the given height-1 from the node and ingests its metadata (Block,evidence)
// into the database. It also queries the next block to access the commits and stores the missed signatures.
func (ex Exporter) process(height int64) error {
	Block, err := ex.getBlock(height)
	if err != nil {
		return err
	}

	evidenceInfo, err := ex.getEvidenceInfo(height)
	if err != nil {
		return err
	}

	genesisValsInfo, missInfo, accumMissInfo, missDetailInfo, err := ex.getValidatorSetInfo(height)
	if err != nil {
		return err
	}

	transactionInfo, voteInfo, depositInfo, proposalInfo, validatorSetInfo, err := ex.getTransactionInfo(height)
	if err != nil {
		return err
	}

	// Insert data into database
	err = ex.db.InsertExportedData(Block, evidenceInfo, genesisValsInfo, missInfo, accumMissInfo,
		missDetailInfo, transactionInfo, voteInfo, depositInfo, proposalInfo, validatorSetInfo)

	if err != nil {
		return err
	}

	return nil
}
