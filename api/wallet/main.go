package main

import (
	"github.com/cosmostation/cosmostation-cosmos/api/wallet/config"
	"github.com/cosmostation/cosmostation-cosmos/api/wallet/exporter"
)

func main() {
	// Configuration in config.yaml
	config := config.NewAPIConfig()

	// Start syncing tasks using goroutines
	exporter := exporter.NewChainExporterService(config)
	exporter.OnStart()

	// 문제발생!
	// 동일한 main에서 API서버와 exporter를 돌릴 수 없는 것 같다
	// API server app
	// app := &app.App{}
	// app.NewApp(config)
	// app.Run(":" + config.Web.Port)
}
