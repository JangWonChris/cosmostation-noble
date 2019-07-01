package api

import (
	"fmt"
	"log"
	"net/http"

	gaiaApp "github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/controllers"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/databases"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/rpc/client"
)

// App wraps up the required variables that are needed in this app
type App struct {
	Codec     *codec.Codec
	Config    *config.Config
	DB        *pg.DB
	Router    *mux.Router
	RPCClient *client.HTTP
}

// NewApp initializes the app with predefined configuration
func (a *App) NewApp(config *config.Config) {
	// Configuration
	a.Config = config

	// Connect to Tendermint RPC client through websocket
	a.RPCClient = client.NewHTTP(a.Config.Node.GaiadURL, "/websocket")

	// Connect to PostgreSQL
	a.DB = databases.ConnectDatabase(config)

	// Register Cosmos SDK codecs
	a.Codec = gaiaApp.MakeCodec()

	// Register routers
	a.setRouters()
}

// Sets the all required routers
func (a *App) setRouters() {
	a.Router = mux.NewRouter()
	a.Router = a.Router.PathPrefix("/v1").Subrouter()
	controllers.AccountController(a.Router, a.RPCClient, a.DB, a.Config)
	controllers.BlockController(a.Router, a.RPCClient, a.DB, a.Config)
	controllers.DistributionController(a.Config, a.DB, a.Router, a.RPCClient)
	controllers.GovernanceController(a.Router, a.RPCClient, a.DB, a.Config)
	controllers.MintingController(a.Router, a.RPCClient, a.DB, a.Config)
	controllers.TransactionController(a.Router, a.RPCClient, a.DB, a.Codec, a.Config)
	controllers.ValidatorController(a.Router, a.RPCClient, a.DB, a.Codec, a.Config)
	controllers.StatusController(a.Router, a.RPCClient, a.DB, a.Config)
	controllers.StatsController(a.Router, a.RPCClient, a.DB, a.Config)
	// a.Router.NotFoundHandler = http.HandleFunc("/", notFound)
}

// Run the app
func (a *App) Run(host string) {
	fmt.Print("Server is Starting on http://localhost", host, "\n")
	log.Fatal(http.ListenAndServe(host, a.Router))
}
