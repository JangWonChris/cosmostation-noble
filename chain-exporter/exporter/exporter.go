package exporter

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	ceCodec "github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/databases"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/lcd"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	resty "gopkg.in/resty.v1"
)

// ChainExporterService implemnts a wrapper around configuration for this project
type ChainExporterService struct {
	codec     *codec.Codec
	config    *config.Config
	db        *databases.Database
	wsCtx     context.Context
	wsOut     <-chan ctypes.ResultEvent
	rpcClient *client.HTTP
}

// NewChainExporterService initializes the required config
func NewChainExporterService(config *config.Config) *ChainExporterService {
	ces := &ChainExporterService{
		codec:     ceCodec.Codec, // register Cosmos SDK codecs
		config:    config,
		db:        databases.Connect(config), // connect to PostgreSQL
		wsCtx:     context.Background(),
		rpcClient: client.NewHTTP(config.Node.GaiadURL, "/websocket"), // connect to Tendermint RPC client
	}

	ces.db.CreateSchema()

	resty.SetTimeout(5 * time.Second)
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	return ces
}

// OnStart is an override method for BaseService, which starts a service
func (ces ChainExporterService) OnStart() error {
	ces.rpcClient.OnStart()

	// Store data initially
	lcd.SaveBondedValidators(ces.db, ces.config)
	lcd.SaveUnbondingAndUnBondedValidators(ces.db, ces.config)
	lcd.SaveProposals(ces.db, ces.config)

	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		for {
			fmt.Println("start - sync blockchain")
			err := ces.sync()
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
			lcd.SaveBondedValidators(ces.db, ces.config)
			lcd.SaveUnbondingAndUnBondedValidators(ces.db, ces.config)
			lcd.SaveProposals(ces.db, ces.config)
			fmt.Println("finish - ", msg1)
		case msg2 := <-c2:
			fmt.Println("start - ", msg2)
			ces.SaveValidatorKeyBase()
			fmt.Println("finish - ", msg2)
		}
	}
}

// OnStop is an override method for BaseService, which stops a service
func (ces ChainExporterService) OnStop() {
	ces.rpcClient.OnStop()
}

// sync synchronizes the block data from connected full node
func (ces ChainExporterService) sync() error {
	var blocks []schema.BlockInfo
	err := ces.db.Model(&blocks).
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
	status, err := ces.rpcClient.Status()
	if err != nil {
		return err
	}
	maxHeight := status.SyncInfo.LatestBlockHeight

	if currentHeight == 1 {
		currentHeight = 0
	}

	// ingest all blocks up to the best height
	for i := currentHeight + 1; i <= maxHeight; i++ {
		err = ces.process(i)
		if err != nil {
			return err
		}
		fmt.Printf("synced block %d/%d \n", i, maxHeight)
	}
	return nil
}

// sync queries the block at the given height-1 from the node and ingests its metadata (blockinfo,evidence)
// into the database. It also queries the next block to access the commits and stores the missed signatures.
func (ces ChainExporterService) process(height int64) error {
	blockInfo, err := ces.getBlockInfo(height)
	if err != nil {
		return err
	}

	evidenceInfo, err := ces.getEvidenceInfo(height)
	if err != nil {
		return err
	}

	genesisValsInfo, missInfo, accumMissInfo, missDetailInfo, err := ces.getValidatorSetInfo(height)
	if err != nil {
		return err
	}

	transactionInfo, voteInfo, depositInfo, proposalInfo, validatorSetInfo, err := ces.getTransactionInfo(height)
	if err != nil {
		return err
	}

	// Insert data into database
	err = ces.db.SaveExportedData(blockInfo, evidenceInfo, genesisValsInfo, missInfo, accumMissInfo,
		missDetailInfo, transactionInfo, voteInfo, depositInfo, proposalInfo, validatorSetInfo)

	if err != nil {
		return err
	}

	return nil
}
