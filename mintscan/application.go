package main

import (
	app "github.com/cosmostation/cosmostation-cosmos/mintscan/api"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
)

func main() {
	// configuration for this app in config.yaml
	config := config.NewConfig()

	// starting the server
	a := app.NewApp(config)
	a.Run(":" + config.Web.Port)
}
