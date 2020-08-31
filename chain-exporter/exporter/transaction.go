package exporter

import (
	"fmt"

	ceCodec "github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// getTxs decodes transactions in a block and return a format of database transaction.
func (ex *Exporter) getTxs(block *tmctypes.ResultBlock, txResp []*sdk.TxResponse) ([]schema.Transaction, error) {
	txs := make([]schema.Transaction, 0)

	if len(txResp) <= 0 {
		return txs, nil
	}

	for _, tx := range txResp {
		stdTx, ok := tx.Tx.(auth.StdTx)
		if !ok {
			return txs, fmt.Errorf("unsupported tx type: %s", tx.Tx)
		}

		msgsBz, err := ceCodec.Codec.MarshalJSON(stdTx.GetMsgs())
		if err != nil {
			return txs, fmt.Errorf("failed to unmarshal transaction messages: %s", err)
		}

		feeBz, err := ceCodec.Codec.MarshalJSON(stdTx.Fee)
		if err != nil {
			return txs, fmt.Errorf("failed to unmarshal tx fee: %s", err)
		}

		sigs := make([]auth.StdSignature, len(stdTx.GetSignatures()), len(stdTx.GetSignatures()))
		for i, s := range stdTx.GetSignatures() {
			sigs[i].Signature = s.Signature
			// }
			// for i, pk := range stdTx.GetPubKeys() {
			// sigs[i].PubKey = pk
			sigs[i].PubKey = s.PubKey
		}

		sigsBz, err := ceCodec.Codec.MarshalJSON(sigs)
		if err != nil {
			return txs, fmt.Errorf("failed to unmarshal tx signatures: %s", err)
		}

		logsBz, err := ceCodec.Codec.MarshalJSON(tx.Logs)
		if err != nil {
			return txs, fmt.Errorf("failed to unmarshal tx logs: %s", err)
		}

		t := &schema.Transaction{
			ChainID:    block.Block.ChainID,
			Height:     tx.Height,
			Code:       tx.Code,
			TxHash:     tx.TxHash,
			GasWanted:  tx.GasWanted,
			GasUsed:    tx.GasUsed,
			Messages:   string(msgsBz),
			Fee:        string(feeBz),
			Signatures: string(sigsBz),
			Logs:       string(logsBz),
			Memo:       stdTx.GetMemo(),
			Timestamp:  tx.Timestamp,
		}

		txs = append(txs, *t)
	}

	return txs, nil
}
