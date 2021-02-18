package exporter

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"

	//internal
	"github.com/cosmostation/cosmostation-cosmos/chain-config/custom"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/client"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/db"

	// mbl
	"github.com/cosmostation/mintscan-backend-library/config"
	"github.com/cosmostation/mintscan-backend-library/db/schema"
	"github.com/cosmostation/mintscan-backend-library/types"

	// sdk
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
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
	rawdb  *db.RawDatabase
}

// NewExporter returns new Exporter instance
func NewExporter() *Exporter {
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)
	defer l.Sync()

	config := config.ParseConfig()

	client := client.NewClient(&config.Client)

	database := db.Connect(&config.DB)
	err := database.Ping()
	if err != nil {
		zap.L().Error("failed to ping database", zap.Error(err))
		return &Exporter{}
	}

	rawdb := db.RawDBConnect(&config.RAWDB)
	err = rawdb.Ping()
	if err != nil {
		zap.L().Error("failed to ping database", zap.Error(err))
		return &Exporter{}
	}

	// Create database tables if not exist already
	database.CreateTablesAndIndexes()
	rawdb.CreateTablesAndIndexes()

	return &Exporter{config, client, database, rawdb}
}

// Start starts to synchronize blockchain data
func (ex *Exporter) Start(initialHeight int64, op int) {
	zap.S().Info("Starting Chain Exporter...")
	zap.S().Infof("Version: %s | Commit: %s", Version, Commit)

	//close grpc
	// defer ex.client.Close()

	tick7Sec := time.NewTicker(time.Second * 7)
	tick5Min := time.NewTicker(time.Minute * 5)
	tick20Min := time.NewTicker(time.Minute * 20)

	done := make(chan struct{})
	defer close(done)

	go func() {
		for {
			zap.S().Info("start - sync blockchain")
			err := ex.sync(initialHeight, op)
			if err != nil {
				zap.S().Infof("error - sync blockchain: %s\n", err)
			}
			zap.S().Info("finish - sync blockchain")

			time.Sleep(time.Second)
		}
	}()

	if op == BASIC_MODE {
		go func() {
			ex.saveValidators()
			ex.saveProposals()
			for {
				select {
				case <-tick7Sec.C:
					zap.S().Infof("start sync governance and validators")
					ex.saveValidators()
					ex.saveProposals()
					zap.S().Infof("finish sync governance and validators")
				case <-tick5Min.C:
					ex.SaveStatsMarket5M()
					zap.S().Info("successfully saved market data @every 5m ")
				case <-tick20Min.C:
					zap.S().Infof("start sync validators keybase identities")
					ex.saveValidatorsIdentities()
					zap.S().Infof("finish sync validators keybase identities")
				case <-done:
					return
				}
			}
		}()
	}

	<-done // implement gracefully shutdown when signal received
	zap.S().Infof("shutdown signal received")
}

func (ex *Exporter) ReproducePowerEventHistory(op int) error {
	// Query latest block height saved in database
	dbHeight, err := ex.db.QueryLatestBlockHeight()
	if dbHeight == -1 {
		return fmt.Errorf("unexpected error in database: %s", err)
	}

	zap.S().Infof("dst db %d \n", dbHeight)

	i := int64(1)
	endHeight := int64(200)
	for ; i <= dbHeight; i += endHeight {
		zap.S().Info("working height : ", i)
		rawTxs, err := ex.db.QueryTxForPowerEventHistory(i, i+endHeight) // 1 <= x < 201, 201 <= x < 201+200
		if err != nil {
			return err
		}

		rawTxsLen := len(rawTxs)
		if rawTxsLen > 0 {

			txs := make([]*sdktypes.TxResponse, len(rawTxs))
			for i := range rawTxs {
				zap.S().Info("height : ", i, ", num_txs : ")
				tx := new(sdktypes.TxResponse)
				if err := custom.AppCodec.UnmarshalJSON([]byte(rawTxs[i].Chunk), tx); err != nil {
					return err
				}
				txs[i] = tx
			}

			exportData := new(schema.ExportData)
			exportData.ResultValidatorsPowerEventHistory, err = ex.getPowerEventHistoryNew(txs)
			if err != nil {
				return err
			}
			return ex.db.InsertExportedData(exportData)

		}

	}
	return nil

}

func (ex *Exporter) Refine(op int) error {
	// Query latest block height saved in database
	srcDBHeight, err := ex.rawdb.QueryLatestBlockHeight()
	if srcDBHeight == -1 {
		return fmt.Errorf("unexpected error in database: %s", err)
	}
	dstDBHeight, err := ex.db.QueryLatestBlockHeight()
	if dstDBHeight == -1 {
		return fmt.Errorf("unexpected error in database: %s", err)
	}

	zap.S().Infof("src db %d, dst db %d \n", srcDBHeight, dstDBHeight)

	for i := dstDBHeight + 1; i <= srcDBHeight; {
		zap.S().Info("working height : ", i)
		bs, err := ex.rawdb.GetBlocks(i)
		if err != nil {
			return err
		}
		for _, b := range bs {
			block := new(tmctypes.ResultBlock)
			if err := json.Unmarshal([]byte(b.Chunk), &block); err != nil {
				return err
			}
			var txs []*sdktypes.TxResponse
			if b.NumTxs != 0 {
				txs = make([]*sdktypes.TxResponse, b.NumTxs)
				zap.S().Info("height : ", b.Height, ", num_txs : ", b.NumTxs)
				ts, err := ex.rawdb.GetTransactions(b.Height)
				if err != nil {
					return err
				}
				for i, t := range ts {
					tx := new(sdktypes.TxResponse)
					if err := custom.AppCodec.UnmarshalJSON([]byte(t.Chunk), tx); err != nil {
						return err
					}
					txs[i] = tx
				}

			}
			ex.process(block, txs, op)
		}

		i += int64(len(bs))
	}
	return nil

}

// sync compares block height between the height saved in your database and
// the latest block height on the active chain and calls process to start ingesting data.
func (ex *Exporter) sync(initialHeight int64, op int) error {
	// Query latest block height saved in database
	dbHeight, err := ex.db.QueryLatestBlockHeight()
	if dbHeight == -1 {
		return fmt.Errorf("unexpected error in database: %s", err)
	}
	rawDBHeight, err := ex.rawdb.QueryLatestBlockHeight()
	if rawDBHeight == -1 {
		return fmt.Errorf("unexpected error in database: %s", err)
	}

	// Query latest block height on the active network
	latestBlockHeight, err := ex.client.RPC.GetLatestBlockHeight()
	if latestBlockHeight == -1 {
		return fmt.Errorf("failed to query the latest block height on the active network: %s", err)
	}

	if dbHeight == 0 && initialHeight != 0 {
		dbHeight = initialHeight - 1
		rawDBHeight = initialHeight - 1
		zap.S().Info("initial Height set : ", initialHeight)
	}

	beginHeight := dbHeight
	if dbHeight > rawDBHeight || op == RAW_MODE {
		beginHeight = rawDBHeight
	}

	zap.S().Infof("dbHeight %d, rawHeight %d \n", dbHeight, rawDBHeight)

	for i := beginHeight + 1; i <= latestBlockHeight; i++ {
		block, err := ex.client.RPC.GetBlock(i)
		if err != nil {
			return fmt.Errorf("failed to query block: %s", err)
		}
		txs, err := ex.client.CliCtx.GetTxs(block)
		if err != nil {
			return fmt.Errorf("failed to get transactions for block: %s", err)
		}
		switch op {
		case BASIC_MODE:
			if i > dbHeight {
				err = ex.process(block, txs, op)
				if err != nil {
					return err
				}
			}
			fallthrough //continue to case RAW_MODE
		case RAW_MODE:
			if i > rawDBHeight {
				err = ex.rawProcess(block, txs)
				if err != nil {
					return err
				}
			}
		case REFINE_MODE:
		default:
			zap.S().Info("unknown mode = ", op)
			os.Exit(1)
		}
		zap.S().Infof("synced block %d/%d", i, latestBlockHeight)
	}
	return nil
}

func (ex *Exporter) rawProcess(block *tmctypes.ResultBlock, txs []*sdktypes.TxResponse) (err error) {
	exportRawData := new(schema.ExportRawData)

	exportRawData.ResultBlockJSONChunk, err = ex.getBlockJSONChunk(block)
	if err != nil {
		return fmt.Errorf("failed to get block: %s", err)
	}
	exportRawData.ResultTxsJSONChunk, err = ex.getTxsJSONChunk(block, txs)
	if err != nil {
		return fmt.Errorf("failed to get txs: %s", err)
	}
	return ex.rawdb.InsertExportedData(exportRawData)
}

// process ingests chain data, such as block, transaction, validator, evidence information and
// save them in database.
func (ex *Exporter) process(block *tmctypes.ResultBlock, txs []*sdktypes.TxResponse, op int) (err error) {
	exportData := new(schema.ExportData)

	exportData.ResultBlock, err = ex.getBlock(block)
	if err != nil {
		return fmt.Errorf("failed to get block: %s", err)
	}

	if op != REFINE_MODE {
		exportData.ResultEvidence, err = ex.getEvidence(block)
		if err != nil {
			return fmt.Errorf("failed to get evidence: %s", err)
		}

		if block.Block.LastCommit.Height != 0 {
			prevBlock, err := ex.client.RPC.GetBlock(block.Block.LastCommit.Height)
			if err != nil {
				return fmt.Errorf("failed to query previous block: %s", err)
			}

			vals, err := ex.client.RPC.GetValidatorsInHeight(block.Block.LastCommit.Height, types.DefaultQueryValidatorsPage, types.DefaultQueryValidatorsPerPage)
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
	}

	if exportData.ResultBlock.NumTxs > 0 {
		if op != REFINE_MODE {
			// exportData.ResultAccounts, err = ex.getAccounts(block, txs)
			// if err != nil {
			// 	return fmt.Errorf("failed to get accounts: %s", err)
			// }
			// exportData.ResultAccountCoin, err = ex.getAccounts(block, txs)
			// if err != nil {
			// 	return fmt.Errorf("failed to get accounts: %s", err)
			// }
			exportData.ResultProposals, exportData.ResultDeposits, exportData.ResultVotes, err = ex.getGovernance(block, txs)
			if err != nil {
				return fmt.Errorf("failed to get governance: %s", err)
			}
			// exportData.ResultValidatorsPowerEventHistory, err = ex.getPowerEventHistory(block, txs)
			exportData.ResultValidatorsPowerEventHistory, err = ex.getPowerEventHistoryNew(txs)
			if err != nil {
				return fmt.Errorf("failed to get transactions: %s", err)
			}
		}
		exportData.ResultTxs, err = ex.getTxs(block, txs)
		if err != nil {
			return fmt.Errorf("failed to get txs: %s", err)
		}
		exportData.ResultTxsAccount, err = ex.transactionAccount(block.Block.ChainID, txs)
		if err != nil {
			return fmt.Errorf("failed to get account by each tx message: %s", err)
		}
	}

	// TODO: is this right place to be?
	if ex.config.Alarm.Switch {
		ex.handlePushNotification(block, txs)
	}

	return ex.db.InsertExportedData(exportData)
}
