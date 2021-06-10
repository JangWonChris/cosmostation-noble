package exporter

import (
	"fmt"
	"os"
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdktypestx "github.com/cosmos/cosmos-sdk/types"
	legacytx "github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/cosmostation/cosmostation-cosmos/app"
)

var (
	ex *Exporter
)

func TestMain(m *testing.M) {
	chainEx := app.NewApp("chain-exporter")
	ex = NewExporter(chainEx)

	os.Exit(m.Run())
}

func commonTxParser(txHash string) (*sdkTypes.TxResponse, sdktypestx.Tx, error) {
	txResponse, err := ex.Client.CliCtx.GetTx(txHash)
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
