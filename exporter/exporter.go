package exporter

import (
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"

	//internal

	"github.com/cosmostation/cosmostation-cosmos/client"
	"github.com/cosmostation/cosmostation-cosmos/db"

	// mbl
	mblconfig "github.com/cosmostation/mintscan-backend-library/config"
	mbltypes "github.com/cosmostation/mintscan-backend-library/types"
	mdschema "github.com/cosmostation/mintscan-database/schema"

	// sdk
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var (
	// Version is a project's version string.
	Version = "Development"

	// Commit is commit hash of this project.
	Commit = ""

	controler     = make(chan struct{}, 20)
	wg            = new(sync.WaitGroup)
	ChainNumMap   = map[int]string{}
	ChainIDMap    = map[string]int{}
	ChainID       string
	InitialHeight = int64(0)
)

// Exporter is
type Exporter struct {
	config *mblconfig.Config
	client *client.Client
	db     *db.Database
	rawdb  *db.RawDatabase
}

// NewExporter returns new Exporter instance
func NewExporter(op int) *Exporter {
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)
	defer l.Sync()

	fileBaseName := "chain-exporter"
	config := mblconfig.ParseConfig(fileBaseName)

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

	// switch op {
	// case REFINE_MODE:
	// 	mdschema.SetCommonSchema(config.DB.CommonSchema)
	// 	// database.CreateRefineTables()
	// case BASIC_MODE, RAW_MODE, GENESIS_MODE:
	// 	mdschema.SetCommonSchema(config.DB.CommonSchema)
	// 	// database.CreateTables()
	// 	// rawdb.CreateTables()
	// default:
	// 	zap.S().Panic("Unknow operator type :", op)
	// 	os.Exit(1)
	// }

	return &Exporter{config, client, database, rawdb}
}

// preProcess 는 실제 프로세스 수행 전, 필요한 설정 환경 등을 동적으로 설정
func (ex *Exporter) PreProcess(initialHeight int64) {

	var err error

	// chain id map 설정
	/*
		1. 현재 체인 ID를 Database에 조회 후 없으면 저장
		2. (1)이 완료되면, database에 저장된 모든 체인ID를 가져온다.
		ChainNumMap map[int]string : chainNum을 통해 chain-id를 가져옴
		ChainIDMap map[string]int : chainid를 통해 chainnum을 가져옴(chainNum은 chain-info의 id 컬럼)
	*/

	// ChainID, err := ex.client.GetNetworkChainID()
	// if err != nil {
	// 	panic(err)
	// }

	ChainID = "cosmoshub-4"
	InitialHeight = initialHeight

	exist, err := ex.db.ExistChainID(ChainID)
	if err != nil {
		panic(err)
	}

	if !exist {
		// insert db
		if err := ex.db.InsertChainID(ChainID); err != nil {
			panic(err)
		}
	}

	chainInfo, err := ex.db.QueryChainInfo()
	if err != nil {
		panic(err)
	}

	for _, c := range chainInfo {
		ChainNumMap[int(c.ID)] = c.ChainID
		ChainIDMap[c.ChainID] = int(c.ID)
	}

	// bonded denom 설정 (1회 호출하면, MBL에서 캐싱 됨)
	// ex.client.GetStakingDenom() // legacy 체인 client를 연결해놓지 않으면 문제가 될 수 있음. 최신 체인에서만 사용

	// 마지막 chain-id의 마지막 블록 넘버 설정(1회 호출해서, block id를 받아옴)
	// chaininfoid <= 1 이면, LastBlockID == 0임
	// if len(chainInfo) > 1 {
	// 	LastBlockID, err = ex.db.GetPrevChainLastBlockID(ChainIDMap[ChainID])
	// }

	// fmt.Println("LastBlockID :", LastBlockID)
	fmt.Println("ChainIDMap :", ChainIDMap)
	fmt.Println("ChainNumMap :", ChainNumMap)
}

// Start starts to synchronize blockchain data
func (ex *Exporter) Start(op int) {
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
			err := ex.sync(op)
			if err != nil {
				zap.S().Infof("error - sync blockchain: %s\n", err)
			}
			zap.S().Info("finish - sync blockchain")

			time.Sleep(time.Second)
		}
	}()

	if op == BASIC_MODE {
		go func() {
			ex.SaveStatsMarket5M()
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

// sync compares block height between the height saved in your database and
// the latest block height on the active chain and calls process to start ingesting data.
func (ex *Exporter) sync(op int) error {
	// Query latest block height saved in database
	dbHeight, err := ex.db.QueryLatestBlockHeight(ChainIDMap[ChainID])
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

	if dbHeight == 0 && InitialHeight != 0 {
		dbHeight = InitialHeight - 1
		rawDBHeight = InitialHeight - 1
		zap.S().Info("initial Height set : ", InitialHeight)
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
		retryFlag := false
		zap.S().Infof("number of Transactions : %d", len(block.Block.Txs))
		txList := block.Block.Txs
		txs := make([]*sdktypes.TxResponse, len(block.Block.Txs))

		for idx, tx := range txList {
			hex := fmt.Sprintf("%X", tx.Hash())
			controler <- struct{}{}
			wg.Add(1)
			go func(i int, hex string) {
				zap.S().Info(i, hex)
				defer func() {
					<-controler
					wg.Done()
				}()

				txs[i], err = ex.client.CliCtx.GetTx(hex)
				if err != nil {
					zap.S().Error("Error while getting tx ", hex)
					retryFlag = true
					return
				}
			}(idx, hex)
		}
		wg.Wait()

		if retryFlag {
			zap.S().Error("can not get all of tx, retry get tx in block height = ", i)
			i--
			time.Sleep(1 * time.Second)
			continue
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
	rawData := new(mdschema.RawData)

	rawData.Block, err = ex.getRawBlock(block)
	if err != nil {
		return fmt.Errorf("failed to get block: %s", err)
	}
	rawData.Transactions, err = ex.getRawTransactions(block, txs)
	if err != nil {
		return fmt.Errorf("failed to get txs: %s", err)
	}
	return ex.rawdb.InsertExportedData(rawData)
}

// process ingests chain data, such as block, transaction, validator, evidence information and
// save them in database.
func (ex *Exporter) process(block *tmctypes.ResultBlock, txs []*sdktypes.TxResponse, op int) (err error) {
	basic := new(mdschema.BasicData)

	basic.Block, err = ex.getBlock(block)
	if err != nil {
		return fmt.Errorf("failed to get block: %s", err)
	}

	basic.Evidence, err = ex.getEvidence(block)
	if err != nil {
		return fmt.Errorf("failed to get evidence: %s", err)
	}

	if block.Block.LastCommit.Height != 0 {
		prevBlock, err := ex.client.RPC.GetBlock(block.Block.LastCommit.Height)
		if err != nil {
			return fmt.Errorf("failed to query previous block: %s", err)
		}

		vals, err := ex.client.RPC.GetValidatorsInHeight(block.Block.LastCommit.Height, mbltypes.DefaultQueryValidatorsPage, mbltypes.DefaultQueryValidatorsPerPage)
		if err != nil {
			return fmt.Errorf("failed to query validators: %s", err)
		}

		basic.GenesisValidatorsSet, err = ex.getGenesisValidatorsSet(block, vals)
		if err != nil {
			return fmt.Errorf("failed to get genesis validator set: %s", err)
		}
		basic.MissBlocks, basic.AccumulatedMissBlocks, basic.MissDetailBlocks, err = ex.getValidatorsUptime(prevBlock, block, vals)
		if err != nil {
			return fmt.Errorf("failed to get missing blocks: %s", err)
		}
	}

	if basic.Block.NumTxs > 0 {
		basic.Proposals, basic.Deposits, basic.Votes, err = ex.getGovernance(block, txs)
		if err != nil {
			return fmt.Errorf("failed to get governance: %s", err)
		}
		// exportData.ValidatorsPowerEventHistory, err = ex.getPowerEventHistory(block, txs)
		basic.ValidatorsPowerEventHistory, err = ex.getPowerEventHistoryNew(txs)
		if err != nil {
			return fmt.Errorf("failed to get transactions: %s", err)
		}

		// 시작
		// block-id 추출을 위해 사용
		list := make(map[int64]*mdschema.Block)
		list[block.Block.Height] = basic.Block
		// 종료

		basic.Transactions, err = ex.getTxs(block.Block.ChainID, list, txs)
		if err != nil {
			return fmt.Errorf("failed to get txs: %s", err)
		}
		basic.TMAs = ex.disassembleTransaction(txs)
	}

	// TODO: is this right place to be?
	if ex.config.Alarm.Switch {
		ex.handlePushNotification(block, txs)
	}

	return ex.db.InsertExportedData(basic)
}
