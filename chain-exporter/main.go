package main

import (
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/exporter"
)

func main() {
	// Configuration in config.yaml
	config := config.NewConfig()

	// Start exporting data from blockchain
	exporter := exporter.NewChainExporter(config)
	exporter.OnStart()
}
