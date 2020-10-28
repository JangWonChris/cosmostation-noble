package exporter

import (
	"encoding/json"
	"fmt"

	ceCodec "github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdktypestx "github.com/cosmos/cosmos-sdk/types/tx"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// getTxs decodes transactions in a block and return a format of database transaction.
func (ex *Exporter) getTxs(block *tmctypes.ResultBlock, txResps []*sdkTypes.TxResponse) ([]schema.Transaction, error) {
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
