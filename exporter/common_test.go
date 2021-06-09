package exporter

import (
	"fmt"
	"os"
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdktypestx "github.com/cosmos/cosmos-sdk/types"
	legacytx "github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
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
