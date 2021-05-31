package exporter

import (
	"fmt"
	"log"

	// internal
	"github.com/cosmostation/cosmostation-cosmos/chain-config/custom"

	// core
	"github.com/cosmostation/mintscan-backend-library/types"
	"github.com/cosmostation/mintscan-database/schema"

	// sdk
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// getTxs decodes transactions in a block and return a format of database transaction.
func (ex *Exporter) getTxs(chainID string, list map[int64]*schema.Block, txResp []*sdktypes.TxResponse) ([]schema.Transaction, error) {
	txs := make([]schema.Transaction, 0)

	if len(txResp) <= 0 {
		return txs, nil
	}

	codec := custom.AppCodec

	for i := range txResp {

		chunk, err := codec.MarshalJSON(txResp[i])
		if err != nil {
			return txs, fmt.Errorf("failed to marshal tx : %s", err)
		}

		t := schema.Transaction{
			ChainInfoID: ChainIDMap[chainID],
			BlockID:     list[txResp[i].Height].ID,
			Height:      txResp[i].Height,
			Code:        txResp[i].Code,
			Hash:        txResp[i].TxHash,
			Chunk:       chunk,
			Timestamp:   list[txResp[i].Height].Timestamp,
		}

		txs = append(txs, t)
	}

	return txs, nil
}

// getTxsChunk decodes transactions in a block and return a format of database transaction.
func (ex *Exporter) getRawTransactions(block *tmctypes.ResultBlock, txResps []*sdktypes.TxResponse) ([]schema.RawTransaction, error) {
	txChunk := make([]schema.RawTransaction, len(txResps), len(txResps))
	if len(txResps) <= 0 {
		return txChunk, nil
	}

	for i, txResp := range txResps {
		chunk, err := custom.AppCodec.MarshalJSON(txResp)
		if err != nil {
			log.Println(err)
			return txChunk, fmt.Errorf("failed to marshal tx : %s", err)
		}
		txChunk[i].ChainID = block.Block.ChainID
		txChunk[i].Height = txResp.Height
		txChunk[i].TxHash = txResp.TxHash
		txChunk[i].Chunk = chunk
		// show result
		// fmt.Println(jsonString[i])
	}

	return txChunk, nil
}

func (ex *Exporter) disassembleTransaction(txResps []*sdktypes.TxResponse) (uniqTransactionMessageAccounts []schema.TMA) {

	// 별도의 스키마 필요
	// id, account, hash, timestamp(불필요, 조인하면 되기 때문)
	// id, tx_hash(id는 불가능, db에 저장할 때 알 수 없음)
	if len(txResps) <= 0 {
		return nil
	}

	for _, txResp := range txResps {
		msgs := txResp.GetTx().GetMsgs()

		txHash := txResp.TxHash

		uniqueMsgAccount := make(map[string]map[string]struct{}) // tx 내 동일 메세지에 대한 유일한 어카운트 저장

		for _, msg := range msgs {

			msgType, accounts := types.AccountExporterFromCosmosTxMsg(&msg)
			// 어떤 msg 타입에 대해서도 signer를 이용해 accounts를 확보하면, 모든 메세지를 파싱할 수 있다.
			signers := getSignerAddress(msg.GetSigners())
			accounts = append(accounts, signers...)

			if msgType == "" {
				msgType, accounts = custom.AccountExporterFromCustomTxMsg(&msg, txHash)
			}

			for i := range accounts {
				ma, ok := uniqueMsgAccount[msgType]
				if !ok {
					ma = make(map[string]struct{})
					uniqueMsgAccount[msgType] = ma
				}
				ma[accounts[i]] = struct{}{}
			}
		} // end msgs for loop

		// msg 별 유일 어카운트 수집
		tma := parseTransactionMessageAccount(txHash, uniqueMsgAccount)
		uniqTransactionMessageAccounts = append(uniqTransactionMessageAccounts, tma...)
	} // 모든 tx 완료

	return uniqTransactionMessageAccounts
}

// msg - account 매핑 unique
func parseTransactionMessageAccount(txHash string, msgAccount map[string]map[string]struct{}) []schema.TMA {
	tma := make([]schema.TMA, 0)
	for msg := range msgAccount {
		for acc := range msgAccount[msg] {
			ta := schema.TMA{
				TxHash:         txHash,
				MsgType:        msg,
				AccountAddress: acc,
			}
			tma = append(tma, ta)
		}
	}
	return tma
}

func getSignerAddress(accAddrs []sdktypes.AccAddress) (address []string) {
	for _, addr := range accAddrs {
		address = append(address, addr.String())
	}

	return address
}
