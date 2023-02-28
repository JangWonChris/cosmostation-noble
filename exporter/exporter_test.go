package exporter

import (
	"os"
	"testing"

	"github.com/cosmostation/cosmostation-noble/app"
)

var (
	ex *Exporter
)

func TestMain(m *testing.M) {
	chainEx := app.NewApp("chain-exporter")
	ex = NewExporter(chainEx)

	os.Exit(m.Run())
}
