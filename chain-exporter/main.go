package main

import (
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/exporter"
	"go.uber.org/zap"
)

var (
	// Version is a project's version string.
	Version = ""

	// Commit is commit hash of this project.
	Commit = ""
)

func main() {
	zap.S().Info("Starting Chain Exporter...")
	zap.S().Infof("Version: %s Commit: %s", Version, Commit)

	exporter := exporter.NewExporter()
	exporter.Start()
}
