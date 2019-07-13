package main

import (
	"fmt"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/exporter"
)

func main() {
	fmt.Println("Stats Exporter")

	// Configuration in config.yaml
	config := config.NewConfig()

	// Start exporting data from blockchain
	exporter := exporter.NewStatsExporterService(config)
	exporter.OnStart()
}
