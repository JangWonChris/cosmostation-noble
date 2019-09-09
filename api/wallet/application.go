package main

import (
	app "github.com/cosmostation/cosmostation-cosmos/api/wallet/api"
	"github.com/cosmostation/cosmostation-cosmos/api/wallet/config"
)

func main() {
	// configuration for this app in config.yaml
	config := config.NewConfig()

	// 동일한 main에서 API서버와 exporter는 돌리려면 따로 쓰레드로 돌려야 될 것 같다.
	// starting the server
	app := &app.App{}
	app.NewApp(config)
	app.Run(":" + config.Web.Port)
}
