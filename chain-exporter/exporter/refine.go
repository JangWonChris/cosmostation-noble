package exporter

import (
	"fmt"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmostation/cosmostation-cosmos/chain-config/custom"
	"github.com/cosmostation/mintscan-database/schema"
	"go.uber.org/zap"
)

func (ex *Exporter) Refine(op int) error {
	var chainID string
	aminoUnmarshal := custom.AppCodec.UnmarshalJSON

	if true { // 블록 가공 시작
		rawBlockIDMax, err := ex.rawdb.QueryBlockIDMax()
		if rawBlockIDMax == -1 {
			return fmt.Errorf("fail to get block max(id) in database: %s", err)
		}

		zap.S().Infof("total count of raw blocks : %d\n", rawBlockIDMax)
		for i := int64(1); i <= rawBlockIDMax; i++ {
			zap.S().Info("block working id : ", i)
			rb, err := ex.rawdb.GetBlockByID(i)
			if err != nil {
				return err
			}

			if rb == nil {
				return nil
			}
			i = rb[len(rb)-1].ID // 최종적으로 받아온 id 값으로 다음 loop를 시작한다.

			if err := ex.refineRawBlocks(rb); err != nil {
				return err
			}
		}
	}

	if true { // 트랜잭션 가공 시작
		rawTxIDMax, err := ex.rawdb.QueryTxIDMax()
		if rawTxIDMax == -1 {
			return fmt.Errorf("fail to get transaction max(id) in database: %s", err)
		}
		zap.S().Infof("total count of raw_transaction : %d\n", rawTxIDMax)
		for i := int64(1); i <= rawTxIDMax; i++ {
			zap.S().Info("transaction working id : ", i)

			ts, err := ex.rawdb.GetTransactionsByID(i)
			if err != nil {
				return err
			}
			if ts == nil {
				return nil
			}
			txs := make([]*sdktypes.TxResponse, len(ts))
			for j, t := range ts {
				tx := new(sdktypes.TxResponse)
				if err := aminoUnmarshal(t.Chunk, tx); err != nil {
					return err
				}
				txs[j] = tx
				chainID = t.ChainID
				i = t.ID // 최종적으로 받아온 id 값으로 다음 loop를 시작한다.
			}

			if err := ex.refineRawTransactions(chainID, txs); err != nil {
				return err
			}
		}
	}
	return nil
}
func (ex *Exporter) refineRawTransactions(chainID string, txs []*sdktypes.TxResponse) (err error) {
	refineData := new(schema.RefineData)

	//cosmoshub-1에서는 txResponse에 timestamp를 가지고 있지 않기 때문에, 블록으로부터 가져온다.
	// 시작
	var begin, end int64
	if len(txs) > 0 {
		begin = txs[0].Height
		end = txs[len(txs)-1].Height

	}
	blockList, err := ex.getBlocksHasTxs(chainID, begin, end)
	if err != nil {
		return fmt.Errorf("failed to get block list: %s", err)
	}
	// 추가 로직 종료

	refineData.Transactions, err = ex.getTxs(chainID, blockList, txs)
	if err != nil {
		return fmt.Errorf("failed to get txs: %s", err)
	}

	refineData.SourceTransactionMessageAccounts = ex.disassembleTransaction(txs)

	if true {
		return ex.db.InsertExportedRefineData(refineData)
	}
	return fmt.Errorf("currently, disabled to store data into database\n")

}

// getBlockHasTxs()는 database로부터 chainID, 시작 블록 높이, 종료 블록 높이의 범위 내에서 []블록 슬라이스를 리턴(id, height, timestamp)를 리턴한다.
// schema.Block의 나머지 구조체는 nil 임에 주의한다.
func (ex *Exporter) getBlocksHasTxs(chainID string, begin, end int64) (list map[int64]*schema.Block, err error) {
	blocks, err := ex.db.GetBlockHasTxs(ChainIDMap[chainID], begin, end)
	if err != nil {
		return list, err
	}

	list = map[int64]*schema.Block{} // 초기화

	for i := range blocks {
		list[blocks[i].Height] = &blocks[i]
	}
	return list, nil
}

func (ex *Exporter) refineRawBlocks(block []schema.RawBlock) (err error) {
	refineData := new(schema.RefineData)

	refineData.Block, err = ex.getBlockFromDB(block)
	if err != nil {
		return fmt.Errorf("failed to get block: %s", err)
	}

	if true {
		return ex.db.InsertExportedRefineData(refineData)
	}
	return fmt.Errorf("currently do not store any data\n")

}
