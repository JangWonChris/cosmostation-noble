package main

import (
	"github.com/cosmostation/cosmostation-cosmos/alarm-notification/config"
	"github.com/cosmostation/cosmostation-cosmos/alarm-notification/exporter"
)

func main() {
	// Configuration in config.yaml
	config := config.NewConfig()

	// Start exporting data from blockchain
	exporter := exporter.NewChainExporterService(config)
	exporter.OnStart()
}
