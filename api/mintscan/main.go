package main

import (
	app "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
)

func main() {
	// Configuration in config.yaml
	config := config.NewConfig()

	// API server app
	app := &app.App{}
	app.NewApp(config)
	app.Run(":" + config.Web.Port)
}
