package exporter

import (
	"fmt"
	"os"
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdktypestx "github.com/cosmos/cosmos-sdk/types"
	legacytx "github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/cosmostation/cosmostation-cosmos/chain-config/custom"
	mdschema "github.com/cosmostation/mintscan-database/schema"
	"go.uber.org/zap"
)

var (
	ex *Exporter
)

func TestMain(m *testing.M) {
	ex = NewExporter(BASIC_MODE)

	os.Exit(m.Run())
}

func commonTxParser(txHash string) (*sdkTypes.TxResponse, sdktypestx.Tx, error) {
	txResponse, err := ex.client.CliCtx.GetTx(txHash)
	if err != nil {
		return &sdkTypes.TxResponse{}, legacytx.StdTx{}, err
	}

	if txResponse.Empty() {
		return &sdkTypes.TxResponse{}, legacytx.StdTx{}, nil
	}

	stdTx, ok := txResponse.GetTx().(sdktypestx.Tx)
	// stdTx, ok := txResponse.Tx.(legacytx.StdTx)
	if !ok {
		return &sdkTypes.TxResponse{}, legacytx.StdTx{}, fmt.Errorf("unsupported tx type: %s", txResponse.Tx)
	}

	return txResponse, stdTx, nil
}

func TestReproducePowerEventHistory(t *testing.T) {
	// Query latest block height saved in database
	dbHeight, err := ex.db.QueryLatestBlockHeight(ChainIDMap[ChainID])
	if dbHeight == -1 {
		fmt.Errorf("unexpected error in database: %s", err)
		return
	}

	zap.S().Infof("dst db %d \n", dbHeight)

	i := int64(1)
	endHeight := int64(200)
	for ; i <= dbHeight; i += endHeight {
		zap.S().Info("working height : ", i)
		rawTxs, err := ex.db.QueryTxForPowerEventHistory(i, i+endHeight) // 1 <= x < 201, 201 <= x < 201+200
		if err != nil {
			zap.S().Info("get query error : ", err)
			return
		}

		rawTxsLen := len(rawTxs)
		if rawTxsLen > 0 {

			txs := make([]*sdktypes.TxResponse, len(rawTxs))
			for i := range rawTxs {
				zap.S().Info("height : ", i, ", num_txs : ")
				tx := new(sdktypes.TxResponse)
				if err := custom.AppCodec.UnmarshalJSON([]byte(rawTxs[i].Chunk), tx); err != nil {
					zap.S().Info("unmarshal error")
					return
				}
				txs[i] = tx
			}

			exportData := new(mdschema.BasicData)
			exportData.ValidatorsPowerEventHistory, err = ex.getPowerEventHistoryNew(txs)
			if err != nil {
				return
			}
			if err := ex.db.InsertExportedData(exportData); err != nil {
				return
			}

		}

	}
	return

}
