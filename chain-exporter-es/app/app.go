package app

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmostation-cosmos/chain-exporter-es/elastic"
	"github.com/cosmostation-cosmos/chain-exporter-es/model"
)

type App struct {
	cmn.BaseService
	IsCaughtUp    bool
	Client        *client.HTTP
	WsCtx         context.Context
	WsOut         <-chan ctypes.ResultEvent
	ElasticSearch *elastic.ElasticSearch
}

func (a *App) NewContext() *Context {
	return &Context{
		Logger:        log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		ElasticSearch: a.ElasticSearch,
	}
}

func NewApp(network string, env string) (app *App, err error) {
	app = &App{}

	elasticConfig, err := elastic.InitConfig(network, env)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	appConfig, err := InitConfig(network, env)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	appContext := app.NewContext()

	app.ElasticSearch, err = elastic.NewElastic(elasticConfig)
	app.Client = client.NewHTTP(appConfig.RPCEndPoint, "/websocket")
	app.Client.SetLogger(appContext.Logger)
	app.Client.WSEvents.SetLogger(appContext.Logger)
	app.BaseService = *cmn.NewBaseService(appContext.Logger, "ES Crawler Server", app)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return app, nil
}

func (a *App) OnStart() error {
	err := a.BaseService.OnStart()
	if err != nil {
		return err
	}

	err = a.Client.OnStart()
	if err != nil {
		return err
	}

	a.WsCtx = context.Background()
	// ctypes.ResultEvent 타입의 채널을 생성한다.
	a.WsOut = make(chan ctypes.ResultEvent)
	a.ElasticSearch.CreateIndex(a.WsCtx)
	a.WsOut, err = a.Client.Subscribe(a.WsCtx, "new block", "tm.event = 'NewBlock'", 1)
	if err != nil {
		return err
	}

	// 고루틴 실행
	go a.routine()
	return nil
}

// Todo : txid parsing,insert 하던 중간에 동기화 멈췄다가 다시 실행할 경우, tx_index에 해당블록높이의 txs 중 어디까지 들어갔는지 체크하는 로직이 필요하겠다.
func (a *App) routine() {
	a.IsCaughtUp = false
	height, err := a.ElasticSearch.GetCurrHeight(a.WsCtx)
	a.Logger.Info("sync start at height : ", height)
	if err != nil {
		a.Logger.Error(fmt.Sprintf("Error on GetCurrHeightFromES: %s", err.Error()), "height", height)
	}

	a.ElasticSearch.SetCurrHeight(height)
	for {
		select {
		//v, ok := <-ch
		//ok is false if there are no more values to receive and the channel is closed.
		case i, ok := <-a.WsOut:
			if ok {
				a.wsNewBlockRoutine(i)
			}
		default:
			a.syncRoutine()
		}
	}
}

func (a *App) syncRoutine() {
	if a.IsCaughtUp {
		time.Sleep(100 * time.Microsecond)
		return
	}

	height, err := a.ElasticSearch.GetNextHeight(a.WsCtx)
	if err != nil {
		a.Logger.Error(fmt.Sprintf("Error on GetNextHeight: %s", err.Error()), "height", height)
		time.Sleep(1 * time.Second)
		return
	}

	block, err := a.Client.Block(&height)
	if err != nil {
		a.Logger.Error(fmt.Sprintf("Error on pulling: %s", err.Error()), "height", height)
		time.Sleep(1 * time.Second)
		return
	}

	for index, tx := range block.Block.Txs {
		var TxType sdkTypes.Tx
		// Use tx codec to unmarshal binary length prefix
		err := a.ElasticSearch.UnmarshalBinaryLengthPrefixed([]byte(tx), &TxType)
		if err != nil {
			time.Sleep(1 * time.Second)
			return
		}

		txJson, err := a.ElasticSearch.MarsharJson(TxType)
		if err != nil {
			time.Sleep(1 * time.Second)
			return
		}

		// Transaction Hash
		txByte := crypto.Sha256(tx)
		txHash := hex.EncodeToString(txByte)

		fmt.Println("syncHeight : ", block.Block.Height)
		fmt.Println("syncTxHash : ", index, txHash)

		txResult, err := a.Client.Tx([]byte(tx.Hash()), false)
		if err != nil {
			time.Sleep(1 * time.Second)
			return
		}

		var tags []model.Tag
		for _, tag := range txResult.TxResult.Tags {
			tagBytes, err := json.Marshal(tag)
			if err != nil {
				time.Sleep(1 * time.Second)
				return
			}
			var tagObj sdkTypes.KVPair
			err = json.Unmarshal(tagBytes, &tagObj)

			tags = append(tags, model.Tag{string(tagObj.Key), string(tagObj.Value)})
		}
		tagsBytes, err := json.Marshal(tags)
		if err != nil {
			time.Sleep(1 * time.Second)
			return
		}

		elasticTx := &model.ElasticsearchTxInfo{
			Hash:   txHash,
			Height: block.Block.Height,
			Time:   block.Block.Time,
			Tx:     txJson,
			Result: &model.TxResultInfo{
				GasWanted: txResult.TxResult.GasWanted,
				GasUsed:   txResult.TxResult.GasUsed,
				Log:       json.RawMessage(txResult.TxResult.Log),
				Tags:      json.RawMessage(tagsBytes),
			},
		}
		//logrus.Info(elasticTx)
		a.ElasticSearch.InsertTx(a.WsCtx, elasticTx)
	}

	a.ElasticSearch.SetCurrHeight(block.Block.Height)
}

func (a *App) wsNewBlockRoutine(i ctypes.ResultEvent) {
	resultEvent := ctypes.ResultEvent{}
	resultEvent = i

	byteValue, err := json.Marshal(resultEvent.Data)
	if err != nil {
		return
	}
	var newBlock tmtypes.EventDataNewBlock
	json.Unmarshal(byteValue, &newBlock)

	nextHeight, err := a.ElasticSearch.GetNextHeight(a.WsCtx)
	if err != nil {
		time.Sleep(1 * time.Second)
		return
	}
	// isCaughtUp이 ture일 때만, websocket으로 받는 것으로!
	if newBlock.Block.Height == nextHeight {
		if a.IsCaughtUp == false {
			a.IsCaughtUp = true
			a.Logger.Info("Catch up!", "height", newBlock.Block.Height)
		}
	} else {
		fmt.Println("===== 동기화 아직 NONO =====")
		a.Logger.Info("Get block from ws, but not expected", "expected", nextHeight, "but", newBlock.Block.Height)
		if a.IsCaughtUp {
			a.IsCaughtUp = false
			a.Logger.Info("Lost catching up", "expected", nextHeight, "but", newBlock.Block.Height)
		}
		return
	}

	for index, tx := range newBlock.Block.Txs {
		var TxType sdkTypes.Tx
		// Use tx codec to unmarshal binary length prefix
		err := a.ElasticSearch.UnmarshalBinaryLengthPrefixed([]byte(tx), &TxType)
		if err != nil {
			time.Sleep(1 * time.Second)
			return
		}

		txJson, err := a.ElasticSearch.MarsharJson(TxType)
		if err != nil {
			time.Sleep(1 * time.Second)
			return
		}

		// Transaction Hash
		txByte := crypto.Sha256(tx)
		txHash := hex.EncodeToString(txByte)

		fmt.Println("syncHeight : ", newBlock.Block.Height)
		fmt.Println("syncTxHash : ", index, txHash)

		txResult, err := a.Client.Tx([]byte(tx.Hash()), false)
		if err != nil {
			time.Sleep(1 * time.Second)
			return
		}

		var tags []model.Tag
		for _, tag := range txResult.TxResult.Tags {
			tagBytes, err := json.Marshal(tag)
			if err != nil {
				time.Sleep(1 * time.Second)
				return
			}
			var tagObj sdkTypes.KVPair
			err = json.Unmarshal(tagBytes, &tagObj)

			tags = append(tags, model.Tag{string(tagObj.Key), string(tagObj.Value)})
		}

		tagsBytes, err := json.Marshal(tags)
		if err != nil {
			time.Sleep(1 * time.Second)
			return
		}

		elasticTx := &model.ElasticsearchTxInfo{
			Hash:   txHash,
			Height: newBlock.Block.Height,
			Time:   newBlock.Block.Time,
			Tx:     txJson,
			Result: &model.TxResultInfo{
				GasWanted: txResult.TxResult.GasWanted,
				GasUsed:   txResult.TxResult.GasUsed,
				Log:       json.RawMessage(txResult.TxResult.Log),
				Tags:      json.RawMessage(tagsBytes),
			},
		}

		a.ElasticSearch.InsertTx(a.WsCtx, elasticTx)
		//logrus.Info(elasticTx)
	}
	a.ElasticSearch.SetCurrHeight(newBlock.Block.Height)
}

func (a *App) OnStop() {
	a.BaseService.OnStop()
	a.Client.OnStop()
}
