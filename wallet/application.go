package main

import (
	app "github.com/cosmostation/cosmostation-cosmos/wallet/api"
	"github.com/cosmostation/cosmostation-cosmos/wallet/api/config"
)

func main() {
	config := config.NewConfig()

	app := &app.App{}
	app.NewApp(config)
	app.Run(":" + config.Web.Port)
}
