package exporter

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/databases"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/lcd"
	dtypes "github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	gaiaApp "github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	"github.com/cosmos/cosmos-sdk/codec"

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

// Monitor wraps Tendermint RPC client and PostgreSQL database
type ChainExporterService struct {
	cmn.BaseService
	Codec     *codec.Codec
	Config    *config.Config
	DB        *pg.DB
	WsCtx     context.Context
	WsOut     <-chan ctypes.ResultEvent
	RPCClient *client.HTTP
}

// Initializes all the required configs
func NewChainExporterService(config *config.Config) *ChainExporterService {
	ces := &ChainExporterService{
		Codec:     gaiaApp.MakeCodec(), // Register Cosmos SDK codecs
		Config:    config,
		DB:        databases.ConnectDatabase(config), // Connect to PostgreSQL
		WsCtx:     context.Background(),
		RPCClient: client.NewHTTP(config.Node.GaiadURL, "/websocket"), // Connect to Tendermint RPC client
	}
	// Setup database schema
	databases.CreateSchema(ces.DB)

	// Register a service that can be started, stopped, and reset
	ces.BaseService = *cmn.NewBaseService(logger, "ChainExporterService", ces)

	// SetTimeout method sets timeout for request.
	resty.SetTimeout(5 * time.Second)
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // Test locally

	return ces
}

// Override method for BaseService, which starts a service
func (ces *ChainExporterService) OnStart() error {
	// OnStart both service and rpc client
	ces.BaseService.OnStart()
	ces.RPCClient.OnStart()

	// Initialize private fields and start subroutines, etc.
	// ces.WsOut, _ = ces.RPCClient.Subscribe(ces.WsCtx, "new block", "tm.event = 'NewBlock'", 1)
	ces.WsOut, _ = ces.RPCClient.Subscribe(ces.WsCtx, "new tx", "tm.event = 'Tx'", 1)

	// Store data initially
	lcd.SaveGovernance(ces.DB, ces.Config)
	lcd.SaveBondedValidators(ces.DB, ces.Config)
	lcd.SaveUnbondedAndUnbodingValidators(ces.DB, ces.Config)

	// Start the syncing task
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

	// Allow graceful closing of the governance loop
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	for {
		select {
		case <-time.Tick(7 * time.Second):
			fmt.Println("start - sync LCD governance & validators")
			lcd.SaveGovernance(ces.DB, ces.Config)
			lcd.SaveBondedValidators(ces.DB, ces.Config)
			lcd.SaveUnbondedAndUnbodingValidators(ces.DB, ces.Config)
			fmt.Println("finish - sync LCD governance & validators")
		case <-signalCh:
			return nil
			// Push Notification 을 위함
			// case eventData, ok := <-ces.WsOut:
			// 	fmt.Println("start new tx subscription from full node")
			// 	if ok {
			// 		fmt.Println("===Event===================================================")
			// 		fmt.Println(eventData)
			// 	}
			// 	fmt.Println("finish tx subscription from full node")
			// case <-signalCh:
			// 	return nil
		}
	}
}

// Override method for BaseService, which stops a service
func (ces *ChainExporterService) OnStop() {
	ces.BaseService.OnStop()
	ces.RPCClient.OnStop()
}

// Synchronizes the block data from connected full node
func (ces *ChainExporterService) sync() error {
	// Check current height in db
	var blocks []dtypes.BlockInfo
	err := ces.DB.Model(&blocks).
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

	// Query the node for its height
	status, err := ces.RPCClient.Status()
	if err != nil {
		return err
	}
	maxHeight := status.SyncInfo.LatestBlockHeight

	if currentHeight == 1 {
		currentHeight = 0
	}

	// Ingest all blocks up to the best height
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
func (ces *ChainExporterService) process(height int64) error {
	blockInfo, err := ces.getBlockInfo(height)
	if err != nil {
		return err
	}

	evidenceInfo, err := ces.getEvidenceInfo(height)
	if err != nil {
		return err
	}

	validatorSetInfo, missInfo, accumMissInfo, missDetailInfo, err := ces.getValidatorSetInfo(height)
	if err != nil {
		return err
	}

	transactionInfo, voteInfo, depositInfo, proposalInfo, err := ces.getTransactionInfo(height)
	if err != nil {
		return err
	}

	// Insert data in PostgreSQL database
	err = ces.DB.RunInTransaction(func(tx *pg.Tx) error {
		if len(blockInfo) > 0 {
			err = tx.Insert(&blockInfo)
			if err != nil {
				return err
			}
		}

		if len(validatorSetInfo) > 0 {
			err = tx.Insert(&validatorSetInfo)
			if err != nil {
				return err
			}
		}

		if len(evidenceInfo) > 0 {
			err = tx.Insert(&evidenceInfo)
			if err != nil {
				return err
			}
		}

		if len(missInfo) > 0 {
			err = tx.Insert(&missInfo)
			if err != nil {
				return err
			}
		}

		if len(missDetailInfo) > 0 {
			err = tx.Insert(&missDetailInfo)
			if err != nil {
				return err
			}
		}

		if len(transactionInfo) > 0 {
			err = tx.Insert(&transactionInfo)
			if err != nil {
				return err
			}
		}

		if len(depositInfo) > 0 {
			err = tx.Insert(&depositInfo)
			if err != nil {
				return err
			}
		}

		// Update accumulative missing block info
		var tempMissInfo dtypes.MissInfo
		if len(accumMissInfo) > 0 {
			for i := 0; i < len(accumMissInfo); i++ {
				_, err = tx.Model(&tempMissInfo).
					Set("address = ?", accumMissInfo[i].Address).
					Set("start_height = ?", accumMissInfo[i].StartHeight).
					Set("end_height = ?", accumMissInfo[i].EndHeight).
					Set("missing_count = ?", accumMissInfo[i].MissingCount).
					Set("start_time = ?", accumMissInfo[i].StartTime).
					Set("end_time = ?", blockInfo[0].Time).
					Where("end_height = ? AND address = ?", accumMissInfo[i].EndHeight-int64(1), accumMissInfo[i].Address).
					Update()
				if err != nil {
					return err
				}
			}
		}

		// Insert vote tx info
		if len(voteInfo) > 0 {
			var tempVoteInfo dtypes.VoteInfo
			for i := 0; i < len(voteInfo); i++ {
				// Check if a validator already voted
				count, _ := tx.Model(&tempVoteInfo).
					Where("proposal_id = ? AND voter = ?", voteInfo[i].ProposalID, voteInfo[i].Voter).
					Count()
				if count > 0 {
					_, err = tx.Model(&tempVoteInfo).
						Set("height = ?", voteInfo[i].Height).
						Set("option = ?", voteInfo[i].Option).
						Set("tx_hash = ?", voteInfo[i].TxHash).
						Set("gas_wanted = ?", voteInfo[i].GasWanted).
						Set("gas_used = ?", voteInfo[i].GasUsed).
						Set("time = ?", voteInfo[i].Time).
						Where("proposal_id = ? AND voter = ?", voteInfo[i].ProposalID, voteInfo[i].Voter).
						Update()
					if err != nil {
						return err
					}
				} else {
					err = tx.Insert(&voteInfo)
					if err != nil {
						return err
					}
				}
			}
		}

		// Exist and update proposerInfo
		if len(proposalInfo) > 0 {
			var tempProposalInfo dtypes.ProposalInfo
			for i := 0; i < len(proposalInfo); i++ {
				// Check if a validator already voted
				count, _ := tx.Model(&tempProposalInfo).
					Where("id = ?", proposalInfo[i].ID).
					Count()

				if count > 0 {
					// Save and update proposalInfo
					_, err = tx.Model(&tempProposalInfo).
						Set("tx_hash = ?", proposalInfo[i].TxHash).
						Set("proposer = ?", proposalInfo[i].Proposer).
						Set("initial_deposit_amount = ?", proposalInfo[i].InitialDepositAmount).
						Set("initial_deposit_denom = ?", proposalInfo[i].InitialDepositDenom).
						Where("id = ?", proposalInfo[i].ID).
						Update()
					if err != nil {
						return err
					}
				} else {
					err = tx.Insert(&proposalInfo)
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	})

	// Roll back
	if err != nil {
		return err
	}

	return nil
}
