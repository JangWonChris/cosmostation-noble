package main

import (
	// app "github.com/cosmostation/cosmostation-cosmos/api/wallet/api"
	"github.com/cosmostation/cosmostation-cosmos/api/wallet/config"
	"github.com/cosmostation/cosmostation-cosmos/api/wallet/exporter"
)

func main() {
	// Configuration in config.yaml
	config := config.NewAPIConfig()

	// API server app
	// app := &app.App{}
	// app.NewApp(config)
	// app.Run(":" + config.Web.Port)

	// Start syncing tasks using goroutines
	exporter := exporter.NewChainExporterService(config)
	exporter.OnStart()
}
