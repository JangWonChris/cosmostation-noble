package main

import (
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter-revive/exporter"
)

func main() {
	exporter := exporter.NewExporter()
	exporter.Start()
}
