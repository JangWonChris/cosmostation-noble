package main

import (
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/exporter"
)

func main() {
	exporter := exporter.NewExporter()
	exporter.Start()
}
