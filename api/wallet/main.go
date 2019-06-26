package main

import (
	app "github.com/cosmostation/cosmostation-cosmos/api/wallet/api"
	"github.com/cosmostation/cosmostation-cosmos/api/wallet/config"
)

func main() {
	// Configuration in config.yaml
	config := config.NewAPIConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":" + string(config.Web.Port))
}
