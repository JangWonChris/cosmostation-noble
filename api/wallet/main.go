package main

import (
	"github.com/cosmostation/cosmostation-cosmos/api/wallet/app"
	"github.com/cosmostation/cosmostation-cosmos/api/wallet/app/config"
)

func main() {

	// Deployment Type
	// config := config.GetMainnetConfig()
	config := config.GetMainnetDevConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":5000")
}
