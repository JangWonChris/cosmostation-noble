package server

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"time"

	"github.com/cosmostation/cosmostation-cosmos/api/wallet/app/chain-exporter/types"
)

// tx 동기화&pulling
// 텐더민트와 블록 높이가 동기화되어 있지 않으면 pulling
// 높이가 같아지면 catch up 상태가 되고 웹소켓으로 블록을 실시간
type SubServer struct {
	cmn.BaseService
	IsCaughtUp bool
	RPCClient  *client.HTTP
	WSCtx      context.Context
	WSOut      <-chan ctypes.ResultEvent
	Subscriber types.Subscriber
}

func NewSubServer(logger log.Logger, subscriber types.Subscriber, remote, wsEndpoint string) *SubServer {
	subServer := &SubServer{}
	subServer.Subscriber = subscriber
	subServer.RPCClient = client.NewHTTP(remote, wsEndpoint)
	subServer.RPCClient.SetLogger(logger)
	subServer.RPCClient.WSEvents.SetLogger(logger)
	subServer.BaseService = *cmn.NewBaseService(logger, "Subscribe-Server", subServer)

	return subServer
}

func (s *SubServer) OnStart() error {
	err := s.BaseService.OnStart()
	if err != nil {
		return err
	}

	err = s.RPCClient.OnStart()
	if err != nil {
		return err
	}

	s.WSCtx = context.Background()
	//s.wsOut = make(chan interface{})
	// ctypes.ResultEvent 타입의 채널을 생성한다.
	s.WSOut = make(chan ctypes.ResultEvent)
	s.Subscriber.CreateIndex(s.WSCtx)

	//cannot assign <-chan ctypes.ResultEvent to s.wsOut (type chan ctypes.ResultEvent) in multiple assignment less... (⌘F1)
	//Inspection info: Reports incompatible types in binary and unary expressions.
	s.WSOut, err = s.RPCClient.Subscribe(s.WSCtx, "new block", "tm.event = 'NewBlock'", 1)
	if err != nil {
		return err
	}

	// 비동기로 실행
	go s.routine()
	return nil
}

// Todo : txid parsing,insert 하던 중간에 동기화 멈췄다가 다시 실행할 경우, tx_index에 해당블록높이의 txs 중 어디까지 들어갔는지 체크하는 로직이 필요하겠다.

func (s *SubServer) routine() {
	s.IsCaughtUp = false
	height, err := s.Subscriber.GetCurrHeightFromES(s.WSCtx)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error on GetCurrHeightFromES: %s", err.Error()), "height", height)
	}
	s.Subscriber.SetCurrentHeight(height)
	for {
		select {
		case i, ok := <-s.WSOut:
			if ok {
				s.wsNewBlockRoutine(i)
			}
		default:
			s.syncRoutine()
		}
	}
}

func (s *SubServer) syncRoutine() {

	if s.IsCaughtUp {
		time.Sleep(100 * time.Microsecond)
		return
	}

	height, err := s.Subscriber.GetNextHeight(s.WSCtx)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error on GetNextHeight: %s", err.Error()), "height", height)
		time.Sleep(1 * time.Second)
		return
	}

	block, err := s.RPCClient.Block(&height)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error on pulling: %s", err.Error()), "height", height)
		time.Sleep(1 * time.Second)
		return
	}

	for index, tx := range block.Block.Txs {
		var TxType sdk.Tx
		// Use tx codec to unmarshal binary length prefix
		err := s.Subscriber.UnmarshalBinaryLengthPrefixed([]byte(tx), &TxType)
		if err != nil {
			time.Sleep(1 * time.Second)
			return
		}

		txJson, err := s.Subscriber.MarsharJson(TxType)
		if err != nil {
			time.Sleep(1 * time.Second)
			return
		}

		// Transaction Hash
		txByte := crypto.Sha256(tx)
		txHash := hex.EncodeToString(txByte)

		fmt.Println("syncHeight : ", block.Block.Height)
		fmt.Println("syncTxHash : ", index, txHash)

		txResult, err := s.RPCClient.Tx([]byte(tx.Hash()), false)
		if err != nil {
			time.Sleep(1 * time.Second)
			return
		}

		var tags []types.Tag
		for _, tag := range txResult.TxResult.Tags {
			tagBytes, err := json.Marshal(tag)
			if err != nil {
				time.Sleep(1 * time.Second)
				return
			}
			var tagObj sdk.KVPair
			err = json.Unmarshal(tagBytes, &tagObj)

			tags = append(tags, types.Tag{string(tagObj.Key), string(tagObj.Value)})
		}

		tagsBytes, err := json.Marshal(tags)
		if err != nil {
			time.Sleep(1 * time.Second)
			return
		}

		elasticTx := &types.ElasticsearchTxInfo{
			Hash:   txHash,
			Height: block.Block.Height,
			Time:   block.Block.Time,
			Tx:     txJson,
			Result: &types.TxResultInfo{
				GasWanted: txResult.TxResult.GasWanted,
				GasUsed:   txResult.TxResult.GasUsed,
				Log:       json.RawMessage(txResult.TxResult.Log),
				Tags:      json.RawMessage(tagsBytes),
			},
		}
		s.Subscriber.Commit(s.WSCtx, elasticTx)
	}

	s.Subscriber.SetCurrentHeight(block.Block.Height)
}

func (s *SubServer) wsNewBlockRoutine(i ctypes.ResultEvent) {

	resultEvent := ctypes.ResultEvent{}
	resultEvent = i

	byteValue, err := json.Marshal(resultEvent.Data)
	if err != nil {
		return
	}

	var newBlock tmtypes.EventDataNewBlock
	json.Unmarshal(byteValue, &newBlock)

	nextHeight, err := s.Subscriber.GetNextHeight(s.WSCtx)
	if err != nil {
		time.Sleep(1 * time.Second)
		return
	}
	// isCaughtUp이 ture일 때만, websocket으로 받는 것으로!
	if newBlock.Block.Height == nextHeight {
		if s.IsCaughtUp == false {
			s.IsCaughtUp = true
			s.Logger.Info("Catch up!", "height", newBlock.Block.Height)
		}
	} else {
		fmt.Println("===== 동기화 아직 NONO =====")
		s.Logger.Info("Get block from ws, but not expected", "expected", nextHeight, "but", newBlock.Block.Height)
		if s.IsCaughtUp {
			s.IsCaughtUp = false
			s.Logger.Info("Lost catching up", "expected", nextHeight, "but", newBlock.Block.Height)
		}
		return
	}

	for index, tx := range newBlock.Block.Txs {
		var TxType sdk.Tx
		// Use tx codec to unmarshal binary length prefix
		err := s.Subscriber.UnmarshalBinaryLengthPrefixed([]byte(tx), &TxType)
		if err != nil {
			time.Sleep(1 * time.Second)
			return
		}

		txJson, err := s.Subscriber.MarsharJson(TxType)
		if err != nil {
			time.Sleep(1 * time.Second)
			return
		}

		// Transaction Hash
		txByte := crypto.Sha256(tx)
		txHash := hex.EncodeToString(txByte)

		fmt.Println("syncHeight : ", newBlock.Block.Height)
		fmt.Println("syncTxHash : ", index, txHash)

		txResult, err := s.RPCClient.Tx([]byte(tx.Hash()), false)
		if err != nil {
			time.Sleep(1 * time.Second)
			return
		}

		var tags []types.Tag
		for _, tag := range txResult.TxResult.Tags {
			tagBytes, err := json.Marshal(tag)
			if err != nil {
				time.Sleep(1 * time.Second)
				return
			}
			var tagObj sdk.KVPair
			err = json.Unmarshal(tagBytes, &tagObj)

			tags = append(tags, types.Tag{string(tagObj.Key), string(tagObj.Value)})
		}

		tagsBytes, err := json.Marshal(tags)
		if err != nil {
			time.Sleep(1 * time.Second)
			return
		}

		elasticTx := &types.ElasticsearchTxInfo{
			Hash:   txHash,
			Height: newBlock.Block.Height,
			Time:   newBlock.Block.Time,
			Tx:     txJson,
			Result: &types.TxResultInfo{
				GasWanted: txResult.TxResult.GasWanted,
				GasUsed:   txResult.TxResult.GasUsed,
				Log:       json.RawMessage(txResult.TxResult.Log),
				Tags:      json.RawMessage(tagsBytes),
			},
		}
		s.Subscriber.Commit(s.WSCtx, elasticTx)
	}

	s.Subscriber.SetCurrentHeight(newBlock.Block.Height)
}

func (s *SubServer) OnStop() {
	s.BaseService.OnStop()
	s.RPCClient.OnStop()

	if err := s.Subscriber.Stop(); err != nil {
		s.Logger.Error(fmt.Sprintf("Err on stop subscriber: %s", err.Error()))
	}
}
