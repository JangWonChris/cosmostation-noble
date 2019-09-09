package exporter

import (
	"context"
	"crypto/tls"
	"os"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/alarm-notification/config"
	"github.com/cosmostation/cosmostation-cosmos/alarm-notification/databases"

	"github.com/cosmos/cosmos-sdk/codec"
	gaiaApp "github.com/cosmos/gaia/app"

	"github.com/go-pg/pg"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/rpc/client"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	resty "gopkg.in/resty.v1"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
)

// ChainExporterService wraps below params
type ChainExporterService struct {
	cmn.BaseService
	codec     *codec.Codec
	config    *config.Config
	db        *pg.DB
	wsCtx     context.Context
	wsOut     <-chan ctypes.ResultEvent
	rpcClient *client.HTTP
}

// NewChainExporterService initializes the required config
func NewChainExporterService(config *config.Config) *ChainExporterService {
	ces := &ChainExporterService{
		codec:     gaiaApp.MakeCodec(), // register Cosmos SDK codecs
		config:    config,
		db:        databases.ConnectDatabase(config), // connect to PostgreSQL
		wsCtx:     context.Background(),
		rpcClient: client.NewHTTP(config.Node.GaiadURL, "/websocket"), // connect to Tendermint RPC client
	}

	// setup database schema
	databases.CreateSchema(ces.db)

	// Register a service that can be started, stopped, and reset
	ces.BaseService = *cmn.NewBaseService(logger, "ChainExporterService", ces)

	// sets timeout for request.
	resty.SetTimeout(5 * time.Second)
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // test locally

	return ces
}

// OnStart is an override method for BaseService, which starts a service
func (ces *ChainExporterService) OnStart() error {
	ces.BaseService.OnStart()
	ces.rpcClient.OnStart()

	// Initialize private fields and start subroutines, etc.
	// https://godoc.org/github.com/tendermint/tendermint/types#pkg-constants
	// ces.wsOut, _ = ces.rpcClient.Subscribe(ces.wsCtx, "new block", "tm.event = 'NewBlock'", 1)
	ces.wsOut, _ = ces.rpcClient.Subscribe(ces.wsCtx, "new tx", "tm.event = 'Tx'", 1)
	// ces.wsOut, _ = ces.rpcClient.Subscribe(ces.wsCtx, "new tx", "tm.event = 'ValidatorSetUpdates'", 1) // Works

	ces.startSubscription()

	return nil
}

// OnStop is an override method for BaseService, which stops a service
func (ces *ChainExporterService) OnStop() {
	ces.BaseService.OnStop()
	ces.rpcClient.OnStop()
}

// sync synchronizes the block data from connected full node
// func (ces *ChainExporterService) sync() error {
// 	var blocks []dtypes.BlockInfo
// 	err := ces.db.Model(&blocks).
// 		Order("height DESC").
// 		Limit(1).
// 		Select()
// 	if err != nil {
// 		return err
// 	}

// 	currentHeight := int64(1)
// 	if len(blocks) > 0 {
// 		currentHeight = blocks[0].Height
// 	}

// 	// query current height
// 	status, err := ces.rpcClient.Status()
// 	if err != nil {
// 		return err
// 	}
// 	maxHeight := status.SyncInfo.LatestBlockHeight

// 	fmt.Println(maxHeight)

// 	if currentHeight == 1 {
// 		currentHeight = 0
// 	}

// 	// ingest all blocks up to the best height
// 	for i := currentHeight + 1; i <= maxHeight; i++ {
// 		err = ces.process(i)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Printf("synced block %d/%d \n", i, maxHeight)
// 	}
// 	return nil
// }

// sync queries the block at the given height-1 from the node and ingests its metadata (blockinfo,evidence)
// into the database. It also queries the next block to access the commits and stores the missed signatures.
// func (ces *ChainExporterService) process(height int64) error {

// 	transactionInfo, voteInfo, depositInfo, proposalInfo, validatorSetInfo, err := ces.getTransactionInfo(height)
// 	if err != nil {
// 		return err
// 	}

// 	// insert data in PostgreSQL database
// 	err = ces.db.RunInTransaction(func(tx *pg.Tx) error {
// 		if len(transactionInfo) > 0 {
// 			err = tx.Insert(&transactionInfo)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})

// 	// roll back
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
