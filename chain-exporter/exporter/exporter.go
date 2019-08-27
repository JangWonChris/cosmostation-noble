package exporter

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/databases"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/lcd"
	dtypes "github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/utils"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/go-pg/pg"
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
		codec:     utils.MakeCodec(), // register Cosmos SDK codecs
		config:    config,
		db:        databases.ConnectDatabase(config), // connect to PostgreSQL
		wsCtx:     context.Background(),
		rpcClient: client.NewHTTP(config.Node.GaiadURL, "/websocket"), // connect to Tendermint RPC client
	}

	// setup database schema
	databases.CreateSchema(ces.db)

	// sets timeout for request.
	resty.SetTimeout(5 * time.Second)
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // test locally

	return ces
}

// OnStart is an override method for BaseService, which starts a service
func (ces *ChainExporterService) OnStart() error {
	// OnStart rpc client
	ces.rpcClient.OnStart()

	// Store data initially
	lcd.SaveBondedValidators(ces.db, ces.config)
	lcd.SaveUnbondingValidators(ces.db, ces.config)
	lcd.SaveUnbondedValidators(ces.db, ces.config)
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
		case msg2 := <-c1:
			fmt.Println("start - ", msg2)
			lcd.SaveBondedValidators(ces.db, ces.config)
			lcd.SaveUnbondingValidators(ces.db, ces.config)
			lcd.SaveUnbondedValidators(ces.db, ces.config)
			lcd.SaveProposals(ces.db, ces.config)
			fmt.Println("finish - ", msg2)
		case msg3 := <-c2:
			fmt.Println("start - ", msg3)
			ces.SaveValidatorKeyBase()
			fmt.Println("finish - ", msg3)
		}
	}
}

// OnStop is an override method for BaseService, which stops a service
func (ces *ChainExporterService) OnStop() {
	ces.rpcClient.OnStop()
}

// sync synchronizes the block data from connected full node
func (ces *ChainExporterService) sync() error {
	// Check current height in db
	var blocks []dtypes.BlockInfo
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

	// Query the node for its height
	status, err := ces.rpcClient.Status()
	if err != nil {
		return err
	}
	maxHeight := status.SyncInfo.LatestBlockHeight

	fmt.Println(maxHeight)

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

	genesisValidatorsInfo, missInfo, accumMissInfo, missDetailInfo, err := ces.getValidatorSetInfo(height)
	if err != nil {
		return err
	}

	transactionInfo, voteInfo, depositInfo, proposalInfo, validatorSetInfo, err := ces.getTransactionInfo(height)
	if err != nil {
		return err
	}

	// Insert data in PostgreSQL database
	err = ces.db.RunInTransaction(func(tx *pg.Tx) error {
		if len(blockInfo) > 0 {
			err = tx.Insert(&blockInfo)
			if err != nil {
				return err
			}
		}

		if len(genesisValidatorsInfo) > 0 {
			err = tx.Insert(&genesisValidatorsInfo)
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
