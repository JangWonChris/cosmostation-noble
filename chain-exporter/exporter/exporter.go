package exporter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"

	ceCodec "github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/db"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/lcd"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	resty "gopkg.in/resty.v1"
)

// ChainExporter implemnts a wrapper around configuration for this project
type ChainExporter struct {
	codec     *codec.Codec
	config    *config.Config
	db        *db.Database
	wsCtx     context.Context
	wsOut     <-chan ctypes.ResultEvent
	rpcClient *client.HTTP
}

// NewChainExporter initializes the required config
func NewChainExporter(config *config.Config) *ChainExporter {
	ce := &ChainExporter{
		codec:     ceCodec.Codec, // register Cosmos SDK codecs
		config:    config,
		db:        db.Connect(config), // connect to PostgreSQL
		wsCtx:     context.Background(),
		rpcClient: client.NewHTTP(config.Node.RPCNode, "/websocket"), // connect to Tendermint RPC client
	}

	// Ping database to verify connection is succeeded
	err := ce.db.Ping()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to ping database."))
	}

	// Setup database tables
	ce.db.CreateSchema()

	resty.SetTimeout(5 * time.Second) // sets timeout for request.

	return ce
}

// OnStart is an override method for BaseService, which starts a service
func (ce ChainExporter) OnStart() error {
	ce.rpcClient.OnStart()

	// Store data initially
	lcd.SaveBondedValidators(ce.db, ce.config)
	lcd.SaveUnbondingAndUnBondedValidators(ce.db, ce.config)
	lcd.SaveProposals(ce.db, ce.config)

	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		for {
			fmt.Println("start - sync blockchain")
			err := ce.sync()
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
			lcd.SaveBondedValidators(ce.db, ce.config)
			lcd.SaveUnbondingAndUnBondedValidators(ce.db, ce.config)
			lcd.SaveProposals(ce.db, ce.config)
			fmt.Println("finish - ", msg1)
		case msg2 := <-c2:
			fmt.Println("start - ", msg2)
			ce.SaveValidatorKeyBase()
			fmt.Println("finish - ", msg2)
		}
	}
}

// OnStop is an override method for BaseService, which stops a service
func (ce ChainExporter) OnStop() {
	ce.rpcClient.OnStop()
}

// sync synchronizes the block data from connected full node
func (ce ChainExporter) sync() error {
	var blocks []schema.BlockInfo
	err := ce.db.Model(&blocks).
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
	status, err := ce.rpcClient.Status()
	if err != nil {
		return err
	}
	maxHeight := status.SyncInfo.LatestBlockHeight

	if currentHeight == 1 {
		currentHeight = 0
	}

	// ingest all blocks up to the best height
	for i := currentHeight + 1; i <= maxHeight; i++ {
		err = ce.process(i)
		if err != nil {
			return err
		}
		fmt.Printf("synced block %d/%d \n", i, maxHeight)
	}
	return nil
}

// sync queries the block at the given height-1 from the node and ingests its metadata (blockinfo,evidence)
// into the database. It also queries the next block to access the commits and stores the missed signatures.
func (ce ChainExporter) process(height int64) error {
	blockInfo, err := ce.getBlockInfo(height)
	if err != nil {
		return err
	}

	evidenceInfo, err := ce.getEvidenceInfo(height)
	if err != nil {
		return err
	}

	genesisValsInfo, missInfo, accumMissInfo, missDetailInfo, err := ce.getValidatorSetInfo(height)
	if err != nil {
		return err
	}

	transactionInfo, voteInfo, depositInfo, proposalInfo, validatorSetInfo, err := ce.getTransactionInfo(height)
	if err != nil {
		return err
	}

	// Insert data into database
	err = ce.db.InsertExportedData(blockInfo, evidenceInfo, genesisValsInfo, missInfo, accumMissInfo,
		missDetailInfo, transactionInfo, voteInfo, depositInfo, proposalInfo, validatorSetInfo)

	if err != nil {
		return err
	}

	return nil
}
