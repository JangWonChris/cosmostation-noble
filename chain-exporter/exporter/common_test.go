package exporter

import (
	"fmt"
	"os"
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

var (
	ex *Exporter
)

func TestMain(m *testing.M) {
	ex = NewExporter()

	os.Exit(m.Run())
}

func commonTxParser(txHash string) (sdkTypes.TxResponse, auth.StdTx, error) {
	txResponse, err := ex.client.GetTx(txHash)
	if err != nil {
		return sdkTypes.TxResponse{}, auth.StdTx{}, err
	}

	if txResponse.Empty() {
		return sdkTypes.TxResponse{}, auth.StdTx{}, nil
	}

	stdTx, ok := txResponse.Tx.(auth.StdTx)
	if !ok {
		return sdkTypes.TxResponse{}, auth.StdTx{}, fmt.Errorf("unsupported tx type: %s", txResponse.Tx)
	}

	return txResponse, stdTx, nil
}
