package exporter

import (
	"encoding/json"
	"fmt"
	"log"

	ceCodec "github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdktypestx "github.com/cosmos/cosmos-sdk/types/tx"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// getTxs decodes transactions in a block and return a format of database transaction.
func (ex *Exporter) getTxs(block *tmctypes.ResultBlock, txResps []*sdktypes.TxResponse) ([]schema.Transaction, error) {
	txs := make([]schema.Transaction, 0)

	if len(txResps) <= 0 {
		return txs, nil
	}

	for _, txResp := range txResps {
		txI := txResp.GetTx()
		tx, ok := txI.(*sdktypestx.Tx)
		if !ok {
			return txs, fmt.Errorf("unsupported type")
		}

		// for _, msg := range getMsgs {
		// 	msgz, err := ceCodec.AppCodec.MarshalJSON(msg)
		// }
		msgsBz, err := ceCodec.AppCodec.MarshalJSON(tx.GetBody())
		if err != nil {
			return txs, fmt.Errorf("failed to unmarshal transaction messages: %s", err)
		}

		feeBz, err := ceCodec.AppCodec.MarshalJSON(tx.GetAuthInfo().GetFee())
		if err != nil {
			return txs, fmt.Errorf("failed to unmarshal tx fee: %s", err)
		}

		type SIG struct {
			Signatures []byte
		}

		tx.GetSigners()

		sigs := make([]SIG, len(tx.GetSignatures()), len(tx.GetSignatures()))
		for i, s := range tx.GetSignatures() {
			sigs[i].Signatures = s
		}

		sigsBz, err := json.Marshal(sigs)
		if err != nil {
			return txs, fmt.Errorf("failed to unmarshal tx signatures: %s", err)
		}

		logsBz, err := json.Marshal(txResp.Logs.String())
		if err != nil {
			return txs, fmt.Errorf("failed to unmarshal tx logs: %s", err)
		}

		t := &schema.Transaction{
			ChainID:    block.Block.ChainID,
			Height:     txResp.Height,
			Code:       txResp.Code,
			TxHash:     txResp.TxHash,
			GasWanted:  txResp.GasWanted,
			GasUsed:    txResp.GasUsed,
			Messages:   string(msgsBz),
			Fee:        string(feeBz),
			Signatures: string(sigsBz),
			Logs:       string(logsBz),
			Memo:       tx.GetBody().Memo,
			Timestamp:  txResp.Timestamp,
		}

		txs = append(txs, *t)
	}

	return txs, nil
}

// getTxsChunk decodes transactions in a block and return a format of database transaction.
func (ex *Exporter) getTxsJSONChunk(txResps []*sdktypes.TxResponse) ([]schema.TransactionJSONChunk, error) {
	txChunk := make([]schema.TransactionJSONChunk, len(txResps), len(txResps))
	if len(txResps) <= 0 {
		return txChunk, nil
	}

	for i, txResp := range txResps {
		chunk, err := ceCodec.AppCodec.MarshalJSON(txResp)
		if err != nil {
			log.Println(err)
			return txChunk, fmt.Errorf("failed to marshal tx : %s", err)
		}
		txChunk[i].Height = txResp.Height
		txChunk[i].Chunk = string(chunk)
		// show result
		// fmt.Println(jsonString[i])
	}

	return txChunk, nil
}
